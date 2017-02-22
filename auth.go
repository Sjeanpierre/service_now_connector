package main

import (
	//"encoding/json"
	"io/ioutil"
	"bytes"
	"net/http"
	"strings"
	"log"
	"net/url"
	"encoding/json"
	"crypto/tls"
)



type credentials struct {
	clientID     string
	clientSecret string
	username     string
	password     string
}

type oauthPayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	Expires      int `json:"expires_in"`
}

func (pl oauthPayload) bearerToken() string{
	return strings.Join([]string{"Bearer",pl.AccessToken}," ")
}

func (pl oauthPayload) valid() bool {
	if (oauthPayload{}) == pl {
		return false
	}
	return true
}

func (c credentials) oauthRequest(grantType string) oauthPayload {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool}}}
	uri := strings.Join([]string{host, "oauth_token.do"}, "/")
	v := url.Values{"grant_type": {grantType},
		"client_id": {c.clientID},
		"client_secret":{c.clientSecret},
		"username":{c.username},
		"password":{c.password}}
	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(v.Encode()))
	if err != nil {
		log.Fatalln("An error was encountered while building oauth token request to Service Now\n", err)
	}
	req.Header.Add("accept", "json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	if err != nil {
		log.Fatalln("An error was encountered retrieving bearer token from Service Now\n", err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatalf("A non-200 status code was returned for oauth call\n %+v", response)
	}
	defer response.Body.Close()
	ResponseText, error := ioutil.ReadAll(response.Body)
	if error != nil {
		log.Fatalln("An error was encountered reading response data from bearer token request")
	}
	result := oauthPayload{}
	err = json.Unmarshal([]byte(ResponseText), &result)
	if err != nil || !result.valid() {
		log.Fatalf("Could not unmarshall response body to valid" +
			" object for oauth request. error: %v, result: %v", err, result)
	}
	return result
}

func (c credentials) oauthToken() string {
	//Check for existing valid token based on TTL, return
	//Get new token from oauth function, return
	payload := c.oauthRequest("password")
	log.Printf("%+v", payload)
	token := strings.Join([]string{"Bearer", payload.AccessToken}, " ")
	return token
}