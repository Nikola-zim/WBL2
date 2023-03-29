package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func findAnagramSets(words []string) map[string][]string {
	// Создаем карту, где ключ - отсортированное по алфавиту слово,
	// значение - срез строк, содержащий все анаграммы этого слова
	anagramMap := make(map[string][]string)
	for _, word := range words {
		word = strings.ToLower(word)
		sortedWord := sortString(word)
		anagramMap[sortedWord] = append(anagramMap[sortedWord], word)
	}

	// Создаем мапу для хранения множеств анаграмм
	anagramSets := make(map[string][]string)
	// Обходим карту и добавляем все срезы строк со значениями,
	// содержащими более одного элемента, в мапу множеств анаграмм
	for _, words := range anagramMap {
		if len(words) > 1 {
			// убираем дубликаты с помощью анонимной функции
			words = func(arr []string) []string {
				m, uniq := make(map[string]struct{}), make([]string, 0, len(arr))
				for _, v := range arr {
					if _, ok := m[v]; !ok {
						m[v], uniq = struct{}{}, append(uniq, v)
					}
				}
				return uniq
			}(words)
			// сортируем слайсы слов
			sort.Strings(words)
			anagramSets[words[0]] = words
		}
	}
	return anagramSets
}

// Функция для сортировки строки по алфавиту
func sortString(str string) string {
	chars := strings.Split(str, "")
	sort.Strings(chars)
	return strings.Join(chars, "")
}

func main() {
	words := []string{"Столик", "пятак", "листок", "eat", "tea", "tan", "ate", "nat", "bat", "пятка", "слиток", "тяпка", "кот", "слиток"}

	fmt.Println(findAnagramSets(words))
}
