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
	IdToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6Ijc5M2Y3N2Q0N2ViOTBiZjRiYTA5YjBiNWFkYzk2ODRlZTg1NzJlZTYiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiVGhhbmggQsOsbmggTmd1eeG7hW4iLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUNnOG9jSnBKVnZ5MWRKd2JkN05WbG9OSVpSZlBRNThiQkhoc3NocGFFaldTRlZlSWJCZ3FvTT1zOTYtYyIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS90ZW5rLWhvdXJzLXNsZWVwIiwiYXVkIjoidGVuay1ob3Vycy1zbGVlcCIsImF1dGhfdGltZSI6MTcyMDAxODg5NSwidXNlcl9pZCI6Im1UNWJwWFk0M2JhR0RVbzd6bmsxVE9JWkl1QzMiLCJzdWIiOiJtVDVicFhZNDNiYUdEVW83em5rMVRPSVpJdUMzIiwiaWF0IjoxNzIwMDE4ODk1LCJleHAiOjE3MjAwMjI0OTUsImVtYWlsIjoibnRiaW5oMjQzQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7Imdvb2dsZS5jb20iOlsiMTA5NTEwNjI4MDI0ODQ3OTI3Mzg4Il0sImVtYWlsIjpbIm50YmluaDI0M0BnbWFpbC5jb20iXX0sInNpZ25faW5fcHJvdmlkZXIiOiJnb29nbGUuY29tIn19.XlJDSR0vkBqvwoMbKuR_4ywo-1cJnhhfHl-hFqfLaonc3bNzVdTc2qeuClg73nSMzMstJ84B2PrV5id8oIpjY2gp9ABgvOrfoSPEdzILJ92HMov2Ae28wFmSH5jSuHsmgPbnT1YwTuCvoYoGURuj20L_lYlf3eGXiodRkWInvu-lJPSWr8sQDXJg4lhMkViea36b-VbYbRZYUT2stJcrPHZ_Ll1E_lm9hBjJzQSuvBavWB-Ymiq1qPV8UX7u1GOUH31lAfQmgQRpa052hNIWE94i19FIIPHlrgqbJ4jinSVgDVVic4VhUoC3nEI-pLJk6wHxTthjJt4-cuV9Ewu7Pw"
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
