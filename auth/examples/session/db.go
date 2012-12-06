package main

type cred struct {
	Name string
	Hash []byte
	Salt []byte
}

var db = make([]cred, 0)

func LookupUser(id interface{}) (hash, salt []byte) {
	for _, u := range db {
		if u.Name == id.(string) {
			hash = u.Hash
			salt = u.Salt
			break
		}
	}
	return
}

func SetUserPass(id interface{}, hash, salt []byte) {
	for _, u := range db {
		if u.Name == id.(string) {
			u.Hash = hash
			u.Salt = salt
			return
		}
	}
	u := cred{id.(string), hash, salt}
	db = append(db, u)
}
