package pattern

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
)

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

/*
Паттерн Состояние позволяет конструировать объект, способный иметь набор дискретных состояний и ведущий
себя по-разному в зависимости от состояния, в котором находится.
Паттерн Состояние используется, когда:
 - нужно менять поведение объекта в зависимости от его состояния;
Пример использования:
синтаксический разбор
обработка заказов
AI в играх
*/

type State int

const (
	New State = iota
	Paid
	Reserved
	Delivered
)

type Order struct {
	state State
	id    string
	data  map[string]string
}

func NewOrder(data map[string]string) *Order {
	id := rand.Int63()
	return &Order{state: New, id: strconv.FormatInt(id, 10), data: data}
}

func (o *Order) Pay() error {
	if o.state != New {
		return fmt.Errorf("wrong state %v", o.state)
	}
	transId := rand.Int63()
	fmt.Printf("order %s was paid\n", o.id)
	o.data["transaction id"] = strconv.FormatInt(transId, 10)
	o.state = Paid
	return nil
}

func (o *Order) Reserve() error {
	if o.state != Paid {
		return fmt.Errorf("wrong state %v", o.state)
	}
	reservationId := rand.Int63()
	fmt.Printf("order %s was reserved\n", o.id)
	o.data["reservation id"] = strconv.FormatInt(reservationId, 10)
	o.state = Reserved
	return nil
}

func (o *Order) Deliver() error {
	if o.state != Reserved {
		return fmt.Errorf("wrong state %v", o.state)
	}
	deliveryId := rand.Int63()
	fmt.Printf("order %s was delivered\n", o.id)
	o.data["delivery id"] = strconv.FormatInt(deliveryId, 10)
	o.state = Delivered
	return nil
}

func main() {
	order := NewOrder(map[string]string{"product id": "123"})
	err := order.Pay()
	if err != nil {
		log.Fatal(err)
	}
	err = order.Reserve()
	if err != nil {
		log.Fatal(err)
	}
	err = order.Deliver()
	if err != nil {
		log.Fatal(err)
	}
	err = order.Pay()
	if err != nil {
		log.Fatal(err)
	}
}
