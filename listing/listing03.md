Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:
```
<nil>
false

В первом print выводится <nil>, так как os.PathError был присвоен nil. Однако интерфейс 
error имеет 2 поля: (itab и data). Интерфейс равен nil, только если и тип, 
и значение равны nil, поэтому во втором случае получаем false.

```
