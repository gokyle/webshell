// Package auth provides PBKDF2 authentication for webapps.
// This package is designed so that users should really only have to supply
// their own means of procuring an AuthProvider in order to make use of the
// provided functionality.
//
// To return a hash, salt from a password, the HashPass function is provided.
// The end user can store these safely in their database.
package auth

import "bitbucket.org/taruti/pbkdf2"

// AuthProviders are functions that takes a user ID and returns a salt and hash.
type AuthProvider func(user interface{}) (salt []byte, hash []byte)

// Authenticators take a user and password and return true if authentication succeeds.
type Authenticator func(user interface{}, password string) bool

var (
	// LookupCredentials should be set to a valid AuthProvider function
	LookupCredentials AuthProvider = DefaultAuthProvider

	// Authenticate should be set to a valid Authenticator
	Authenticate Authenticator = DefaultAuthenticator
)

// CheckPass compares a password to the salt/hash combination; returns true
// if they match. It may be used in custom 
func CheckPass(password string, salt, hash []byte) bool {
	if len(salt) == 0 || len(hash) == 0 {
		return false
	}
	ph := pbkdf2.PasswordHash{salt, hash}
	return pbkdf2.MatchPassword(password, ph)
}

// DefaultAuthenticator calls LookupCredentials to retrieve the user's
// credentials, and then validates those against the PBKDF2 standard.
func DefaultAuthenticator(user interface{}, password string) bool {
	salt, hash := LookupCredentials(user)
	return CheckPass(password, salt, hash)
}

// DefaultAuthProvider is an empty authentication provider that returns empty.
// Authentication fails by default.
func DefaultAuthProvider(user interface{}) (salt, hash []byte) {
	return
}

// HashPass is provided to return a salt and string for a password.
func HashPass(password string) (salt, hash []byte) {
	ph := pbkdf2.HashPassword(password)
	hash = ph.Hash
	salt = ph.Salt
	return
}
