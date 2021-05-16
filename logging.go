package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/natefinch/lumberjack"
)

var (
	Trace *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func StartLogging(Config Config) {
	fmt.Printf("\n****  Starting logger in [%-20s] **** \n", Config.Logfile)

	llogger := &lumberjack.Logger{
		Filename:   Config.Logfile,
		MaxSize:    Config.LogMaxSize, // megabytes
		MaxBackups: Config.LogMaxBackups,
		MaxAge:     Config.LogMaxAge, //days
		Compress:   false,            // disabled by default
	}

	var this io.Writer

	//Trace logger
	if Config.Loglevel >= 4 {
		this = llogger
	} else {
		this = ioutil.Discard
	}
	Trace = log.New(this,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	//Info logger
	if Config.Loglevel >= 3 {
		this = llogger
	} else {
		this = ioutil.Discard
	}
	Info = log.New(this,
		"INFO:  ",
		log.Ldate|log.Ltime|log.Lshortfile)

	//Error Logger
	this = llogger
	Error = log.New(this,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(this,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Trace.Print("Trace Message test")
	Info.Print("Info Message test")
	Error.Print("Error Message test")
	Info.Print("**** Logging Service is started ****")

}
