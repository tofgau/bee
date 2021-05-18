package main

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

//
// registering of httpx handler
//

func (C Config) setupHttpHandlersx(myRouter *mux.Router) {
	myRouter.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { C.replyPing(w, r) })
	myRouter.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) { C.update(w, r) })
	//myRouter.HandleFunc("/updatexmlrpc", func(w http.ResponseWriter, r *http.Request) { C.replyPing(w, r) })

}

//
// replyPing handler
//
func (C Config) replyPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bee is alive")
	return
}

//
// update handler.
// This is the main bee handler. POST OR GET. Other method returns an http error
//
func (C Config) update(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Parses the request body

	//Extract the UID and, If the config flag is set, alter the case
	var uid string
	var err error
	if uid, err = GetParameterAsString(r, "UID"); err != nil {
		Error.Printf(" Error :  Unable to find UID  : %s\n", err)
		returnCode500(w, r, "[ERROR:HANDLER::ERRSEARCHINGUID] Error while searching UID")
		return
	}
	if uid == "" {
		Info.Printf(" Error :  UID is not given. Stopping %s", r)
		returnCode400(w, r, "[ERROR:HANDLER::UIDNOTGIVEN] UID is not GiVEN")
		return
	}
	fmt.Println(C.RewriteUIDCase)
	switch C.RewriteUIDCase {
	case "ToUpper":
		uid = strings.ToUpper(uid)

	case "ToLower":
		uid = strings.ToLower(uid)
	}

	//Load the object or create a new one

	var beeObjReceived *beeObj

	ce := make(chan error)
	cr := make(chan *beeObj)
	chanLoad <- loadQuery{uid, cr, ce}
	select {
	case err := <-ce:
		Error.Printf("Error loading : %s", err)
		returnCode500(w, r, "[ERROR:HANDLER::LOADOBJ] Error loading bee object")
		return

	case beeObjReceived = <-cr:
		//Info.Printf("Obj loaded :  %s", beeObjReceived.String())
		break
	}

	//Extract genral bee data field

	beeObjReceived.XsourceIP = GetIP(r)
	beeObjReceived.XlastUpdateTime = time.Now()

	reflexT := reflect.TypeOf(beeObjReceived).Elem()
	reflexV := reflect.ValueOf(beeObjReceived).Elem()

	Info.Printf("update received from %s", beeObjReceived.XsourceIP)

	//Validate AuthKey
	switch beeObjReceived.AuthKey {
	// AuthKey is empty -> we access it. Perhaps it will be populated later
	case "":
		Trace.Printf("AuthKey is empty for %s. Passing", beeObjReceived.UID)
		break
		//AuthKey is not empty, then we control and perhaps we rise an error
	default:
		if v, _ := GetParameterAsString(r, "AuthKey"); v != beeObjReceived.AuthKey {
			Error.Printf("Authkey error for : %s", beeObjReceived.UID)
			returnCode401(w, r, "401 [ERROR:HANDLER::AUTHKEY] Authkey Mismatch ")
			return
		}
	}

	// curl --insecure -X POST -d 'UID=6635'  https://192.168.9.114:444/update
	// curl -X POST -d 'UID=6635'  http://192.168.9.114:81/update

	//
	// For all non private beeObj Field, check theirs  presence in the GET then POST parameters
	//

	for i := 0; i < reflexT.NumField(); i++ {
		if _, ok := reflexT.Field(i).Tag.Lookup("Synced"); !ok {
			continue
		} else {

			if v, err := GetParameterAsString(r, reflexT.Field(i).Name); err != nil {
				Trace.Printf(" 400  wrong HTPP METHOD : %s\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				// If the parameter exist but is empty
				if v == "" {
					continue
				}

				Trace.Printf("Searching field %s = [%s]", reflexT.Field(i).Name, v)
				//fmt.Printf("%s ", reflexV.Field(i).Type().String())
				switch typex := reflexV.Field(i).Kind(); typex {
				//
				case reflect.String:
					Trace.Printf("Assign field type string  %s = [%s]", reflexT.Field(i).Name, v)
					reflexV.Field(i).SetString(v)
				//
				case reflect.Int64:
					n, err := strconv.ParseInt(v, 10, 64)
					if err != nil {

						text := fmt.Sprintf(" 400 [ERROR:HANDLER::CONVERT] wrong unable to convert %s value %s\n", reflexT.Field(i).Name, v)
						Error.Printf(text)
						returnCode400(w, r, text)
						return
					}

					Trace.Printf("Assign field type int64 %s = [%s]", reflexT.Field(i).Name, n)
					reflexV.Field(i).SetInt(n)
				//
				case reflect.Int32:
					n, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						text := fmt.Sprintf(" 400 [ERROR:HANDLER::CONVERT] wrong unable to convert %s value %s\n", reflexT.Field(i).Name, v)
						Error.Printf(text)
						returnCode400(w, r, text)

						return
					}

					Trace.Printf("Assign field type int64 %s = [%s]", reflexT.Field(i).Name, n)
					reflexV.Field(i).SetInt(n)
				//
				default:
					text := fmt.Sprintf(" 500 [ERROR:HANDLER::CONVERT] DataType not allowed %s : %s", reflexT.Field(i).Name, reflexV.Field(i).Kind())
					Error.Printf(text)
					returnCode500(w, r, text)
					return
				}
			}
		}
	}
	//Save

	{
		ce := make(chan error)
		chanStore <- storeQuery{beeObjReceived, ce}

		select {
		case err := <-ce:
			if err != nil {
				text := fmt.Sprintf(" 500 [ERROR:HANDLER::STORE] Error while storing beeObj %s", err)
				returnCode500(w, r, text)
				return
			}
		}
	}

	returnCode200(w, r, beeObjReceived.String())
	return
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

//
//This function retrieve a tcp parameter from POST or GET
//
func GetParameterAsString(r *http.Request, name string) (string, error) {
	//Trace.Printf("searching %s", name)
	switch r.Method {
	case "GET":
		x := r.URL.Query().Get(name)
		//Trace.Printf("returned for GET %s", x)
		return x, nil

	case "POST":
		err := r.ParseForm()
		if err != nil {
			return "", err
		}

		x := r.Form.Get(name)
		//Trace.Printf("returned for POST  %s", x)
		return x, nil

	}
	return "", errors.New("Sorry, only GET and POST methods are supported")
}
