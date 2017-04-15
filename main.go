package main

import (
	"log"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
	"github.com/sjeanpierre/service_now_connector/lib/servicenow/snclient"
)



func main() {
	router := mux.NewRouter().StrictSlash(true)
	registerHandlers(router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Println("Validating Config")
	snclient.NewClient() //used to check if environment variables are defined
	log.Println("Started: Ready to serve")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter)) //todo, refactor to make port dynamic
}

func registerHandlers(r *mux.Router)  {
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	r.HandleFunc(`/incidents/{incident:INC\d{7,10}}`, IncidentHandler)
	r.HandleFunc("/incidents/{option:count}/{team}", IncidentTeamHandler)
	r.HandleFunc("/incidents/{option:list}/{team}", IncidentTeamHandler)
}
