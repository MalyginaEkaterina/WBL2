package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

/*
Взаимодействие с ОС


Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:


- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).

*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран.
Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).

*/

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter command")
	fmt.Print(">> ")
	// Читаем строки из STDIN
	for scanner.Scan() {
		str := scanner.Text()
		str = strings.TrimSpace(str) // Удаляем пробелы
		args := strings.SplitN(str, " ", 2)

		switch args[0] {
		case "quit":
			return
		case "cd":
			if len(args) < 2 {
				fmt.Println("Usage: cd <dir>")
			} else {
				err := os.Chdir(args[1])
				if err != nil {
					fmt.Printf("Change dir error: %v\n", err)
				}
			}
		case "pwd":
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Printf("Get current dir error: %v\n", err)
			} else {
				fmt.Printf("Current dir: %s\n", currentDir)
			}
		case "echo":
			if len(args) < 2 {
				fmt.Println("Usage: echo <args>")
			} else {
				fmt.Println(args[1])
			}
		case "kill":
			if len(args) < 2 {
				fmt.Println("Usage: kill <id>")
			} else {
				err := killProcess(args[1])
				if err != nil {
					fmt.Printf("Kill process error: %v\n", err)
				} else {
					fmt.Printf("Process %s was killed\n", args[1])
				}
			}
		case "ps":
			err := getProcessInfo()
			if err != nil {
				fmt.Printf("Get process info error: %v\n", err)
			}
		default:
			err := execCmd(str)
			if err != nil {
				fmt.Printf("Run error: %v\n", err)
			}
		}
		fmt.Print(">> ")
	}

	// Проверяем на наличие ошибки при чтении
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Read command error")
		os.Exit(1)
	}
}

func execCmd(commands string) error {
	signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()
	errGroup, ctx := errgroup.WithContext(signalCtx)

	commandList := strings.Split(commands, "|")
	cmdList := make([]*exec.Cmd, len(commandList))
	for i, command := range commandList {
		command = strings.TrimSpace(command)
		args := strings.Split(command, " ")
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		cmdList[i] = cmd
	}

	for i := 0; i < len(cmdList); i++ {
		cmdList[i].Stderr = os.Stderr
		if i == 0 {
			cmdList[i].Stdin = os.Stdin
		} else {
			stdout, err := cmdList[i-1].StdoutPipe()
			if err != nil {
				return err
			}
			cmdList[i].Stdin = stdout
		}
		if i == len(cmdList)-1 {
			cmdList[i].Stdout = os.Stdout
		}
	}

	for i, cmd := range cmdList {
		command := commandList[i]
		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("start command %q error: %w", command, err)
		}
	}

	for i, cmd := range cmdList {
		command := commandList[i]
		cmd := cmd
		errGroup.Go(func() error {
			err := cmd.Wait()
			if err != nil {
				return fmt.Errorf("wait command %q error: %w", command, err)
			}
			return nil
		})
	}

	err := errGroup.Wait()
	if err != nil {
		if !errors.Is(signalCtx.Err(), context.Canceled) {
			return err
		}
	}
	return nil
}

func killProcess(pid string) error {
	var command *exec.Cmd

	if runtime.GOOS == "windows" {
		command = exec.Command("taskkill", "/F", "/PID", pid)
	} else if runtime.GOOS == "linux" {
		command = exec.Command("kill", pid)
	} else {
		return fmt.Errorf("unsupported OS")
	}

	return command.Run()
}

func getProcessInfo() error {
	var command *exec.Cmd

	if runtime.GOOS == "windows" {
		command = exec.Command("tasklist")
	} else if runtime.GOOS == "linux" {
		command = exec.Command("ps")
	} else {
		return fmt.Errorf("unsupported OS")
	}

	output, err := command.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
