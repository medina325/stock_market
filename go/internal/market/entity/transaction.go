package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/medina325/stock_market/go/internal/market/enums"
)

type Transaction struct {
	ID           string
	SellingOrder *Order
	BuyingOrder  *Order
	Shares       int
	Price        float64
	Total        float64
	DateTime     time.Time
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	total := price * float64(shares)

	return &Transaction{
		ID:           uuid.New().String(),
		SellingOrder: sellingOrder,
		BuyingOrder:  buyingOrder,
		Shares:       shares,
		Price:        price,
		Total:        total,
		DateTime:     time.Now(),
	}
}

func (t *Transaction) LiquidateBuyPendingShares() {
	t.BuyingOrder.PendingShares -= t.Shares
}

func (t *Transaction) LiquidateSellPendingShares() {
	t.SellingOrder.PendingShares -= t.Shares
}

func (t *Transaction) UpdateSellOrderAssetPosition() {
	t.SellingOrder.Investor.UpdateAssetPosition(t.SellingOrder.Asset.ID, -t.Shares)
}

func (t *Transaction) UpdateBuyOrderAssetPosition() {
	t.BuyingOrder.Investor.UpdateAssetPosition(t.BuyingOrder.Asset.ID, t.Shares)
}

func (t *Transaction) UpdateBuyOrderStatus() {
	if t.BuyingOrder.PendingShares == 0 {
		t.BuyingOrder.Status = enums.Closed
	}
}

func (t *Transaction) UpdateSellOrderStatus() {
	if t.SellingOrder.PendingShares == 0 {
		t.SellingOrder.Status = enums.Closed
	}
}
