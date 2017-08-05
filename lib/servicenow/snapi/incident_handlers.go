package snapi

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"github.com/sjeanpierre/service_now_proxy/lib/servicenow/snclient"
	"strconv"
	"log"
)

func IncidentHandler(w http.ResponseWriter, r *http.Request,isGuid bool) {
	singleIncidentParams := snclient.IncidentParams{}
	vars := mux.Vars(r)
	v := fmt.Sprintf("%+v", vars)
	incidentID := vars["incident"]
	if isGuid {
		singleIncidentParams = snclient.IncidentParams{Active:false, IncidentGUID:incidentID, Limit:"100"}

	} else {
		singleIncidentParams = snclient.IncidentParams{Active:false, IncidentID:incidentID, Limit:"100"}
	}

	serviceNow := snclient.NewClient()
	singleIncident := serviceNow.Incidents(singleIncidentParams)
	log.Println("%+v",singleIncident)
	ret := Response{Type:"response",Message:v, Data:singleIncident}
	if singleIncident.DataPresent() {
		JSONResponseHandler(w, ret)
		return
	}
	resourceNotFoundHandler(w,r)
}

func IncidentFromNumber(w http.ResponseWriter, r *http.Request) {IncidentHandler(w,r,false)}

func IncidentFromGUID(w http.ResponseWriter, r *http.Request) {IncidentHandler(w,r,true)}

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
