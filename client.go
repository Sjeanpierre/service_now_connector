package main

import (
	"net/http"
	"strings"
	"log"
	"io/ioutil"
	"crypto/tls"
)

var (
	INCIDENTLISTPATH = "api/now/v2/table/incident"
	USERPATH = "api/now/table/sys_user"
	USERGROUPPATH = "api/now/table/sys_user_grmember"
)

type returnData []byte

type client struct {
	creds oauthPayload
}

type getParams struct {
	params map[string]string
	path   string
	Client client
}

func NewClientwCreds(c credentials) client {
	//todo, cache client
	oauthCreds := c.oauthRequest("password")
	return client{creds:oauthCreds}
}

func NewClient() client {
	//todo, cache client
	var c = credentials{}
	creds := credentials{snClientID, snClientSecret, snUsername, snPassword}
	if creds == c  {
		log.Fatalf("Error: Environment variables for credentials are not set\n Exiting...")
	}
	return NewClientwCreds(creds)
}

func (c getParams) buildURL(path string) string {
	return strings.Join([]string{host, path}, "/")
}

func (gp getParams) Get() returnData {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool}}}
	uri := gp.buildURL(gp.path)
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		log.Fatalln("An error was encountered while building get request", err)
	}
	req.Header.Add("Authorization", gp.Client.creds.bearerToken())
	req.Header.Add("Accept", "application/json")
	params := req.URL.Query()
	for name, value := range gp.params {
		params.Add(name, value)
	}
	req.URL.RawQuery = params.Encode()
	response, err := client.Do(req)

	if err != nil {
		log.Fatalln("An error was encountered while performing get request", err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatalf("A non-200 status code was returned for oauth call\n %+v", response)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln("an error was encountered reading response data from request", err)
	}
	a := returnData{}
	a = responseBody
	return a
}


