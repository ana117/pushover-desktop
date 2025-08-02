package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

type Notification struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Icon     string `json:"icon,omitempty"`
	IconByte []byte `json:"-,omitempty"`
}

type AppConfig struct {
	Port    string
	AppName string
}

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendNotification(title, message string, icon any) error {
	if err := beeep.Notify(title, message, icon); err != nil {
		return fmt.Errorf("Failed to send notification: %w", err)
	}
	return nil
}

func (s *NotificationService) SendAlert(title, message string, icon any) error {
	if err := beeep.Alert(title, message, icon); err != nil {
		return fmt.Errorf("Failed to send alert: %w", err)
	}
	return nil
}

type APIError struct {
	Message string `json:"message"`
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

func parseNotificationRequest(r *http.Request, notification *Notification,
) error {
	if err := json.NewDecoder(r.Body).Decode(notification); err != nil {
		return fmt.Errorf("invalid request payload: %w", err)
	}

	if notification.Title == "" || notification.Message == "" {
		return fmt.Errorf("title and message are required")
	}

	if notification.Icon != "" {
		data, err := base64.StdEncoding.DecodeString(notification.Icon)
		if err == nil {
			notification.IconByte = data
		}
	}

	return nil
}

func (s *NotificationService) handleNotification(w http.ResponseWriter, r *http.Request) {
	var notification Notification
	if err := parseNotificationRequest(r, &notification); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, APIError{Message: err.Error()})
		return
	}

	if err := s.SendNotification(notification.Title, notification.Message, notification.IconByte); err != nil {
		writeJSONResponse(
			w,
			http.StatusInternalServerError,
			APIError{Message: "Failed to send notification: " + err.Error()},
		)
		return
	}

	writeJSONResponse(
		w,
		http.StatusOK,
		map[string]string{"message": "Notification sent successfully"},
	)
}

func (s *NotificationService) handleAlert(w http.ResponseWriter, r *http.Request) {
	var notification Notification
	if err := parseNotificationRequest(r, &notification); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, APIError{Message: err.Error()})
		return
	}

	if err := s.SendAlert(notification.Title, notification.Message, notification.IconByte); err != nil {
		writeJSONResponse(
			w,
			http.StatusInternalServerError,
			APIError{Message: "Failed to send alert: " + err.Error()},
		)
		return
	}

	writeJSONResponse(
		w,
		http.StatusOK,
		map[string]string{"message": "Alert sent successfully"},
	)
}

func loadConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "Pushover Desktop Notification Service"
	}
	beeep.AppName = appName

	return &AppConfig{
		Port:    port,
		AppName: appName,
	}, nil
}

func setupRouter(ns *NotificationService) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSONResponse(
			w,
			http.StatusOK,
			map[string]string{"message": "Pushover Desktop Notification Service is running"},
		)
	})

	r.Post("/notification", ns.handleNotification)
	r.Post("/alert", ns.handleAlert)

	return r
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	notificationService := NewNotificationService()
	router := setupRouter(notificationService)

	address := fmt.Sprintf(":%s", config.Port)
	log.Printf("Server starting on %s...", address)

	err = http.ListenAndServe(address, router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
