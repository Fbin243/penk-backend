package test

import (
	"encoding/json"
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
	IdToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6Ijc5M2Y3N2Q0N2ViOTBiZjRiYTA5YjBiNWFkYzk2ODRlZTg1NzJlZTYiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiVGhhbmggQsOsbmggTmd1eeG7hW4iLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUNnOG9jSnBKVnZ5MWRKd2JkN05WbG9OSVpSZlBRNThiQkhoc3NocGFFaldTRlZlSWJCZ3FvTT1zOTYtYyIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS90ZW5rLWhvdXJzLXNsZWVwIiwiYXVkIjoidGVuay1ob3Vycy1zbGVlcCIsImF1dGhfdGltZSI6MTcxOTkzNzg1OSwidXNlcl9pZCI6Im1UNWJwWFk0M2JhR0RVbzd6bmsxVE9JWkl1QzMiLCJzdWIiOiJtVDVicFhZNDNiYUdEVW83em5rMVRPSVpJdUMzIiwiaWF0IjoxNzE5OTM3ODU5LCJleHAiOjE3MTk5NDE0NTksImVtYWlsIjoibnRiaW5oMjQzQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7Imdvb2dsZS5jb20iOlsiMTA5NTEwNjI4MDI0ODQ3OTI3Mzg4Il0sImVtYWlsIjpbIm50YmluaDI0M0BnbWFpbC5jb20iXX0sInNpZ25faW5fcHJvdmlkZXIiOiJnb29nbGUuY29tIn19.IivQk-YKjnvNMNWXGQ-F4KsxBynqGVSRnL62K2kJLIyXqUUvM_24fRJQeBshuiXAMpO_4KFlJESUunbSTLvjrYrCBQWoL9RqPdF84yvD7DldoJmP6BJMRJPr_n8PZXw36qPKeWLnF0F56mCSt7ksr89QpSkXFyl_bqx01cdBWMjxAXTdLmcKgf5LOe_H1YgyTta1XKwhG0nVnJJ02wEc51BBxJ_SxK8DnfNet3n5Nc9q5a9O4BHUT0aKq2cr-q4kvew2C6VstvC-MsU0iGMrGRwJN-F5QWM8yb702MjoNAGGNGwvKu8Dql-NyO8GRxDuz2ynQWcCX0OrHqZntoSbQQ"
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
