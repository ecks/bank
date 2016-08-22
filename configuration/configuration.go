package configuration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"github.com/kardianos/osext"
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v3"
)

type Configuration struct {
	TimeZone     string
	MySQLUser    string
	MySQLPass    string
	MySQLHost    string
	MySQLPort    string
	MySQLDB      string
	RedisHost    string
	RedisPort    string
	PasswordSalt string
	FQDN         string
	HttpPort     string
	Db           *sql.DB
	Redis        *redis.Client
	PushEnv      string
}

// Initialization of the working directory. Needed to load asset files.
var ImportPath = os.Getenv("GOPATH") + "/src/github.com/bvnk/bank/"

//configPath must be an absolute path
var configPath = ImportPath + "config.json"

func LoadConfig() (configuration Configuration, err error) {

	workingDir, err := os.Getwd()
	if err != nil {
		return Configuration{}, errors.New("configuration.LoadConfig: Could not load config. " + err.Error())
	}

	//Go uses this path during tests
	tempDir := os.TempDir()
	//fmt.Println("TempDir: ", workingDir)

	//Check if testing
	isTest := strings.Index(workingDir, tempDir)

	if isTest != -1 {
		/*
			We're testing, so we must be in Go's temp directory
			We can assume the import path is github.com/[user]/bank
			We will replace anything before github with GOPATH/src
			We will also remove anything after bank
			@TODO Implement a more generic workaround
		*/
		workingDir = regexp.MustCompile(".*github").ReplaceAllString(workingDir, "github")
		ImportPath = os.Getenv("GOPATH") + "/src/" + regexp.MustCompile("/bank/.*").ReplaceAllString(workingDir, "/bank/")
	}

	var configPath = ImportPath + "config.json"

	// Get config
	file, _ := os.Open(configPath)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)

	if err != nil {
		return Configuration{}, errors.New("configuration.LoadConfig: Could not load config. " + err.Error())
	}

	// Load MySQL
	err = loadMySQL(&configuration)
	if err != nil {
		return Configuration{}, errors.New("configuration.LoadConfig: Could not load MySQL. " + err.Error())
	}
	// Load Redis
	loadRedis(&configuration)

	return
}

func loadMySQL(configuration *Configuration) (err error) {
	configuration.Db, err = sql.Open("mysql", configuration.MySQLUser+":"+configuration.MySQLPass+"@tcp("+configuration.MySQLHost+":"+configuration.MySQLPort+")/"+configuration.MySQLDB)
	if err != nil {
		return errors.New("configuration.loadMySQL: Could not connect to database")
	}

	return
}

func loadRedis(configuration *Configuration) {
	configuration.Redis = redis.NewClient(&redis.Options{
		Addr:     configuration.RedisHost + ":" + configuration.RedisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}

func determineWorkingDirectory() string {
	var customPath string

	// Check if a custom path has been provided by the user.
	flag.StringVar(&customPath, "p", "", "Specify a custom path to the asset files. This needs to be an absolute path.")
	flag.Parse()

	// Get the absolute path this executable is located in.
	executablePath, err := osext.ExecutableFolder()

	if err != nil {
		log.Fatal("Error: Couldn't determine working directory: " + err.Error())
	}
	// Set the working directory to the path the executable is located in.
	os.Chdir(executablePath)

	// Return the user-specified path. Empty string if no path was provided.
	if customPath != "" {
		return customPath + "/"
	}
	return customPath
}
