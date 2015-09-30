package generator

import (
	"strings"
	"testing"
)

var charsTest = []struct {
	given    string
	expected string
}{
	{"zxcv", "`~]'"},
	{"123", "`~]'"},
	{"a", "!&^#"},
	{"e", "!&^#"},
	{"i", "!&^#"},
	{"o", "!&^#"},
	{"u", "!&^#"},
	{"y", "!&^#"},
	{"ae", ")(*$"},
	{"axe", ")(*$"},
	{"aei", "[ -="},
	{"aeio", "@%.;"},
	{"aeiou", "<,}+"},
	{"aeiouy", "`~]'"},
	{"aeiouya", "!&^#"},
}

var genPassTests = []struct {
	Settings
	expected string
}{
	{Settings{MasterPhrase: "secret", Domain: "localhost", AdditionalInfo: "", PasswordLength: 12, AddSpecialCharacters: true}, "B8MYkTQT`~]'"},
	{Settings{MasterPhrase: "secret", Domain: "localhost", AdditionalInfo: "", PasswordLength: 12, AddSpecialCharacters: false}, "B8MYkTQTtUwW"},
	{Settings{MasterPhrase: "secret", Domain: "google.com", AdditionalInfo: "", PasswordLength: 12, AddSpecialCharacters: true}, "ODejwny3!&^#"},
	{Settings{MasterPhrase: "terces", Domain: "google.com", AdditionalInfo: "", PasswordLength: 12, AddSpecialCharacters: true}, "cLJk0Cnq!&^#"},
	{Settings{MasterPhrase: "terces", Domain: "google.com", AdditionalInfo: "", PasswordLength: 20, AddSpecialCharacters: true}, "cLJk0CnqwfDqjv4Y!&^#"},
	{Settings{MasterPhrase: "terces", Domain: "google.com", AdditionalInfo: "addInfo", PasswordLength: 12, AddSpecialCharacters: true}, "SKkSa4NN)(*$"},
	{Settings{MasterPhrase: "terces", Domain: "google.com", AdditionalInfo: "addInfo", PasswordLength: 12, AddSpecialCharacters: false}, "SKkSa4NN+5Xo"},
}

var invalidSettings = []struct {
	Settings
	expectedErr error
}{
	{Settings{}, ErrorEmptyPass},
	{Settings{MasterPhrase: "secret"}, ErrorEmptyDomain},
	{Settings{Domain: "google.com"}, ErrorEmptyPass},
	{Settings{MasterPhrase: "secret", Domain: "google.com"}, nil},
}

func TestSpecialCharactersGenerationLowerCase(t *testing.T) {
	for _, inOut := range charsTest {
		actual := getSpecialCharacters(inOut.given)
		if actual != inOut.expected {
			t.Errorf("getSpecialCharacters(%s): expected %s, actual %s", inOut.given, inOut.expected, actual)
		}
	}
}

func TestSpecialCharactersGenerationUpperCase(t *testing.T) {
	for _, inOut := range charsTest {
		given := strings.ToUpper(inOut.given)
		actual := getSpecialCharacters(given)
		if actual != inOut.expected {
			t.Errorf("getSpecialCharacters(%s): expected %s, actual %s", given, inOut.expected, actual)
		}
	}
}

func TestSpecialCharactersGenerationMixedCase(t *testing.T) {
	for _, inOut := range charsTest {
		given := toMixedCase(inOut.given)
		actual := getSpecialCharacters(given)
		if actual != inOut.expected {
			t.Errorf("getSpecialCharacters(%s): expected %s, actual %s", given, inOut.expected, actual)
		}
	}
}

// masterPhrase, domain, additionalInfo string, passLength int, addSpecialChars bool
func TestGeneratePassword(t *testing.T) {
	for _, inOut := range genPassTests {
		actual, err := GeneratePassword(inOut.Settings)
		if err != nil {
			t.Errorf("GeneratePassword returned an error %s", err)
		}
		actualStr := string(actual)
		if actualStr != inOut.expected {
			t.Errorf("GeneratePassword(%v): expected %s, actual %s", inOut.Settings, inOut.expected, actualStr)
		}
	}
}

func TestDefaultSettings(t *testing.T) {

	defaultPasswordLength := 12

	defaultSettings := DefaultSettings()
	if defaultSettings.Domain != "" || defaultSettings.MasterPhrase != "" || defaultSettings.AdditionalInfo != "" {
		t.Errorf("Default settings contain data, that must be empty! Domain: [%s], MasterPhrase: [%s], AdditionalInfo: [%s]\n", defaultSettings.Domain, defaultSettings.MasterPhrase, defaultSettings.AdditionalInfo)
	}

	if !defaultSettings.AddSpecialCharacters || defaultSettings.PasswordLength != defaultPasswordLength {
		t.Errorf("Default settings contain unexpected data! AddSpecialCharacters: expected [%s], actual [%s]; PasswordLength: expected [%d], actual [%d]\n", true, defaultSettings.AddSpecialCharacters, defaultPasswordLength, defaultSettings.PasswordLength)
	}
}

func TestNegTooLongPassRequirement(t *testing.T) {
	settings := Settings{MasterPhrase: "qwerty", Domain: "qwerty", PasswordLength: 1000}
	_, err := GeneratePassword(settings)
	if err == nil {
		t.Errorf("Expected error after requesting a too long password. Got nothing instead (error was nil)")
	}
}

func TestNegInvalidSettings(t *testing.T) {
	for _, inOut := range invalidSettings {
		_, actualErr := GeneratePassword(inOut.Settings)
		if actualErr != inOut.expectedErr {
			t.Errorf("GeneratePassword(%s): expected error %s, actual %s", inOut.Settings, inOut.expectedErr, actualErr)
		}
	}

}

func toMixedCase(input string) string {
	result := ""
	for i, c := range input {
		if i%2 == 0 {
			result += strings.ToUpper(string(c))
		} else {
			result += string(c)
		}
	}
	return result
}
