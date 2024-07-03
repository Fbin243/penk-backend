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
	IdToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6Ijc5M2Y3N2Q0N2ViOTBiZjRiYTA5YjBiNWFkYzk2ODRlZTg1NzJlZTYiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiQsOsbmggTmd1eeG7hW4iLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUNnOG9jS3dXQnNrdTFRU18weXVPdTJOS3d4M2Z6LW8xRzlYMmZqcUlpcmdkZi1FVngxWFk5az1zOTYtYyIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS90ZW5rLWhvdXJzLXNsZWVwIiwiYXVkIjoidGVuay1ob3Vycy1zbGVlcCIsImF1dGhfdGltZSI6MTcxOTk4NjY0OSwidXNlcl9pZCI6IlZxSVRTTXJicVFWM0QzRDlzV2ZCMTE3V0xobzIiLCJzdWIiOiJWcUlUU01yYnFRVjNEM0Q5c1dmQjExN1dMaG8yIiwiaWF0IjoxNzE5OTg2NjQ5LCJleHAiOjE3MTk5OTAyNDksImVtYWlsIjoibnRiaW5oMjQzLmRldkBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJnb29nbGUuY29tIjpbIjExMDQ0Mjg1MDE1Njk5ODAzMTE1MyJdLCJlbWFpbCI6WyJudGJpbmgyNDMuZGV2QGdtYWlsLmNvbSJdfSwic2lnbl9pbl9wcm92aWRlciI6Imdvb2dsZS5jb20ifX0.UPs5IDu5cJKJPOIwlY5sIpZcuYQT7rMymbe5EkizsGHdih9-jYWtcj6PVQGVncWBZEQ7LstBTCmdhrw2csa3kDCia0u0JbeZDBH4LbdGbiDSds6JK_cAQTn0XEcX0E5AwLHl-5vX-jKN48WV7OibVrXydkfpPmTVGOdY3zvS45a3sel_JKnymwjjr4RmW9U5UNoeQ3xaxtbBxU8zPJSqqw9eyH0O6jXqmXJF5hoGzTZyrGWoy4wQIobuq4TRlIx6Vrf2KYERhmHu5GvNjd4Q8ZecYTFsSLI4seCuaHdDzGO5Hfmt8HEQqSRe87grrCkU4u5J3vDi7WyL2JHYyi20Rg"
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
