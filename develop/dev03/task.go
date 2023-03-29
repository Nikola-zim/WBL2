package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// функция для удаления повторяющихся строк
func uniqueStrings(lines [][]string) [][]string {
	var rowsConcatenate []string
	for _, row := range lines {
		rowsConcatenate = append(rowsConcatenate, strings.Join(row, " "))
	}
	result := make([][]string, len(lines))
	seen := make(map[string]bool)
	for _, line := range rowsConcatenate {
		if !seen[line] {
			seen[line] = true
			result = append(result, rowsConcatenate)
		}
	}
	for _, row := range rowsConcatenate {
		result = append(lines, strings.Split(row, " "))
	}
	return result
}

func sortLines(lines [][]string, column int, reverse bool, numberSort bool) [][]string {
	// результирующий двумерный срез, имеет максимальную длину чтобы избежать паники
	result := make([][]string, len(lines))
	var max int
	for _, row := range lines {
		if max < len(row) {
			max = len(row)
		}
	}
	//Проверка что индекс колонки есть в файле
	if max < column {
		panic("такой колонки нет")
	}
	for i := range result {
		result[i] = make([]string, max, max)
		for j := range lines[i] {
			result[i][j] = lines[i][j]
		}
	}
	// опции сортировки
	var sortFunc func(i, j int) bool
	if !numberSort {
		sortFunc = func(i, j int) bool {
			if reverse {
				return result[i][column] > result[j][column]
			}
			return result[i][column] < result[j][column]
		}
	} else {
		sortFunc = func(i, j int) bool {
			ii, err := strconv.Atoi(result[i][column])
			if err != nil {
				panic(err)
			}
			jj, err := strconv.Atoi(result[j][column])
			if err != nil {
				panic(err)
			}
			if reverse {
				return ii > jj
			}
			return ii < jj
		}
	}

	if reverse {
		sort.Slice(result, sortFunc)
	} else {
		sort.Slice(result, sortFunc)
	}
	return result
}

func main() {
	// определяем флаги командной строки
	filename := flag.String("f", "", "имя файла")
	reverse := flag.Bool("r", false, "сортировать в обратном порядке")
	unique := flag.Bool("u", false, "не выводить повторяющиеся строки")
	numberSort := flag.Bool("n", false, "сортировать по числовому значению")
	column := flag.Int("k", 0, "указание колонки для сортировки")
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
		lines = append(lines, strings.Split(input.Text(), " "))
	}

	// отбрасывание повторяющихся строк при соответствующем флаге
	if *unique {
		lines = uniqueStrings(lines)
	}
	// сортировка
	lines = sortLines(lines, *column, *reverse, *numberSort)

	// выводим отсортированные строки
	for _, line := range lines {
		fmt.Printf("%s\n", line)
	}
}
