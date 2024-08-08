package test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

var (
	cookieJar, _ = cookiejar.New(nil)
	cli          = &http.Client{
		Timeout: time.Second * 20,
		Jar:     cookieJar,
	}
	url     = "http://localhost:8080/graphql"
	IdToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjBjYjQyNzQyYWU1OGY0ZGE0NjdiY2RhZWE0Yjk1YTI5ZmJhMGM1ZjkiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiTmd1eeG7hW4gVGhhbmggQsOsbmggKDIxMTI3MjMyKSIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQ2c4b2NKaGZlMDBhZHJJTWM4Z1lfN3J1UG96SlNqUFpxaXd1QjJFZmF4X2pfZUZudlg3VEJFPXM5Ni1jIiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL3RlbmstaG91cnMtc2xlZXAiLCJhdWQiOiJ0ZW5rLWhvdXJzLXNsZWVwIiwiYXV0aF90aW1lIjoxNzIxOTk3NjYxLCJ1c2VyX2lkIjoicDRvcXRmdXZ1VVVtVFg1eUhxY2NJRnFwODlBMiIsInN1YiI6InA0b3F0ZnV2dVVVbVRYNXlIcWNjSUZxcDg5QTIiLCJpYXQiOjE3MjE5OTc2NjEsImV4cCI6MTcyMjAwMTI2MSwiZW1haWwiOiJudGJpbmgyMUBjbGMuZml0dXMuZWR1LnZuIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsiZ29vZ2xlLmNvbSI6WyIxMDc3MDI0NzEwMDQyNDQxNDAzMzAiXSwiZW1haWwiOlsibnRiaW5oMjFAY2xjLmZpdHVzLmVkdS52biJdfSwic2lnbl9pbl9wcm92aWRlciI6Imdvb2dsZS5jb20ifX0.Riok5TY_Yq2bwoTBGKsM_rt7dZfedsK9iB0z3UnvOwwz-O9UcELO3EaaXh93HiVSXAM1pEHP93lS_STGvwMNmpwf_Hr4ItOzubOdY3gN2c5V3V284b0NZth_WAgYb7RYo5uBNMK0_2K9sSLY55q8925k8YARQjtALsmr0f1EzKm_JoUIigbyle9kpZAlk69VlFXG1Ak9H33FKj-HJWEgqe1Vn5PnAbGLFNLZ-bxjiZTOX1ORiBRTKR4RTdz9o7Pl8L4TdQZ0MyMd0ujohDY4ISyDsjid4292VFpeUhceKmGbkkXdxmx3H87LE6KTNex4xn0HG_EbOe0T3sBKKl_djg"
)

type Map map[string]interface{}

func (r *Map) getFieldValue(fieldString string) interface{} {
	fields := strings.Split(fieldString, ".")
	log.Println(fields)
	n := len(fields)
	rMap := *r
	for i := 0; i < n-1; i++ {
		rMap = rMap[fields[i]].(map[string]interface{})
	}

	return rMap[fields[n-1]]
}

func (r *Map) log() error {
	jsonData, err := json.MarshalIndent(*r, "", "  ")
	if err != nil {
		return err
	}

	log.Printf("--> Response: %v\n", string(jsonData))
	return nil
}
