package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"os/signal"
	"syscall"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

var ntpServer = "0.beevik-ntp.pool.ntp.org"

func main() {
	//Обработка паники и возврат ненулевого кода выхода в OS
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "Error:", r)
			os.Exit(1)
		}
	}()
	//Выход из программы
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go func() {
		<-c
		fmt.Println("Got SIGINT signal")
		// Дополнительный код обработки сигнала
		os.Exit(0)
	}()
	//Получение времени от ntp сервера
	flagT := "f"
	for {
		fmt.Println("Для получения точного времени введите t:")
		if fmt.Scan(&flagT); flagT == "t" {
			time, err := ntp.Time(ntpServer)
			if err != nil {
				panic(err)
			} else {
				fmt.Println(time)
			}
			flagT = "f"
		}
	}

}
