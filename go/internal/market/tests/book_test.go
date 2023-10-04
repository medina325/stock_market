package tests

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/medina325/stock_market/go/internal/market/entity"
	"github.com/medina325/stock_market/go/internal/market/enums"
	"github.com/stretchr/testify/assert"
)

func TestBuyAsset(t *testing.T) {
	sellInvestor := entity.NewInvestor(uuid.NewString())
	buyInvestor := entity.NewInvestor(uuid.NewString())

	a := entity.NewAsset(uuid.NewString(), "Asset 1", 1000)

	assetPosition := entity.NewInvestorAssetPosition(a.ID, 20)
	sellInvestor.AddAssetPosition(assetPosition)

	inputChannel := make(chan *entity.Order)
	outputChannel := make(chan *entity.Order)
	wg := sync.WaitGroup{}

	book := entity.NewBook(inputChannel, outputChannel, &wg)

	// Acionar book.Trade()
	go book.Trade()

	wg.Add(1)

	// Criar orders de venda e compra
	o1 := entity.NewOrder(
		uuid.NewString(), sellInvestor, a, 20, 10, enums.Sell,
	)
	o2 := entity.NewOrder(
		uuid.NewString(), buyInvestor, a, 20, 10, enums.Buy,
	)

	// Alimentar channel de entrada com Orders
	inputChannel <- o2
	inputChannel <- o1

	wg.Wait()

	// Fazer assert
	// no status das orders,
	// na posição do investor,
	// transações criadas
	assert := assert.New(t)

	assert.Equal(enums.Closed, o1.Status, "Order 1 should be closed")
	assert.Equal(enums.Closed, o2.Status, "Order 2 should be closed")
	assert.Equal(0, o1.PendingShares, "Order 1 should not have any pending shares")
	assert.Equal(0, o2.PendingShares, "Order 2 should not have any pending shares")

	assert.Equal(0, sellInvestor.GetAssetPosition(a.ID).Shares, "Sell Investor should have 0 shares")
	assert.Equal(20, buyInvestor.GetAssetPosition(a.ID).Shares, "Buy Investor should have 20 shares")

	assert.Equal(1, o1.TransactionsCount(), "There should be 1 transaction for Order 1")
	assert.Equal(1, o2.TransactionsCount(), "There should be 1 transaction for Order 2")
	assert.Equal(200.0, o1.Transactions[0].Total, "Transaction value of Order 1 should be of 200.00")
	assert.Equal(200.0, o2.Transactions[0].Total, "Transaction value of Order 1 should be of 200.00")
}
