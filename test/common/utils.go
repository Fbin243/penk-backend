package common

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

const (
	GatewayUrl       = "http://localhost:8070/graphql"
	CoreUrl          = "http://localhost:8080/graphql"
	TimeTrackingsUrl = "http://localhost:8082/graphql"
	AnalyticsUrl     = "http://localhost:8083/graphql"
)

var (
	cookieJar, _ = cookiejar.New(nil)
	cli          = &http.Client{
		Timeout: time.Second * 20,
		Jar:     cookieJar,
	}
	IdToken  string
	DeviceId string
)
