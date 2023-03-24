package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/
/*
Есть робот-уборщик который может убирать разные части помещения: пылесосить, вытирать пыль и мыть окна.
Каждая из этих задач выполняется отдельным модулем, в определенном заданном порядке
*/

// Интерфейс обработчика
type RobotModule interface {
	SetNext(module RobotModule)
	Execute(*roomTaskInfo)
}

// Информация о задачах
type roomTaskInfo struct {
	room               uint
	VacuumDone         bool
	DustingDone        bool
	WindowCleaningDone bool
}

// Конкретный обработчик 1
type VacuumModule struct {
	next RobotModule
}

func (m *VacuumModule) Execute(tasks *roomTaskInfo) {
	if tasks.VacuumDone {
		fmt.Println("Помещение уже пропылесошено")
		m.next.Execute(tasks)
		return
	}
	fmt.Println("Работа пылесоса")
	tasks.VacuumDone = true
	m.next.Execute(tasks)
}

func (m *VacuumModule) SetNext(next RobotModule) {
	m.next = next
}

// Конкретный обработчик 2
type WindowCleaningModule struct {
	next RobotModule
}

func (m *WindowCleaningModule) Execute(tasks *roomTaskInfo) {
	if tasks.WindowCleaningDone {
		fmt.Println("Окна уже помыты")
		m.next.Execute(tasks)
		return
	}
	fmt.Println("Мытье окон")
	tasks.VacuumDone = true
	m.next.Execute(tasks)
}

func (m *WindowCleaningModule) SetNext(next RobotModule) {
	m.next = next
}

// Конкретный обработчик 3
type DustingModule struct {
	next RobotModule
}

func (m *DustingModule) Execute(tasks *roomTaskInfo) {
	if tasks.DustingDone {
		fmt.Println("Пыль уже вытерта")
		m.next.Execute(tasks)
		return
	}
	fmt.Println("Вытирание пыли")
	tasks.VacuumDone = true
	//m.next.Execute(tasks)
}
func (m *DustingModule) SetNext(next RobotModule) {
	m.next = next
}

// Конкретный обработчик 4
type startClean struct {
	next RobotModule
}

func (m *startClean) Execute(tasks *roomTaskInfo) {
	fmt.Println("\nСтарт уборки")
	m.next.Execute(tasks)
}

func (m *startClean) SetNext(next RobotModule) {
	m.next = next
}

func chainOfResponseWork() {
	//Создаём экземпляры обработчиков
	cleaning := &startClean{}
	dusting := &DustingModule{}
	windowClean := &WindowCleaningModule{}
	vacuuming := &VacuumModule{}

	//Определение порядока выполнения
	cleaning.SetNext(vacuuming)
	vacuuming.SetNext(windowClean)
	windowClean.SetNext(dusting)

	//Экземпляры комнат для уборки
	room1 := &roomTaskInfo{
		room: 12,
	}
	room2 := &roomTaskInfo{
		room:       21,
		VacuumDone: true,
	}

	cleaning.Execute(room1)
	cleaning.Execute(room2)
}
