package cards

import (
	"context"
	"database/sql"
)

type CardService struct {
	storage *CardStorage
}

func NewCardService(storage *CardStorage) *CardService {
	return &CardService{storage: storage}
}

func (service *CardService) GetCard(c context.Context, cardID int) (*CardInfo, error) {
	return service.storage.FindOne(c, cardID)
}

func (service *CardService) GetListCards(c context.Context, params *FilterParams) (*Pagination, error) {
	return service.storage.FindMany(c, params)
}

func (service *CardService) AddCard(c context.Context, params *AddCardRequestParams) (*int64, error) {
	isExist, err := service.storage.isExistUser(c, params.UserID)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, err
	}

	return service.storage.AddCardItem(c, params)
}

func (service *CardService) UpdateCard(c context.Context, params *UpdateCardRequestParams) (bool, error) {
	err := service.storage.UpdateCardItem(c, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (service *CardService) DeleteCard(c context.Context, cardID int) (bool, error) {
	err := service.storage.DeleteCardItem(c, cardID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (service *CardService) RefillCard(c context.Context, refillParams *RefillCardRequestParams) (bool, error) {
	err := service.storage.RefillCard(c, refillParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (service *CardService) TransferBalanceCard(c context.Context, params *TransferBalanceCardRequestParams) (bool, error) {
	err := service.storage.TransferBalanceCard(c, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
