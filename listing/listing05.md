Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
Выведет error.
Т.к. переменная err инициализируется как интерфейс error, то нулевой
указатель *customError оборачивается в интерфейс и в err лежит
указатель на то, как тип *customError реализует интерфейс error, поэтому
err не равен nil.
