package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

//
// TO DO : should lock a record after a read (=checkout) and un lock it (=checkin) after a write.
//    NOt a major issue because should not have concurrent issue except with GUI
//

//
// Take a bee Object and return the marshalled JSON string
//
func (beeObj beeObj) marshal() ([]byte, error) {
	var jsonData []byte
	jsonData, err := json.MarshalIndent(beeObj, "", " ")
	if err != nil {
		return jsonData, err
	}
	fmt.Println("jsonData")
	return jsonData, nil
}

//
// Take a bee []byte and return beeObj, err]
//
func unMarshal(in []byte) (beeObj, error) {
	var z beeObj
	//err=json.Unmarshal([]byte(marshal), &z)
	return z, nil
}

type loadQuery struct {
	UID     string
	chanRet chan *beeObj
	err     chan error
}

type storeQuery struct {
	beeObj *beeObj
	err    chan error
}

func (C Config) setupStorage() (chanLoad chan loadQuery, chanStore chan storeQuery, chanList chan chan string, done chan interface{}) {
	chanLoad = make(chan loadQuery)
	chanStore = make(chan storeQuery)
	chanList = make(chan chan string)
	done = make(chan interface{})
	fmt.Println(C)
	go func(C Config) {
		Trace.Println("Starting abstraction storage")
		for {
			select {
			//
			// Load
			//
			case loadQuery := <-chanLoad:
				Trace.Printf("Store:: Load of %s", loadQuery.UID)

				if _, err := os.Stat(C.BeeObjectPath + "/" + loadQuery.UID); os.IsNotExist(err) {
					Info.Printf("Store:: %s do not exist. Return new beeObj", loadQuery.UID)
					ret := beeObj{UID: loadQuery.UID, XcreationTime: time.Now(), XS_lastState: "DS", S_currentState: "DS"}
					ret.AddNotification("this job is born in the Hive")
					loadQuery.chanRet <- &ret
					break
				}

				jsonFile, err := os.Open(C.BeeObjectPath + "/" + loadQuery.UID)

				if err != nil {
					loadQuery.err <- errors.New(fmt.Sprintf("Store::Error opening  %s : %s", loadQuery.UID, err))
					break
				}
				byteValue, _ := ioutil.ReadAll(jsonFile)
				jsonFile.Close()

				var beeObjRet = beeObj{}
				err = json.Unmarshal(byteValue, &beeObjRet)
				if err != nil {
					loadQuery.err <- errors.New(fmt.Sprintf("Store::Error unmarshalling %s : %s", loadQuery.UID, err))
					break

				}
				loadQuery.chanRet <- &beeObjRet
				break

				//
				//Store
				//
			case storeQuery := <-chanStore:
				Trace.Printf("Store:: request storage of %s", storeQuery.beeObj)
				var marshall, err = storeQuery.beeObj.marshal()
				if err != nil {
					storeQuery.err <- errors.New(fmt.Sprintf("Store::Error marshalling %s : %s", storeQuery.beeObj.UID, err))
					break
				}

				err = ioutil.WriteFile(C.BeeObjectPath+"/"+storeQuery.beeObj.UID, marshall, 644)
				if err != nil {
					storeQuery.err <- errors.New(fmt.Sprintf("Store::Error writing beeobj %s disk full or unauthorized char (%s)", storeQuery.beeObj, err))
					break
				}

				close(storeQuery.err)
				break

				//
				//List
				//

			case chanPathString := <-chanList:
				go func(C Config, chanPathString chan<- string) {
					files, err := ioutil.ReadDir(C.BeeObjectPath)
					if err != nil {
						Error.Println(err)
					}
					for _, file := range files {
						chanPathString <- file.Name()
					}

					close(chanPathString)
					return
				}(C, chanPathString)
			case q := <-done:
				_ = q
				Trace.Println("Terminating abstraction storage")
				return
			}
		}
	}(C)

	return chanLoad, chanStore, chanList, done

}
