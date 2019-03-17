package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sjeanpierre/service_now_proxy/lib/servicenow/snapi"
	"github.com/sjeanpierre/service_now_proxy/lib/servicenow/snclient"
	"log"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	snapi.RegisterHandlers(router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Println("Validating Service Now proxy Config")
	snclient.NewClient() //used to check if environment variables are defined
	log.Println("Started: Ready to serve")
	go snclient.HydrateServiceCache()
	log.Fatal(http.ListenAndServe(":8080", loggedRouter)) //todo, refactor to make port dynamic
}
