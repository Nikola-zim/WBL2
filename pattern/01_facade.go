package pattern

import (
	"errors"
	"fmt"
	"math"
	"time"
)

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/
// Фасад для работы с роботом
type RobotFacade struct {
	robotAuth  *robotAuth
	kinematics *kinematics
	sensors    *Sensors
}

func NewRobotFacade(adminLogin string, adminPassword string) *RobotFacade {
	robotAuth := newRobotAuth(adminLogin, adminPassword)
	kinematics := newKinematics(50, 10, robotAuth)
	return &RobotFacade{
		robotAuth:  robotAuth,
		kinematics: kinematics,
		sensors:    NewSensors(robotAuth, kinematics),
	}
}

func (rf *RobotFacade) MoveToPoint(coordinateX float64, coordinateY float64) error {
	err := rf.kinematics.singularityCheck(coordinateX, coordinateY)
	if err != nil {
		return err
	}
	// Преобразование и отправка координат
	rf.kinematics.sendCoordinates(rf.kinematics.coordinateConverter(coordinateX, coordinateY))
	// Проверка достижения точки
	realX, realY, err := rf.sensors.readEncoders()
	if err != nil {
		return err
	}
	if realX == coordinateX && realY == coordinateY {
		fmt.Println("Координаты достигнуты")
		return nil
	} else {
		return errors.New("координаты не достигнуты")
	}
}

// Подсистема робота
// Часть с авторизацией пользователя
type robotAuth struct {
	isAuthorized bool
	users        map[string]string
}

func newRobotAuth(adminLogin string, adminPassword string) *robotAuth {
	users := make(map[string]string)
	users[adminLogin] = adminPassword
	return &robotAuth{
		users: users,
	}
}

func (ru *robotAuth) robotUserLogin(login string, password string) bool {
	if ru.users[login] == password {
		fmt.Println("authorized on robot")
		ru.isAuthorized = true
		// Автоматический выход
		go func() {
			time.Sleep(30 * time.Second)
			ru.isAuthorized = false
		}()
		return true
	}
	return false
}

// Часть с управлением
type kinematics struct {
	*robotAuth
	theta1           float64
	theta2           float64
	wheelR           float64
	singularityParam float64
}

func newKinematics(singularityParam float64, wheelR float64, auth *robotAuth) *kinematics {
	if wheelR == 0 {
		return nil
	}
	robotAuth := auth
	return &kinematics{
		singularityParam: singularityParam,
		wheelR:           wheelR,
		robotAuth:        robotAuth,
	}
}

// Часть с проверкой на достигаемость заданных координат
func (k *kinematics) singularityCheck(coordinateX float64, coordinateY float64) error {
	module := math.Pow(math.Pow(coordinateX, 2)+math.Pow(coordinateY, 2), 1/2)
	if module >= k.singularityParam {
		return errors.New("неверные координаты, выход за допустимую область")
	} else {
		return nil
	}
}

// Часть с перевод координат в координаты приводов
func (k *kinematics) coordinateConverter(coordinateX float64, coordinateY float64) (float64, float64) {
	return coordinateX / k.wheelR, coordinateY / k.wheelR
}

// Отправка команды
func (k *kinematics) sendCoordinates(theta1 float64, theta2 float64) bool {
	if k.robotAuth.isAuthorized {
		fmt.Println("Отправка команды на движение")
		// Робот двигается
		k.theta1 = theta1
		k.theta2 = theta2
		return true
	} else {
		fmt.Println("Ошибка авторизации при отправке координат ")
		return true
	}
}

// Часть для работы с датчиками
type Sensors struct {
	*robotAuth
	*kinematics
}

func NewSensors(robotAuth *robotAuth, kinematics *kinematics) *Sensors {
	return &Sensors{
		kinematics: kinematics,
		robotAuth:  robotAuth,
	}
}

func (s *Sensors) readEncoders() (float64, float64, error) {
	if s.kinematics.robotAuth.isAuthorized {
		realX := s.kinematics.theta1 * s.wheelR
		realY := s.kinematics.theta2 * s.wheelR
		return realX, realY, nil
	}
	return 0, 0, errors.New("Нет доступа к датчикам, ошибка авториззации ")
}

// Работа "клиента" с роботом при помощи фасада
func FacadeWork() error {
	robotFacade := NewRobotFacade("admin", "admin")
	robotFacade.robotAuth.robotUserLogin("admin", "admin")
	err := robotFacade.MoveToPoint(3, 4)
	if err != nil {
		return err
	}
	return nil
}
