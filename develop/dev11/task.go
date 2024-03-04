package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

/*
HTTP-сервер

Реализовать HTTP-сервер для работы с календарем.
В рамках задания необходимо работать строго со стандартной HTTP-библиотекой.

В рамках задания необходимо:
Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
Реализовать middleware для логирования запросов

Методы API:
POST /create_event
POST /update_event
POST /delete_event
GET /events_for_day
GET /events_for_week
GET /events_for_month


Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON-документ содержащий либо {"result": "..."}
в случае успешного выполнения метода, либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
Реализовать все методы.
Бизнес логика НЕ должна зависеть от кода HTTP сервера.
В случае ошибки бизнес-логики сервер должен возвращать HTTP 503.
В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400.
В случае остальных ошибок сервер должен возвращать HTTP 500.
Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
*/

// Config содержит описание конфигурационного файла сервера
type Config struct {
	Address string `json:"server_address"`
}

// Status соответствует статусу события
type Status int

const (
	// Created соответствует созданному событию
	Created Status = iota
	// Updated соответствует обновленному событию
	Updated
	// Deleted соответствует удаленному событию
	Deleted
)

// PostResult возвращается в API модификации событий
type PostResult struct {
	ID     string `json:"id"`
	Status Status `json:"status"`
}

// EventResult возвращается в API поиска событий
type EventResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`
}

// Response это формат ответа API модификации событий
type Response struct {
	Result any `json:"result"`
}

// ErrorResponse это формат ошибочного ответа API
type ErrorResponse struct {
	Error string `json:"error"`
}

// ValidationError структура для ошибки валидации параметров
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ParseEvent разбирает переданные параметры event и возвращает ссылку на Event, date - строка в формате 2019-09-09
func ParseEvent(v url.Values) (*Event, error) {
	event := Event{}
	var err error
	for key, value := range v {
		switch key {
		case "user_id":
			event.UserID = value[0]
		case "id":
			event.ID = value[0]
		case "name":
			event.Name = value[0]
		case "date":
			event.Date, err = time.Parse("2006-01-02", value[0])
			if err != nil {
				return nil, fmt.Errorf("date parse error: %w", err)
			}
		}
	}
	return &event, nil
}

// ParseUserAndDate парсит id пользователя и дату события из query
func ParseUserAndDate(v url.Values) (userID string, date time.Time, err error) {
	for key, value := range v {
		switch key {
		case "user_id":
			userID = value[0]
		case "date":
			date, err = time.Parse("2006-01-02", value[0])
			if err != nil {
				return "", time.Time{}, fmt.Errorf("date parse error: %w", err)
			}
		}
	}
	return
}

// ParseUserAndMonth парсит id пользователя и год и месяц события из query
func ParseUserAndMonth(v url.Values) (userID string, year int, month time.Month, err error) {
	for key, value := range v {
		switch key {
		case "user_id":
			userID = value[0]
		case "year":
			year, err = strconv.Atoi(value[0])
			if err != nil {
				return "", 0, 0, fmt.Errorf("year parse error: %w", err)
			}
		case "month":
			monthNum, err := strconv.Atoi(value[0])
			if err != nil {
				return "", 0, 0, fmt.Errorf("month parse error: %w", err)
			}
			if monthNum < 1 || monthNum > 12 {
				return "", 0, 0, fmt.Errorf("month parse error")
			}
			month = time.Month(monthNum)
		}
	}
	return
}

// Event это внутреннее представление события
type Event struct {
	UserID string
	ID     string
	Name   string
	Date   time.Time
}

// UserCalendar хранит события одного пользователя
type UserCalendar map[string]Event

// Storage хранит календари в памяти
type Storage struct {
	// user -> event ID -> event
	events map[string]UserCalendar
}

// NewStorage возвращает новый storage
func NewStorage() *Storage {
	return &Storage{events: make(map[string]UserCalendar)}
}

// Create создает новое событие
func (s *Storage) Create(event *Event) (*Event, error) {
	if event.UserID == "" || event.Name == "" || event.Date.IsZero() {
		return nil, &ValidationError{Message: "empty parameters"}
	}
	event.ID = uuid.New().String()
	calendar, ok := s.events[event.UserID]
	if !ok {
		calendar = make(map[string]Event)
		s.events[event.UserID] = calendar
	}
	calendar[event.ID] = *event
	return event, nil
}

// Update обновляет существующее событие
func (s *Storage) Update(event *Event) (*Event, error) {
	if event.ID == "" || event.UserID == "" || event.Name == "" || event.Date.IsZero() {
		return nil, &ValidationError{Message: "empty parameters"}
	}
	calendar, ok := s.events[event.UserID]
	if !ok {
		return nil, &ValidationError{Message: "UserID does not exist"}
	}
	_, ok = calendar[event.ID]
	if !ok {
		return nil, &ValidationError{Message: "Event does not exist"}
	}
	calendar[event.ID] = *event
	return event, nil
}

// Delete удаляет существующее событие
func (s *Storage) Delete(event *Event) (*Event, error) {
	if event.ID == "" || event.UserID == "" {
		return nil, &ValidationError{Message: "empty parameters"}
	}
	calendar, ok := s.events[event.UserID]
	if !ok {
		return nil, &ValidationError{Message: "UserID does not exist"}
	}
	_, ok = calendar[event.ID]
	if !ok {
		return nil, &ValidationError{Message: "Event does not exist"}
	}
	delete(calendar, event.ID)
	return event, nil
}

// GetEventsPerDay возвращает события в заданный день
func (s *Storage) GetEventsPerDay(userID string, date time.Time) ([]Event, error) {
	if userID == "" || date.IsZero() {
		return nil, &ValidationError{Message: "empty parameters"}
	}
	var res []Event
	calendar, ok := s.events[userID]
	if !ok {
		return nil, &ValidationError{Message: "UserID does not exist"}
	}
	for _, event := range calendar {
		if event.Date == date {
			res = append(res, event)
		}
	}
	return res, nil
}

// GetEventsPerWeek возвращает события в заданную неделю
func (s *Storage) GetEventsPerWeek(userID string, startDate time.Time) ([]Event, error) {
	if userID == "" || startDate.IsZero() {
		return nil, &ValidationError{Message: "empty parameters"}
	}
	var res []Event
	calendar, ok := s.events[userID]
	if !ok {
		return nil, &ValidationError{Message: "UserID does not exist"}
	}
	endDate := startDate.Add(time.Hour * 24 * 7)
	for _, event := range calendar {
		if (!event.Date.Before(startDate)) && (event.Date.Before(endDate)) {
			res = append(res, event)
		}
	}
	return res, nil
}

// GetEventsPerMonth возвращает события в заданный месяц
func (s *Storage) GetEventsPerMonth(userID string, year int, month time.Month) ([]Event, error) {
	if userID == "" || year == 0 || month == 0 {
		return nil, &ValidationError{Message: "empty parameters"}
	}
	var res []Event
	calendar, ok := s.events[userID]
	if !ok {
		return nil, &ValidationError{Message: "UserID does not exist"}
	}
	for _, event := range calendar {
		if event.Date.Year() == year && event.Date.Month() == month {
			res = append(res, event)
		}
	}
	return res, nil
}

type loggingResponseWriter struct {
	w          http.ResponseWriter
	statusCode int
}

func (l *loggingResponseWriter) Header() http.Header {
	return l.w.Header()
}

func (l *loggingResponseWriter) Write(bytes []byte) (int, error) {
	return l.w.Write(bytes)
}

func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	l.statusCode = statusCode
	l.w.WriteHeader(statusCode)
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		nextW := &loggingResponseWriter{w: w}
		next.ServeHTTP(nextW, r)
		fmt.Fprintf(os.Stdout, "Request %s %s processed in %v, code %d\n", r.Method, r.URL.String(), time.Since(start), nextW.statusCode)
	})
}

func validatePostRequest(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		writeErrorMessage(w, http.StatusMethodNotAllowed, "Wrong method")
		return false
	}
	// Парсим тело запроса
	err := r.ParseForm()
	if err != nil {
		writeErrorMessage(w, http.StatusBadRequest, "Failed to parse form")
		return false
	}
	return true
}

func createEvent(w http.ResponseWriter, r *http.Request, storage *Storage) {
	if !validatePostRequest(w, r) {
		return
	}
	event, err := ParseEvent(r.PostForm)
	if err != nil {
		writeErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	event, err = storage.Create(event)
	if err != nil {
		writeError(w, err)
		return
	}
	response := Response{Result: PostResult{
		ID:     event.ID,
		Status: Created,
	}}
	marshalResponseAndWrite(w, http.StatusOK, response)
}

func updateEvent(w http.ResponseWriter, r *http.Request, storage *Storage) {
	if !validatePostRequest(w, r) {
		return
	}
	event, err := ParseEvent(r.PostForm)
	if err != nil {
		writeErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	event, err = storage.Update(event)
	if err != nil {
		writeError(w, err)
		return
	}
	response := Response{Result: PostResult{
		ID:     event.ID,
		Status: Updated,
	}}
	marshalResponseAndWrite(w, http.StatusOK, response)
}

func deleteEvent(w http.ResponseWriter, r *http.Request, storage *Storage) {
	if !validatePostRequest(w, r) {
		return
	}
	event, err := ParseEvent(r.PostForm)
	if err != nil {
		writeErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	event, err = storage.Delete(event)
	if err != nil {
		writeError(w, err)
		return
	}
	response := Response{Result: PostResult{
		ID:     event.ID,
		Status: Deleted,
	}}
	marshalResponseAndWrite(w, http.StatusOK, response)
}

func getEventsPerDay(w http.ResponseWriter, r *http.Request, storage *Storage) {
	if r.Method != http.MethodGet {
		writeErrorMessage(w, http.StatusMethodNotAllowed, "Wrong method")
		return
	}

	userID, date, err := ParseUserAndDate(r.URL.Query())
	if err != nil {
		writeErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	events, err := storage.GetEventsPerDay(userID, date)
	if err != nil {
		writeError(w, err)
		return
	}
	writeEventsResponse(w, events)
}

func getEventsPerWeek(w http.ResponseWriter, r *http.Request, storage *Storage) {
	if r.Method != http.MethodGet {
		writeErrorMessage(w, http.StatusMethodNotAllowed, "Wrong method")
		return
	}

	userID, date, err := ParseUserAndDate(r.URL.Query())
	if err != nil {
		writeErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	events, err := storage.GetEventsPerWeek(userID, date)
	if err != nil {
		writeError(w, err)
		return
	}
	writeEventsResponse(w, events)
}

func getEventsPerMonth(w http.ResponseWriter, r *http.Request, storage *Storage) {
	if r.Method != http.MethodGet {
		writeErrorMessage(w, http.StatusMethodNotAllowed, "Wrong method")
		return
	}

	userID, year, month, err := ParseUserAndMonth(r.URL.Query())
	if err != nil {
		writeErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	events, err := storage.GetEventsPerMonth(userID, year, month)
	if err != nil {
		writeError(w, err)
		return
	}
	writeEventsResponse(w, events)
}

func writeEventsResponse(w http.ResponseWriter, events []Event) {
	respEvents := make([]EventResult, len(events))
	for i, e := range events {
		respEvents[i] = EventResult{
			ID:   e.ID,
			Name: e.Name,
			Date: e.Date.Format("2006-01-02"),
		}
	}
	response := Response{Result: respEvents}
	marshalResponseAndWrite(w, http.StatusOK, response)
}

func marshalResponseAndWrite(w http.ResponseWriter, status int, response any) {
	respJSON, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while serializing response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(respJSON)
}

func writeError(w http.ResponseWriter, err error) {
	if err, ok := err.(*ValidationError); ok {
		writeErrorMessage(w, http.StatusBadRequest, err.Message)
	} else {
		fmt.Fprintf(os.Stderr, "Internal error while processing request: %v\n", err)
		writeErrorMessage(w, http.StatusServiceUnavailable, "Service unavailable")
	}
}

func writeErrorMessage(w http.ResponseWriter, statusCode int, message string) {
	resp := ErrorResponse{Error: message}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while serializing response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(respJSON)
}

func getHandler() http.Handler {
	storage := NewStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("/create_event/", func(w http.ResponseWriter, r *http.Request) {
		createEvent(w, r, storage)
	})
	mux.HandleFunc("/update_event/", func(w http.ResponseWriter, r *http.Request) {
		updateEvent(w, r, storage)
	})
	mux.HandleFunc("/delete_event/", func(w http.ResponseWriter, r *http.Request) {
		deleteEvent(w, r, storage)
	})
	mux.HandleFunc("/events_for_day/", func(w http.ResponseWriter, r *http.Request) {
		getEventsPerDay(w, r, storage)
	})
	mux.HandleFunc("/events_for_week/", func(w http.ResponseWriter, r *http.Request) {
		getEventsPerWeek(w, r, storage)
	})
	mux.HandleFunc("/events_for_month/", func(w http.ResponseWriter, r *http.Request) {
		getEventsPerMonth(w, r, storage)
	})
	return loggingHandler(mux)
}

func main() {
	// Определяем название конфиг файла
	configName := flag.String("c", "conf.json", "name of config file")
	flag.Parse()

	confData, err := os.ReadFile(*configName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading config file: %v\n", err)
		os.Exit(1)
	}

	var cfg Config
	err = json.Unmarshal(confData, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while parsing config file: %v\n", err)
		os.Exit(1)
	}

	handler := getHandler()
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: handler,
	}

	//Обрабатываем сигналы для корректного завершения
	signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()
	go func() {
		<-signalCtx.Done()
		if er := server.Shutdown(context.Background()); er != nil {
			fmt.Fprintln(os.Stderr, "Failed to shutdown server: ", er)
		}
	}()

	fmt.Fprintln(os.Stdout, "Started server")
	if er := server.ListenAndServe(); er != http.ErrServerClosed {
		fmt.Fprintln(os.Stderr, "HTTP server ListenAndServe: ", er)
	}
	fmt.Fprintln(os.Stdout, "Stopped server")
}
