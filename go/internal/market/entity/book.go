package entity

import (
	"container/heap"
	"sync"

	"github.com/medina325/stock_market/go/internal/market/enums"
)

// Book represents a bookkeeping structure.
type Book struct {
	Orders       []*Order
	Transactions []*Transaction
	OrdersChanIn chan *Order
	OrderChanOut chan *Order
	Wg           *sync.WaitGroup
}

func NewBook(orderChanIn chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Orders:       []*Order{},
		Transactions: []*Transaction{},
		OrdersChanIn: orderChanIn,
		OrderChanOut: orderChanOut,
		Wg:           wg,
	}
}

// Returns number of shares to be applied in a transaction, i.e., the number
// of shares to be removed from the selling order, and the number of shares
// to be added to the buying order.
// The number of shares need to be the minimum value between the number of shares
// being sold and the number of shares being bought.
func getTransactionShares(sellingShares, buyingShares int) int {
	if sellingShares < buyingShares {
		return sellingShares
	}
	return buyingShares
}

func (b *Book) Trade() {
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range b.OrdersChanIn {
		if order.OrderType == enums.Buy {
			buyOrders.Push(order)

			thereAreNoSellOrders := sellOrders.Len() == 0

			if thereAreNoSellOrders {
				continue
			}

			thereAreNoSellOrderMatch := (*sellOrders)[0].Price > order.Price
			orderMatchHasNoPendingShares := (*sellOrders)[0].PendingShares == 0

			if thereAreNoSellOrderMatch || orderMatchHasNoPendingShares {
				continue
			}

			sellOrder := sellOrders.Pop().(*Order)

			transactionShares := getTransactionShares(sellOrder.PendingShares, order.PendingShares)

			// Duvidas:
			// - como eu sei se passo o preço da order de compra ou venda
			transaction := NewTransaction(sellOrder, order, transactionShares, sellOrder.Price)
			b.ExecuteTransaction(transaction)

			sellOrder.AddTransaction(transaction)
			order.AddTransaction(transaction)

			b.OrderChanOut <- sellOrder
			b.OrderChanOut <- order

			// REFACTOR
			if sellOrder.PendingShares > 0 {
				sellOrders.Push(sellOrder)
			}
		} else if order.OrderType == enums.Sell {
			sellOrders.Push(order)

			thereAreNoBuyOrder := buyOrders.Len() == 0

			if thereAreNoBuyOrder {
				continue
			}

			thereAreNoBuyOrderMatch := order.Price > (*buyOrders)[0].Price
			orderMatchHasNoPendingShares := (*buyOrders)[0].PendingShares == 0

			if thereAreNoBuyOrderMatch || orderMatchHasNoPendingShares {
				continue
			}

			buyOrder := buyOrders.Pop().(*Order)

			transactionShares := getTransactionShares(order.PendingShares, buyOrder.PendingShares)
			transaction := NewTransaction(order, buyOrder, transactionShares, buyOrder.Price)

			b.ExecuteTransaction(transaction)

			buyOrder.AddTransaction(transaction)
			order.AddTransaction(transaction)

			b.OrderChanOut <- buyOrder
			b.OrderChanOut <- order

			if buyOrder.PendingShares > 0 {
				buyOrders.Push(buyOrder)
			}
		}
	}
}

func (b *Book) ExecuteTransaction(t *Transaction) {
	defer b.Wg.Done()

	t.UpdateSellOrderAssetPosition()
	t.LiquidateSellPendingShares(t.Shares)
	t.UpdateSellOrderStatus()

	t.UpdateBuyOrderAssetPosition()
	t.LiquidateBuyPendingShares(t.Shares)
	t.UpdateBuyOrderStatus()

	b.Transactions = append(b.Transactions, t)
}
