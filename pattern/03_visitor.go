package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
Посетитель позволяет отвязать функциональность от объекта. Новые методы добавляются не для каждого типа из семейства,
а для промежуточного объекта visitor, аккумулирующего функциональность.
Типам семейства добавляется только один метод visit(visitor). Так проще добавлять операции к существующей
базе кода без особых изменений и страха всё сломать.
Паттерн Посетитель используется, когда:
 - нужно применить одну и ту же операцию к объектам разных типов;
 - часто добавляются новые операции для объектов;
 - требуется добавить новый функционал, но избежать усложнения кода объекта.
Пример использования:
обход сложных структур(подсчет статистик, сериализация/десериализация)
*/

type SizeVisitor struct {
	size int
}

func (v *SizeVisitor) Add(len int) {
	v.size += len
}

type Visited interface {
	Visit(visitor *SizeVisitor)
}

type MyString string
type MySlice []Visited
type MyMap map[string]Visited

func (s MyString) Visit(visitor *SizeVisitor) {
	visitor.Add(len(s))
}

func (s MySlice) Visit(visitor *SizeVisitor) {
	for _, v := range s {
		v.Visit(visitor)
	}
}

func (m MyMap) Visit(visitor *SizeVisitor) {
	for k, v := range m {
		visitor.Add(len(k))
		v.Visit(visitor)
	}
}

func main() {
	s := SizeVisitor{}
	mstr := MyString("test")
	msl := MySlice{MyString("one"), MyString("two")}
	mm := MyMap{"1": MyString("zero"), "2": msl}
	mstr.Visit(&s)
	mm.Visit(&s)
	fmt.Println(s.size)
}
