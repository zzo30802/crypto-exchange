package main

import (
	"fmt"
	"reflect"
	"testing"
)

// func assert(t *testing.T, a, b any) {
// 	t.Helper()
// 	if !reflect.DeepEqual(a, b) {
// 		t.Errorf("%+v != %+v", a, b)
// 	}
// }

func assertEqual(t *testing.T, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

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

// func TestOrderbook(t *testing.T) {
// 	ob := NewOrderbook()

// 	buyOrder1 := NewOrder(true, 10)
// 	ob.PlaceOrder(18000, buyOrder1)
// 	sellOrder := NewOrder(false, 5)
// 	ob.PlaceOrder(18000, sellOrder)
// 	buyOrder2 := NewOrder(true, 20)
// 	ob.PlaceOrder(19000, buyOrder2)

// 	fmt.Println(ob.Bids) // buy
// 	fmt.Println(ob.Asks) // sell
// }

func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrderA := NewOrder(false, 10)
	sellOrderB := NewOrder(false, 5)
	ob.PlaceLimitOrder(10000, sellOrderA)
	ob.PlaceLimitOrder(9000, sellOrderB)

	buyOrderA := NewOrder(true, 5)
	ob.PlaceLimitOrder(10000, buyOrderA)

	// 判斷sell limit數量
	assertEqual(t, len(ob.Asks), 2)
	// 判斷buy limit數量
	assertEqual(t, len(ob.Bids), 0)
}
