package session

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type JWTClaims struct {
	Sub    string                 `json:"sub"`
	Exp    int64                  `json:"exp"`
	Hasura map[string]interface{} `json:"https://hasura.io/jwt/claims"`
}

func DecodeClaims(token string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid token format")
	}
	payload := parts[1]
	data, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}
	var claims JWTClaims
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

func ExtractUserID(claims *JWTClaims) (int, error) {
	if claims == nil {
		return 0, fmt.Errorf("missing claims")
	}

	if claims.Hasura != nil {
		if value, ok := claims.Hasura["x-hasura-user-id"]; ok {
			switch v := value.(type) {
			case string:
				if v != "" {
					return strconv.Atoi(v)
				}
			case float64:
				return int(v), nil
			}
		}
	}

	if claims.Sub != "" {
		return strconv.Atoi(claims.Sub)
	}

	return 0, fmt.Errorf("user id not present in claims")
}
