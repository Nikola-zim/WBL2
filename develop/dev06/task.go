package main

import (
	"flag"
	"fmt"
	"io"
	"log"
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

func main() {
	// обработка флагов
	myFlags := flags{}
	flag.IntVar(&myFlags.fields, "f", 0, "выбрать поля (колонки)")
	flag.StringVar(&myFlags.delimiter, "d", " ", "использовать другой разделитель")
	flag.BoolVar(&myFlags.separated, "s", false, "только строки с разделителем")
	flag.Parse()

	filename := "test_text.txt"
	args := flag.Args()
	if len(args) > 1 {
		log.Fatalln("qwer")
	}
	if len(args) == 1 {
		filename = args[0]
	}

	input := ""
	if filename != "" {
		bytes, err := os.ReadFile(filename)
		if err != nil {
			log.Fatalln(err)
		}
		input = string(bytes)
	} else {
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}
		input = string(bytes)
	}

	lines := strings.Split(input, "\n")

	for _, line := range lines {
		cols := strings.Split(line, f.D)

		if f.F > 0 && f.F <= len(cols) {
			fmt.Println(cols[f.F-1])
		} else if !f.S {
			fmt.Println(line)
		}
	}
}
