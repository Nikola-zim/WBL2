package pattern

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern

Стратегия — это поведенческий паттерн, выносит набор алгоритмов в собственные классы и делает их взаимозаменимыми.
В данном примере определены две конкретные стратегии: атака и защита. Контекст использует интерфейс стратегии для выполнения действий в игре.
При использовании контекста можно легко изменять стратегию, не затрагивая клиентский код.
В зависимости от текущей ситуации в игре, можно выбирать соответствующую стратегию.
*/
import "fmt"

// Интерфейс стратегии
type Strategy interface {
	Attack()
	Defend()
}

// Конкретная стратегия для атаки
type AttackStrategy struct{}

func (a AttackStrategy) Attack() {
	fmt.Println("Атакую врага!")
}

func (a AttackStrategy) Defend() {
	fmt.Println("Защищаюсь от врага!")
}

// Конкретная стратегия для защиты
type DefendStrategy struct{}

func (d DefendStrategy) Attack() {
	fmt.Println("Не могу атаковать, защищаюсь!")
}

func (d DefendStrategy) Defend() {
	fmt.Println("Укрепляю оборону!")
}

// Контекст, использующий стратегию
type Context struct {
	strategy Strategy
}

func (c *Context) SetStrategy(strategy Strategy) {
	c.strategy = strategy
}

func (c *Context) Attack() {
	c.strategy.Attack()
}

func (c *Context) Defend() {
	c.strategy.Defend()
}

// Пример использования
func StrategyWork() {
	context := Context{}

	// Используем стратегию атаки
	context.SetStrategy(AttackStrategy{})
	context.Attack()
	context.Defend()

	// Используем стратегию защиты
	context.SetStrategy(DefendStrategy{})
	context.Attack()
	context.Defend()
}
