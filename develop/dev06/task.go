package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// структура с параметрами
type flags struct {
	fields    int
	delimiter string
	separated bool
}

func cutUtil(lines [][]string, myFlags flags) [][]string {
	res := make([][]string, len(lines))
	for i, line := range lines {
		for j, word := range line {
			if myFlags.fields == j {
				res[i] = append(res[i], word)
			}
		}
	}
	return res
}

func main() {
	// обработка флагов
	myFlags := flags{}
	filename := flag.String("n", "", "имя файла")
	flag.IntVar(&myFlags.fields, "f", 0, "выбрать поля (колонки)")
	flag.StringVar(&myFlags.delimiter, "d", "    ", "использовать другой разделитель")
	flag.BoolVar(&myFlags.separated, "s", false, "только строки с разделителем")
	flag.Parse()
	// считываем данные из stdin или файла
	var input *bufio.Scanner
	if *filename != "" {
		file, err := os.Open(*filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()
		input = bufio.NewScanner(file)
	} else {
		input = bufio.NewScanner(os.Stdin)
	}
	fmt.Println(*filename)
	// считываем строки и добавляем их в слайс
	var lines [][]string
	for input.Scan() {
		if newRow := strings.Split(input.Text(), myFlags.delimiter); (myFlags.separated == false || len(newRow) > 1) && len(newRow) != 0 {
			lines = append(lines, newRow)
		}
	}
	// Применение функции обрезки и вывод в консоль результата
	fmt.Println(cutUtil(lines, myFlags))
}
