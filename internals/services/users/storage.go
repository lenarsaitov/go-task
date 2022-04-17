package users

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math"
	"strings"
	"sync/atomic"
)

type UserStorage struct {
	db atomic.Value
}

type QueryResult interface {
	Scan(dest ...interface{}) error
}

var (
	_ QueryResult = &sql.Rows{}
	_ QueryResult = &sql.Row{}
)

func NewUserStorage(db *sqlx.DB) *UserStorage {
	res := &UserStorage{}
	res.db.Store(db)
	return res
}

func (s *UserStorage) getDB() *sqlx.DB {
	return s.db.Load().(*sqlx.DB)
}

func (s *UserStorage) FindOne(ctx context.Context, id int) (*UserInfo, error) {
	query := `SELECT user_id, user_full_name, create_time
	          FROM users
	          WHERE user_id = $1;`

	row := s.getDB().QueryRowContext(ctx, query, id)

	m, err := s.readUserInfo(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return m, nil
}

func (s *UserStorage) readUserInfo(r QueryResult) (*UserInfo, error) {
	userInfo := &UserInfo{}
	err := r.Scan(&userInfo.UserID, &userInfo.UserName, &userInfo.CreateTime)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (s *UserStorage) buildFindManyWhereClause(filter *FilterParams, pos int) (clause string, args []interface{}) {
	var predicates []string

	if len(filter.UserName) != 0 {
		predicates = append(predicates,
			fmt.Sprintf("user_full_name = $%d", pos))
		args = append(args, filter.UserName)
		pos++
	}

	clause = strings.Join(predicates, " and ")
	if len(clause) > 0 {
		clause = "where " + clause
	}

	return
}

func (s *UserStorage) FindMany(ctx context.Context, filter *FilterParams) (*Pagination, error) {
	limit := filter.Size
	if limit <= 0 {
		limit = math.MaxInt32
	}

	offset := 0
	if filter.Page > 0 {
		offset = (filter.Page - 1) * filter.Size
	}

	paginationArgs := []interface{}{limit, offset}

	whereClause, whereArgs := s.buildFindManyWhereClause(filter, 3)

	template := `SELECT user_id, user_full_name, create_time
	             FROM users %s
	             ORDER BY (user_id)
	             LIMIT $1 OFFSET $2;`

	query := fmt.Sprintf(template, whereClause)
	rows, err := s.getDB().QueryContext(ctx, query, append(paginationArgs, whereArgs...)...)
	if err != nil {
		return nil, fmt.Errorf("Cant query services: %w", err)
	}
	defer rows.Close()

	items := make([]*UserInfo, 0)
	for rows.Next() {
		userInfo, err := s.readUserInfo(rows)
		if err != nil {
			return nil, fmt.Errorf("Cannot read user info: %w", err)
		}
		items = append(items, userInfo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Query error: %w", err)
	}

	countWhereClause, countWhereArgs := s.buildFindManyWhereClause(filter, 1)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users %s;`, countWhereClause)
	row := s.getDB().QueryRowContext(ctx, countQuery, countWhereArgs...)
	var count int
	if err := row.Scan(&count); err != nil {
		return nil, fmt.Errorf("Items count query error: %w", err)
	}

	pc := count / limit
	if count%limit > 0 {
		pc++
	}

	return &Pagination{
		Page:       filter.Page,
		Size:       filter.Size,
		PagesCount: pc,
		ItemsCount: count,
		Items:      items,
	}, nil
}

func (s *UserStorage) readCardsInfo(r QueryResult) (*CardInfo, error) {
	cardInfo := &CardInfo{}
	err := r.Scan(&cardInfo.CardID)
	if err != nil {
		return nil, err
	}
	return cardInfo, nil
}

func (s *UserStorage) isExistCards(ctx context.Context, userID int) (bool, error) {
	query := `SELECT card_id
	          FROM cards
	          WHERE user_id = $1;`

	row := s.getDB().QueryRowContext(ctx, query, userID)

	_, err := s.readCardsInfo(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *UserStorage) AddUserItem(ctx context.Context, req *AddUserRequestParams) (*int64, error) {
	row := s.getDB().QueryRowContext(ctx,
		`INSERT INTO users
		              (user_full_name)
		        VALUES ($1)
		        RETURNING user_id;`, req.UserName)

	var requestID int64
	err := row.Scan(&requestID)
	if err != nil {
		return nil, err
	}
	return &requestID, nil
}

func (s *UserStorage) UpdateUserItem(ctx context.Context, req *UpdateUserRequestParams) error {
	_, err := s.getDB().ExecContext(ctx, `UPDATE users SET user_full_name = $2 WHERE user_id = $1;`, req.UserID, req.UserName)

	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) DeleteUserItem(ctx context.Context, userID int) error {
	tx, err := s.getDB().BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("Cannot start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	res, err := tx.ExecContext(ctx, "DELETE FROM users WHERE user_id=$1", userID)
	if err == nil {
		_, err := res.RowsAffected()
		if err != nil {
			return err
		}
	}

	return nil
}
