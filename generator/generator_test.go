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

type generatorInput struct {
	pass            string
	domain          string
	additionalInfo  string
	passLength      int
	addSpecialChars bool
}

var genPassTests = []struct {
	generatorInput
	expected string
}{
	{generatorInput{pass: "secret", domain: "localhost", additionalInfo: "", passLength: 12, addSpecialChars: true}, "B8MYkTQT`~]'"},
	{generatorInput{pass: "secret", domain: "localhost", additionalInfo: "", passLength: 12, addSpecialChars: false}, "B8MYkTQTtUwW"},
	{generatorInput{pass: "secret", domain: "google.com", additionalInfo: "", passLength: 12, addSpecialChars: true}, "ODejwny3!&^#"},
	{generatorInput{pass: "terces", domain: "google.com", additionalInfo: "", passLength: 12, addSpecialChars: true}, "cLJk0Cnq!&^#"},
	{generatorInput{pass: "terces", domain: "google.com", additionalInfo: "", passLength: 20, addSpecialChars: true}, "cLJk0CnqwfDqjv4Y!&^#"},
	{generatorInput{pass: "terces", domain: "google.com", additionalInfo: "addInfo", passLength: 12, addSpecialChars: true}, "SKkSa4NN)(*$"},
	{generatorInput{pass: "terces", domain: "google.com", additionalInfo: "addInfo", passLength: 12, addSpecialChars: false}, "SKkSa4NN+5Xo"},
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
		actual, err := GeneratePassword(inOut.pass, inOut.domain, inOut.additionalInfo, inOut.passLength, inOut.addSpecialChars)
		if err != nil {
			t.Errorf("GeneratePassword returned an error %s", err)
		}
		actualStr := string(actual)
		if actualStr != inOut.expected {
			t.Errorf("GeneratePassword(%s, %s, %s, %d, %v): expected %s, actual %s", inOut.pass, inOut.domain, inOut.additionalInfo, inOut.passLength, inOut.addSpecialChars, inOut.expected, actualStr)
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
