package order_book

import (
	"crypto-balance-service/transaction_log"
	"fmt"
	"sort"
	"sync"
	"time"
)

type Order struct {
	ID             string
	UserID         string
	Cryptocurrency string
	Type           string // "buy" or "sell"
	Price          float64
	Amount         float64
	Status         string // "open" or "filled"
}

type OrderBook struct {
	buys      map[string]Order
	sells     map[string]Order
	buyPrice  []Order
	sellPrice []Order
	mutex     sync.RWMutex
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		buys:      make(map[string]Order),
		sells:     make(map[string]Order),
		buyPrice:  make([]Order, 0),
		sellPrice: make([]Order, 0),
	}
}

func (ob *OrderBook) PlaceOrder(order Order) error {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	if order.Type != "buy" && order.Type != "sell" {
		return fmt.Errorf("Invalid order type: %s", order.Type)
	}

	if order.Amount <= 0 {
		return fmt.Errorf("Order amount must be greater than zero")
	}

	if order.Price <= 0 {
		return fmt.Errorf("Order price must be greater than zero")
	}

	// Place the order in the appropriate map
	if order.Type == "buy" {
		if _, ok := ob.buys[order.ID]; ok {
			return fmt.Errorf("Order with ID %s already exists", order.ID)
		}
		ob.buys[order.ID] = order
		ob.addToBuyPriceList(order)
	} else {
		if _, ok := ob.sells[order.ID]; ok {
			return fmt.Errorf("Order with ID %s already exists", order.ID)
		}
		ob.sells[order.ID] = order
		ob.addToSellPriceList(order)
	}

	err := transaction_log.SaveTransaction(&transaction_log.Transaction{
		ID:           order.ID,
		Operation:    "placeOrder",
		UserID:       order.UserID,
		OrderID:      order.ID,
		OrderDetails: &transaction_log.OrderDetails{order.Cryptocurrency, order.Type, order.Price, order.Amount},
		Timestamp:    time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("error logging transaction: %v", err)
	}

	return nil
}

func (ob *OrderBook) MatchOrders() {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	for len(ob.buyPrice) > 0 && len(ob.sellPrice) > 0 && ob.buyPrice[0].Price >= ob.sellPrice[0].Price {
		buyOrder := ob.buyPrice[0]
		sellOrder := ob.sellPrice[0]

		minAmount := buyOrder.Amount
		if sellOrder.Amount < minAmount {
			minAmount = sellOrder.Amount
		}

		buyOrder.Amount -= minAmount
		sellOrder.Amount -= minAmount

		if buyOrder.Amount == 0 {
			ob.removeFromBuyPriceList(buyOrder.ID)
			delete(ob.buys, buyOrder.ID)
			buyOrder.Status = "filled"
		}

		if sellOrder.Amount == 0 {
			ob.removeFromSellPriceList(sellOrder.ID)
			delete(ob.sells, sellOrder.ID)
			sellOrder.Status = "filled"
		}
	}
}

func (ob *OrderBook) GetOrders() ([]Order, []Order) {
	ob.mutex.RLock()
	defer ob.mutex.RUnlock()
	return ob.buyPrice, ob.sellPrice
}

func (ob *OrderBook) addToBuyPriceList(order Order) {
	for i, o := range ob.buyPrice {
		if o.Price < order.Price {
			ob.buyPrice = append(ob.buyPrice[:i], ob.buyPrice[i+1:]...)
			ob.buyPrice = append(ob.buyPrice, order)
			sort.Slice(ob.buyPrice, func(i, j int) bool { return ob.buyPrice[i].Price > ob.buyPrice[j].Price })
			return
		}
	}
	ob.buyPrice = append(ob.buyPrice, order)
	sort.Slice(ob.buyPrice, func(i, j int) bool { return ob.buyPrice[i].Price > ob.buyPrice[j].Price })
}

func (ob *OrderBook) addToSellPriceList(order Order) {
	for i, o := range ob.sellPrice {
		if o.Price > order.Price {
			ob.sellPrice = append(ob.sellPrice[:i], ob.sellPrice[i+1:]...)
			ob.sellPrice = append(ob.sellPrice, order)
			sort.Slice(ob.sellPrice, func(i, j int) bool { return ob.sellPrice[i].Price < ob.sellPrice[j].Price })
			return
		}
	}
	ob.sellPrice = append(ob.sellPrice, order)
	sort.Slice(ob.sellPrice, func(i, j int) bool { return ob.sellPrice[i].Price < ob.sellPrice[j].Price })
}

func (ob *OrderBook) removeFromBuyPriceList(orderID string) {
	for i, order := range ob.buyPrice {
		if order.ID == orderID {
			ob.buyPrice = append(ob.buyPrice[:i], ob.buyPrice[i+1:]...)
			return
		}
	}
}

func (ob *OrderBook) removeFromSellPriceList(orderID string) {
	for i, order := range ob.sellPrice {
		if order.ID == orderID {
			ob.sellPrice = append(ob.sellPrice[:i], ob.sellPrice[i+1:]...)
			return
		}
	}
}
