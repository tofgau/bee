package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const tickerPIDFileRewriteMilliSec = 1000
const maxAgePIDFileMilliSec = 2000

func (C Config) startWatchdog() chan interface{} {
	mainDone := make(chan interface{})

	fileInfo, err := os.Stat(C.PIDFile)
	if err != nil || (time.Now().Sub(fileInfo.ModTime()) > time.Millisecond*maxAgePIDFileMilliSec) {

		go func(mainDone chan interface{}) {
			ticker := time.NewTicker(tickerPIDFileRewriteMilliSec * time.Millisecond)
			firstLoop := true
			for {
				select {
				case <-mainDone:
					Trace.Printf("PID::Deleting PID file and stopping observer")
					err := os.Remove(C.PIDFile)
					if err != nil {
						Error.Println(err)
					}
					return
				case _ = <-ticker.C:
					if !firstLoop {

						_, err := os.Stat(C.PIDFile)

						if err != nil {
							Info.Printf("PID::pid file removed : requesting to terminate now")

							close(mainDone)
							return
						}
					}
					firstLoop = false

					pid := os.Getpid()
					d1 := []byte(strconv.Itoa(pid))
					err := ioutil.WriteFile(C.PIDFile, d1, 0644)
					if err != nil {
						Error.Printf("Unable to write PID file : %s, %s", C.PIDFile, err)
					}

				}

			}
		}(mainDone)

		// catch CTRL C
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func(mainDone chan interface{}) {
			for _ = range c {
				Info.Println("*** Catching SIGTERM - terminate signal has been laucnhed ***")
				close(mainDone)
				return
			}
		}(mainDone)

		return mainDone
	}
	return nil // if we are already running
}
