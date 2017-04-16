package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/sjeanpierre/service_now_connector/lib/servicenow/snclient"
)

type Response struct {
	Type    string
	Message string
	Data    interface{}
}

func JSONResponseHandler(w http.ResponseWriter, returnval interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnval)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Type:"error", Message:fmt.Sprintf("Route %s not found, " +
		"please check request and try again", r.URL.Path)})
}

func IncidentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	v := fmt.Sprintf("%+v", vars)
	incidentID := vars["incident"]
	serviceNow := snclient.NewClient()
	singleIncidentParams := snclient.IncidentParams{Active:false, IncidentID:incidentID, Limit:"100"}
	singleIncident := serviceNow.Incidents(singleIncidentParams)
	ret := Response{Type:"response",Message:v, Data:singleIncident}
	JSONResponseHandler(w, ret) //todo, return 404 if incident lookup did not yield results
}

func IncidentTeamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	v := fmt.Sprintf("%+v", vars)
	teamID := vars["team"]
	serviceNow := snclient.NewClient()
	teamIncidentListParams := snclient.IncidentParams{Active:true, TeamID: teamID, Limit:"100"}
	teamIncidentList := serviceNow.Incidents(teamIncidentListParams)
	if vars["option"] == "count" {
		ret := Response{Type:"response",Message:v, Data:map[string]string{"count":strconv.Itoa(teamIncidentList.Count)}}
		JSONResponseHandler(w, ret)
		return
	}
	if vars["option"] == "list" {
		ret := Response{Type:"response",Message:v, Data:teamIncidentList}
		JSONResponseHandler(w, ret)
		return
	}
	  notFoundHandler(w,r)
}
