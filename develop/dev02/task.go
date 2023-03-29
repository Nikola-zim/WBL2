package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func UnpackString(s string) string {
	var result string
	var prev string
	var repeat int
	for _, char := range s {
		if unicode.IsDigit(char) {
			digit, _ := strconv.Atoi(string(char))
			repeat = digit
		} else {
			if repeat > 0 {
				result += strings.Repeat(prev, repeat-1)
				repeat = 0
			}
			result += string(char)
			prev = string(char)
		}
	}
	if repeat > 0 {
		result += strings.Repeat(prev, repeat-1)
	}
	return result
}

func main() {
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
	inputStr := ""
	for {
		fmt.Scan(&inputStr)
		resStr := UnpackString(inputStr)
		fmt.Println(resStr)
	}
}
