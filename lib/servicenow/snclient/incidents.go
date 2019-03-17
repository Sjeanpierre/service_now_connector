package snclient

import (
	"encoding/json"
	"log"
	"math"
	"strings"
	"time"
)

const IncidentTimeFormat = "2006-01-02 15:04:05"

var (
	UrgencyOptions    = map[string]string{"1": "critical", "2": "high", "3": "medium", "4": "low", "5": "very low"}
	PriorityOptions   = map[string]string{"1": "critical", "2": "high", "3": "moderate", "4": "low", "5": "requests"}
	productCacheStore = make(map[string]string)
)

type IncidentResult struct {
	Incidents []Incident `json:"result"`
	Count     int        `json:"incident_count"`
}

type SNLink struct {
	Link  string `json:"link"`
	Value string `json:"value"`
}

type Incident struct {
	Number                string          `json:"number"`
	SysCreatedBy          string          `json:"sys_created_by"`
	UIncidentType         string          `json:"u_incident_type"`
	IncidentState         string          `json:"incident_state"`
	Impact                string          `json:"impact"`
	Active                string          `json:"active"`
	Priority              string          `json:"priority"`
	ShortDescription      string          `json:"short_description"`
	TicketID              string          `json:"sys_id"`
	ClosedBy              string          `json:"closed_by"`
	AssignedToRaw         json.RawMessage `json:"assigned_to,omitempty"` //todo,demote this to unexported value
	ULsmCustomerImpacting string          `json:"u_lsm_customer_impacting"`
	UResolvedOn           string          `json:"u_resolved_on"`
	UCategoryTier1        string          `json:"u_category_tier_1"`
	SysUpdatedBy          string          `json:"sys_updated_by"`
	UCategoryTier3        string          `json:"u_category_tier_3"`
	UCategoryTier2        string          `json:"u_category_tier_2"`
	SysCreatedOn          string          `json:"sys_created_on"`
	USLA                  string          `json:"u_sla"`
	AssignmentGroup       SNLink          `json:"assignment_group,omitempty"`
	Urgency               string          `json:"urgency"`
	Severity              string          `json:"severity"`
	LSMAssigned           User            `json:"lsm_assigned"`
	ProductName           string          `json:"product_name"`
	ProductRaw            json.RawMessage `json:"u_live_services_product,omitempty"`
}

type IncidentParams struct {
	Limit        string
	Active       bool
	TeamID       string
	IncidentID   string
	IncidentGUID string
	Query        string
}

func (c Client) Incidents(p IncidentParams) IncidentResult {
	gp := make(map[string]string)
	if p.TeamID != "" {
		gp["assignment_group"] = p.TeamID
	}
	gp["sysparm_limit"] = p.Limit
	if p.Limit == "" {
		gp["sysparm_limit"] = "100"
	}
	if p.Active {
		gp["sysparm_query"] = "active=true"
	}
	if p.IncidentID != "" {
		gp["number"] = p.IncidentID
	}

	if p.IncidentGUID != "" {
		gp["sys_id"] = p.IncidentGUID
	}

	if gp["sys_id"] == "" && gp["assignment_group"] == "" && gp["number"] == "" {
		log.Fatal("either teamID or incidentID must be provided")
	}
	IncidentRequest := getParams{params: gp, path: INCIDENTLISTPATH, Client: c}
	return IncidentRequest.Get().IncidentsData(c)
}

func (i Incident) AssignedUser(c Client) User {
	if string(i.AssignedToRaw) == `""` {
		return noUser
	}

	u := SNLink{}
	err := json.Unmarshal(i.AssignedToRaw, &u)
	if err != nil {
		log.Printf("Could not parse Assigned to details, %+v", string(i.AssignedToRaw))
		return noUser
	}
	return c.User(u.Value)
}

func (i Incident) ImpactedProduct(c Client) string {
	if string(i.ProductRaw) == `""` {
		return ""
	}

	u := SNLink{}
	err := json.Unmarshal(i.ProductRaw, &u)
	if err != nil {
		log.Printf("Could not parse Assigned to details, %+v", string(i.ProductRaw))
		return ""
	}
	return c.Service(u.Value)
}

func (ir IncidentResult) DataPresent() bool {
	if ir.Count > 0 {
		return true
	}
	return false
}

func (rd returnData) IncidentsData(c Client) (res IncidentResult) {
	err := json.Unmarshal(rd, &res)
	if err != nil {
		log.Printf("Could not unmarshall Incident response to struct - %+v\n", err)
		return
	}
	res.Count = len(res.Incidents)
	for index, incident := range res.Incidents {
		res.Incidents[index].ProductName = incident.ImpactedProduct(c)
	}
	for index, incident := range res.Incidents {
		res.Incidents[index].LSMAssigned = incident.AssignedUser(c)
	}
	return
}

func (ir IncidentResult) Aggregate() map[string]int {
	currentTime := time.Now()
	var m = make(map[string]int)
	var d []time.Duration
	var ad []time.Duration
	for _, incident := range ir.Incidents {
		if string(incident.AssignedToRaw) == `""` {
			//+1 unassigned
			m["count.unassigned"] += 1
			//+1 incident.SysCreatedOn added to time slice
			createdTime, err := time.Parse(IncidentTimeFormat, incident.SysCreatedOn)
			if err != nil {
				log.Println("Could not parse created date for incident", incident.Number)
				createdTime = currentTime
			}
			d = append(d, currentTime.Sub(createdTime))
			//+1 category - increment counter map[string]int?
			m[strings.Join([]string{"product.unassigned", incident.ProductName}, ".")] += 1
			//+1 priority - increment counter map[int]int
			m[strings.Join([]string{"priority.unassigned", incident.Priority}, ".")] += 1
			//+1 urgency - increment counter map[int]int
			m[strings.Join([]string{"urgency.unassigned", incident.Urgency}, ".")] += 1
		} else {
			//+1 unassigned
			m["count.assigned"] += 1
			//+1 incident.SysCreatedOn added to time slice
			createdTime, err := time.Parse(IncidentTimeFormat, incident.SysCreatedOn)
			if err != nil {
				log.Println("Could not parse created date for incident", incident.Number)
				createdTime = currentTime
			}
			ad = append(ad, currentTime.Sub(createdTime))
			//+1 category - increment counter map[string]int?
			m[strings.Join([]string{"product.assigned", incident.ProductName}, ".")] += 1
			//+1 priority - increment counter map[int]int
			m[strings.Join([]string{"priority.assigned", incident.Priority}, ".")] += 1
			//+1 urgency - increment counter map[int]int
			m[strings.Join([]string{"urgency.assigned", incident.Urgency}, ".")] += 1
		}
	}
	m["average_age.unassigned"] = averageDurationInDays(d)
	m["average_age.assigned"] = averageDurationInDays(ad)
	return m
}

func averageDurationInDays(durations []time.Duration) int {
	var total int64
	for _, duration := range durations {
		total += int64(duration)
	}
	durationHours := time.Duration(total).Hours()
	avg := float64(durationHours) / float64(len(durations))
	return int(math.RoundToEven(avg))
}
