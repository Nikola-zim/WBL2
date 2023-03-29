package pattern

import (
	"errors"
	"fmt"
	"log"
)

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern

Состояние — это поведенческий паттерн, позволяющий динамически изменять поведение объекта при смене его состояния.
Очень важным нюансом, отличающим этот паттерн от Стратегии, является то,
что и контекст, и сами конкретные состояния могут знать друг о друге и инициировать переходы от одного состояния к другому.
*/

// GameHero Контекст
type GameHero struct {
	attack       State
	defend       State
	currentState State
	bulletCount  int
}

func newGameHero(ammo int) *GameHero {
	gH := &GameHero{bulletCount: ammo}
	gH.attack = &AttackState{gameHero: gH}
	gH.defend = &DefendState{gameHero: gH}
	gH.setState(gH.attack)
	return gH
}

func (gH *GameHero) setState(s State) {
	gH.currentState = s
}
func (gH *GameHero) fireShot(shots int) error {
	return gH.currentState.fireShot(shots)
}
func (gH *GameHero) getCover() error {
	return gH.currentState.getCover()
}
func (gH *GameHero) findBullets() error {
	return gH.currentState.findBullets()
}

// State Интерфейс состояния
type State interface {
	fireShot(shots int) error
	getCover() error
	findBullets() error
}

// AttackState Конкретное атакующее состояние
type AttackState struct {
	gameHero *GameHero
}

func (as *AttackState) fireShot(shots int) error {
	if bulletsRest := as.gameHero.bulletCount - shots; bulletsRest >= 0 {
		as.gameHero.bulletCount = bulletsRest
		fmt.Printf("персонаж пострелял. Осталось: %d пуль \n", bulletsRest)
		return nil
	} else {
		return errors.New("недостаточно пуль")
	}
}

func (as *AttackState) getCover() error {
	as.gameHero.setState(as.gameHero.defend)
	fmt.Println("персонаж в укрытии")
	return nil
}

func (as *AttackState) findBullets() error {
	return errors.New("персонаж атакует, пули сейчас найти невозможно")
}

// DefendState Состояние защиты
type DefendState struct {
	gameHero *GameHero
}

func (as *DefendState) fireShot(shots int) error {
	return errors.New("персонаж в укрытии и не может атаковать")
}

func (as *DefendState) getCover() error {
	return errors.New("персонаж уже в укрытии")
}

func (as *DefendState) findBullets() error {
	as.gameHero.bulletCount += 30
	as.gameHero.setState(as.gameHero.attack)
	fmt.Println("Взят новый магазин")
	return nil
}

func StateWork() {
	newHero := newGameHero(30)
	err := newHero.fireShot(5)
	if err != nil {
		log.Println(err)
	}
	err = newHero.fireShot(30)
	if err != nil {
		log.Println(err)
	}
	err = newHero.findBullets()
	if err != nil {
		log.Println(err)
	}
	err = newHero.getCover()
	if err != nil {
		log.Println(err)
	}
	err = newHero.findBullets()
	if err != nil {
		log.Println(err)
	}
	err = newHero.fireShot(30)
	if err != nil {
		log.Println(err)
	}
}
