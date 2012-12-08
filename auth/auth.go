// Package auth provides PBKDF2 authentication for webapps.
// This package is designed so that users should really only have to supply
// their own means of procuring an AuthProvider in order to make use of the
// provided functionality.
//
// To return a hash, salt from a password, the HashPass function is provided.
// The end user can store these safely in their database.
package auth

import (
        "fmt"
        "github.com/gokyle/pbkdf2"
        "github.com/gokyle/uuid"
        "net/http"
        "time"
)

var (
        DefaultCheck time.Duration
        DefaultExpire time.Duration
)

func init() {
        var err error
        DefaultCheck, err = time.ParseDuration("2h")
        if err != nil {
                panic("auth - error parsing duration: " + err.Error())
        }

        DefaultExpire, err = time.ParseDuration("1h")
        if err != nil {
                panic("auth - error parsing duration: " + err.Error())
        }
}

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
// if they match. It may be used in custom Authenticators.
func CheckPass(password string, salt, hash []byte) bool {
	if len(salt) == 0 || len(hash) == 0 {
		return false
	}
	ph := &pbkdf2.PasswordHash{hash, salt}
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

// Type SessionStore provides basic session support via authentication for Go
// web apps. SessionStores should be created with CreateSessionStore.
type SessionStore struct {
        Name            string
        Sessions        map[string]*time.Time
        Check           time.Duration
        Secure          bool
}

// CreateSessionStore creates new SessionStores properly.
// Name is the name new cookies are given; this should be something that
// identifies the cookie as a session cookie for your application. Secure
// should be true if the app itself runs as a TLS server. If the app
// runs behind an SSL server (i.e. an nginx proxy), it should be set to false.
// dur is a time.Duration that specifies how long cookies will live.
func CreateSessionStore(name string, secure bool, dur *time.Duration) *SessionStore {
        store := &SessionStore{name, nil, DefaultCheck, secure}
        store.Sessions = make(map[string]*time.Time, 0)
        if dur != nil {
                store.Check = *dur
        }
        go store.CheckExpired()
        return store
}

func (store *SessionStore) _checkExpired() {
        for k, t := range store.Sessions {
                if t != nil && t.After(time.Now()) {
                        delete(store.Sessions, k)
                }
        }
}

// CheckExpired typically runs in the background, launched in
// CreateSessionStore; it may be called to manually / force check the store
// for expired cookies.
func (store *SessionStore) CheckExpired() {
        for {
                <-time.After(store.Check)
                store._checkExpired()
        }
}

// NewSession creates a new session without authenticating. This should
// be used if you are authentication in some way other than
// Authenticate. t should be a duration specifying when the session
// should expire, and no_expire should be true if it should never
// expire. If no_expire is false and t is nil, the default expiration
// will be used.
func (store *SessionStore) NewSession() (c *http.Cookie, err error) {
        var session_id string
        session_id, err = uuid.GenerateV4String()
        if err != nil {
                return
        }
        t := time.Now().Add(DefaultExpire)
        store.Sessions[session_id] = &t

        c = new(http.Cookie)
        c.Name = store.Name
        c.Value = session_id
        c.Path = "/"
        c.Secure = false
        c.HttpOnly = true
        return
}

// NewPSession creates a persistent session, i.e. one that survives the
// current session.
func (store *SessionStore) NewPSession(age string) (c *http.Cookie, err error) {
        var session_id string
        session_id, err = uuid.GenerateV4String()
        if err != nil {
                return
        }
        dur, err := time.ParseDuration(age)
        if err != nil {
                return
        }
        t := time.Now().Add(dur)
        store.Sessions[session_id] = &t

        c = new(http.Cookie)
        c.Name = store.Name
        c.Value = session_id
        c.Path = "/"
        c.Expires = *store.Sessions[session_id]
        c.Secure = false
        c.HttpOnly = true
        return
}


// AuthSession attempts to authenticate the user; if successful, it returns a
// cookie that can be sent to the client to set up a session. If
// persistent is false, the cookie will last for the current session
// only.
func (store *SessionStore) AuthSession(id interface{}, pass string, persist bool, age string) (c *http.Cookie, err error) {
        if !Authenticate(id, pass) {
                return
        }

        if !persist {
                c, err = store.NewSession()
        } else {
                c, err = store.NewPSession(age)
        }
        return
}

// CheckSession checks a client request to see if it contains a valid session
// cookie. If this fails, the client should be re-authenticated.
func (store *SessionStore) CheckSession(r *http.Request) bool {
        for _, c := range r.Cookies() {
                if c.Name != store.Name {
                        continue
                }
                sid := c.Value
                t, valid := store.Sessions[sid]
                if !valid {
                        return false
                }
                if t == nil || time.Now().After(*t) {
                        return false
                }
                return true
        }
        fmt.Println("I couldn't find any valid cookies!")
        return false
}

// DestroySession should be called for any logout action or when the session
// should not be maintained anymore.
func (store *SessionStore) DestroySession(r *http.Request) bool {
        for _, c := range r.Cookies() {
                if c.Name == store.Name && c.Domain == r.URL.Host {
                        delete(store.Sessions, c.Value)
                        return true
                }
        }
        return false
}
