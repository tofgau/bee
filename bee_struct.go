//
// bee_Struct - propbably the most important part of bee Hive
//
// Bee represent job by a struct savec on disk as JSON.
// This struct is defined here.
// From it it, and from its TAGs, comes automatically storage and sync with http update datas.
//
//
package main

import (
	"fmt"
	"reflect"
	"time"
)

// validState is a public string array that define all possible value (bee job state) for the field   S_currentState
var validState = [...]string{"OK", "WARNING", "CRITICAL", "TIMEOUT", "DISABLED"}

//type BeeNotif is a timetag and a string.
//This are stored in the array XbeeEngineNotification
//The idea is to be able to store some event happened to the bee job
// use *beeObj.AddNotification to add entries
type BeeNotif struct {
	Time time.Time `json:"Time"`
	Text string    `json:"text"`
}

//BeeNotif.String()
func (B BeeNotif) String() string {
	return fmt.Sprintf("\n[%s %s ]", B.Time, B.Text)
}

//
// type beeObj
// This is bee task object representation.  Synced fields are retrieved from GET/POST http(s) update
//    Onbly  string, int32 and int64 are valid types
//

type beeObj struct {
	UID                    string       `json:"UID" ` // must be private because it is the first field read at the beginning and it has a case rewrite
	AuthKey                string       `json:"AuthKey"  Synced:"1"`
	Heartbeat              int32        `json:"Heartbeat"  Synced:"1"`
	Zombeat                int32        `json:"Zombeat"  Synced:"1"`
	RunningLocation        string       `json:"RunningLocation"  Synced:"1"`
	RunningPath            string       `json:"RunningPath"  Synced:"1"`
	XsourceIP              string       `json:"SourceIP" `
	XcreationTime          time.Time    `json:"CreationTime" `
	XlastUpdateTime        time.Time    `json:"LastupdateTime" `
	XlastStateChangeTime   time.Time    `json:"LastupdateTime" `
	XbeeEngineNotification [10]BeeNotif `json:"S_beeEngineNotification" `

	S_currentState         string `json:"SI_currentState"  Synced:"1"`
	S_currentStateText     string `json:"SI_currentStateText"  Synced:"1"`
	XS_lastState           string `json:"S_lastState" `
	XS_lastStateChangeTime string `json:"S_lastStateChangeTime" `

	XS_eventSend int `json:"S_eventSend" `

	XS_inWAcount   int32 `json:"S_inWAcount" `
	XS_inCRcount   int32 `json:"S_inCRcount" `
	XS_inTOcount   int32 `json:"S_inTOcount" `
	XS_inOKcount   int32 `json:"S_inOKcount" `
	S_WAgraceCount int32 `json:"SI_WAgraceCount"  Synced:"1"`
	S_CRgraceCount int32 `json:"SI_CRgraceCount"  Synced:"1"`
	S_OKgraceCount int32 `json:"SI_OKgraceCount"  Synced:"1"`

	XS_transitionCount     int32 `json:"S_transitionCount" `
	S_FLlockDuration       int32 `json:"SI_FLlockDuration"  Synced:"1"`
	S_transitionThreshold  int32 `json:"SI_transitionThreshold"  Synced:"1"`
	S_transitionDecaySpeed int32 `json:"SI_transitionDecaySpeed"  Synced:"1"`

	DOC_teamID      string `json:"DOC_teamID"  Synced:"1"`
	DOC_cmdbService string `json:"DOC_cmdbService"  Synced:"1"`
	DOC_revision    string `json:"DOC_revision"  Synced:"1"`

	DOC_inlineText string `json:"DOC_inlineText"  Synced:"1"`
	DOC_URI        string `json:"DOC_URI"  Synced:"1"`
	DOC_URL        string `json:"DOC_URL"  Synced:"1"`
}

//*beeObj.String()
//  reflexion is used to present a comprehensioble view of bee job!
func (beeObj *beeObj) String() string {
	//*beeObj.String()
	//  reflexion is used to present a comprehensioble view of bee job!

	v := reflect.ValueOf(*beeObj)
	ret := ""
	ret = ret + "\n[BEE::\n"
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Name == "AuthKey" {
			continue
		}
		ret = ret + fmt.Sprintf("%22s", v.Type().Field(i).Name)
		ret = ret + ":"
		ret = ret + fmt.Sprint(v.Field(i).Interface())
		ret = ret + "\n"
	}
	ret = ret + "\n]\n"
	return ret

}

//*beeJob.AddNotification(text string) (*beeJob)
//    This is used to add a notification
//    New notification are always in position 0 for simplicity (then 0 become 1, ...)
func (beeObj *beeObj) AddNotification(text string) (ret_beeObj *beeObj) {

	// I prefer to move all the data in place of playing with start index. more easy to read
	for i := len(beeObj.XbeeEngineNotification) - 2; i >= 0; i-- {
		beeObj.XbeeEngineNotification[i+1] = beeObj.XbeeEngineNotification[i]
	}

	beeObj.XbeeEngineNotification[0] = BeeNotif{Time: time.Now(), Text: text}

	return beeObj

}

//Errors
/*
type errorInvalidState struct{}

func (m *errorInvalidState) Error() string {
	return "You tried to configure an invalid bee state"
}
//*beeObj.Validate()
// *Not Used*  This functions performs Validation on a bee
func (beeObj *beeObj) Validate() (*beeObj, error) {
	// check that the state of the object is well the one defined
	for _, b := range validState {
		if b == beeObj.S_currentState {
			return beeObj, nil
		}
	}
	return beeObj, &errorInvalidState{}
}
*/
