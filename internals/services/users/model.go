package users

type UserInfo struct {
	UserID     int64  `json:"user_id"`
	UserName   string `json:"user_full_name"`
	CreateTime string `json:"create_time"`
}

type CardInfo struct {
	CardID     int64  `json:"card_id"`
	Balance    int64  `json:"balance"`
	UserID     int64  `json:"user_id"`
	UserName   string `json:"user_full_name"`
	CreateTime string `json:"create_time"`
}

type Pagination struct {
	Page       int         `json:"page,omitempty"`
	Size       int         `json:"size,omitempty"`
	PagesCount int         `json:"pagesCount"`
	ItemsCount int         `json:"itemsCount"`
	Items      []*UserInfo `json:"items"`
}

type FilterParams struct {
	UserName string
	Page     int `validate:"gte=1"`
	Size     int `validate:"gte=1,lte=50"`
}

type AddUserRequestParams struct {
	UserName string
}

type UpdateUserRequestParams struct {
	UserID   int
	UserName string
}
