package pattern

import "fmt"

/*
	Реализовать паттерн «команда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

/*
Паттерн Команда преобразует все параметры операции или события в объект-команду.
Впоследствии можно выполнить эту операцию, вызвав соответствующий метод объекта.
Объект-команда заключает в себе всё необходимое для проведения операции, поэтому её легко выполнять, логировать и отменять.
Применяется, когда:
 - нужно преобразовать запросы в объекты, что позволяет передавать, обрабатывать и хранить их;
 - нужно создать очередь операций и работать с ней;
 - нужно откатить выполненные действия.
Пример использования на практике:
Команда применяется в работе с базами данных.
В стандартной библиотеке Go есть пример SQL-инструкции sql.Stmt.
Такую заранее подготовленную инструкцию можно многократно выполнять методом Stmt.Exec,
не задумываясь о её внутренней структуре.
sql.Stmt, выполненную в рамках транзакции Tx.Stmt(), легко откатить с помощью Tx.Rollback().
*/

// command - объект, который знает все методы receiver и вызывает один из них,
// а также хранит параметры для вызова. Команда выполняется вызовом своего метода execute()
type command interface {
	execute()
}

// receiver — непосредственный исполнитель команды. Методы этого объекта и совершают фактическую работу.
type receiver interface {
	action(message string)
}

// invoker — объект, который знает интерфейс команд, умеет их вызывать, пользуясь только этим интерфейсом.
// Не зависит от внутреннего устройства команд.
type invoker struct {
	commands map[string]command
}

func newInvoker() *invoker {
	i := new(invoker)
	i.commands = make(map[string]command)
	return i
}

func (i *invoker) do(c string) {
	i.commands[c].execute()
}

// реализация command
type helloPrinter struct {
	receiver receiver
}

func (c *helloPrinter) execute() {
	c.receiver.action("Hello")
}

type goodbyePrinter struct {
	receiver receiver
}

func (c *goodbyePrinter) execute() {
	c.receiver.action("Goodbye")
}

// реализация receiver
type rcvr struct {
	prefix string
}

func (r *rcvr) action(message string) {
	fmt.Println(r.prefix + message)
}

func main() {
	// клиентский код
	var rcvr1 = rcvr{"John: "}
	var cmd1 = helloPrinter{&rcvr1}
	var rcvr2 = rcvr{"Another: "}
	var cmd2 = goodbyePrinter{&rcvr2}
	invkr := newInvoker()
	invkr.commands["print_hello"] = &cmd1
	invkr.commands["print_goodbye"] = &cmd2
	// применение команд
	invkr.do("print_hello")
	invkr.do("print_goodbye")
}
