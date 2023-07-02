package auth

import (
	"fmt"
	"strings"

	"github.com/vtno/zypher/internal/store"
)

const defaultAuthBucket = "auth"

type AuthGuard struct {
	store store.Store
}

func splitAuthHeader(authHeader string) (string, string) {
	s := strings.Split(authHeader, " ")
	if len(s) > 0 {
		return s[0], s[1]
	}
	return "", ""
}

func NewAuthGuard(store store.Store) *AuthGuard {
	return &AuthGuard{store: store}
}

func (a *AuthGuard) Guard(authHeader string) bool {
	if authHeader == "" {
		return false
	}
	_, token := splitAuthHeader(authHeader)
	result, err := a.store.GetByBucket(defaultAuthBucket, token)
	if err != nil {
		fmt.Printf("error getting token from store: %v\n", err)
		return false
	}

	return result != ""
}
