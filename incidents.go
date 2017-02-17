package main

import (
	"encoding/json"
	"log"
)




//list incidents belonging to team
//List incidents belonging to invidual member on the team
//Retrieve incident by ID

type IncidentResult struct {
	Incidents []Incident `json:"result"`
	count int
}

type Incident struct {
	Number string `json:"number"`
	SysCreatedBy string `json:"sys_created_by"`
	UIncidentType string `json:"u_incident_type"`
	Impact string `json:"impact"`
	Active string `json:"active"`
	Priority string `json:"priority"`
	ShortDescription string `json:"short_description"`
	ClosedBy string `json:"closed_by"`
	AssignedTo struct {
		       Link string `json:"link"`
		       Value string `json:"value"`
	       } `json:"assigned_to"`
	ULsmCustomerImpacting string `json:"u_lsm_customer_impacting"`
	UResolvedOn string `json:"u_resolved_on"`
	UCategoryTier1 string `json:"u_category_tier_1"`
	SysUpdatedBy string `json:"sys_updated_by"`
	UCategoryTier3 string `json:"u_category_tier_3"`
	UCategoryTier2 string `json:"u_category_tier_2"`
	SysCreatedOn string `json:"sys_created_on"`
	USLA string `json:"u_sla"`
	AssignmentGroup struct {
		       Link string `json:"link"`
		       Value string `json:"value"`
	       } `json:"assignment_group"`
	Urgency string `json:"urgency"`
	Severity string `json:"severity"`
}

type IncidentParams struct {
	limit string
	active bool
	teamID string
	incidentID string
	query string
}


func (c client) Incidents(p IncidentParams) (IncidentResult){
	gp := make(map[string]string)
	if p.teamID != "" {
		gp["assignment_group"] = p.teamID
	}
	gp["sysparm_limit"] = p.limit
	if p.limit != "" {
		gp["sysparm_limit"] = "100"
	}
	if p.active {
		gp["sysparm_query"] = "active=true"
	}
	if p.incidentID != "" {
		gp["number"] = p.incidentID
	}

	if gp["assignment_group"] == "" && gp["number"] == "" {
		log.Fatal("either teamID or incidentID must be provided")
	}

	IncidentRequest := getParams{}
	IncidentRequest.params = gp
	IncidentRequest.path = INCIDENTLISTPATH
	IncidentRequest.Client = c
	return IncidentRequest.Get().IncidentsData()
}

func (i Incident) AssignedUser() User {
	if i.AssignedTo.Value != "" {
		userInfo := serviceNow.User(i.AssignedTo.Value)
		if len(userInfo) > 0 {
			return userInfo[0]
		}
	}
	var u = User{"N/A","N/A","N/A","N/A","N/A","N/A","N/A","N/A","N/A"}
	return u
}

func (ir IncidentResult) DataPresent() bool {
	if ir.count > 0 {
		return true
	}
	return false
}

func (d returnData) IncidentsData() (res IncidentResult){
	err := json.Unmarshal(d, &res)
	if err != nil {
		log.Fatal("Could not unmarshall Incident response to struct",err)
	}
	res.count = len(res.Incidents)
        return
}