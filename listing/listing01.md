Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
)

func main() {
	a := [5]int{76, 77, 78, 79, 80}
	var b []int = a[1:4]
	fmt.Println(b)
}
```

Ответ:
[77 78 79]
В переменной b будет слайс, в который входят элементы 
с 1го по 3й включительно из массива a.