package routes

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"singel-sondag/config"
	"singel-sondag/domain"
	"strconv"
)

type registerDTO struct {
	Name    string            `json:"name"`
	Email   string            `json:"email"`
	Answers map[string]string `json:"answers"`
}

func InviteHandler(regRepo domain.RegisterRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := regRepo.GetEvent(r.Context())
		if err != nil {
			http.Error(w, "could not load event", http.StatusInternalServerError)
			return
		}
		tmpl, err := template.ParseFiles("web/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, event); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func RegisterHandler(regRepo domain.RegisterRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto registerDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if dto.Name == "" || dto.Email == "" {
			http.Error(w, "name and email are required", http.StatusBadRequest)
			return
		}
		event, err := regRepo.GetEvent(r.Context())
		if err != nil {
			http.Error(w, "could not validate registration", http.StatusInternalServerError)
			return
		}
		allowed := make(map[string]bool, len(event.Questions))
		for _, question := range event.Questions {
			allowed[question.Label] = true
			if question.Required && dto.Answers[question.Label] == "" {
				http.Error(w, "all required questions must be answered", http.StatusBadRequest)
				return
			}
		}
		for answer := range dto.Answers {
			if !allowed[answer] {
				delete(dto.Answers, answer)
			}
		}
		if err := regRepo.Register(r.Context(), dto.Name, dto.Email, dto.Answers); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func GetEventHandler(regRepo domain.RegisterRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := regRepo.GetEvent(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(event)
	}
}

func UpdateEventHandler(regRepo domain.RegisterRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event domain.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "invalid event", http.StatusBadRequest)
			return
		}
		if event.Title == "" {
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}
		if err := regRepo.UpdateEvent(r.Context(), event); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		GetEventHandler(regRepo)(w, r)
	}
}

func GetRegistrationsHandler(regRepo domain.RegisterRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registrations, err := regRepo.GetRegistrations(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(registrations)
	}
}

func DeleteRegistrationHandler(regRepo domain.RegisterRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if err := regRepo.DeleteRegistration(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteAllRegistrationsHandler(regRepo domain.RegisterRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := regRepo.DeleteRegistrations(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func AdminPageHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.AdminPassword == "" {
			http.Error(w, "admin password is not configured", http.StatusNotImplemented)
			return
		}
		tmpl, err := template.ParseFiles("web/admin.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, map[string]bool{"Authenticated": validSession(r, cfg)})
	}
}
