package main

import (
	"fmt"
	"strings"
	"unicode"
)

/*
Задача на распаковку
Создать Go-функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы/руны, например:
"a4bc2d5e" => "aaaabccddddde"
"abcd" => "abcd"
"45" => "" (некорректная строка)
"" => ""

Дополнительно
Реализовать поддержку escape-последовательностей.
Например:
qwe\4\5 => qwe45 (*)
qwe\45 => qwe44444 (*)
qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка, функция должна возвращать ошибку. Написать unit-тесты.
*/

func main() {
	s, err := Unpack("a4bc2d5e")
	fmt.Println(err)
	fmt.Println(s)
}

// Unpack распаковывает сжатую строку
func Unpack(s string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}
	b := strings.Builder{}
	var prev rune
	var prevEscaped bool
	// Идем по символам и в зависимости от текущего символа, предыдущего символа и флага(заэскейплен ли предыдущий
	// символ) добавляем предыдущий символ в результирующую строку, либо возвращаем ошибку
	for i, r := range s {
		if i == 0 {
			if unicode.IsDigit(r) {
				return "", fmt.Errorf("wrong string")
			}
		} else {
			if unicode.IsDigit(r) {
				if (unicode.IsDigit(prev) && prevEscaped) || (prev == '\\' && prevEscaped) || isNotDigitAndSlash(prev) {
					for j := 0; j < int(r-'0'); j++ {
						b.WriteRune(prev)
					}
					prevEscaped = false
				} else if prev == '\\' {
					prevEscaped = true
				} else {
					return "", fmt.Errorf("wrong string")
				}
			} else if r == '\\' {
				if prev == '\\' {
					prevEscaped = true
				} else if !(unicode.IsDigit(prev) && !prevEscaped) {
					b.WriteRune(prev)
					prevEscaped = false
				}
			} else {
				if prev == '\\' && !prevEscaped {
					return "", fmt.Errorf("wrong string")
				}
				if !(unicode.IsDigit(prev) && !prevEscaped) {
					b.WriteRune(prev)
					prevEscaped = false
				}
			}
		}
		prev = r
	}
	// обрабатываем последний символ
	if prev == '\\' && !prevEscaped {
		return "", fmt.Errorf("wrong string")
	}
	if ((unicode.IsDigit(prev) || prev == '\\') && prevEscaped) || isNotDigitAndSlash(prev) {
		b.WriteRune(prev)
	}
	return b.String(), nil
}

func isNotDigitAndSlash(r rune) bool {
	return !unicode.IsDigit(r) && !(r == '\\')
}
