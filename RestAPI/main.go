package main

import (
	"Golang-Templates/RestAPI/constants"
	"Golang-Templates/RestAPI/models"
	"Golang-Templates/RestAPI/storage"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	dbhelper "github.com/JojiiOfficial/GoDBHelper"
	"gopkg.in/alecthomas/kingpin.v2"

	_ "github.com/go-sql-driver/mysql"
)

const version = "0.0.1a"

var (
	app         = kingpin.New("server", "A Rest server")
	appLogLevel = app.Flag("log-level", "Enable debug mode").HintOptions(constants.LogLevels...).Envar(getEnVar(EnVarLogLevel)).Short('l').Default(constants.LogLevels[2]).String()
	appNoColor  = app.Flag("no-color", "Disable colors").Envar(getEnVar(EnVarNoColor)).Bool()
	appYes      = app.Flag("yes", "Skips confirmations").Short('y').Envar(getEnVar(EnVarYes)).Bool()
	appCfgFile  = app.
			Flag("config", "the configuration file for the server").
			Envar(getEnVar(EnVarConfigFile)).
			Short('c').String()

	//Server commands
	//Server start
	serverCmd      = app.Command("server", "Commands for the server")
	serverCmdStart = serverCmd.Command("start", "Start the server")

	//Config commands
	//Config create
	configCmd           = app.Command("config", "Commands for the config file")
	configCmdCreate     = configCmd.Command("create", "Create config file")
	configCmdCreateName = configCmdCreate.Arg("name", "Config filename").Default(models.GetDefaultConfig()).String()
)

var (
	config  *models.Config
	db      *dbhelper.DBhelper
	isDebug bool = false
)

//Env vars
const (
	//EnVarPrefix prefix of all used env vars
	EnVarPrefix = "REST"

	EnVarLogLevel   = "LOG_LEVEL"
	EnVarNoColor    = "NO_COLOR"
	EnVarYes        = "SKIP_CONFIRM"
	EnVarConfigFile = "CONFIG"
)

//Return the variable using the server prefix
func getEnVar(name string) string {
	return fmt.Sprintf("%s_%s", EnVarPrefix, name)
}

func main() {
	//Set app attributes
	app.HelpFlag.Short('h')
	app.Version(version)

	//parsing the args
	parsed := kingpin.MustParse(app.Parse(os.Args[1:]))

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.Stamp,
		FullTimestamp:    true,
		ForceColors:      !*appNoColor,
		DisableColors:    *appNoColor,
	})

	log.Infof("LogLevel: %s\n", *appLogLevel)

	//set app logLevel
	switch *appLogLevel {
	case constants.LogLevels[0]:
		//Debug
		log.SetLevel(log.DebugLevel)
		isDebug = true
	case constants.LogLevels[1]:
		//Info
		log.SetLevel(log.InfoLevel)
	case constants.LogLevels[2]:
		//Warning
		log.SetLevel(log.WarnLevel)
	case constants.LogLevels[3]:
		//Error
		log.SetLevel(log.ErrorLevel)
	default:
		fmt.Println("LogLevel not found!")
		os.Exit(1)
		return
	}

	if parsed != configCmdCreate.FullCommand() {
		var shouldExit bool
		config, shouldExit = models.InitConfig(*appCfgFile, false)
		if shouldExit {
			return
		}

		if !config.Check() {
			log.Info("Exiting")
			return
		}

		var err error
		db, err = storage.ConnectDB(config, isDebug, *appNoColor)
		err = nil
		if err != nil {
			log.Fatalln(err.Error())
			return
		}

	}

	switch parsed {
	//Server --------------------
	case serverCmdStart.FullCommand():
		{
			startAPI()
		}

	//Config --------------------
	case configCmdCreate.FullCommand():
		{
			//whsub config create
			models.InitConfig(*configCmdCreateName, true)
		}
	}
}
