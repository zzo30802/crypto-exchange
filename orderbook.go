package main

import (
	"fmt"
	"sort"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

// 一筆訂單
type Order struct {
	Size      float64 // 數量
	Bid       bool    //buy or sell
	Limit     *Limit  // 限價
	Timestamp int64   //下單時間，用來決定同價位的先後順序
}

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size: size, //數量
		Bid:  bid,  //true->buy, false->sell
		// 下單時間，因為同價格可能會有多個人的訂單，以時間為撮合的優先順序
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("[size: %.2f]", o.Size)
}

// 代表price這個價位有多少Order集合
// buy和sell的Order並不會放在同個Orders中
// Orderbook中會分開存放這兩種訂單
// Order.Bid的用途是辨認新訂單是BIds or Asks
type Limit struct {
	Price       float64 //價位
	Orders      Orders  //掛在這個價位的所有訂單
	TotalVolume float64 //這個價位的總數量
}

type Limits []*Limit

// for Asks
type ByBestAsk struct{ Limits }

func (a ByBestAsk) Len() int           { return len(a.Limits) }
func (a ByBestAsk) Swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestAsk) Less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

// for Bids
type ByBestBid struct{ Limits }

func (b ByBestBid) Len() int           { return len(b.Limits) }
func (b ByBestBid) Swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestBid) Less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }

func (l *Limit) String() string {
	return fmt.Sprintf("[price: %.2f | volume: %.2f]", l.Price, l.TotalVolume)
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrder(o *Order) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
			o.Limit = nil
			l.TotalVolume -= o.Size
			break
		}
	}
	// o.Limit = nil
	// l.TotalVolume -= o.Size
	sort.Sort(l.Orders)
}

type Orderbook struct {
	Asks []*Limit // 想賣的人，開價由低到高
	Bids []*Limit // 想買的人，開價由高到低

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Asks:      []*Limit{},
		Bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (ob *Orderbook) PlaceLimitOrder(price float64, o *Order) {
	var limit *Limit
	if o.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}
	if limit == nil {
		limit = NewLimit(price)
		if o.Bid {
			ob.Bids = append(ob.Bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.Asks = append(ob.Asks, limit)
			ob.AskLimits[price] = limit
		}
	}

	limit.AddOrder(o)
}
