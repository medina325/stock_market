package entity

type OrderQueue []*Order

func (o OrderQueue) Len() int {
	return len(o)
}

func (o OrderQueue) Less(i, j int) bool {
	return o[i].Price < o[j].Price
}

func (o *OrderQueue) Swap(i, j int) {
	(*o)[i], (*o)[j] = (*o)[j], (*o)[i]
}

// Content of "o" receives its own content appended with "x".
// Notice that the interface method Push needs to be able
// to receive any data type, so "x" is generic - hence its
// interface{} type, hence our need to cast it to *Order
func (o *OrderQueue) Push(x interface{}) {
	*o = append(*o, x.(*Order))
}

func (o *OrderQueue) Pop() interface{} {
	old := *o
	length := len(old)

	// last := old[length - 1]
	// Removing last element
	*o = old[0 : length-1]

	return old[length-1]
}

// Praticar essas coisas de container/heap no meu repo de go

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}
