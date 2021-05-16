package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// global vars users to be able to use the play function

var chanLoad chan loadQuery
var chanStore chan storeQuery
var chanList chan chan string
var storageDone chan interface{}

func main() {

	// Load the configuration
	fmt.Println("**** LOADING BEE CONFIGURATION FILE ****")
	Config := LoadConfiguration("config.cfg")
	fmt.Println("Configuration loaded. ")

	//SetupLogging. Loggers are Debug Info and Error
	fmt.Println("**** STRATING LOGGING COMPONENT ****")
	StartLogging(Config)
	Info.Println("**** LOGGING SERVICE HAS BEEN STARTED - STARTING BEE ****")
	Info.Print(Config.String())

	//setup storage abstraction layer

	chanLoad, chanStore, chanList, storageDone = Config.setupStorage()

	_ = chanLoad
	_ = chanStore
	_ = chanList
	_ = storageDone

	//starting web servers

	serverDone := Config.setupWebServers()

	//Play during creation
	//	play()
	Info.Println("PLAY IS FINISH")
	//
	//
	//
	//MAIN LOOP **********************************************************************
	mainDone := Config.startWatchdog() // mainDone is the global termination channel
	if mainDone == nil {
		fmt.Println("Already running, silent exit")
		os.Exit(0)
	}
	for {
		select {
		case <-mainDone:
			Info.Println("MAIN::sending httpx servers termination request")
			close(serverDone)
			Info.Println("MAIN::sending storage termination request")
			close(storageDone)
			time.Sleep(500 * time.Millisecond)
			Info.Println("MAIN::exit now")
			os.Exit(0)
		}
	}

}

func play() {

	//play **********************************************************************

	var x = &beeObj{UID: "MYSELF", AuthKey: "AZZEERT"}
	x.AddNotification("coucou")
	x.AddNotification("caca")

	Info.Println(x)
	fmt.Println(x)
	var marshal, err = x.marshal()
	fmt.Println(marshal)
	_ = err
	fmt.Println("Do a storage request now..")
	fmt.Println(x)
	ce := make(chan error)
	chanStore <- storeQuery{x, ce}
	fmt.Println("End of storage request")
	fmt.Println(x)

	{
		fmt.Println("Do a load  request now..")
		ce := make(chan error)
		cr := make(chan beeObj)
		chanLoad <- loadQuery{"myself", cr, ce}
		select {
		case err := <-ce:
			fmt.Printf("Error loading : %s", err)
			break
		case obj := <-cr:
			fmt.Printf("Obj loaded :  %s", obj)
			break
		}

		fmt.Println("End of load request")
	}

	Info.Println("Do a list request")
	cs := make(chan string)
	chanList <- cs
	for file := range cs {
		Info.Printf(" list request received %s", file)
	}

	Info.Print("End of list request")

	var z beeObj

	err = json.Unmarshal([]byte(marshal), &z)

	fmt.Println(z)

}
