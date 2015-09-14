package generator

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

const vowels = "aeiouy"

// TODO - provide a way for users to define their own suffixes (either via file with suffix per row or a single suffix as param). Beware that changing the (order of) suffixes will result in different password. There must be a better way...
var specialCharsGroups = []string{"`~]'", "!&^#", ")(*$", "[ -=", "@%.;", "<,}+"}

// getSpecialCharacters selects from a predefined set of special characters and returns them. For now the "algorithm" is based on the number of vowels in the encrypted string.
func getSpecialCharacters(encrypted string) string {
	vowelCount := 0
	for _, v := range vowels {
		vowelCount += strings.Count(strings.ToLower(encrypted), string(v))
	}
	return specialCharsGroups[vowelCount%len(specialCharsGroups)]
}

// GeneratePassword generates the domain specific password.
func GeneratePassword(masterPhrase, domain, additionalInfo string, passLength int, addSpecialChars bool) ([]byte, error) {
	passToEncrypt := fmt.Sprintf("%s:%s:%s", masterPhrase, domain, additionalInfo)

	encrypted := sha256.Sum256([]byte(passToEncrypt))
	fullEncryptHash := base64.StdEncoding.EncodeToString(encrypted[:])

	if len(fullEncryptHash) < passLength {
		return nil, fmt.Errorf("Cannot generate password with so many symbols. The current limit is [%d]. Please lower the -password-length value.", len(fullEncryptHash))
	}

	trimmedPass := fullEncryptHash[:passLength]

	if addSpecialChars {
		charsToAdd := getSpecialCharacters(fullEncryptHash)
		trimmedPass = trimmedPass[:len(trimmedPass)-len(charsToAdd)] + charsToAdd
	}

	return []byte(trimmedPass), nil
}
