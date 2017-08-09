package snapi

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sjeanpierre/service_now_proxy/lib/servicenow/snclient"
	"log"
	"net/http"
)

func ChangeHandler(w http.ResponseWriter, r *http.Request, isGuid bool) {
	singleChangeParams := snclient.ChangeParams{}
	vars := mux.Vars(r)
	v := fmt.Sprintf("%+v", vars)
	ChangeID := vars["change"]
	if isGuid {
		singleChangeParams = snclient.ChangeParams{ChangeGUID: ChangeID}

	} else {
		singleChangeParams = snclient.ChangeParams{ChangeID: ChangeID}
	}

	serviceNow := snclient.NewClient()
	singleChange := serviceNow.Changes(singleChangeParams)
	log.Println("%+v", singleChange)
	ret := Response{Type: "response", Message: v, Data: singleChange}
	if singleChange.DataPresent() {
		JSONResponseHandler(w, ret)
		return
	}
	resourceNotFoundHandler(w, r)
}

func ChangeFromNumber(w http.ResponseWriter, r *http.Request) { ChangeHandler(w, r, false) }

func ChangeFromGUID(w http.ResponseWriter, r *http.Request) { ChangeHandler(w, r, true) }
