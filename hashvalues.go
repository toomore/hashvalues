package hashvalues

import (
	"crypto/hmac"
	"errors"
	"fmt"
	"hash"
	"net/url"
)

type HashValues struct {
	Values   url.Values
	hashfunc func() hash.Hash
	hashkey  []byte
}

func New(hashkey []byte, hashfunc func() hash.Hash) *HashValues {
	return &HashValues{
		Values:   url.Values{},
		hashfunc: hashfunc,
		hashkey:  hashkey,
	}
}

func (h *HashValues) Set(key, value string) {
	h.Values.Set(key, value)
}

func (h *HashValues) Add(key, value string) {
	h.Values.Add(key, value)
}

func (h *HashValues) Del(key string) {
	h.Values.Del(key)
}

func (h *HashValues) Get(key string) string {
	return h.Values.Get(key)
}

func (h *HashValues) Decode(key []byte, message string) error {
	var err error
	var hashed = hmac.New(h.hashfunc, h.hashkey)

	hashed.Write([]byte(message))

	if fmt.Sprintf("%x", hashed.Sum(nil)) == fmt.Sprintf("%s", key) {
		h.Values, err = url.ParseQuery(message)
	} else {
		err = errors.New("wrong key!")
	}
	return err
}

func (h *HashValues) Encode() ([]byte, string) {
	var value = h.Values.Encode()

	hm := hmac.New(h.hashfunc, h.hashkey)
	hm.Write([]byte(value))

	return hm.Sum(nil), value
}
