package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

/*
Утилита telnet

Реализовать простейший telnet-клиент.

Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123


Требования:
Программа должна подключаться к указанному хосту (ip или доменное имя + порт) по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s)
При нажатии Ctrl+D программа должна закрывать сокет и завершаться.
Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout

*/

func main() {
	// Определяем флаг timeout
	timeout := flag.Int("timeout", 10, "timeout in seconds")
	flag.Parse()

	// Определяем хост и порт
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: host port")
		os.Exit(1)
	}
	host := args[0]
	port := args[1]

	// устанавливаем соединение
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), time.Duration(*timeout)*time.Second)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Connection error: ", err)
		os.Exit(1)
	}

	defer conn.Close()

	// канал для сигнала о закрытии соединения со стороны сервера
	closed := make(chan interface{})

	// направляем на Stdout то, что приходит от сервера, после закрытия сокета закрываем канал closed
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Read from socket error: ", err)
		}
		close(closed)
	}()

	// канал для перенаправления ввода в Stdin. Требуется, чтобы было возможно обработать сигнал о закрытии сокета вовремя
	in := make(chan []byte)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		// Читаем строки из STDIN
		for scanner.Scan() {
			str := scanner.Text()
			in <- []byte(str + "\n")
		}

		// Проверяем на наличие ошибки при чтении
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Read from stdin error")
		}
		close(in)
	}()

	for {
		select {
		case <-closed:
			fmt.Fprintln(os.Stdout, "Socket was closed")
			os.Exit(0)
		case data, ok := <-in:
			if !ok {
				fmt.Fprintln(os.Stdout, "Stdin was closed")
				os.Exit(0)
			}
			_, err = conn.Write(data)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Write into socket error")
				os.Exit(1)
			}
		}
	}
}
