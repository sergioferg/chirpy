package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestJWT(t *testing.T) {
	user1ID := uuid.New()
	secretToken1 := "ReallySecretToken"
	secretToken2 := "EvenMoreSecretToken"

	oneSecDur, _ := time.ParseDuration("1s")
	quickDuration, _ := time.ParseDuration("1ms")

	token1, _ := MakeJWT(user1ID, secretToken1, oneSecDur)

	tests := []struct {
		name        string
		token       string
		secretToken string
		duration    time.Duration
		wantErr     bool
	}{
		{
			name:        "Valid token",
			token:       token1,
			secretToken: secretToken1,
			duration:    quickDuration,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			token:       "someothertoken",
			secretToken: secretToken1,
			duration:    quickDuration,
			wantErr:     true,
		},
		{
			name:        "Incorrect secret token",
			token:       token1,
			secretToken: secretToken2,
			duration:    quickDuration,
			wantErr:     true,
		},
		{
			name:        "Timeout",
			token:       token1,
			secretToken: secretToken1,
			duration:    oneSecDur,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(tt.duration)
			_, err := ValidateJWT(tt.token, tt.secretToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetToken(t *testing.T) {
	goodReq, err := http.NewRequest("GET", "/test-path", nil)
	if err != nil {
		t.Fatal(err)
	}
	badReq, err := http.NewRequest("GET", "/test-path", nil)
	if err != nil {
		t.Fatal(err)
	}
	noHeaderReq, err := http.NewRequest("GET", "/test-path", nil)
	if err != nil {
		t.Fatal(err)
	}

	goodReq.Header.Set("Authorization", "Bearer jwtToken")
	badReq.Header.Set("Authorization", "Bearer")

	tests := []struct {
		name    string
		request http.Request
		wantErr bool
	}{
		{
			name:    "Valid request header",
			request: *goodReq,
			wantErr: false,
		},
		{
			name:    "Invalid request header",
			request: *badReq,
			wantErr: true,
		},
		{
			name:    "No header in request",
			request: *noHeaderReq,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBearerToken(tt.request.Header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
