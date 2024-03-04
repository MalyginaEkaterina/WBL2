Что выведет программа? Объяснить вывод программы. 
Объяснить внутреннее устройство интерфейсов и 
их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:
<nil>
false
У fmt.Println есть явная обработка случая вызова для нулевого указателя,
обернутого в интерфейс, в этом случае он выводит <nil>.
Проверка err == nil дает false, т.к. нулевой указатель, обернутый в интерфейс 
не равен nil, т.к. err в main содержит указатель на то, как тип 
*os.PathError реализует интерфейс error.