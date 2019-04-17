package secret

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/acharyab/gophercises/ex17-secrets-api/secret/cipher"
)

// Vault contains an encodingKey and the a map of keyValues
type Vault struct {
	encodingKey string
	filepath    string
	mutex       sync.Mutex
	keyValues   map[string]string
}

// File returns a Vault from memory
func File(encodingKey, filepath string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    filepath,
	}
}

func (v *Vault) load() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()
	r, err := cipher.DecryptReader(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.readKeyValues(r)
}

func (v *Vault) readKeyValues(r io.Reader) error {
	dec := json.NewDecoder(r) // -> decryptReader -> file
	return dec.Decode(&v.keyValues)
}

func (v *Vault) save() error {
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	w, err := cipher.EncryptWriter(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.writeKeyValues(w)
}

func (v *Vault) writeKeyValues(w io.Writer) error {
	enc := json.NewEncoder(w) // -> encrypt writer -> file
	return enc.Encode(v.keyValues)
}

// Get returns the value of the encoding key
func (v *Vault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.load()
	if err != nil {
		return "", err
	}
	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}
	return value, nil
}

// Set will set the value of the encoding key
func (v *Vault) Set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.load()
	if err != nil {
		return err
	}
	v.keyValues[key] = value
	err = v.save()
	return err
}
