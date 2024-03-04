package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
Утилита sort
Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание и основные параметры): на входе подается файл из несортированными строками, на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:

-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительно

Реализовать поддержку утилитой следующих ключей:

-b — игнорировать хвостовые пробелы
*/

import (
	"fmt"
)

// Comparable это интерфейс для сравнения элементов друг с другом
type Comparable interface {
	Less(other Comparable) bool
}

// IntKey это обертка над строкой, реализующая Comparable
type IntKey int

var _ Comparable = IntKey(0)

// Less реализует Comparable
func (ik IntKey) Less(other Comparable) bool {
	return int(ik) < int(other.(IntKey))
}

// StringKey это обертка над строкой, реализующая Comparable
type StringKey string

var _ Comparable = StringKey("")

// Less реализует Comparable
func (sk StringKey) Less(other Comparable) bool {
	return strings.Compare(string(sk), string(other.(StringKey))) < 0
}

// ReverseKey инвертирует порядок сортировки
type ReverseKey struct {
	inner Comparable
}

// Less реализует Comparable
func (rk ReverseKey) Less(other Comparable) bool {
	return !rk.inner.Less(other.(ReverseKey).inner)
}

// Pair это строка и ключ сортировки этой строки
type Pair struct {
	Row string
	Key Comparable
}

// KeyExtractor достает из строки ключ сортировки
type KeyExtractor func(string) Comparable

// ColumnExtractor берет из строки нужный столбец
func ColumnExtractor(prev KeyExtractor, column int) KeyExtractor {
	return func(s string) Comparable {
		k := string(prev(s).(StringKey))
		return StringKey(strings.Fields(k)[column])
	}
}

// IgnoreSpacesExtractor убирает пробелы спереди и сзади
func IgnoreSpacesExtractor(prev KeyExtractor) KeyExtractor {
	return func(s string) Comparable {
		k := string(prev(s).(StringKey))
		return StringKey(strings.TrimSpace(k))
	}
}

// IntExtractor преобразует строку в число, в случае если не получилось преобразовать устанавливает в значение 0
func IntExtractor(prev KeyExtractor) KeyExtractor {
	return func(s string) Comparable {
		k := string(prev(s).(StringKey))
		n, err := strconv.Atoi(k)
		if err != nil {
			n = 0
		}
		return IntKey(n)
	}
}

// ReverseExtractor используется для сортировки в обратном порядке
func ReverseExtractor(prev KeyExtractor) KeyExtractor {
	return func(s string) Comparable {
		k := prev(s)
		return ReverseKey{inner: k}
	}
}

// Params это параметры сортировки
type Params struct {
	IgnoreSpaces bool
	ColumnNum    int
	IsNumber     bool
	Reverse      bool
	NoDuplicates bool
}

// PrepareExtractor подготавливает KeyExtractor на основании параметров
func PrepareExtractor(params Params) KeyExtractor {
	extractor := func(s string) Comparable {
		return StringKey(s)
	}
	if params.IgnoreSpaces {
		extractor = IgnoreSpacesExtractor(extractor)
	}
	if params.ColumnNum > -1 {
		extractor = ColumnExtractor(extractor, params.ColumnNum)
	}
	if params.IsNumber {
		extractor = IntExtractor(extractor)
	}
	if params.Reverse {
		extractor = ReverseExtractor(extractor)
	}
	return extractor
}

// Sort сортирует строки из reader и записывает результат в writer на основании параметров.
func Sort(params Params, reader io.Reader, writer io.Writer) error {
	// На основании флагов формируем функцию, которой будем получать из строки ключ, по которому дальше будем сортировать
	extractor := PrepareExtractor(params)

	var rows []Pair
	scanner := bufio.NewScanner(reader)
	// Читаем файл построчно, применяем extractor для получения ключа и формируем слайс пар(ключ, строка)
	for scanner.Scan() {
		s := scanner.Text()
		p := Pair{
			Row: s,
			Key: extractor(s),
		}
		rows = append(rows, p)
	}

	// Проверяем на наличие ошибки при чтении файла
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read inputFile error: %w", err)
	}

	// Сортируем слайс по ключам
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Key.Less(rows[j].Key)
	})

	// Записываем в файл без дубликатов(ключей)
	if params.NoDuplicates {
		for i := 0; i < len(rows); i++ {
			if i == 0 || rows[i-1].Key != rows[i].Key {
				_, err := fmt.Fprintln(writer, rows[i].Row)
				if err != nil {
					return fmt.Errorf("write into output inputFile error: %w", err)
				}
			}
		}
	} else {
		// Записываем в файл все строки
		for i := 0; i < len(rows); i++ {
			_, err := fmt.Fprintln(writer, rows[i].Row)
			if err != nil {
				return fmt.Errorf("write into output inputFile error: %w", err)
			}
		}
	}
	return nil
}

func main() {
	// Определяем флаги
	k := flag.Int("k", -1, "column number")
	n := flag.Bool("n", false, "sort key is number")
	r := flag.Bool("r", false, "reverse order")
	u := flag.Bool("u", false, "delete duplicates")
	b := flag.Bool("b", false, "ignore spaces")
	flag.Parse()

	params := Params{
		IgnoreSpaces: *b,
		ColumnNum:    *k,
		IsNumber:     *n,
		Reverse:      *r,
		NoDuplicates: *u,
	}

	// Получаем имена файлов входного и выходного
	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: input-inputFile output-inputFile")
		os.Exit(1)
	}
	input := args[0]
	output := args[1]

	// Открываем файл для чтения
	inputFile, err := os.Open(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open inputFile error: %v", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	// Открываем или создаем файл для записи
	outputFile, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open output inputFile error: %v", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	err = Sort(params, inputFile, outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Sort error: %v", err)
		os.Exit(1)
	}
}
