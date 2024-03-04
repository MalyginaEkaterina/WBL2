Что выведет программа? Объяснить вывод программы. 
Объяснить как работают defer’ы и порядок их вызовов.

```go
package main

import (
	"fmt"
)

func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```

Ответ:
2
1
defer выполняется перед выходом из функции, но после стейтмента
return.
Получается, что в случае неименованных возвращаемых параметров
defer на возвращаемое значение не повлияет, а в случае
именованных возвращаемых параметров в defer значение для них может 
быть изменено.
