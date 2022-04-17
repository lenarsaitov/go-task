package cards

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math"
	"strings"
	"sync/atomic"
)

type CardStorage struct {
	db atomic.Value
}

type QueryResult interface {
	Scan(dest ...interface{}) error
}

var (
	_ QueryResult = &sql.Rows{}
	_ QueryResult = &sql.Row{}
)

func NewCardStorage(db *sqlx.DB) *CardStorage {
	res := &CardStorage{}
	res.db.Store(db)
	return res
}

func (s *CardStorage) getDB() *sqlx.DB {
	return s.db.Load().(*sqlx.DB)
}

func (s *CardStorage) FindOne(ctx context.Context, cardID int) (*CardInfo, error) {
	query := `SELECT cards.card_id, users.user_id, users.user_full_name, cards.balance, cards.create_time
	          FROM cards INNER JOIN users
			  ON 	cards.user_id = users.user_id
	          WHERE cards.card_id = $1;`

	row := s.getDB().QueryRowContext(ctx, query, cardID)

	m, err := s.readCardInfo(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return m, nil
}

func (s *CardStorage) readCardInfo(r QueryResult) (*CardInfo, error) {
	cardInfo := &CardInfo{}
	err := r.Scan(&cardInfo.CardID, &cardInfo.UserID, &cardInfo.UserName, &cardInfo.Balance, &cardInfo.CreateTime)
	if err != nil {
		return nil, err
	}
	return cardInfo, nil
}

func (s *CardStorage) buildFindManyWhereClause(filter *FilterParams, pos int) (clause string, args []interface{}) {
	var predicates []string

	if filter.UserID != 0 {
		predicates = append(predicates,
			fmt.Sprintf("user_id = $%d", pos))
		args = append(args, filter.UserID)
		pos++
	}

	clause = strings.Join(predicates, " and ")
	if len(clause) > 0 {
		clause = "where " + clause
	}

	return
}

func (s *CardStorage) FindMany(ctx context.Context, filter *FilterParams) (*Pagination, error) {
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

	template := `SELECT cards.card_id, users.user_id, users.user_full_name, cards.balance, cards.create_time
	             FROM cards INNER JOIN users 
				 ON cards.user_id = users.user_id
				 %s
	             ORDER BY (cards.card_id)
	             LIMIT $1 OFFSET $2;`

	query := fmt.Sprintf(template, whereClause)
	rows, err := s.getDB().QueryContext(ctx, query, append(paginationArgs, whereArgs...)...)
	if err != nil {
		return nil, fmt.Errorf("Cant query cards: %w", err)
	}
	defer rows.Close()

	items := make([]*CardInfo, 0)
	for rows.Next() {
		userInfo, err := s.readCardInfo(rows)
		if err != nil {
			return nil, fmt.Errorf("Cannot read card info: %w", err)
		}
		items = append(items, userInfo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Query error: %w", err)
	}

	countWhereClause, countWhereArgs := s.buildFindManyWhereClause(filter, 1)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM cards %s;`, countWhereClause)
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

func (s *CardStorage) readUserInfo(r QueryResult) (*UserInfo, error) {
	deviceInfo := &UserInfo{}
	err := r.Scan(&deviceInfo.UserID, &deviceInfo.UserName, &deviceInfo.CreateTime)
	if err != nil {
		return nil, err
	}
	return deviceInfo, nil
}

func (s *CardStorage) isExistUser(ctx context.Context, id int) (bool, error) {
	query := `SELECT user_id, user_full_name, create_time
	          FROM users
	          WHERE user_id = $1;`

	row := s.getDB().QueryRowContext(ctx, query, id)

	_, err := s.readUserInfo(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *CardStorage) AddCardItem(ctx context.Context, req *AddCardRequestParams) (*int64, error) {
	row := s.getDB().QueryRowContext(ctx,
		`INSERT INTO cards
		              (user_id, balance)
		        VALUES ($1, $2)
		        RETURNING card_id;`, req.UserID, req.Balance)

	var requestID int64
	err := row.Scan(&requestID)
	if err != nil {
		return nil, err
	}
	return &requestID, nil
}

func (s *CardStorage) UpdateCardItem(ctx context.Context, req *UpdateCardRequestParams) error {
	tx, err := s.getDB().BeginTx(ctx, nil)
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

	cardRow := tx.QueryRowContext(ctx, `SELECT balance FROM cards WHERE card_id = $1 FOR UPDATE;`, req.CardID)
	var balance int
	err = cardRow.Scan(&balance)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE cards SET balance = $2 WHERE card_id = $1;`, req.CardID, req.Balance)
	if err != nil {
		return err
	}

	return nil
}

func (s *CardStorage) DeleteCardItem(ctx context.Context, cardID int) error {
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

	cardRow := tx.QueryRowContext(ctx, `SELECT balance FROM cards WHERE card_id = $1 FOR UPDATE;`, cardID)
	var balance int
	err = cardRow.Scan(&balance)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, "DELETE FROM cards WHERE card_id=$1", cardID)
	if err == nil {
		_, err := res.RowsAffected()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *CardStorage) RefillCard(ctx context.Context, req *RefillCardRequestParams) error {
	tx, err := s.getDB().BeginTx(ctx, nil)
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

	cardRow := tx.QueryRowContext(ctx, `SELECT balance FROM cards WHERE card_id = $1 FOR UPDATE;`, req.CardID)
	var balance int
	err = cardRow.Scan(&balance)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE cards SET balance = $2 WHERE card_id = $1;`, req.CardID, balance+req.AddBalance)
	if err != nil {
		return err
	}

	return nil
}

func (s *CardStorage) TransferBalanceCard(ctx context.Context, req *TransferBalanceCardRequestParams) error {
	tx, err := s.getDB().BeginTx(ctx, nil)
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

	cardRow := tx.QueryRowContext(ctx, `SELECT balance FROM cards WHERE card_id = $1 FOR UPDATE;`, req.CardFrom)
	var balanceFrom int
	err = cardRow.Scan(&balanceFrom)
	if err != nil {
		return err
	}

	cardRow = tx.QueryRowContext(ctx, `SELECT balance FROM cards WHERE card_id = $1 FOR UPDATE;`, req.CardTo)
	var balanceTo int
	err = cardRow.Scan(&balanceTo)
	if err != nil {
		return err
	}

	if balanceFrom-req.AddBalance < 0 {
		return errors.New("Dont exist essential balance")
	}

	_, err = tx.ExecContext(ctx, `UPDATE cards SET balance = $2 WHERE card_id = $1;`, req.CardFrom, balanceFrom-req.AddBalance)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE cards SET balance = $2 WHERE card_id = $1;`, req.CardTo, balanceTo+req.AddBalance)
	if err != nil {
		return err
	}

	return nil
}
