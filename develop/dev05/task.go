package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
Утилита grep


Реализовать утилиту фильтрации по аналогии с консольной утилитой (man grep — смотрим описание и основные параметры).


Реализовать поддержку утилитой следующих ключей:
-A - "after" печатать +N строк после совпадения
-B - "Before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "Count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "Fixed", точное совпадение со строкой, не паттерн
-n - "line num", напечатать номер строки
*/

// Params это параметры grep
type Params struct {
	After      int
	Before     int
	Count      bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	LineNum    bool
}

// Matcher проверяет строку на соответствие условиям
type Matcher func(s string) bool

func grepStdin(expr string, params Params) {
	err := Grep(expr, params, "<stdin>", os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Grep error: %v", err)
		os.Exit(1)
	}
}

func grepInput(input string, expr string, params Params) {
	// Открываем файл для чтения
	inputFile, err := os.Open(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open inputFile error: %v", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	err = Grep(expr, params, input, inputFile, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Grep error: %v", err)
		os.Exit(1)
	}
}

// Grep записывает согласно параметрам строки в writer из reader, подходящие под фильтр.
func Grep(expr string, params Params, readerName string, reader io.Reader, writer io.Writer) error {
	// Подготавливаем функцию, которая будет проверять, подходит ли строка под фильтр
	matcher, err := PrepareMatcher(params, expr)
	if err != nil {
		return fmt.Errorf("prepare matcher error: %w", err)
	}

	scanner := bufio.NewScanner(reader)

	// если нужно только количество строк, попадающих под фильтр
	if params.Count {
		count := 0
		for scanner.Scan() {
			s := scanner.Text()
			if matcher(s) {
				count++
			}
		}

		// Проверяем на наличие ошибки при чтении файла
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("read inputFile error: %w", err)
		}

		_, err = fmt.Fprintln(writer, strconv.Itoa(count))
		if err != nil {
			return fmt.Errorf("write error: %w", err)
		}
		return nil
	}

	// номер строки
	lineNum := 0

	// слайс для хранения строк контекста "до"
	beforeBuf := make([]string, 0, params.Before)
	// параметр для определения, сколько строк контекста "после" осталось вывести
	after := 0

	for scanner.Scan() {
		s := scanner.Text()
		lineNum++
		if matcher(s) {
			// записываем все строки до
			for _, str := range beforeBuf {
				_, err = fmt.Fprintln(writer, str)
				if err != nil {
					return fmt.Errorf("write error: %w", err)
				}
			}
			beforeBuf = beforeBuf[:0]

			if params.LineNum {
				s = fmt.Sprintf("%s:%d:%s", readerName, lineNum, s)
			}
			_, err = fmt.Fprintln(writer, s)
			if err != nil {
				return fmt.Errorf("write error: %w", err)
			}

			after = params.After
		} else {
			// если нужно вывести контекст "после", выводим
			if after > 0 {
				if params.LineNum {
					s = fmt.Sprintf("%s:%d:%s", readerName, lineNum, s)
				}
				_, err = fmt.Fprintln(writer, s)
				if err != nil {
					return fmt.Errorf("write error: %w", err)
				}
				after--
			} else {
				// если нужно выводить контекст "до", то сохраняем в слайс, чтобы затем вывести, когда встретится строка,
				// попадающая под фильтр
				if params.Before > 0 {
					if len(beforeBuf) == params.Before {
						beforeBuf = beforeBuf[1:]
					}
					if params.LineNum {
						s = fmt.Sprintf("%s:%d:%s", readerName, lineNum, s)
					}
					beforeBuf = append(beforeBuf, s)
				}
			}
		}
	}

	// Проверяем на наличие ошибки при чтении файла
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read inputFile error: %w", err)
	}

	return nil
}

// PrepareMatcher подготавливает на основании параметров функцию, определяющую, подходит ли строка под фильтр
func PrepareMatcher(params Params, expr string) (Matcher, error) {
	if params.IgnoreCase {
		expr = strings.ToLower(expr)
	}
	if params.Fixed {
		if params.Invert {
			if params.IgnoreCase {
				return func(s string) bool {
					return !strings.Contains(strings.ToLower(s), expr)
				}, nil
			}
			return func(s string) bool {
				return !strings.Contains(s, expr)
			}, nil
		}
		if params.IgnoreCase {
			return func(s string) bool {
				return strings.Contains(strings.ToLower(s), expr)
			}, nil
		}
		return func(s string) bool {
			return strings.Contains(s, expr)
		}, nil
	}

	// добавляем в регулярное выражение игнорирование регистра
	if params.IgnoreCase {
		expr = "(?i)" + expr
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("regexp compile error: %w", err)
	}

	if params.Invert {
		return func(s string) bool {
			return !re.MatchString(s)
		}, nil
	}
	return func(s string) bool {
		return re.MatchString(s)
	}, nil
}

func main() {
	// Определяем флаги
	a := flag.Int("A", 0, "print N lines after match")
	b := flag.Int("B", 0, "print N lines Before match")
	ab := flag.Int("C", 0, "print N lines after and Before match")
	c := flag.Bool("c", false, "print just Count of matches")
	i := flag.Bool("i", false, "ignore case")
	v := flag.Bool("v", false, "not to print matches")
	f := flag.Bool("F", false, "not pattern")
	n := flag.Bool("n", false, "print line number")

	flag.Parse()

	params := Params{
		After:      *a,
		Before:     *b,
		Count:      *c,
		IgnoreCase: *i,
		Invert:     *v,
		Fixed:      *f,
		LineNum:    *n,
	}
	if *ab != 0 {
		params.After = *ab
		params.Before = *ab
	}

	// Получаем expression и имена файлов
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: expression [file1 ...]")
		os.Exit(1)
	}
	expr := args[0]
	inputs := args[1:]

	if len(inputs) == 0 {
		grepStdin(expr, params)
	} else {
		for _, input := range inputs {
			grepInput(input, expr, params)
		}
	}
}
