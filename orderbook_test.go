package main

import (
	"fmt"
	"testing"
)

func TestLimit(t *testing.T) {
	l := NewLimit(10000)
	buyOrderA := NewOrder(true, 5) // 買
	buyOrderB := NewOrder(true, 8)
	buyOrderC := NewOrder(true, 10)

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.DeleteOrder(buyOrderB)
	fmt.Println(l)
}

func TestOrderbook(t *testing.T) {
	ob := NewOrderbook()

	buyOrder1 := NewOrder(true, 10)
	ob.PlaceOrder(18000, buyOrder1)
	sellOrder := NewOrder(false, 5)
	ob.PlaceOrder(18000, sellOrder)
	buyOrder2 := NewOrder(true, 20)
	ob.PlaceOrder(19000, buyOrder2)

	fmt.Println(ob.Bids) // buy
	fmt.Println(ob.Asks) // sell
}
