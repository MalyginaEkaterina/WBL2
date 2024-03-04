package pattern

import (
	"fmt"
	"log"
	"os"
)

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/*
Фабричный метод — это паттерн проектирования, который определяет интерфейс для создания объектов определённого типа.
Паттерн предлагает создавать объекты не напрямую, а через вызов специального фабричного метода.
Все объекты должны иметь общий интерфейс, который отвечает за их создание.
Паттерн используется, когда:
 - заранее неизвестны типы объектов — фабричный метод отделяет код создания объектов от остального кода, где они используются;
 - нужна возможность расширять части существующей системы.
Пример использования на практике:
Создание различных типов логирования (например, файл, консоль, база данных).
*/

type Logger interface {
	Log(message string)
	Close()
}

var _ Logger = (*FileLogger)(nil)

type FileLogger struct {
	file *os.File
}

func newFileLogger(name string) (*FileLogger, error) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &FileLogger{file: file}, nil
}

func (l *FileLogger) Log(message string) {
	_, err := fmt.Fprintln(l.file, message)
	if err != nil {
		fmt.Println(err)
	}
}

func (l *FileLogger) Close() {
	l.file.Close()
}

var _ Logger = (*ConsoleLogger)(nil)

type ConsoleLogger struct{}

func newConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Log(message string) {
	fmt.Println(message)
}

func (l *ConsoleLogger) Close() {
}

// NewLogger реализует фабричный метод
func NewLogger(t string) (Logger, error) {
	switch t {
	case "console":
		return newConsoleLogger(), nil
	case "file":
		return newFileLogger("test.log")
	default:
		return nil, fmt.Errorf("unknown logger type")
	}
}

func main() {
	fileLogger, err := NewLogger("file")
	if err != nil {
		log.Fatal(err)
	}
	defer fileLogger.Close()
	fileLogger.Log("file logger was created")
	fileLogger.Log("write something")

	consoleLogger, _ := NewLogger("console")
	consoleLogger.Log("console logger created")
}
