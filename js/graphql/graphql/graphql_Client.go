package graphql

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"zone01-dashboard/session"
)

const (
	defaultBaseAPI = "https://platform.zone01.gr/api"
	signinPath     = "/auth/signin"
	graphqlPath    = "/graphql-engine/v1/graphql"
)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

func SignIn(identifier, password string) (string, *session.JWTClaims, int) {
	payload := identifier + ":" + password
	encoded := base64.StdEncoding.EncodeToString([]byte(payload))

	url := defaultBaseAPI + signinPath
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", nil, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", "Basic "+encoded)

	res, err := httpClient.Do(req)
	if err != nil {
		return "", nil, http.StatusInternalServerError
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return "", nil, res.StatusCode
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", nil, res.StatusCode
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil, http.StatusInternalServerError
	}
	body := string(bodyBytes)

	token := strings.TrimSpace(body)
	token = strings.Trim(token, `"`)
	if token == "" {
		return "", nil, http.StatusInternalServerError
	}

	claims, err := session.DecodeClaims(token)
	if err != nil {
		return "", nil, http.StatusInternalServerError
	}

	return token, claims, http.StatusOK
}

func Query(query string, variables map[string]interface{}, token string) (map[string]interface{}, int, error) {
	url := defaultBaseAPI + graphqlPath

	respData := make(map[string]interface{})

	if strings.TrimSpace(token) == "" {
		return respData, http.StatusUnauthorized, fmt.Errorf("missing authentication token")
	}

	args := make(map[string]interface{})
	for k, v := range variables {
		args[k] = v
	}
	// Request body
	body := map[string]interface{}{
		"query":     query,
		"variables": args,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return respData, http.StatusInternalServerError, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return respData, http.StatusInternalServerError, fmt.Errorf("failed to create graphql request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return respData, http.StatusInternalServerError, fmt.Errorf("failed to execute graphql request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return respData, resp.StatusCode, fmt.Errorf("graphql request returned status %d", resp.StatusCode)
	}

	// Decode JSON response
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return respData, http.StatusInternalServerError, err
	}
	for k := range respData {
		if k == "errors" {
			return respData, http.StatusBadRequest, fmt.Errorf("error response from platform. validation failed due to bad request")
		}
	}

	return respData, http.StatusOK, nil
}

func GetCollabs(userId int, token string) (map[string]interface{}, int, error) {
	url := defaultBaseAPI + graphqlPath

	respData := make(map[string]interface{})

	if strings.TrimSpace(token) == "" {
		return respData, http.StatusUnauthorized, fmt.Errorf("missing authentication token")
	}

	// get ids
	ids, code, err := GetIds(userId, token)
	if err != nil {
		return respData, code, err
	}
	// Request body
	body := map[string]interface{}{
		"query":     USER_COLLABORATORS,
		"variables": map[string]interface{}{"userId": userId, "groupIds": ids},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return respData, http.StatusInternalServerError, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return respData, http.StatusInternalServerError, fmt.Errorf("failed to create graphql request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return respData, http.StatusInternalServerError, fmt.Errorf("failed to execute graphql request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return respData, resp.StatusCode, fmt.Errorf("graphql request returned status %d", resp.StatusCode)
	}

	// Decode JSON response
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return respData, http.StatusInternalServerError, err
	}
	for k := range respData {
		if k == "errors" {
			return respData, http.StatusBadRequest, fmt.Errorf("error response from platform. validation failed due to bad request")
		}
	}

	return respData, http.StatusOK, nil
}

func GetIds(userId int, token string) ([]int, int, error) {
	var ids []int

	idsMap, code, err := Query(USER_GROUP_IDS, map[string]interface{}{"userId": userId}, token)
	if err != nil {
		return ids, code, err
	}

	ids = ExtractGroupIDs(idsMap)
	return ids, code, err
}

func ExtractGroupIDs(data map[string]interface{}) []int {
	groupIDs := []int{}

	if d, ok := data["data"].(map[string]interface{}); ok {
		if gu, ok := d["group_user"].([]interface{}); ok {
			for _, g := range gu {
				if gm, ok := g.(map[string]interface{}); ok {
					if id, ok := gm["groupId"].(float64); ok {
						groupIDs = append(groupIDs, int(id))
					}
				}
			}
		}
	}
	return groupIDs
}
