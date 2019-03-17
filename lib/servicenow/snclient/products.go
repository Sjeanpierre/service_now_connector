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
	ID string `json:"sys_id"`
}

func (c Client) Service(id string) string {
	service, ok := serviceCacheStore[id]
	if ok {
		//log.Printf("Serving Product %s with %s from cache",service,id)
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

func HydrateServiceCache() {
	log.Println("Hydrating Product Cache")
	c := NewClient()
	p := make(map[string]string)
	r := getParams{path: SERVICENAMEPATH, params: p, Client: c}
	jsonResponse := r.Get()
	var data ServiceResult
	if err := json.Unmarshal(jsonResponse, &data); err != nil {
		log.Printf("Could not unmarshall service name from Service now response, %s", err.Error())
	}
	for _, service := range data.Services {
		//log.Printf("Adding %s - %s to cache",service.ID,service.Name)
		serviceCacheStore[service.ID] = service.Name
	}
	log.Printf("Product Cache Hydration completed, %d products added",len(data.Services))
}