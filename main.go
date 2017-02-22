package main

import (
	"log"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
	"crypto/x509"
)

var (
	host = os.Getenv("SERVICE_NOW_HOSTNAME")
        snClientID = os.Getenv("SERVICE_NOW_CLIENT_ID")
	snClientSecret = os.Getenv("SERVICE_NOW_CLIENT_SECRET")
	snUsername = os.Getenv("SERVICE_NOW_USERNAME")
	snPassword = os.Getenv("SERVICE_NOW_PASSWORD")
	serviceNow = client{}
	pool = &x509.CertPool{}
)


func init() {
	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)
	//client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool}}}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	registerHandlers(router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Println("Validating Config")
	NewClient() //used to check if environment variables are defined
	log.Println("Started: Ready to serve")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter)) //todo, refactor to make port dynamic
}

func registerHandlers(r *mux.Router)  {
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	r.HandleFunc(`/incidents/{incident:INC\d{7,10}}`, IncidentHandler)
	r.HandleFunc("/incidents/{option:count}/{team}", IncidentTeamHandler)
	r.HandleFunc("/incidents/{option:list}/{team}", IncidentTeamHandler)
}
