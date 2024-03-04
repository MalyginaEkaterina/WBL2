package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

/*
Фасад — это паттерн, который добавляет простой интерфейс к сложной системе для взаимодействия с ней.
Паттерн Фасад используется, когда:
 - есть сложная система, работу с которой нужно упростить;
 - хочется уменьшить количество зависимостей между клиентом и сложной системой;
 - требуется разбить сложную систему на компоненты — применение паттерна к каждому компоненту упростит взаимодействие между ними.
Пример использования на практике:
Упрощение работы с базой данных
*/

type DbStatement struct {
	query string
}

type DbConnection struct {
}

func (c *DbConnection) Prepare(query string) *DbStatement {
	fmt.Printf("Preparing query %s\n", query)
	return &DbStatement{
		query: query,
	}
}

func (c *DbConnection) Execute(stmt *DbStatement) {
	fmt.Printf("Executing query %s\n", stmt.query)
}

type DbConnectionPool struct {
}

func (p *DbConnectionPool) Get() *DbConnection {
	fmt.Println("Getting new connection")
	return &DbConnection{}
}

func (p *DbConnectionPool) Put(*DbConnection) {
	fmt.Println("Returning connection")
}

type DbFacade struct {
	pool *DbConnectionPool
}

func newDbFacade() *DbFacade {
	return &DbFacade{pool: &DbConnectionPool{}}
}

// Скрываем за фасадом всю работу с пулом коннектов к БД и по подготовке стейтментов
func (f *DbFacade) Execute(query string) {
	conn := f.pool.Get()
	defer f.pool.Put(conn)
	preparedStmt := conn.Prepare(query)
	conn.Execute(preparedStmt)
}

func main() {
	db := newDbFacade()
	db.Execute("SELECT * FROM table")
}
