package entity

import "github.com/medina325/stock_market/go/internal/market/enums"

type Order struct {
	ID            string
	Investor      *Investor
	Asset         *Asset
	Shares        int
	PendingShares int
	Price         float64
	OrderType     int
	Status        int
	Transactions  []*Transaction
}

func (o *Order) NewOrder(orderID string, investor *Investor, asset *Asset, shares int, price float64, orderType int) *Order {
	return &Order{
		ID:            orderID,
		Investor:      investor,
		Asset:         asset,
		Shares:        shares,
		PendingShares: shares,
		Price:         price,
		OrderType:     orderType,
		Status:        enums.Open,
		Transactions:  []*Transaction{},
	}
}
