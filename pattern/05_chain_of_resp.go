package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
Паттерн Цепочка выстраивает конвейер обработчиков для поступающих запросов.
Объект «обработчик» выполняет свою часть процессинга и передаёт запрос дальше по цепочке.
Обработчики не влияют друг на друга и не меняют состояние друг друга.
Поэтому их легко писать, отлаживать и переносить между проектами.
Паттерн Цепочка обязанностей используется, когда:
 - нужно иметь несколько обработчиков, которые будут вызываться в определённом порядке;
 - нужно обрабатывать разные типы запросов разными обработчиками.
Пример применения в Go — цепочка middleware-обработчиков http.Request
*/

type Request struct {
	kind string
	body string
}

type Handler interface {
	Handle(r Request)
	SetNext(handler Handler)
}

type CheckHandler struct {
	next Handler
}

func (c *CheckHandler) Handle(r Request) {
	if r.kind == "wrong" {
		return
	}
	if c.next != nil {
		c.next.Handle(r)
	}
}

func (c *CheckHandler) SetNext(handler Handler) {
	c.next = handler
}

type LogHandler struct {
	next Handler
}

func (l *LogHandler) Handle(r Request) {
	fmt.Printf("Handle request %s\n", r.kind)
	if l.next != nil {
		l.next.Handle(r)
	}
}

func (l *LogHandler) SetNext(handler Handler) {
	l.next = handler
}

type MainHandler struct {
	next Handler
}

func (m *MainHandler) Handle(r Request) {
	fmt.Printf("Request body: %s\n", r.body)
	if m.next != nil {
		m.next.Handle(r)
	}
}

func (m *MainHandler) SetNext(handler Handler) {
	m.next = handler
}

type Handlers struct {
	main Handler
}

func NewHandlers(handlers ...Handler) *Handlers {
	h := Handlers{main: handlers[0]}
	for i := 1; i < len(handlers); i++ {
		handlers[i-1].SetNext(handlers[i])
	}
	return &h
}

func (h *Handlers) Handle(r Request) {
	h.main.Handle(r)
}

func main() {
	logger := &LogHandler{}
	check := &CheckHandler{}
	mainHandler := &MainHandler{}
	handlers := NewHandlers(logger, check, mainHandler)
	handlers.Handle(Request{kind: "wrong", body: "something wrong"})
	handlers.Handle(Request{kind: "right", body: "something right"})
}
