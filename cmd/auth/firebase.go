package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/coredb"

	"github.com/tidwall/gjson"
)

const baseUri = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=%s"

func GetCustomTokenByEmail(email string) (string, error) {
	profileRepo := coredb.NewProfilesRepo(db.GetDBManager().DB)
	profile, err := profileRepo.GetProfileByEmail(email)
	if err != nil {
		return "", err
	}

	authClient := auth.GetFirebaseManager().Client

	// Create a custom token from the Firebase UID
	customToken, err := authClient.CustomToken(context.Background(), profile.FirebaseUID)
	if err != nil {
		return "", err
	}

	apiKey := os.Getenv("WEB_API_KEY")
	url := fmt.Sprintf(baseUri, apiKey)
	requestBody := map[string]any{
		"token":             customToken,
		"returnSecureToken": true,
	}

	json, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(json))
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return gjson.Get(string(body), "idToken").String(), nil
}
