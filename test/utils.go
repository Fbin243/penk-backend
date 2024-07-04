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
	IdToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjU2OTFhMTk1YjI0MjVlMmFlZDYwNjMzZDdjYjE5MDU0MTU2Yjk3N2QiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiVGhhbmggQsOsbmggTmd1eeG7hW4iLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUNnOG9jSnBKVnZ5MWRKd2JkN05WbG9OSVpSZlBRNThiQkhoc3NocGFFaldTRlZlSWJCZ3FvTT1zOTYtYyIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS90ZW5rLWhvdXJzLXNsZWVwIiwiYXVkIjoidGVuay1ob3Vycy1zbGVlcCIsImF1dGhfdGltZSI6MTcyMDExNjg5OCwidXNlcl9pZCI6Im1UNWJwWFk0M2JhR0RVbzd6bmsxVE9JWkl1QzMiLCJzdWIiOiJtVDVicFhZNDNiYUdEVW83em5rMVRPSVpJdUMzIiwiaWF0IjoxNzIwMTE2ODk4LCJleHAiOjE3MjAxMjA0OTgsImVtYWlsIjoibnRiaW5oMjQzQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7Imdvb2dsZS5jb20iOlsiMTA5NTEwNjI4MDI0ODQ3OTI3Mzg4Il0sImVtYWlsIjpbIm50YmluaDI0M0BnbWFpbC5jb20iXX0sInNpZ25faW5fcHJvdmlkZXIiOiJnb29nbGUuY29tIn19.U31g-Isimwx09f9gZYiH_5sjNqhjBYeyqplottPlkOfQfmVuebMMKJn8wYyP9HhyjutjqgOq6BfY1MRQzrz89bETRzroCIQTdIYKki4wQnmNKsJGxFJtM2eIf2QkWCIXdb7QzVN10-unUGcRzF1QWfLViAAUVXz2sGwCo1rp5pMS09z1EFBhN0Uf33dJqpbB90JSQ7y8NXlpnW06Dxno2T1MUbg6UYqMdP3_aOxRSfFduf96n7GOZanubxZrSUW6aTbE364y1CAPV4xiAKehOsBKAiK4S8bXymJNZkWH753AqPJ0nyx8TcpyEwgwaGV1pMC1wgPk2iF7-CMUMMzRlw"
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
