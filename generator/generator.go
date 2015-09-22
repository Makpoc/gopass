package generator

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// Settings holds all properties for the password, which will be generated.
type Settings struct {
	MasterPhrase         string /// TODO - add tags
	Domain               string
	AdditionalInfo       string
	PasswordLength       int
	AddSpecialCharacters bool
}

var (
	ErrorEmptyPass   = errors.New("Empty master password")
	ErrorEmptyDomain = errors.New("Empty domain")
)

const vowels = "aeiouy"

func DefaultSettings() Settings {
	return Settings{PasswordLength: 12, AddSpecialCharacters: true}
}

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

func validateSettings(settings Settings) error {
	if settings.MasterPhrase == "" {
		return ErrorEmptyPass
	}
	if settings.Domain == "" {
		return ErrorEmptyDomain
	}

	return nil
}

// GeneratePassword generates the domain specific password.
func GeneratePassword(settings Settings) ([]byte, error) {
	if err := validateSettings(settings); err != nil {
		return nil, err
	}

	passToEncrypt := fmt.Sprintf("%s:%s:%s", settings.MasterPhrase, settings.Domain, settings.AdditionalInfo)

	encrypted := sha256.Sum256([]byte(passToEncrypt))
	fullEncryptHash := base64.StdEncoding.EncodeToString(encrypted[:])

	if len(fullEncryptHash) < settings.PasswordLength {
		return nil, fmt.Errorf("Cannot generate password with so many symbols. The current limit is [%d]. Please lower the -password-length value.", len(fullEncryptHash))
	}

	trimmedPass := fullEncryptHash[:settings.PasswordLength]

	if settings.AddSpecialCharacters {
		charsToAdd := getSpecialCharacters(fullEncryptHash)
		trimmedPass = trimmedPass[:len(trimmedPass)-len(charsToAdd)] + charsToAdd
	}

	return []byte(trimmedPass), nil
}
