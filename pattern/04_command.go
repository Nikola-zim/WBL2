package pattern

import "fmt"

/*
	Реализовать паттерн «комманда».

Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.

	https://en.wikipedia.org/wiki/Command_pattern
*/
/*
 Код представляет собой управление роботом Omnibot при помощи 3 кнопок: moveForward, moveBackward, stop.
*/
// Интерфейс команды
type Command interface {
	execute()
}

// Отправитель
type Button struct {
	command Command
}

func (b *Button) press() {
	b.command.execute()
}

// Интерфейс получателя
type OmniBot interface {
	moveForward()
	moveBackward()
	stop()
}

// Конкретная команда
type forwardCommand struct {
	robot OmniBot
}

func (c *forwardCommand) execute() {
	c.robot.moveForward()
}

// Конкретная команда
type backwardCommand struct {
	robot OmniBot
}

func (c *backwardCommand) execute() {
	c.robot.moveBackward()
}

// Конкретная команда
type stopCommand struct {
	robot OmniBot
}

func (c *stopCommand) execute() {
	c.robot.stop()
}

// Конкретный получатель
type myOmniBot struct {
	direction string
}

func (ob *myOmniBot) moveForward() {
	ob.direction = "forward"
	fmt.Println("bot is moving forward")
}

func (ob *myOmniBot) moveBackward() {
	ob.direction = "backward"
	fmt.Println("bot is moving backward")
}
func (ob *myOmniBot) stop() {
	ob.direction = "stop"
	fmt.Println("bot stopped")
}

func CommandWork() {
	robot := &myOmniBot{}
	forwardCommand := &forwardCommand{
		robot: robot,
	}
	backwardCommand := &backwardCommand{
		robot: robot,
	}
	stopCommand := &stopCommand{
		robot: robot,
	}

	//Экземпляры "кнопок"
	forwardButton := &Button{
		command: forwardCommand,
	}
	backwardButton := &Button{
		command: backwardCommand,
	}
	stopButton := &Button{
		command: stopCommand,
	}

	//Отработка нажатий кнопок
	forwardButton.press()
	backwardButton.press()
	stopButton.press()
}
