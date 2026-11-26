package basicauth

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"strconv"
)

type account struct {
	username, password [sha256.Size]byte
}

// NewVerifier snapshots the allowed accounts. An empty account set rejects every request.
func NewVerifier(credentials map[string]string) func(*http.Request) (string, bool) {
	accounts := make([]account, 0, len(credentials))
	for username, password := range credentials {
		accounts = append(accounts, account{
			username: sha256.Sum256([]byte(username)),
			password: sha256.Sum256([]byte(password)),
		})
	}

	return func(r *http.Request) (string, bool) {
		username, password, ok := r.BasicAuth()
		if !ok {
			return "", false
		}

		userHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		matched := 0
		for _, account := range accounts {
			matched |= subtle.ConstantTimeCompare(userHash[:], account.username[:]) &
				subtle.ConstantTimeCompare(passwordHash[:], account.password[:])
		}
		return username, matched == 1
	}
}

// Challenge returns the WWW-Authenticate value for a Basic Auth realm.
func Challenge(realm string) string {
	return "Basic realm=" + strconv.Quote(realm)
}
