package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
Утилита cut

Реализовать утилиту аналог консольной команды cut (man cut).
Утилита должна принимать строки через STDIN, разбивать по разделителю (TAB) на колонки и выводить запрошенные.

Реализовать поддержку утилитой следующих ключей:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

*/

// Params это параметры cut
type Params struct {
	Fields    []int
	Delimiter string
	Separated bool
}

// CutString обрезает строку в соотвествии с параметрами
func CutString(s string, params Params) string {
	if len(params.Fields) == 0 {
		return s
	}
	arr := strings.Split(s, params.Delimiter)
	// если строка не содержит разделителя
	if len(arr) == 1 {
		if params.Separated {
			return ""
		}
		return s
	}
	var builder strings.Builder
	// в params.fields лежат номера столбцов, которые нужно вывести, отсортированные по возрастанию
	// невалидные номера столбцов просто игнорируем
	for i := 0; i < len(params.Fields); i++ {
		if params.Fields[i] >= len(arr) {
			break
		}
		if params.Fields[i] >= 0 {
			if builder.Len() > 0 {
				// также добавляем разделитель
				builder.WriteString(params.Delimiter)
			}
			builder.WriteString(arr[params.Fields[i]])
		}
	}
	return builder.String()
}

func main() {
	// Определяем флаги
	f := flag.String("f", "", "column numbers")
	d := flag.String("d", "\t", "delimiter")
	s := flag.Bool("s", false, "do not print strings without delimiter")
	flag.Parse()

	var columns []int
	if *f != "" {
		cols := strings.Split(*f, ",")
		for _, cn := range cols {
			n, err := strconv.Atoi(cn)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Usage: -fields=1,2,3")
				os.Exit(1)
			}
			columns = append(columns, n)
		}
	}

	// Если указать колонки, например, 4,2, то cut выводит в порядке 2,4, поэтому сортируем слайс колонок
	sort.Slice(columns, func(i, j int) bool {
		return columns[i] < columns[j]
	})

	params := Params{
		Fields:    columns,
		Delimiter: *d,
		Separated: *s,
	}

	scanner := bufio.NewScanner(os.Stdin)
	// Читаем строки из STDIN
	for scanner.Scan() {
		in := scanner.Text()
		cutStr := CutString(in, params)
		fmt.Println(cutStr)
	}

	// Проверяем на наличие ошибки при чтении
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Reading error")
		os.Exit(1)
	}
}
