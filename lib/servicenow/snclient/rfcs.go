package snclient

import (
	"encoding/json"
	"log"
)

var serviceCacheStore = make(map[string]string)

type ServiceResult struct {
	Services []ServiceName `json:"result"`
}

type ServiceName struct {
	Name string `json:"u_product_name"`
}
type ChangeResult struct {
	Changes []Change `json:"result"`
	Count   int      `json:"change_count"`
}

type Change struct {
	Approval         string `json:"approval"`
	EndDate          string `json:"end_date"`
	Number           string `json:"number"`
	ShortDescription string `json:"short_description"`
	StartDate        string `json:"start_date"`
	State            string `json:"state"`
	SysCreatedBy     string `json:"sys_created_by"`
	SysCreatedOn     string `json:"sys_created_on"`
	SysID            string `json:"sys_id"`
	Type             string `json:"type"`
	UChangeReason    string `json:"u_change_reason"`
	Product          string
	UProductService  struct {
		//todo, lookup service by value
		Link  string `json:"link"`
		Value string `json:"value"`
	} `json:"u_product_service"`
	URfcNumber string `json:"u_rfc_number"`
	UStatus    string `json:"u_status"`
}

//in case I ever decide to clean this up
//https://stackoverflow.com/questions/26303694/json-marshalling-unmarshalling-same-struct-to-different-json-format-in-go

type ChangeParams struct {
	ChangeID   string
	ChangeGUID string
	Query      string
}

func (cr ChangeResult) DataPresent() bool {
	if cr.Count > 0 {
		return true
	}
	return false
}

func (c Client) Service(id string) string {
	//todo add simple caching here
	service, ok := serviceCacheStore[id]
	if ok {
		return service
	}
	p := make(map[string]string)
	p["sys_id"] = id
	r := getParams{path: SERVICENAMEPATH, params: p, Client: c}
	jsonResponse := r.Get()
	data := ServiceResult{}
	if err := json.Unmarshal(jsonResponse, &data); err != nil {
		log.Printf("Could not unmarshall service name from Service now response, %s", err.Error())
		return ""
	}
	service = data.Services[0].Name
	if service != "" {
		serviceCacheStore[id] = service
	}
	return service
}

func (c Client) Changes(p ChangeParams) ChangeResult {
	gp := make(map[string]string)
	if p.ChangeID != "" {
		gp["number"] = p.ChangeID
	}

	if p.ChangeGUID != "" {
		gp["sys_id"] = p.ChangeGUID
	}

	if gp["sys_id"] == "" && gp["number"] == "" {
		log.Fatal("either Change ID or Change Guid must be provided")
	}
	ChangeRequest := getParams{params: gp, path: CHANGEPATH, Client: c}
	return ChangeRequest.Get().ChangesData(c)
}

func (rd returnData) ChangesData(c Client) (res ChangeResult) {
	err := json.Unmarshal(rd, &res)
	if err != nil {
		log.Printf("Could not unmarshall Change response to struct - %+v\n", err)
		return
	}
	res.Count = len(res.Changes)
	for index, change := range res.Changes {
		res.Changes[index].Product = c.Service(change.UProductService.Value)
	}
	return
}
