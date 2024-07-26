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
	IdToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjBjYjQyNzQyYWU1OGY0ZGE0NjdiY2RhZWE0Yjk1YTI5ZmJhMGM1ZjkiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiTmd1eeG7hW4gVGhhbmggQsOsbmggKDIxMTI3MjMyKSIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQ2c4b2NKaGZlMDBhZHJJTWM4Z1lfN3J1UG96SlNqUFpxaXd1QjJFZmF4X2pfZUZudlg3VEJFPXM5Ni1jIiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL3RlbmstaG91cnMtc2xlZXAiLCJhdWQiOiJ0ZW5rLWhvdXJzLXNsZWVwIiwiYXV0aF90aW1lIjoxNzIxOTg5MDQ5LCJ1c2VyX2lkIjoicDRvcXRmdXZ1VVVtVFg1eUhxY2NJRnFwODlBMiIsInN1YiI6InA0b3F0ZnV2dVVVbVRYNXlIcWNjSUZxcDg5QTIiLCJpYXQiOjE3MjE5ODkwNDksImV4cCI6MTcyMTk5MjY0OSwiZW1haWwiOiJudGJpbmgyMUBjbGMuZml0dXMuZWR1LnZuIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsiZ29vZ2xlLmNvbSI6WyIxMDc3MDI0NzEwMDQyNDQxNDAzMzAiXSwiZW1haWwiOlsibnRiaW5oMjFAY2xjLmZpdHVzLmVkdS52biJdfSwic2lnbl9pbl9wcm92aWRlciI6Imdvb2dsZS5jb20ifX0.xNRPaRoAGNOgUd5_9hYEjLD9TnD8SgrIK2jJmIsraw7s0AB6lztxYTKqdGaGLp9HVnRhJdkPZV4LVpMqda02GIY2-2j1Ygd0PDfIxGUrMW89QiKJo3wHoOmtyOzn2fwKMwVHFdVQc41WH632k1Hmc1FHyB_blaB3uwHqDujFhXJ8ofjJaV2VZ35h7sE6MSm-pKaLAsoB81pkOn83qz99f7wvFOYlZfY_jgYtCS4Oy1WCb11TZLNNt0RHs6uzSQ1UL4MpkLY5TNZ17pFZo4bvHyMkcIBJGHi5mDo8ag0BDabP_NqXue8PEBG0nhKnuNXBFvqjElsalyHXuPjK3u8MQg"
)

type Map map[string]interface{}

func (r *Map) getData() Map {
	return (*r)["data"].(Map)
}

func (r *Map) getError() Map {
	return (*r)["errors"].(Map)
}

func (r *Map) log() error {
	jsonData, err := json.MarshalIndent(*r, "", "  ")
	if err != nil {
		return err
	}

	log.Printf("--> Response: %v\n", string(jsonData))
	return nil
}
