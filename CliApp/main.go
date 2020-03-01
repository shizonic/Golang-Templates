package main

import (
	"fmt"
	"os"
	"time"

	"github.com/JojiiOfficial/Golang-Templates/GoLogServer/constants"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	appName = "clapp"
	version = "1.0.0"

	//EnVarPrefix prefix for env vars
	EnVarPrefix = "GOLOG"
)

var (
	app         = kingpin.New(appName, "A Logging server")
	appLogLevel = app.Flag("log-level", "Enable debug mode").HintOptions(constants.LogLevels...).Envar(getEnVar(EnVarLogLevel)).Short('l').Default(constants.LogLevels[2]).String()
	appNoColor  = app.Flag("no-color", "Disable colors").Envar(getEnVar(EnVarNoColor)).Bool()
	appYes      = app.Flag("yes", "Skips confirmations").Short('y').Envar(getEnVar(EnVarYes)).Bool()
	appCfgFile  = app.
			Flag("config", "the configuration file for the app").
			Envar(getEnVar(EnVarConfigFile)).
			Short('c').String()

	//Server commands
	//Server start
	serverCmd      = app.Command("server", "Commands for the server")
	serverCmdStart = serverCmd.Command("start", "Start the server")
)

var (
	//	config  Config
	isDebug = false
)

func main() {
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

	switch parsed {
	case serverCmdStart.FullCommand():

	}
}

//Env vars
const (
	//EnVarPrefix prefix of all used env vars

	EnVarLogLevel   = "LOG_LEVEL"
	EnVarNoColor    = "NO_COLOR"
	EnVarYes        = "SKIP_CONFIRM"
	EnVarConfigFile = "CONFIG"
)

//Return the variable using the server prefix
func getEnVar(name string) string {
	return fmt.Sprintf("%s_%s", EnVarPrefix, name)
}