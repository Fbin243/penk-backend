package test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var (
	cookieJar, _ = cookiejar.New(nil)
	cli          = &http.Client{
		Timeout: time.Second * 1,
		Jar:     cookieJar,
	}
	url     = "http://localhost:8080/graphql"
	IdToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6Ijc5M2Y3N2Q0N2ViOTBiZjRiYTA5YjBiNWFkYzk2ODRlZTg1NzJlZTYiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiVGhhbmggQsOsbmggTmd1eeG7hW4iLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUNnOG9jSnBKVnZ5MWRKd2JkN05WbG9OSVpSZlBRNThiQkhoc3NocGFFaldTRlZlSWJCZ3FvTT1zOTYtYyIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS90ZW5rLWhvdXJzLXNsZWVwIiwiYXVkIjoidGVuay1ob3Vycy1zbGVlcCIsImF1dGhfdGltZSI6MTcxOTk5MDg2MywidXNlcl9pZCI6Im1UNWJwWFk0M2JhR0RVbzd6bmsxVE9JWkl1QzMiLCJzdWIiOiJtVDVicFhZNDNiYUdEVW83em5rMVRPSVpJdUMzIiwiaWF0IjoxNzE5OTkwODYzLCJleHAiOjE3MTk5OTQ0NjMsImVtYWlsIjoibnRiaW5oMjQzQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7Imdvb2dsZS5jb20iOlsiMTA5NTEwNjI4MDI0ODQ3OTI3Mzg4Il0sImVtYWlsIjpbIm50YmluaDI0M0BnbWFpbC5jb20iXX0sInNpZ25faW5fcHJvdmlkZXIiOiJnb29nbGUuY29tIn19.KDxEErAa6ScpRDJ3HKqVVTL_8MMUNXeeVySLbLlWmlQUuOAy-4aJZNGWZwpoIlQK8ws6IogPxVYZnb6lV4aUU6gVe7pAB3sazm3KOUWiz0MxoCpjBUObzQC7qJmClHVmVPOl_Ml-OQ7q9gjMFQDDXNhrJSfB2q-JaKjanIYbtSKACjB2fWpnnwRgF-pQS_-ROUrbHZeeMfboeoVF0gcnUZc600mWs2NUk5pq_EgeQRukrf1Iwlv6zSDy_KRJ8Oo39_HvpMb7JCyZgxbH81KfAHPRxXIfGoOW_ci6h7Yn_BVIH6YMTGAkoZH3c8nLgPVcDOTkRN2priO5AE-sqpc23Q"
)

type TestContext struct {
	IdUser         string `json:"idUser"`
	IdCharacter    string `json:"idCharacter"`
	IdCustomMetric string `json:"idCustomMetric"`
	IdTimeTracking string `json:"idTimeTracking"`
	IdProperty     string `json:"idProperty"`
}

var responseBody map[string]interface{}

func logResponse(responseBody map[string]interface{}) error {
	jsonData, err := json.MarshalIndent(responseBody["data"], "", "  ")
	if err != nil {
		return err
	}

	log.Printf("--> Response: %v\n", string(jsonData))
	return nil
}
