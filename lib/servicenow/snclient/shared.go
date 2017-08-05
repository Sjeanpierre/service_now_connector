package snclient

import (
	"os"
	"crypto/x509"
)

var (
	INCIDENTLISTPATH = "api/now/v2/table/incident"
	USERPATH = "api/now/table/sys_user"
	USERGROUPPATH = "api/now/table/sys_user_grmember"
	host = os.Getenv("SERVICE_NOW_HOSTNAME")
	snClientID = os.Getenv("SERVICE_NOW_CLIENT_ID")
	snClientSecret = os.Getenv("SERVICE_NOW_CLIENT_SECRET")
	snUsername = os.Getenv("SERVICE_NOW_USERNAME")
	snPassword = os.Getenv("SERVICE_NOW_PASSWORD")
	pool = &x509.CertPool{}
	noUser = User{"N/A", "N/A", "N/A", "N/A", "N/A", "N/A", "N/A", "N/A", "N/A", "N/A"}
)

func init() {
	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)
	//client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool}}}
}
