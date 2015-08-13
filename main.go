package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"time"
)

const (
	VOWELS      = "aeiouy"
	GOPASS_HOME = "GOPASS_HOME"
)

var (
	// TODO - provide a way for users to define their own suffixes (either via file with suffix per row or a single suffix as param)
	SPECIAL_CHARS_GROUPS = []string{"`~]'", "!&^#", ")(*$", "[ -=", "@%.;", "<,}+"}

	defaultMasterFileName = "master"
	logFileName           = "domains.log"

	configFolder      string
	masterPhrase      string
	masterPhraseFile  string
	domain            string
	additionalInfo    string
	passLength        int
	addSpecialChars   bool
	logDomainsAndInfo bool
)

// initHome initializes the configFolder to point to the default or user-definied home
func initHome() error {
	if cf := os.Getenv(GOPASS_HOME); cf != "" {
		configFolder = cf
		return nil
	}

	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("Failed to retrieve current user! Error was %s", err)
	}

	configFolder = path.Join(usr.HomeDir, ".gopass")
	return nil
}

// getSpecialCharacters selects from a predefined set of special characters and returns them. For now the "algorithm" is based on the number of vowels in the encrypted string.
func getSpecialCharacters(encrypted string) string {
	vowelCount := 0
	for _, v := range VOWELS {
		vowelCount += strings.Count(strings.ToLower(encrypted), string(v))
	}
	return SPECIAL_CHARS_GROUPS[vowelCount%len(SPECIAL_CHARS_GROUPS)]
}

// logInfo is used to log a "reminder log" containing the domains this tool was used for and any additional information that was provided. It is disabled by default.
func logInfo() {
	file, err := os.OpenFile(path.Join(configFolder, logFileName), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		fmt.Sprintf("Failed to open log file %s. Error was: %s\n", logFileName, err)
	}
	defer file.Close()

	data := fmt.Sprintf("Date: [%s], Domain: [%s], Special Characters: [%v], AdditionalInfo: [%s]\n", time.Now().Format(time.RFC3339), domain, addSpecialChars, additionalInfo)
	_, err = file.WriteString(data)

	if err != nil {
		fmt.Sprintf("Failed to write to log file %s. Error was: %s\n", logFileName, err)
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

func parseMasterPhrase() error {
	if masterPhrase != "" {
		// master phrase provided on cmd
		return nil
	}

	var err error

	if masterPhraseFile != "" {
		// master phrase provided as file
		if masterPhrase, err = parseMasterPhraseFromFile(masterPhraseFile); err != nil {
			return fmt.Errorf("Failed to retrieve the master from file! Error was: %s", err)
		}
	}

	// try to load from the default file
	masterPhrase, err = parseMasterPhraseFromFile(path.Join(configFolder, "master"))
	if err != nil {
		return fmt.Errorf("Failed to retrieve master from default file. Error was: %s", err)
	}

	return nil
}

// parseArgs parses the command line arguments and checks if they are valid
func parseArgs() {
	flag.StringVar(&masterPhrase, "master", "", "The master phrase to use for password generation. Do NOT forget to escape any special characters contained in the master phrase (e.g. $, space etc).")
	flag.StringVar(&masterPhraseFile, "master-file", "", "The path to a file, containing the master phrase.")

	flag.StringVar(&domain, "domain", "", "The domain for which this password is intended")
	flag.StringVar(&additionalInfo, "additional-info", "", "Free text to add (e.g. index/timestamp/username if the previous password was compromized)")
	flag.IntVar(&passLength, "password-length", 12, "Define the length of the password.")
	flag.BoolVar(&addSpecialChars, "special-characters", true, "Whether to add a known set of special characters to the password")
	flag.BoolVar(&logDomainsAndInfo, "log-domain", false, "Whether to log the parameters that were used for generation to a file. Note that the password itself will NOT be stored!")

	flag.Parse()

	err := parseMasterPhrase()
	if err != nil {
		printAndExit(err.Error(), true)
	}

	validateParams()
}

func printAndExit(errorMsg string, printUsage bool) {
	if errorMsg != "" {
		fmt.Println(errorMsg)
	}
	if printUsage {
		flag.Usage()
	}
	os.Exit(1)
}

// validateParams validates all params that are provided on command line
func validateParams() {
	if domain == "" {
		printAndExit("-domain is required!", true)
	}

	if passLength < 1 {
		printAndExit("-password-length must be a positive number!", true)
	}

	if passLength < 8 {
		fmt.Println("WARN: -password-length is too short. We will grant your wish, but this might be a security risk. Consider using longer password.", false)
	}
}

//constructPasswordToEncrypt combines the different strings, put in the password base (domain, masterPhrase, ...)
func constructPasswordToEncrypt() string {
	return fmt.Sprintf("%s:%s:%s", masterPhrase, domain, additionalInfo)
}

// GeneratePassword does the actual work - encrypt the base string, trims the password to size and appends the special characters.
func GeneratePassword() error {
	encrypted := sha256.Sum256([]byte(constructPasswordToEncrypt()))
	fullEncryptHash := base64.StdEncoding.EncodeToString(encrypted[:])

	if len(fullEncryptHash) < passLength {
		return fmt.Errorf("Cannot generate password with so many symbols. The current limit is [%d]. Please lower the -password-length value.", len(fullEncryptHash))
	}

	trimmedPass := fullEncryptHash[:passLength]

	if addSpecialChars {
		charsToAdd := getSpecialCharacters(fullEncryptHash)
		trimmedPass = trimmedPass[:len(trimmedPass)-len(charsToAdd)] + charsToAdd
	}

	fmt.Printf("Your password for %s is: %s\n", domain, trimmedPass)
	return nil

}

func main() {
	if err := initHome(); err != nil {
		printAndExit(fmt.Sprintf("Failed to initialize default settings. Error was %s", err.Error), false)
	}

	parseArgs()

	if err := GeneratePassword(); err != nil {
		printAndExit(err.Error(), false)
	}

	if logDomainsAndInfo {
		logInfo()
	}

}
