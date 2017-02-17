package main

import (
	"encoding/json"
	"log"
	"fmt"
)

var userCacheStore = make(map[string]User)

type UserResult struct {
	Users []User `json:"result"`
}

type UserGroupResult struct {
	UserGroups []UserGroup `json:"result"`
}

type UserGroup struct {
	User struct {
		     ID string `json:"value"`
	     } `json:"user"`
}

type User struct {
	Active         string `json:"active"`
	Email          string `json:"email"`
	EmployeeNumber string `json:"employee_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Phone          string `json:"phone"`
	TimeZone       string `json:"time_zone"`
	Title          string `json:"title"`
	Zip            string `json:"zip"`
	SystemID       string `json:"sys_id"`
}

func (user User) CacheLookup(userID string) ([]User, bool) {
	u, ok := userCacheStore[userID]
	if ok {
		return []User{u}, true
	}
	return []User{}, false
}

func (user User) CacheAdd() {
	if user.SystemID != "" {
		userCacheStore[user.SystemID] = user
	}
}

type userParams struct {
	userID  string
	groupID string
}

func (c client) User(id string) ([]User) {
	u, ok := User{}.CacheLookup(id)
	if ok {
		return u
	}
	gp := make(map[string]string)
	gp["sys_id"] = id
	gp["sysparm_limit"] = "100"
	UserRequest := getParams{}
	UserRequest.path = USERPATH
	UserRequest.params = gp
	UserRequest.Client = c
	return UserRequest.Get().UsersData()
}

func (d returnData) UsersData() (res []User) {
	var r = UserResult{}
	err := json.Unmarshal(d, &r)
	if err != nil {
		log.Fatal("Could not unmarshall User data response to struct", err)
	}
	res = r.Users
	for _,user := range res {
		user.CacheAdd()
		log.Printf("Added user %s to cache", user.Email)
	}
	return
}

func (c client) UserGroup(id string) ([]User) {
	gp := make(map[string]string)
	gp["sysparm_query"] = fmt.Sprintf("%s=%s", "group", id)
	gp["sysparm_limit"] = "100"
	UserGroupRequest := getParams{}
	UserGroupRequest.path = USERGROUPPATH
	UserGroupRequest.params = gp
	UserGroupRequest.Client = c
	groups := UserGroupRequest.Get().UserGroupData()
	var userList []User
	for _, group := range groups.UserGroups {
		u := c.User(group.User.ID)
		userList = append(userList, u[0])
	}
	return userList
}

func (d returnData) UserGroupData() (res UserGroupResult) {
	err := json.Unmarshal(d, &res)
	if err != nil {
		log.Fatal("Could not unmarshall User Group response to struct", err)
	}
	return
}