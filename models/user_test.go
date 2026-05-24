package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUserJSONDoesNotExposePassword(t *testing.T) {
	user := User{
		ID:       1,
		Username: "alice",
		Password: "$2a$10$secret-hash",
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal user: %v", err)
	}

	body := string(data)
	if strings.Contains(body, "password") || strings.Contains(body, "secret-hash") {
		t.Fatalf("user JSON exposed password data: %s", body)
	}
	if !strings.Contains(body, "alice") {
		t.Fatalf("user JSON did not include username: %s", body)
	}
}
