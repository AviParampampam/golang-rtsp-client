package auth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
)

// Digest access authentication
type Digest struct{}

func md5Hash(text string) hash.Hash {
	h := md5.New()
	io.WriteString(h, text)
	return h
}

func string16(h hash.Hash) string {
	return hex.EncodeToString(h.Sum(nil))
}

// MD5String - Generating text to string MD5
func MD5String(text string) string {
	return string16(md5Hash(text))
}

func (d Digest) h1(username, realm, password string) string {
	return MD5String(fmt.Sprintf("%s:%s:%s", username, realm, password))
}

func (d Digest) h2(method, uri string) string {
	return MD5String(fmt.Sprintf("%s:%s", method, uri))
}

// Generating - Generation "Authorization" header Digest
func (d Digest) Generating(username, password, realm, nonce, method, uri string) string {
	HA1 := d.h1(username, realm, password)
	HA2 := d.h2(method, uri)
	return MD5String(fmt.Sprintf("%s:%s:%s", HA1, nonce, HA2))
}
