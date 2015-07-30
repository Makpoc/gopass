package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	VOWELS        = "aeiouy"
	LOG_FILE_NAME = "domains.log"
)

var (
	// TODO - provide a way for users to define their own suffixes (either via file with suffix per row or a single suffix as param)
	SPECIAL_CHARS_GROUPS = []string{"`~]'", "!&^#", ")(*$", "[ -=", "@%.;", "<,}+"}

	masterPhrase      string
	masterPhraseFile  string
	domain            string
	additionalInfo    string
	passLength        int
	addSpecialChars   bool
	logDomainsAndInfo bool
)

// getSpecialCharacters selects from a predefined set of special characters and returns them
func getSpecialCharacters(encrypted string) string {
	vowelCount := 0
	for _, v := range VOWELS {
		vowelCount += strings.Count(strings.ToLower(string(encrypted)), string(v))
	}
	return SPECIAL_CHARS_GROUPS[vowelCount%len(SPECIAL_CHARS_GROUPS)]
}

// logInfo is used to log a "reminder log" containing the domains this tool was used for and any additional information that was provided. It is disabled by default.
func logInfo() {
	file, err := os.OpenFile(LOG_FILE_NAME, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		fmt.Sprintf("Failed to open log file %s. Error was: %s\n", LOG_FILE_NAME, err)
	}
	defer file.Close()

	data := fmt.Sprintf("Domain: [%s], Special Characters: [%v], AdditionalInfo: [%s]\n", domain, addSpecialChars, additionalInfo)
	_, err = file.WriteString(data)

	if err != nil {
		fmt.Sprintf("Failed to write to log file %s. Error was: %s\n", LOG_FILE_NAME, err)
	}
}

// parseMasterPhraseFromFile tries to load the master phrase from the provided file
func parseMasterPhraseFromFile(file string) (string, error) {
	_, err := os.Stat(file)
	if err != nil {
		return "", err
	}

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return strings.Trim(string(fileContent), "\r\n"), nil
}

// parseArgs parses the command line arguments and checks if they are valid
func parseArgs() {
	flag.StringVar(&masterPhrase, "master", "", "The master phrase to use for password generation. Required unless master-file is provided. Do NOT forget to escape any special characters contained in the master phrase (e.g. $, space etc).")
	flag.StringVar(&masterPhraseFile, "master-file", "", "The path to a file, containing the master phrase. Required unless master is provided.")
	flag.StringVar(&domain, "domain", "", "The domain for which this password is intended")
	flag.StringVar(&additionalInfo, "additional-info", "", "Free text to add (e.g. index/timestamp if the previous password was compromized)")
	flag.IntVar(&passLength, "password-length", 12, "Define the length of the password. Default: 12")
	flag.BoolVar(&addSpecialChars, "special-characters", true, "Whether to add a known set of special characters to the password")
	flag.BoolVar(&logDomainsAndInfo, "log-domain", false, "Whether to log the domain and the additional info for each generated password. Note that the password itself will NOT be stored!")

	flag.Parse()

	validateParams()
}

func pringUsageAndExit() {
	flag.Usage()
	os.Exit(1)
}

// validateParams validates all params that are provided on command line
func validateParams() {
	if (masterPhrase == "" && masterPhraseFile == "") || (masterPhrase != "" && masterPhraseFile != "") {
		fmt.Println("Either -master or -master-file must be specified")
		pringUsageAndExit()
	}

	if masterPhraseFile != "" {
		var err error
		if masterPhrase, err = parseMasterPhraseFromFile(masterPhraseFile); err != nil {
			fmt.Printf("Failed to retreive the master phrase from file! Error was: %s\n", err)
			os.Exit(1)
		}
	}

	if domain == "" {
		fmt.Println("-domain is required!")
		pringUsageAndExit()
	}

	if passLength < 1 {
		fmt.Println("-password-length must be a positive number!")
		pringUsageAndExit()
	}

	if passLength < 8 {
		fmt.Println("WARN: -password-length is too short. We will grant your wish, but this might be a security risk. Consider using longer password.")
	}
}

func main() {
	parseArgs()

	if logDomainsAndInfo {
		logInfo()
	}

	pass := []byte(fmt.Sprintf("%s:%s:%s", masterPhrase, domain, additionalInfo))
	encrypted := sha256.Sum256(pass)
	fullPass := base64.StdEncoding.EncodeToString(encrypted[:])

	if len(fullPass) < passLength {
		fmt.Printf("Cannot generate password with so many symbols. The current limit is [%d]. Please lower the -password-length value.\n", len(fullPass))
		os.Exit(1)
	}

	trimmedPass := fullPass[:passLength]

	if addSpecialChars {
		charsToAdd := getSpecialCharacters(fullPass)
		trimmedPass = trimmedPass[:len(trimmedPass)-len(charsToAdd)] + charsToAdd
	}

	fmt.Printf("Your password for %s is: %s\n", domain, trimmedPass)
}
