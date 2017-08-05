package snclient

import (
	"net/http"
	"strings"
	"log"
	"io/ioutil"
	"crypto/tls"
)



type returnData []byte

type Client struct {
	creds oauthPayload
}

type getParams struct {
	params map[string]string
	path   string
	Client Client
}

func NewClientwCreds(c credentials) Client {
	//todo, cache client
	oauthCreds := c.oauthRequest("password")
	return Client{creds:oauthCreds}
}

func NewClient() Client {
	//todo, cache client
	var c = credentials{}
	creds := credentials{snClientID, snClientSecret, snUsername, snPassword}
	if creds == c  {
		log.Println("Error: Environment variables for credentials are not set\n Exiting...")
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
		log.Println("An error was encountered while building get request", err)
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
		log.Println("An error was encountered while performing get request", err)
		return returnData{}
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("A non-200 status code was returned for oauth call\n %+v", response)
		return returnData{}
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("an error was encountered reading response data from request", err)
		return returnData{}
	}
	return responseBody
}


