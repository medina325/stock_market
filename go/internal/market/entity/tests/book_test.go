package entity

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/medina325/stock_market/go/internal/market/entity"
	"github.com/medina325/stock_market/go/internal/market/enums"
	"github.com/stretchr/testify/assert"
)

func TestAssetsTrading(t *testing.T) {
	a := entity.NewAsset(uuid.NewString(), "Asset 1", 1000)

	assetPosition := entity.NewInvestorAssetPosition(a.ID, 20)
	sellInvestor := entity.NewInvestor(uuid.NewString())
	sellInvestor.AddAssetPosition(assetPosition)

	inputChannel := make(chan *entity.Order)
	outputChannel := make(chan *entity.Order)
	wg := sync.WaitGroup{}

	book := entity.NewBook(inputChannel, outputChannel, &wg)

	// Acionar book.Trade()
	go book.Trade()

	wg.Add(1)

	// Criar orders de venda e compra
	o1 := entity.NewOrder(uuid.NewString(), sellInvestor, a, 20, 10, enums.Sell)

	buyInvestor := entity.NewInvestor(uuid.NewString())
	o2 := entity.NewOrder(uuid.NewString(), buyInvestor, a, 20, 10, enums.Buy)

	// Alimentar channel de entrada com Orders
	inputChannel <- o2
	inputChannel <- o1

	wg.Wait()

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

func TestDifferentAssetsTrading(t *testing.T) {
	asset1 := entity.NewAsset(uuid.NewString(), "Asset 1", 750)
	assetPosition1 := entity.NewInvestorAssetPosition(asset1.ID, 10)
	buyInvestor := entity.NewInvestor(uuid.NewString())
	buyInvestor.AddAssetPosition(assetPosition1)

	asset2 := entity.NewAsset(uuid.NewString(), "Asset 2", 650)
	assetPosition2 := entity.NewInvestorAssetPosition(asset2.ID, 13)
	sellInvestor := entity.NewInvestor(uuid.NewString())
	sellInvestor.AddAssetPosition(assetPosition2)

	orderChanIn := make(chan *entity.Order)
	orderChanOut := make(chan *entity.Order)
	wg := sync.WaitGroup{}

	book := entity.NewBook(orderChanIn, orderChanOut, &wg)
	go book.Trade()

	buyOrder := entity.NewOrder(uuid.NewString(), buyInvestor, asset1, 5, 10, enums.Buy)
	orderChanIn <- buyOrder

	sellOrder := entity.NewOrder(uuid.NewString(), sellInvestor, asset2, 3, 10, enums.Sell)
	orderChanIn <- sellOrder

	// realizar asserts
	assert := assert.New(t)
	// status das ordens estão open
	assert.Equal(enums.Open, buyOrder.Status, "Buy order should still be open")
	assert.Equal(enums.Open, sellOrder.Status, "Sell order should still be open")

	assert.Equal(5, buyOrder.PendingShares, "Buy order should still have 5 pending shares")
	assert.Equal(3, sellOrder.PendingShares, "Sell order should still have 3 pending shares")

	assert.Equal(10, assetPosition1.Shares, "Asset position 1 should still have the same 10 shares")
	assert.Equal(13, assetPosition2.Shares, "Asset position 2 should still have the same 13 shares")

	assert.Equal(0, buyOrder.TransactionsCount(), "There should be no transactions for buy order")
	assert.Equal(0, sellOrder.TransactionsCount(), "There should be no transactions for sell order")
}

func TestPartialMatchingAssetsTrading(t *testing.T) {
	a := entity.NewAsset(uuid.NewString(), "Asset 1", 200)

	buyInvestor := entity.NewInvestor(uuid.NewString())
	sellInvestor := entity.NewInvestor(uuid.NewString())

	sellAssetPosition := entity.NewInvestorAssetPosition(a.ID, 10)
	sellInvestor.AddAssetPosition(sellAssetPosition)

	chanIn := make(chan *entity.Order)
	chanOut := make(chan *entity.Order)
	wg := sync.WaitGroup{}

	wg.Add(1)
	book := entity.NewBook(chanIn, chanOut, &wg)
	go book.Trade()

	buyOrder := entity.NewOrder(uuid.NewString(), buyInvestor, a, 10, 5, enums.Buy)
	chanIn <- buyOrder
	sellOrder := entity.NewOrder(uuid.NewString(), sellInvestor, a, 8, 5, enums.Sell)
	chanIn <- sellOrder

	// Não faz sentido imediatamente falar para esperar, tenho que rodar algo antes
	wg.Wait()

	assert := assert.New(t)

	assert.Equal(8, buyInvestor.GetAssetPosition(a.ID).Shares, "Buy investor should have 8 shares")
	assert.Equal(2, sellInvestor.GetAssetPosition(a.ID).Shares, "Sell investor should have 2 shares (since he had 10 and sold 8)")

	assert.Equal(2, buyOrder.PendingShares, "Buy order should have 2 pending shares")
	assert.Equal(0, sellOrder.PendingShares, "Sell order should have 0 pending shares")

	assert.Equal(enums.Open, buyOrder.Status, "Buy order should still be open")
	assert.Equal(enums.Closed, sellOrder.Status, "Sell order should be closed")

	assert.Equal(40.0, buyOrder.Transactions[0].Total, "Transaction value of buy order should be of 40.00")
	assert.Equal(40.0, sellOrder.Transactions[0].Total, "Transaction value of sell order should be of 40.00")
}

func TestMultipleMatches(t *testing.T) {
	a := entity.NewAsset(uuid.NewString(), "Asset 1", 200)

	buyInvestor := entity.NewInvestor(uuid.NewString())
	sellInvestor := entity.NewInvestor(uuid.NewString())

	sellAssetPosition := entity.NewInvestorAssetPosition(a.ID, 10)
	sellInvestor.AddAssetPosition(sellAssetPosition)

	chanIn := make(chan *entity.Order)
	chanOut := make(chan *entity.Order)
	wg := sync.WaitGroup{}

	wg.Add(2)
	book := entity.NewBook(chanIn, chanOut, &wg)
	go book.Trade()

	buyOrder := entity.NewOrder(uuid.NewString(), buyInvestor, a, 10, 5, enums.Buy)
	chanIn <- buyOrder
	sellOrder1 := entity.NewOrder(uuid.NewString(), sellInvestor, a, 5, 5, enums.Sell)
	chanIn <- sellOrder1

	go func() {
		for range chanOut {
		}
	}()

	sellOrder2 := entity.NewOrder(uuid.NewString(), sellInvestor, a, 5, 5, enums.Sell)
	chanIn <- sellOrder2

	go func() {
		for range chanOut {
		}
	}()

	wg.Wait()

	assert := assert.New(t)

	assert.Equal(10, buyInvestor.GetAssetPosition(a.ID).Shares, "Buy investor should have 10 shares")
	assert.Equal(0, sellInvestor.GetAssetPosition(a.ID).Shares, "Sell investor should have 0 shares")

	assert.Equal(0, buyOrder.PendingShares, "Buy order should have 0 pending shares")
	assert.Equal(0, sellOrder1.PendingShares, "Sell order should have 0 pending shares")
	assert.Equal(0, sellOrder2.PendingShares, "Sell order should have 0 pending shares")

	assert.Equal(enums.Closed, buyOrder.Status, "Buy order should be closed")
	assert.Equal(enums.Closed, sellOrder2.Status, "Sell order 2 should be closed")
	assert.Equal(enums.Closed, sellOrder1.Status, "Sell order 1 should be closed")

	assert.Equal(2, buyOrder.TransactionsCount(), "Buy order should have 2 transactions")
	assert.Equal(1, sellOrder1.TransactionsCount(), "Sell order 1 should have 1 transaction")
	assert.Equal(1, sellOrder2.TransactionsCount(), "Sell order 2 should have 1 transaction")

	assert.Equal(25.0, buyOrder.Transactions[0].Total, "Transaction value of buy order should be of 25.00")
	assert.Equal(25.0, buyOrder.Transactions[1].Total, "Transaction value of buy order should be of 25.00")
	assert.Equal(25.0, sellOrder1.Transactions[0].Total, "Transaction value of sell order should be of 25.00")
	assert.Equal(25.0, sellOrder2.Transactions[0].Total, "Transaction value of sell order should be of 25.00")
}

func TestMultiplePartialMatches(t *testing.T) {
	a := entity.NewAsset(uuid.NewString(), "Asset 1", 200)

	buyInvestor := entity.NewInvestor(uuid.NewString())
	sellInvestor := entity.NewInvestor(uuid.NewString())

	sellAssetPosition := entity.NewInvestorAssetPosition(a.ID, 10)
	sellInvestor.AddAssetPosition(sellAssetPosition)

	chanIn := make(chan *entity.Order)
	chanOut := make(chan *entity.Order)
	wg := sync.WaitGroup{}

	wg.Add(1)
	book := entity.NewBook(chanIn, chanOut, &wg)
	go book.Trade()

	buyOrder := entity.NewOrder(uuid.NewString(), buyInvestor, a, 10, 5, enums.Buy)
	chanIn <- buyOrder
	sellOrder1 := entity.NewOrder(uuid.NewString(), sellInvestor, a, 5, 5, enums.Sell)
	chanIn <- sellOrder1

	go func() {
		for range chanOut {
		}
	}()

	wg.Wait()

	assert := assert.New(t)

	assert.Equal(5, buyInvestor.GetAssetPosition(a.ID).Shares, "Buy investor should have 5 shares")
	assert.Equal(5, sellInvestor.GetAssetPosition(a.ID).Shares, "Sell investor should have 5 shares")

	assert.Equal(5, buyOrder.PendingShares, "Buy order should have 5 pending shares")
	assert.Equal(0, sellOrder1.PendingShares, "Sell order 1 should have 0 pending shares")

	assert.Equal(enums.Open, buyOrder.Status, "Buy order should be open")
	assert.Equal(enums.Closed, sellOrder1.Status, "Sell order 1 should be closed")

	assert.Equal(1, buyOrder.TransactionsCount(), "Buy order should have 1 transactions")
	assert.Equal(1, sellOrder1.TransactionsCount(), "Sell order 1 should have 1 transaction")

	assert.Equal(25.0, buyOrder.Transactions[0].Total, "Transaction value of buy order should be of 25.00")
	assert.Equal(25.0, sellOrder1.Transactions[0].Total, "Transaction value of sell order should be of 25.00")

	wg.Add(1)
	sellOrder2 := entity.NewOrder(uuid.NewString(), sellInvestor, a, 5, 5, enums.Sell)
	chanIn <- sellOrder2
	wg.Wait()

	assert.Equal(10, buyInvestor.GetAssetPosition(a.ID).Shares, "Buy investor should have 10 shares")
	assert.Equal(0, sellInvestor.GetAssetPosition(a.ID).Shares, "Sell investor should have 0 shares")

	assert.Equal(0, buyOrder.PendingShares, "Buy order should have 0 pending shares")
	assert.Equal(0, sellOrder2.PendingShares, "Sell order 2 should have 0 pending shares")

	assert.Equal(enums.Closed, buyOrder.Status, "Buy order should be closed")
	assert.Equal(enums.Closed, sellOrder2.Status, "Sell order 2 should be closed")

	assert.Equal(2, buyOrder.TransactionsCount(), "Buy order should have 2 transactions")
	assert.Equal(1, sellOrder2.TransactionsCount(), "Sell order 2 should have 1 transaction")

	assert.Equal(25.0, buyOrder.Transactions[1].Total, "Transaction value of buy order should be of 25.00")
	assert.Equal(25.0, sellOrder2.Transactions[0].Total, "Transaction value of sell order should be of 25.00")
}

func TestNoMatchingOrders(t *testing.T) {
	a := entity.NewAsset(uuid.NewString(), "Asset 1", 200)

	buyInvestor := entity.NewInvestor(uuid.NewString())
	sellInvestor := entity.NewInvestor(uuid.NewString())

	sellAssetPosition := entity.NewInvestorAssetPosition(a.ID, 10)
	sellInvestor.AddAssetPosition(sellAssetPosition)

	chanIn := make(chan *entity.Order)
	chanOut := make(chan *entity.Order)
	wg := sync.WaitGroup{}

	book := entity.NewBook(chanIn, chanOut, &wg)
	go book.Trade()

	buyOrder := entity.NewOrder(uuid.NewString(), buyInvestor, a, 10, 5, enums.Buy)
	chanIn <- buyOrder
	sellOrder := entity.NewOrder(uuid.NewString(), sellInvestor, a, 10, 6, enums.Sell)
	chanIn <- sellOrder

	assert := assert.New(t)

	var nilAssetPosition *entity.InvestorAssetPosition = nil
	assert.Equal(nilAssetPosition, buyInvestor.GetAssetPosition(a.ID), "Buy investor should not have an asset position")
	assert.Equal(10, sellInvestor.GetAssetPosition(a.ID).Shares, "Sell investor should still have 10 shares")

	assert.Equal(10, buyOrder.PendingShares, "Buy order should have 10 pending shares")
	assert.Equal(10, sellOrder.PendingShares, "Sell order should have 10 pending shares")

	assert.Equal(enums.Open, buyOrder.Status, "Buy order should be open")
	assert.Equal(enums.Open, sellOrder.Status, "Sell order should be closed")

	assert.Equal(0, buyOrder.TransactionsCount(), "Buy order should have 0 transactions")
	assert.Equal(0, sellOrder.TransactionsCount(), "Sell order should have 0 transaction")
}
