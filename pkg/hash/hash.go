package hash

import (
	"fmt"

	"github.com/vandi37/vanerrors"
	"golang.org/x/crypto/sha3"
)

// The errors
const (
	ErrorGettingHash = "error getting hash" // error getting hash
)

// The salt
var (
	SALT = ""
)

// A function that's getting the password and returning it's hash (sha-3-256)
func HashPassword(password string) (string, error) {
	// Creating a 256 byte hash as a slice of bytes
	hash := sha3.New256()
	_, err := hash.Write([]byte(password))

	// In case of error returning the error
	if err != nil {
		return "", vanerrors.NewWrap(ErrorGettingHash, err, vanerrors.EmptyHandler)
	}

	// Writing the slice in sha3
	sha3 := hash.Sum([]byte(SALT))

	// Returning the slice as a string
	return fmt.Sprintf("%x", sha3), nil
}

// A function to compare hash of a password and the password
func CompareHash(password string, hash string) (bool, error) {
	// Creating the hash of the password string
	hashedPassword, err := HashPassword(password)

	// In case of error returning the error
	if err != nil {
		return false, err
	}

	// Returning true if the hash are the same. If not false
	return hashedPassword == hash, nil
}
