package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type App struct {
	DB *sql.DB
	Router *mux.Router
}

const (
	host     = "localhost"
	port     = 5432
	user     = "testgo"
	password = "123456"
	dbname   = "kinopoisk"
)

func (a *App) Initialize() {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	/*err = a.DB.Ping()
	if err != nil {
		fmt.Println("ZZZZZ")
		log.Fatal(err)
	} */

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	handlers := handlers.LoggingHandler(os.Stdout, a.Router)
	if err := http.ListenAndServe(addr, handlers); err != nil {
		log.Fatal(err)
	}
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/status", a.getStatus).Methods("GET")
}

func (a *App) getStatus( w http.ResponseWriter, r *http.Request) {
	minMaxData := MinMaxIds{}

	//result := make([]backup,0)
	defer r.Body.Close()

	result, err := minMaxData.getMinMaxIds(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, result)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}



