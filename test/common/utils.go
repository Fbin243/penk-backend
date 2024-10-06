package common

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

const GatewayUrl = "http://localhost:8070/graphql"

var (
	cookieJar, _ = cookiejar.New(nil)
	cli          = &http.Client{
		Timeout: time.Second * 20,
		Jar:     cookieJar,
	}
	IdToken string   
)
