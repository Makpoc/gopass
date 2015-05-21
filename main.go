package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	SPECIAL_CHARS_PREFIX = "!@^#"
	SPECIAL_CHARS_SUFFIX = ")(*$"

	VOWELS = "aeiouy"
)

var (
	masterPhrase     string
	masterPhraseFile string
	domain           string
	additionalInfo   string
	passLength       int
	addSpecialChars  bool
)

// addSpecialCharacters adds a predefined set of special characters to the beginning or the end of the password, replacing the corresponding amount of symbols
func addSpecialCharacters(pass, encrypted string) string {
	vowelCount := 0
	passHash := sha256.Sum256([]byte(pass))
	passHashStr := base64.StdEncoding.EncodeToString(passHash[:])
	for _, v := range VOWELS {
		vowelCount += strings.Count(strings.ToLower(string(passHashStr)), string(v))
	}
	// TODO - make a more complex algorithm
	fmt.Printf("DEBUG: %d vowels found in %s\n", vowelCount, passHashStr)
	if vowelCount%2 == 0 {
		cutToIndex := passLength - len(SPECIAL_CHARS_PREFIX)
		return fmt.Sprintf("%s%s", encrypted[:cutToIndex], SPECIAL_CHARS_SUFFIX)
	} else {
		cutFromIndex := len(SPECIAL_CHARS_PREFIX)
		return fmt.Sprintf("%s%s", SPECIAL_CHARS_PREFIX, encrypted[cutFromIndex:passLength])
	}
}

// parseMasterPhraseFromFile tries to load the master phrase from the provided file
func parseMasterPhraseFromFile(file string) error {
	_, err := os.Stat(file)
	if err != nil {
		return err
	}

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	masterPhrase = string(fileContent)
	return nil
}

// parseArgs parses the command line arguments and checks if they are valid
func parseArgs() {
	flag.StringVar(&masterPhrase, "master", "", "The master phrase to use for password generation. Required unless master-file is provided.")
	flag.StringVar(&masterPhraseFile, "master-file", "", "The path to a file, containing the master phrase. Required unless master is provided.")
	flag.StringVar(&domain, "domain", "", "The domain for which this password is intended")
	flag.StringVar(&additionalInfo, "additional-info", "", "Free text to add (e.g. index/timestamp if the previous password was compromized)")
	flag.IntVar(&passLength, "password-length", 12, "Define the length of the password. Default: 12")
	flag.BoolVar(&addSpecialChars, "special-characters", true, "Whether to add a known set of special characters to the password")

	flag.Parse()

	if (masterPhrase == "" && masterPhraseFile == "") || (masterPhrase != "" && masterPhraseFile != "") {
		log.Fatal("Either master or master-file must be specified")
	}

	if masterPhraseFile != "" {
		if err := parseMasterPhraseFromFile(masterPhraseFile); err != nil {
			log.Fatalf("Failed to retreive the master phrase from file! Error was: %s", err)
		}
	}

	if domain == "" {
		log.Fatal("domain is required!")
	}
}

func clean(pass string) string {
	pass = strings.Replace(pass, "/", "", -1)
	pass = strings.Replace(pass, "+", "", -1)
	return pass
}

func main() {
	parseArgs()

	pass := []byte(fmt.Sprintf("%s:%s:%s", masterPhrase, domain, additionalInfo))
	encrypted := sha256.Sum256(pass)

	finalPass := clean(base64.StdEncoding.EncodeToString(encrypted[:]))
	if addSpecialChars {
		finalPass = addSpecialCharacters(string(pass), finalPass)
	} else {
		finalPass = finalPass[:passLength]
	}

	fmt.Printf("Your password for %s is: %s\n", domain, finalPass)
}
