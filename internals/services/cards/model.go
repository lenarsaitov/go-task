package cards

type CardInfo struct {
	CardID     int    `json:"card_id"`
	Balance    int    `json:"balance"`
	UserID     int    `json:"user_id"`
	UserName   string `json:"user_full_name"`
	CreateTime string `json:"create_time"`
}

type UserInfo struct {
	UserID     int    `json:"user_id"`
	UserName   string `json:"user_full_name"`
	CreateTime string `json:"create_time"`
}

type Pagination struct {
	Page       int         `json:"page,omitempty"`
	Size       int         `json:"size,omitempty"`
	PagesCount int         `json:"pagesCount"`
	ItemsCount int         `json:"itemsCount"`
	Items      []*CardInfo `json:"items"`
}

type FilterParams struct {
	UserID int
	Page   int `validate:"gte=1"`
	Size   int `validate:"gte=1,lte=50"`
}

type AddCardRequestParams struct {
	UserID  int
	Balance int
}

type UpdateCardRequestParams struct {
	CardID  int
	Balance int
}

type RefillCardRequestParams struct {
	CardID     int
	AddBalance int
}

type TransferBalanceCardRequestParams struct {
	CardFrom   int
	CardTo     int
	AddBalance int
}
