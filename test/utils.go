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
	IdToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6Ijc5M2Y3N2Q0N2ViOTBiZjRiYTA5YjBiNWFkYzk2ODRlZTg1NzJlZTYiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiVGhhbmggQsOsbmggTmd1eeG7hW4iLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUNnOG9jSnBKVnZ5MWRKd2JkN05WbG9OSVpSZlBRNThiQkhoc3NocGFFaldTRlZlSWJCZ3FvTT1zOTYtYyIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS90ZW5rLWhvdXJzLXNsZWVwIiwiYXVkIjoidGVuay1ob3Vycy1zbGVlcCIsImF1dGhfdGltZSI6MTcxOTk0NjAxNSwidXNlcl9pZCI6Im1UNWJwWFk0M2JhR0RVbzd6bmsxVE9JWkl1QzMiLCJzdWIiOiJtVDVicFhZNDNiYUdEVW83em5rMVRPSVpJdUMzIiwiaWF0IjoxNzE5OTQ2MDE1LCJleHAiOjE3MTk5NDk2MTUsImVtYWlsIjoibnRiaW5oMjQzQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7Imdvb2dsZS5jb20iOlsiMTA5NTEwNjI4MDI0ODQ3OTI3Mzg4Il0sImVtYWlsIjpbIm50YmluaDI0M0BnbWFpbC5jb20iXX0sInNpZ25faW5fcHJvdmlkZXIiOiJnb29nbGUuY29tIn19.PBkHFd9rcAnvFtNcyZ3lBU6PW3l7BD7LCO2LaktWTYU7-21ToqyrGPJovSUblt0qmjkNbmOkBuAWxJxW1AEYwXhNfe2VNz9AJ1UCfoGFlKnpNnFJPaFUnAiYg6_wzYJ87rjkHEO1ubRI41SA-F8vFG_jzH7JCQ_9XDta4wEtNl3fRlWip5FXtg2l4qcjQg0PJPnsXDiJhyhw72xnxzJryxydZ847dzK5WIG_71JMJzD4VrM10j1aEartgXAK5E00FIgPmWzLDHlSWe1ji6siyeWk46TU9skeLNts1I5AgUUneZkkwuDmM8x31JbrZsnS9EqhIiUUQDcJ8z_ruSBIFg"
)

type TestContext struct {
	IdUser         string `json:"idUser"`
	IdCharacter    string `json:"idCharacter"`
	IdCustomMetric string `json:"idCustomMetric"`
	IdTimeTracking string `json:"idTimeTracking"`
	IdProperty     string `json:"idProperty"`
}

func decodeResponseData(r1 *http.Response) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := json.NewDecoder(r1.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response["data"].(map[string]interface{}), nil
}

func logResponseData(r1 *http.Response, r2 *http.Request) error {
	data, err := decodeResponseData(r1)
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	log.Printf("--> Response data: %v\n", string(jsonData))
	return nil
}
