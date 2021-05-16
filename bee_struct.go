package main

import (
	"fmt"
	"reflect"
	"time"
)

// List of valid state for bee ojects
var validState = [...]string{"OK", "WARNING", "CRITICAL", "TIMEOUT", "DISABLED"}

// Here is defined what ixs a bee entry.

type BeeNotif struct {
	Time time.Time `json:"Time"`
	Text string    `json:"text"`
}

func (B BeeNotif) String() string {
	return fmt.Sprintf("\n[%s %s ]", B.Time, B.Text)
}

type beeObj struct {
	UID                    string       `json:"UID"`
	AuthKey                string       `json:"AuthKey"`
	Heartbeat              int32        `json:"Heartbeat"`
	Zombeat                int32        `json:"Zombeat"`
	RunningLocation        string       `json:"RunningLocation"`
	RunningPath            string       `json:"RunningPath"`
	XsourceIP              string       `json:"SourceIP"`
	XcreationTime          time.Time    `json:"CreationTime" private:"1"`
	XlastUpdateTime        time.Time    `json:"LastupdateTime" private:1`
	XlastStateChangeTime   time.Time    `json:"LastupdateTime" private:1`
	XbeeEngineNotification [10]BeeNotif `json:"S_beeEngineNotification" private:1`

	S_currentState         string `json:"SI_currentState"`
	S_currentStateText     string `json:"SI_currentStateText"`
	XS_lastState           string `json:"S_lastState" private:1`
	XS_lastStateChangeTime string `json:"S_lastStateChangeTime" private:1`

	XS_eventSend int `json:"S_eventSend" private:1`

	XS_inWAcount   int32 `json:"S_inWAcount" private:1`
	XS_inCRcount   int32 `json:"S_inCRcount" private:1`
	XS_inTOcount   int32 `json:"S_inTOcount" private:1`
	XS_inOKcount   int32 `json:"S_inOKcount" private:1`
	S_WAgraceCount int32 `json:"SI_WAgraceCount"`
	S_CRgraceCount int32 `json:"SI_CRgraceCount"`
	S_OKgraceCount int32 `json:"SI_OKgraceCount"`

	XS_transitionCount     int32 `json:"S_transitionCount" private:1`
	S_FLlockDuration       int32 `json:"SI_FLlockDuration"`
	S_transitionThreshold  int32 `json:"SI_transitionThreshold"`
	S_transitionDecaySpeed int32 `json:"SI_transitionDecaySpeed"`

	DOC_teamID      string `json:"DOC_teamID"`
	DOC_cmdbService string `json:"DOC_cmdbService"`
	DOC_revision    string `json:"DOC_revision"`

	DOC_inlineText string `json:"DOC_inlineText"`
	DOC_URI        string `json:"DOC_URI"`
	DOC_URL        string `json:"DOC_URL"`
}

//Errors

type errorInvalidState struct{}

func (m *errorInvalidState) Error() string {
	return "You tried to configure an invalid bee state"
}

func (beeObj beeObj) validate() (beeObj, error) {
	// check that the state of the object is well the one defined
	for _, b := range validState {
		if b == beeObj.S_currentState {
			return beeObj, nil
		}
	}
	return beeObj, &errorInvalidState{}
}

func (beeObj beeObj) String() string {

	v := reflect.ValueOf(beeObj)
	ret := ""
	ret = ret + "\n[BEE::\n"
	for i := 0; i < v.NumField(); i++ {
		ret = ret + fmt.Sprintf("%22s", v.Type().Field(i).Name)
		ret = ret + ":"
		ret = ret + fmt.Sprint(v.Field(i).Interface())
		ret = ret + "\n"
	}
	ret = ret + "\n]\n"
	return ret

}
func (beeObj *beeObj) AddNotification(text string) (ret_beeObj *beeObj) {

	// I prefer to move all the data in place of playing with start index. more easy to read
	for i := len(beeObj.XbeeEngineNotification) - 2; i >= 0; i-- {
		beeObj.XbeeEngineNotification[i+1] = beeObj.XbeeEngineNotification[i]
	}

	beeObj.XbeeEngineNotification[0] = BeeNotif{Time: time.Now(), Text: text}

	return beeObj

}
