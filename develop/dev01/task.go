package main

import (
	"flag"
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)

/*
Базовая задача

Создать программу печатающую точное время с использованием NTP -библиотеки.
Инициализировать как go module. Использовать библиотеку github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Требования:
Программа должна быть оформлена как go module
Программа должна корректно обрабатывать ошибки библиотеки: выводить их в STDERR и возвращать ненулевой код выхода в OS
*/

func main() {
	a := flag.String("address", "0.beevik-ntp.pool.ntp.org", "ntp server address")
	flag.Parse()
	curTime, err := ntp.Time(*a)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ntp error: %v", err)
		os.Exit(1)
	}
	fmt.Println(curTime.Format(time.RFC850))
}
