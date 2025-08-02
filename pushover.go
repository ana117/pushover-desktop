package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type Notification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Icon    string `json:"icon"`
}

func sendNotification(title, message, icon string) error {
	if err := beeep.Notify(title, message, icon); err != nil {
		return fmt.Errorf("Failed to send notification: %w", err)
	}
	return nil
}
func sendAlert(title, message, icon string) error {
	if err := beeep.Alert(title, message, icon); err != nil {
		return fmt.Errorf("Failed to send alert: %w", err)
	}
	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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

	handler(fmt.Sprintf(":%s", port))
}

func handler(address string) error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Pushover Desktop Notification Service is running"))
	})

	r.Post("/notification", func(w http.ResponseWriter, r *http.Request) {
		var notification Notification
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request payload: " + err.Error()})
			return
		}
		if notification.Title == "" || notification.Message == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Title and message are required"})
			return
		}

		if err := sendNotification(notification.Title, notification.Message, notification.Icon); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Failed to send notification: " + err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Notification sent successfully"})
	})

	r.Post("/alert", func(w http.ResponseWriter, r *http.Request) {
		var notification Notification
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request payload: " + err.Error()})
			return
		}
		if notification.Title == "" || notification.Message == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Title and message are required"})
			return
		}

		if err := sendAlert(notification.Title, notification.Message, notification.Icon); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Failed to send alert: " + err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Alert sent successfully"})
	})

	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Printf("Error listening and serving!")
		log.Fatal(err)
		return err
	}
	return nil
}
