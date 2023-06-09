package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

// Структура события
type Event struct {
	UserID      int    `json:"user_id"`
	ID          int    `json:"id"`
	Date        MyDate `json:"date"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Структура для получения конкретного события
type ConcreteEvent struct {
	UserID int `json:"user_id"`
	ID     int `json:"id"`
}

// Обертка для даты
type MyDate struct {
	date time.Time
}

// Хранилище для событий
type EventStore struct {
	m map[int][]Event
	*sync.RWMutex
}

// Собираем вместе сервер, логер и хранилище
type Scope struct {
	srv        *http.ServeMux
	logger     Logger
	eventStore EventStore
}

// Логер
type Logger struct {
	*log.Logger
}

// Инициалзирует scope
func NewScope() *Scope {
	return &Scope{
		srv: http.NewServeMux(),
		logger: Logger{
			log.New(os.Stdout, "logger: ", log.Lshortfile),
		},
		eventStore: EventStore{
			m:       make(map[int][]Event),
			RWMutex: new(sync.RWMutex),
		},
	}
}

// Настраивает роуты и запускает сервер
func (s *Scope) startServer() {
	s.srv.HandleFunc("/create_event", s.CreateEvent)
	s.srv.HandleFunc("/update_event", s.UpdateEvent)
	s.srv.HandleFunc("/delete_event", s.DeleteEvent)

	s.srv.HandleFunc("/events_for_day", s.EventsForDayAdd)
	s.srv.HandleFunc("/events_for_week", s.EventsForWeekAdd)
	s.srv.HandleFunc("/events_for_month", s.EventsForMonthAdd)

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), s.srv))
}

// Отправляет заданную ошибку с заданным кодов статуса
func sendError(w http.ResponseWriter, errorStr string, status int) {
	response := struct {
		Error string `json:"error"`
	}{errorStr}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResponse)
}

// Отправляет результат запроса
func sendResult(w http.ResponseWriter, resStr string, events []Event, status int) {
	response := struct {
		Result string  `json:"result"`
		Events []Event `json:"events"`
	}{resStr, events}
	responseJson, err := json.Marshal(response)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseJson)
}

// Проверяет содержится ли событие в хранилище
func (s *Scope) containsEvent(e Event) bool {
	for _, v := range s.eventStore.m[e.UserID] {
		if v.ID == e.ID {
			return true
		}
	}
	return false
}

// Создает новое событие и сохраняет его в хранилище
func (s *Scope) CreateNewEvent(event Event) error {
	s.eventStore.RWMutex.Lock()
	defer s.eventStore.RWMutex.Unlock()
	if s.containsEvent(event) {
		return errors.New("duplicate event not allowed")
	} else {
		s.eventStore.m[event.UserID] = append(s.eventStore.m[event.UserID], event)
	}
	return nil
}

// Обработчик /create_event
func (s *Scope) CreateEvent(w http.ResponseWriter, r *http.Request) {
	s.logger.Println(r.URL)
	if r.Method != http.MethodPost {
		sendError(w, "Not correct method", http.StatusBadRequest)
		return
	}

	event, err := parseJson(r)

	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.CreateNewEvent(event); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
	}
	sendResult(w, "Success", []Event{event}, http.StatusCreated)
}

// Обновляет существующее событие
func (s *Scope) UpdateEventFunc(e Event) error {
	s.eventStore.RWMutex.Lock()
	defer s.eventStore.RWMutex.Unlock()

	events := s.eventStore.m[e.UserID]
	for ind := 0; ind < len(events); ind++ {
		if events[ind].ID == e.ID {
			events[ind].Title = e.Title
			events[ind].Description = e.Description
			events[ind].Date = e.Date
			return nil
		}
	}
	return errors.New("event not found")
}

// Обработчик /update_event
func (s *Scope) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	s.logger.Println(r.URL)
	if r.Method != http.MethodPost {
		sendError(w, "Not correct method", http.StatusBadRequest)
		return
	}
	event, err := parseJson(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if EventIsValid(event) {
		err = s.UpdateEventFunc(event)
		if err != nil {
			sendError(w, "Event not found", http.StatusInternalServerError)
		} else {
			sendResult(w, "Success", []Event{event}, http.StatusOK)
		}
	} else {
		sendError(w, "Not valid event", http.StatusBadRequest)
	}
}

// Удаляет существующее событие
func (s *Scope) DeleteEventFunc(cEvent ConcreteEvent) error {
	s.eventStore.RWMutex.Lock()
	defer s.eventStore.RWMutex.Unlock()

	events := s.eventStore.m[cEvent.UserID]
	for ind := 0; ind < len(events); ind++ {
		if events[ind].ID == cEvent.ID {
			s.eventStore.m[cEvent.UserID] = append(s.eventStore.m[cEvent.UserID][0:ind], s.eventStore.m[cEvent.UserID][ind+1:]...)
			return nil
		}
	}
	return errors.New("event not found")
}

// Обработчик /delete_event
func (s *Scope) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	s.logger.Println(r.URL)
	if r.Method != http.MethodPost {
		sendError(w, "Not correct method", http.StatusBadRequest)
		return
	}

	var cEvent ConcreteEvent
	err := json.NewDecoder(r.Body).Decode(&cEvent)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.DeleteEventFunc(cEvent)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResult(w, "Success", nil, http.StatusOK)
}

// Возвращает события на заданный день
func (s *Scope) EventsForDayGetter(userID int, date time.Time) ([]Event, error) {
	s.eventStore.RWMutex.RLock()
	s.eventStore.RWMutex.RUnlock()

	var result []Event

	var allUserEvents []Event
	allUserEvents = s.eventStore.m[userID]
	if allUserEvents == nil {
		return nil, errors.New("unknown user_id")
	}

	for _, event := range allUserEvents {
		if event.Date.date.Year() == date.Year() &&
			event.Date.date.Month() == date.Month() &&
			event.Date.date.Day() == date.Day() {
			result = append(result, event)
		}
	}
	return result, nil
}

// Обработчик для /events_for_day
func (s *Scope) EventsForDayAdd(w http.ResponseWriter, r *http.Request) {
	s.logger.Println(r.URL)
	if r.Method != http.MethodGet {
		sendError(w, "Not correct method", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		sendError(w, "Incorrect args", http.StatusBadRequest)
		return
	}
	events, err := s.EventsForDayGetter(userID, date)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendResult(w, "Success", events, http.StatusOK)
}

// Возвращает события с разницей не более чем в неделю от заданной даты
func (s *Scope) EventsForWeekGetter(userID int, date time.Time) ([]Event, error) {
	s.eventStore.RWMutex.RLock()
	s.eventStore.RWMutex.RUnlock()

	var result []Event

	var allUserEvents []Event
	allUserEvents = s.eventStore.m[userID]
	if allUserEvents == nil {
		return nil, errors.New("unknown user_id")
	}

	for _, event := range allUserEvents {
		difference := date.Sub(event.Date.date)
		if difference < 0 {
			difference = -difference
		}
		if difference <= time.Duration(7*24)*time.Hour {
			result = append(result, event)
		}
	}
	return result, nil
}

// Обработчик для /events_for_week
func (s *Scope) EventsForWeekAdd(w http.ResponseWriter, r *http.Request) {
	s.logger.Println(r.URL)
	if r.Method != http.MethodGet {
		sendError(w, "Not correct method", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		sendError(w, "Incorrect args", http.StatusBadRequest)
		return
	}
	events, err := s.EventsForWeekGetter(userID, date)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendResult(w, "Success", events, http.StatusOK)
}

// Возвращает события в заданном месяце
func (s *Scope) EventsForMonthGetter(userID int, date time.Time) ([]Event, error) {
	s.eventStore.RWMutex.RLock()
	s.eventStore.RWMutex.RUnlock()

	var result []Event

	var allUserEvents []Event
	allUserEvents = s.eventStore.m[userID]
	if allUserEvents == nil {
		return nil, errors.New("unknown user_id")
	}

	for _, event := range allUserEvents {
		if event.Date.date.Year() == date.Year() || event.Date.date.Month() == date.Month() {
			result = append(result, event)
		}
	}
	return result, nil
}

// Обработчик для /events_for_month
func (s *Scope) EventsForMonthAdd(w http.ResponseWriter, r *http.Request) {
	s.logger.Println(r.URL)
	if r.Method != http.MethodGet {
		sendError(w, "Not correct method", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		sendError(w, "Incorrect args", http.StatusBadRequest)
		return
	}
	events, err := s.EventsForMonthGetter(userID, date)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendResult(w, "Success", events, http.StatusOK)
}

// Проверяет корректность заданного события
func EventIsValid(event Event) bool {
	if event.ID <= 0 || event.UserID <= 0 || event.Title == "" || event.Description == "" {
		return false
	}
	return true
}

// Парсит событие из тела запроса
func parseJson(r *http.Request) (Event, error) {
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return event, errors.New("cannot decode json")
	}
	return event, nil
}

// Unmarshal() для типа MyDate
func (d *MyDate) UnmarshalJSON(input []byte) error {
	var err error
	d.date, err = time.Parse(`"2006-01-02"`, string(input))
	return err
}

// String() для типа MyDate
func (d MyDate) String() string {
	return d.date.String()
}

// Marhsal для типа MyDate
func (d *MyDate) MarshalJSON() ([]byte, error) {
	dateStr := d.date.Format("2006-01-02")
	return json.Marshal(dateStr)
}

func main() {
	scope := NewScope()
	scope.startServer()
}
