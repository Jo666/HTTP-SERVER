package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Food struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Server struct {
	*mux.Router

	eatingFood []Food
}

// notre constructeur
func NewServer() *Server {
	s := &Server{
		Router:     mux.NewRouter(),
		eatingFood: []Food{},
	}
	s.routes() // initialise les routes
	return s
}

// un peu comme une fct initialiser() qu'on mettra dans un constructeur
func (s *Server) routes() {
	s.HandleFunc("/eating-Food", s.listEatingFood()).Methods("GET")
	s.HandleFunc("/eating-Food", s.createEatingFood()).Methods("POST")
	s.HandleFunc("/eating-Food/{id}", s.removeEatingFood()).Methods("DELETE")

}

func (s *Server) createEatingFood() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i Food
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		i.ID = uuid.New()
		s.eatingFood = append(s.eatingFood, i)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(i); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listEatingFood() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.eatingFood); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) removeEatingFood() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, _ := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		for i, Food := range s.eatingFood {
			if Food.ID == id {
				s.eatingFood = append(s.eatingFood[:i], s.eatingFood[i+1:]...)
				break
			}
		}
	}
}
