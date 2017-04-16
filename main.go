package main

import (
	"log"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
	"github.com/sjeanpierre/service_now_proxy/lib/servicenow/snclient"
	"github.com/sjeanpierre/service_now_proxy/lib/servicenow/snapi"
	"os"
)



func main() {
	router := mux.NewRouter().StrictSlash(true)
	snapi.RegisterHandlers(router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Println("Validating Service Now proxy Config")
	snclient.NewClient() //used to check if environment variables are defined
	log.Println("Started: Ready to serve")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter)) //todo, refactor to make port dynamic
}
