package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"evolvingPhilosophers.local/globalData"
	"evolvingPhilosophers.local/messageServerStack"
	ringBuffer "evolvingPhilosophers.local/ringbuffer"
)

// *
// StateFromAdjacentDp
// *
func (f *Facilitator) StateFromAdjacentDp(w http.ResponseWriter, r *http.Request) {
	response := RequestData{State: f.DpStates.State}

	jsonData, err := json.Marshal(response)
	if err != nil {
		f.LogMessage("3", "StateFromAdjacentDp, Marshall error", err)
		jsonData = []byte(`{"state": "0"}`)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// *
// StateFromRightDp and StateFromLeftDp need to be replaced with
// *
func (f *Facilitator) StateFromRightDp(w http.ResponseWriter, r *http.Request) {
	response := RequestData{State: f.DpStates.State}

	jsonData, err := json.Marshal(response)
	if err != nil {
		f.LogMessage("3", "StateFromRightDp, Marshall error", err)
		jsonData = []byte(`{"state": "0"}`)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	w.Write(jsonData)
}

// *
// GetStateRightAndReceiveStateLeft
// *
func (f *Facilitator) GetStateRightAndReceiveStateLeft(w http.ResponseWriter, r *http.Request) {

	// In state Quiesce2, the dp does not try to change state to dining
	// When the dp on the left requests this dp's state, it is not dining, so the dp on the left
	// may be able to change state to dining
	if f.DpStates.State == Quiesce2 {
		f.LogMessage("1", "GetStateRightAndReceiveStateLeft,f.dpStates.state == Quiesce2", nil)
		w.WriteHeader(http.StatusOK)
		return
	}

	// get right state from query
	leftState := r.URL.Query().Get("state")

	var err error
	f.DpStates.LeftState, err = strconv.Atoi(leftState)
	if err != nil {
		f.LogMessage("1", "GetStateRightAndReceiveStateLeft,could not convert query parameter state to int type", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var stateRight int
	stateRight, err = f.GetStateRight()
	if err != nil {
		f.LogMessage("3", "GetStateRightAndReceiveStateLeft, failure of getStateRight", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f.DpStates.RightState = stateRight

	f.TestIfAbleToDine()
	w.WriteHeader(http.StatusOK)
}

// *
// GetStateLeftAndReceiveStateRight
// *
func (f *Facilitator) GetStateLeftAndReceiveStateRight(w http.ResponseWriter, r *http.Request) {

	// In state Quiesce2, the dp does not try to change state to dining
	// When the dp on the right requests this dp's state, it is not dining, so the dp on the right
	// may be able to change state to dining
	if f.DpStates.State == Quiesce2 {
		f.LogMessage("1", "GetStateLeftAndReceiveStateRight,f.dpStates.state == Quiesce2", nil)
		w.WriteHeader(http.StatusOK)
		return
	}

	// get left state from query
	rightState := r.URL.Query().Get("state")
	var err error
	f.DpStates.RightState, err = strconv.Atoi(rightState)
	if err != nil {
		f.LogMessage("3", "GetStateLeftAndReceiveStateRight, could not convert query parameter state to int type", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// request state of right dp
	var stateLeft int
	stateLeft, err = f.GetStateLeft()
	if err != nil {
		f.LogMessage("3", "GetStateLeftAndReceiveStateRight, failure of getStateLeft", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f.DpStates.LeftState = stateLeft

	f.TestIfAbleToDine()
	w.WriteHeader(http.StatusOK)
}

// *
// Handler quiesce1Left
// *
func (f *Facilitator) Quiesce1Left(w http.ResponseWriter, r *http.Request) {
	f.Quiesce1Chan <- Quiesce1

	leftRequest := fmt.Sprintf("http://%s/quiesce2", f.LeftAddress)

	response, err := http.Get(leftRequest)
	if err != nil {
		text := fmt.Sprintf("Quiesce1Left to the left, Error sending request, http.Get(leftRequest), leftRequest: %s", leftRequest)
		f.LogMessage("3", text, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer response.Body.Close()

	var quiesce2Response QuiesceResponse

	err = json.NewDecoder(response.Body).Decode(&quiesce2Response)

	quiesce1Response := QuiesceResponse{
		Address: f.OwnAddress,
		State:   "Quiesce1",
	}

	quiesce1LeftResponse := make(map[string]QuiesceResponse)
	quiesce1LeftResponse["target"] = quiesce1Response
	quiesce1LeftResponse["leftOfTarget"] = quiesce2Response

	jsonData, err := json.Marshal(quiesce1LeftResponse)
	if err != nil {
		f.LogMessage("3", "Quiesce1Left, Error json.Marshal(quiesce1LeftResponse)", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	text := fmt.Sprintf("Quiesce1Left, address: %s, left: %s, right: %s", f.OwnAddress, f.LeftAddress, f.RightAddress)
	f.LogMessage("1", text, nil)
	w.Write(jsonData)
}

// *
// Handler Quiesce1Right
// *
func (f *Facilitator) Quiesce1Right(w http.ResponseWriter, r *http.Request) {
	f.Quiesce1Chan <- Quiesce1

	rightRequest := fmt.Sprintf("http://%s/quiesce2", f.RightAddress)

	response, err := http.Get(rightRequest)
	if err != nil {
		text := fmt.Sprintf("Quiesce1Right, Error sending request, http.Get(rightRequest), rightRequest: %s", rightRequest)
		f.LogMessage("3", text, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer response.Body.Close()

	var quiesce2Response QuiesceResponse

	err = json.NewDecoder(response.Body).Decode(&quiesce2Response)

	quiesce1Response := QuiesceResponse{
		Address: f.OwnAddress,
		State:   "Quiesce1",
	}

	quiesce1RightResponse := make(map[string]QuiesceResponse)
	quiesce1RightResponse["target"] = quiesce1Response
	quiesce1RightResponse["rightOfTarget"] = quiesce2Response

	jsonData, err := json.Marshal(quiesce1RightResponse)
	if err != nil {
		f.LogMessage("3", "Quiesce1Right, Error json.Marshal(quiesce1RightResponse)", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	text := fmt.Sprintf("Quiesce1Left, address: %s, left: %s, right: %s", f.OwnAddress, f.LeftAddress, f.RightAddress)
	f.LogMessage("1", text, nil)
	w.Write(jsonData)
}

// *
// Handler Quiesce2
// *
func (f *Facilitator) Quiesce2(w http.ResponseWriter, r *http.Request) {
	f.Quiesce2Chan <- Quiesce2

	text := fmt.Sprintf("Quiesce2, address: %s, left: %s, right: %s", f.OwnAddress, f.LeftAddress, f.RightAddress)
	f.LogMessage("1", text, nil)

	// why not a channel to receive the state of this dp, to make sure it has changed to quiesce2
	quiesce2Response := QuiesceResponse{
		Address: f.OwnAddress,
		State:   "Quiesce2",
	}

	jsonData, err := json.Marshal(quiesce2Response)
	if err != nil {
		f.LogMessage("3", "Quiesce2, failure json.Marshal(quiesce2Response)", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	w.Write(jsonData)
}

// *
// Handler DeQuiesce1Left
// This request is received by the dp on the left of the new or removed dp
// *
func (f *Facilitator) DeQuiesce1Left(w http.ResponseWriter, r *http.Request) {
	leftRequestUrl := fmt.Sprintf("http://%s/deQuiesce2", f.LeftAddress)

	response, err := http.Get(leftRequestUrl)
	if err != nil {
		text := fmt.Sprintf("DeQuiesce1Left to the left, Error sending request, http.Get(leftRequest), leftRequest: %s", leftRequestUrl)
		f.LogMessage("3", text, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	f.ExitQuiesce1Chan <- Thinking

	var responseMap map[string]string

	err = json.NewDecoder(response.Body).Decode(&responseMap)
	if err != nil {
		f.LogMessage("3", "DeQuiesce1Left, failed to decode json response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	text := fmt.Sprintf("DeQuiesce1Left, address: %s, left: %s, right: %s", f.OwnAddress, f.LeftAddress, f.RightAddress)
	f.LogMessage("1", text, nil)

	responseMap["Address"] = f.OwnAddress
	responseMap["NewState"] = "Thinking"
	responseMap["Left"] = f.LeftAddress
	responseMap["Right"] = f.RightAddress

	jsonData, err := json.Marshal(responseMap)
	if err != nil {
		f.LogMessage("3", "Error json.Marshal(DeQuiesce1Left)", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// *
// Handler DeQuiesce1Right
// This request is received by the DP on the right of the added or removed dp
// *
func (f *Facilitator) DeQuiesce1Right(w http.ResponseWriter, r *http.Request) {
	rightRequestUrl := fmt.Sprintf("http://%s/deQuiesce2", f.RightAddress)

	response, err := http.Get(rightRequestUrl)
	if err != nil {
		text := fmt.Sprintf("DeQuiesce1Right to the right, Error sending request, http.Get(rightRequest), rightRequest: %s", rightRequestUrl)
		f.LogMessage("3", text, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	f.ExitQuiesce1Chan <- Thinking

	text := fmt.Sprintf("DeQuiesce1Right, address: %s, left: %s, right: %s", f.OwnAddress, f.LeftAddress, f.RightAddress)
	f.LogMessage("1", text, nil)

	var responseMap map[string]string

	err = json.NewDecoder(response.Body).Decode(&responseMap)
	if err != nil {
		f.LogMessage("3", "DeQuiesce1Right, failed to decode json response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseMap["Address"] = f.OwnAddress
	responseMap["NewState"] = "Thinking"
	responseMap["Left"] = f.LeftAddress
	responseMap["Right"] = f.RightAddress

	jsonData, err := json.Marshal(responseMap)
	if err != nil {
		f.LogMessage("3", "DeQuiesce1Right, Error json.Marshal(responseMap)", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// *
// Handler DeQuiesce2
// *
func (f *Facilitator) DeQuiesce2(w http.ResponseWriter, r *http.Request) {
	f.ExitQuiesce2Chan <- Thinking

	text := fmt.Sprintf("DeQuiesce2, address: %s, left: %s, right: %s", f.OwnAddress, f.LeftAddress, f.RightAddress)
	f.LogMessage("1", text, nil)

	responseMap := make(map[string]string)

	responseMap["Address"] = f.OwnAddress
	responseMap["NewState"] = "Thinking"
	responseMap["Left"] = f.LeftAddress
	responseMap["Right"] = f.RightAddress

	jsonData, err := json.Marshal(responseMap)
	if err != nil {
		f.LogMessage("3", "DeQuiesce2, Marshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// *
// The new dp is started with state Quiesce3
// Handler DeQuiesce3
// *
func (f *Facilitator) DeQuiesce3(w http.ResponseWriter, r *http.Request) {
	f.ExitQuiesce3Chan <- Thinking

	text := fmt.Sprintf("DeQuiesce3, address: %s, left: %s, right: %s", f.OwnAddress, f.LeftAddress, f.RightAddress)
	f.LogMessage("1", text, nil)

	responseMap := make(map[string]string)

	responseMap["Address"] = f.OwnAddress
	responseMap["NewState"] = "Thinking"
	responseMap["Left"] = f.LeftAddress
	responseMap["Right"] = f.RightAddress

	jsonData, err := json.Marshal(responseMap)
	if err != nil {
		f.LogMessage("3", "DeQuiesce3, failure to Marshal(quiesce2Response)", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// *
// DirectAddDp
// *
func (f *Facilitator) DirectAddDp(w http.ResponseWriter, r *http.Request) {
	responseMap := make(map[string]string)

	var newPhilosopher NewPhilosopher
	err := json.NewDecoder(r.Body).Decode(&newPhilosopher)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseMap["decode request body"] = fmt.Sprintf("%v", err)
		jsonData, err := json.Marshal(responseMap)
		if err != nil {
			f.LogMessage("3", "DirectAddDp, failure Marshal", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}

	var addressOfAdjacentPhilosopherBeforeNewPhilosopher string
	// for the philosopher that needs to change a philosopher address besides the one doing the work
	// for adding, this is the side of that existing philosopher on which the new philosopher is
	var theNewPhilosopherIsOnThisSideOfTheAdjacentPhilosopher string

	if newPhilosopher.Side == "Right" {
		addressOfAdjacentPhilosopherBeforeNewPhilosopher = f.RightAddress
		f.RightAddress = newPhilosopher.Address
		theNewPhilosopherIsOnThisSideOfTheAdjacentPhilosopher = "Left"
	} else {
		addressOfAdjacentPhilosopherBeforeNewPhilosopher = f.LeftAddress
		f.LeftAddress = newPhilosopher.Address
		theNewPhilosopherIsOnThisSideOfTheAdjacentPhilosopher = "Right"
	}

	// tell the (left or right) adjacent partner to change the address it stores for the left or right dp
	requestUrlChangeOfAddress := fmt.Sprintf("http://%s/changeAddressOfAdjacentDp", addressOfAdjacentPhilosopherBeforeNewPhilosopher)

	requestDataChangeOfAddress := map[string]string{
		"side":    theNewPhilosopherIsOnThisSideOfTheAdjacentPhilosopher,
		"address": newPhilosopher.Address,
	}

	jsonDataHandlerChangeOfAddress, err := json.Marshal(requestDataChangeOfAddress)
	if err != nil {
		f.LogMessage("3", "DirectAddDp marshal request data failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		responseMap["Marshal request for ChangeAddressOfAdjacentDp"] = fmt.Sprintf("%v", err)
		jsonData, err := json.Marshal(responseMap)
		if err != nil {
			f.LogMessage("3", "DirectAddDp, failure Marshal", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}

	requestChangeOfAddress, err := http.NewRequest("POST", requestUrlChangeOfAddress, bytes.NewBuffer(jsonDataHandlerChangeOfAddress))
	if err != nil {
		f.LogMessage("3", "DirectAddDp failure to create request, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		responseMap["Create request for ChangeAddressOfAdjacentDp"] = fmt.Sprintf("%v", err)
		jsonData, err := json.Marshal(responseMap)
		if err != nil {
			f.LogMessage("3", "DirectAddDp, failure Marshal", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}

	requestChangeOfAddress.Header.Set("Content-Type", "application/json")

	// change the adjacent dp's address from the sender to the new dp
	client := &http.Client{Timeout: 120 * time.Second}
	responseChangeOfAddressRequest, err := client.Do(requestChangeOfAddress)
	if err != nil {
		// error may require rollback here
		f.LogMessage("3", "DirectAddDp failure to send requestChangeOfAddress", err)
		w.WriteHeader(http.StatusInternalServerError)
		responseMap["Send POST for ChangeAddressOfAdjacentDp"] = fmt.Sprintf("%v", err)
		jsonData, err := json.Marshal(responseMap)
		if err != nil {
			f.LogMessage("3", "DirectAddDp, failure Marshal", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}

	defer responseChangeOfAddressRequest.Body.Close()

	err = json.NewDecoder(responseChangeOfAddressRequest.Body).Decode(&responseMap)
	if err != nil {
		f.LogMessage("3", "DirectAddDp NewDecoder Decode failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseMap["addressOfAdjacentPhilosopherBeforeNewPhilosopher"] = addressOfAdjacentPhilosopherBeforeNewPhilosopher
	responseMap["theNewPhilosopherIsOnThisSideOfTheAdjacentPhilosopher"] = theNewPhilosopherIsOnThisSideOfTheAdjacentPhilosopher
	jsonDataAddDp, err := json.Marshal(responseMap)
	if err != nil {
		f.LogMessage("3", "DirectAddDp, Marshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonDataAddDp)
}

// *
// ReceiveRequestToRemoveDp
// *
func (f *Facilitator) ReceiveRequestToRemoveDp(w http.ResponseWriter, r *http.Request) {

	var informationToRemoveDp InformationToRemoveDp

	err := json.NewDecoder(r.Body).Decode(&informationToRemoveDp)
	if err != nil {
		f.LogMessage("3", "ReceiveRequestToRemoveDp, failure NewDecoder, Decode", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	terminateSelf := false

	if informationToRemoveDp.DpNumberToBeRemoved == f.DpNumber {
		// add new dp
		err := f.removeDp()
		if err != nil {
			f.LogMessage("3", "ReceiveRequestToRemoveDp, failure to remove Dp", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if informationToRemoveDp.TerminateSelf == "true" {
			terminateSelf = true
		}
		informationToRemoveDp.Done = "true"
		informationToRemoveDp.Result = "Dp that received the request is the Dp deleted"
		f.ReceiveRequestToRemoveDPChannel <- informationToRemoveDp
	} else {
		informationToRemoveDp.FromRequestHander = "true"
		informationToRemoveDp.Done = "false"
		err := f.relayRequestToRemoveDp(&informationToRemoveDp)
		if err != nil {
			f.LogMessage("3", "ReceiveRequestToRemoveDp, f.relayRequestToRemoveDp failure ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	select {
	case informationToRemoveDp = <-f.ReceiveRequestToRemoveDPChannel:
	case <-time.After(360 * time.Second):
		informationToRemoveDp.Result = "ReceiveRequestToRemoveDp did not receive a reply, timed out after 360 seconds"
	}

	jsonData, err := json.Marshal(informationToRemoveDp)
	if err != nil {
		f.LogMessage("3", "ReceiveRequestToRemoveDp, Marshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if terminateSelf == true {
		go func() {
			time.Sleep(time.Second * 30)
			os.Exit(0)
		}()
	}

	w.Write(jsonData)
}

// *
// RelayRemoveDp
// *
func (f *Facilitator) RelayRemoveDp(w http.ResponseWriter, r *http.Request) {

	var informationToRemoveDp InformationToRemoveDp

	terminateSelf := false

	err := json.NewDecoder(r.Body).Decode(&informationToRemoveDp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	informationToRemoveDp.OriginalLeftAddress = f.LeftAddress

	if informationToRemoveDp.Done == "true" {
		err = f.relayRequestToRemoveDp(&informationToRemoveDp)
		if err != nil {
			f.LogMessage("3", "RelayRemoveDp, relayRequestToRemoveDp failure", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if informationToRemoveDp.DpNumberToBeRemoved == f.DpNumber {
		err = f.removeDp()

		informationToRemoveDp.Done = "true"
		if err != nil {
			informationToRemoveDp.Result = "failed to remove Dp"
			f.LogMessage("3", "RelayRemoveDp, relayRequestToRemoveDp failure, Done has been set to true", err)
		} else {
			informationToRemoveDp.Result = "Dp removed successfully"
			f.LogMessage("1", "RelayRemoveDp, removeDp success", nil)
			if informationToRemoveDp.TerminateSelf == "true" {
				terminateSelf = true
			}
		}

		if terminateSelf == true {
			f.LogMessage("1", "RelayRemoveDp, terminate self", nil)
			go func() {
				time.Sleep(time.Second * 30)
				os.Exit(0)
			}()
		}

		err = f.relayRequestToRemoveDp(&informationToRemoveDp)
		if err != nil {
			f.LogMessage("3", "RelayRemoveDp, relayRequestToRemoveDp failure", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		informationToRemoveDp.Done = "false"
		err = f.relayRequestToRemoveDp(&informationToRemoveDp)
		if err != nil {
			f.LogMessage("3", "RelayRemoveDp, relayRequestToRemoveDp failure", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// *
// DirectRemoveDp
// *
func (f *Facilitator) DirectRemoveDp(w http.ResponseWriter, r *http.Request) {

	var informationToDirectRemoveDp InformationToDirectRemoveDp

	err := json.NewDecoder(r.Body).Decode(&informationToDirectRemoveDp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	terminateSelf := informationToDirectRemoveDp.TerminateSelf

	// send quiesce1 to the dp on the left
	leftQuiesce1RequestUrl := fmt.Sprintf("http://%s/quiesce1Left", f.LeftAddress)

	_, err = http.Get(leftQuiesce1RequestUrl)
	if err != nil {
		text := fmt.Sprintf("DirectRemoveDp, Error sending request, http.Get(leftQuiesce1RequestUrl), leftQuiesce1RequestUrl: %s", leftQuiesce1RequestUrl)
		f.LogMessage("3", text, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// send quiesce1 to the dp on the right
	rightQuiesce1RequestUrl := fmt.Sprintf("http://%s/quiesce1Right", f.RightAddress)

	_, err = http.Get(rightQuiesce1RequestUrl)
	if err != nil {
		text := fmt.Sprintf("DirectRemoveDp, Error sending request, http.Get(rightQuiesce1RequestUrl), rightQuiesce1RequestUrl: %s", rightQuiesce1RequestUrl)
		f.LogMessage("3", text, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// change its own state to quiesce3
	f.Quiesce3Chan <- Quiesce3

	// change the right address of the dp on its left to its own right address
	leftChangeAddressRequestUrl := fmt.Sprintf("http://%s/changeAddressOfAdjacentDp", f.LeftAddress)

	leftChangeAddressRequestData := map[string]string{
		"side":    "Right",
		"address": f.RightAddress,
	}

	jsonData, err := json.Marshal(leftChangeAddressRequestData)
	if err != nil {
		f.LogMessage("3", "DirectRemoveDp marshal request data failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	leftChangeAddresstRequest, err := http.NewRequest("POST", leftChangeAddressRequestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		f.LogMessage("3", "DirectRemoveDp failure to create request", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	leftChangeAddresstRequest.Header.Set("Content-Type", "application/json")

	// change the adjacent dp's address from the sender to the new dp
	client := &http.Client{Timeout: 120 * time.Second}
	_, err = client.Do(leftChangeAddresstRequest)
	if err != nil {
		// error may require rollback here
		f.LogMessage("3", "DirectRemoveDp failure to send request", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// change the left address of the dp on its right to its own left address
	rightChangeAddresstRequestUrl := fmt.Sprintf("http://%s/changeAddressOfAdjacentDp", f.RightAddress)

	rightChangeAddresstRequestData := map[string]string{
		"side":    "Left",
		"address": f.LeftAddress,
	}

	jsonData, err = json.Marshal(rightChangeAddresstRequestData)
	if err != nil {
		f.LogMessage("3", "DirectRemoveDp marshal request data failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rightChangeAddresstRequest, err := http.NewRequest("POST", rightChangeAddresstRequestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		f.LogMessage("3", "DirectRemoveDp failure to create request", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rightChangeAddresstRequest.Header.Set("Content-Type", "application/json")

	// change the adjacent dp's address from the sender to the new dp
	client = &http.Client{Timeout: 120 * time.Second}
	_, err = client.Do(rightChangeAddresstRequest)
	if err != nil {
		// error may require rollback here
		f.LogMessage("3", "DirectRemoveDp failure to send request", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send dequiesce1 to the dp on its left
	leftDeQuiesce1RequestUrl := fmt.Sprintf("http://%s/deQuiesce1Left", f.LeftAddress)

	_, err = http.Get(leftDeQuiesce1RequestUrl)
	if err != nil {
		text := fmt.Sprintf("DirectRemoveDp, Error sending request, http.Get(leftDeQuiesce1RequestUrl), leftDeQuiesce1RequestUrl: %s", leftDeQuiesce1RequestUrl)
		f.LogMessage("3", text, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send dequiesce1 to the dp on its right
	rightDeQuiesce1RequestUrl := fmt.Sprintf("http://%s/deQuiesce1Right", f.RightAddress)

	_, err = http.Get(rightDeQuiesce1RequestUrl)
	if err != nil {
		text := fmt.Sprintf("DirectRemoveDp, Error sending request, http.Get(leftDeQuiesce1RequestUrl), leftDeQuiesce1RequestUrl: %s", leftDeQuiesce1RequestUrl)
		f.LogMessage("3", text, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f.RightAddress = ""
	f.LeftAddress = ""

	if terminateSelf == "true" {
		go func() {
			time.Sleep(time.Second * 30)
			os.Exit(0)
		}()
	}

	w.WriteHeader(http.StatusOK)
}

// *
// Handler SetLeftAndRightAddressOfNewDp
// *
func (f *Facilitator) SetLeftAndRightAddressOfNewDp(w http.ResponseWriter, r *http.Request) {

	var leftAndRightAddress NewDpLeftAndRightAddress

	err := json.NewDecoder(r.Body).Decode(&leftAndRightAddress)
	if err != nil {
		f.LogMessage("3", "SetLeftAndRightAddressOfNewDp, NewDecoder Decode failure", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	f.LeftAddress = leftAndRightAddress.LeftAddress
	f.RightAddress = leftAndRightAddress.RightAddress

	w.WriteHeader(http.StatusOK)
}

// *
// ReceiveRequestToAddNewDpToLeftOfTargetDp
// *
func (f *Facilitator) ReceiveRequestToAddNewDpToLeftOfTargetDp(w http.ResponseWriter, r *http.Request) {
	var informationToAddNewDp InformationToAddNewDp

	err := json.NewDecoder(r.Body).Decode(&informationToAddNewDp)
	if err != nil {
		f.LogMessage("3", "ReceiveRequestToAddNewDpToLeftOfTargetDp, NewDecoder Decode failure", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if informationToAddNewDp.TargetDpNumber == f.DpNumber {
		// add new dp
		err = f.addNewDpToLeftOfTargetDp(&informationToAddNewDp)
		if err != nil {
			f.LogMessage("3", "ReceiveRequestToAddNewDpToLeftOfTargetDp, failed", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		informationToAddNewDp.Done = "true"
		informationToAddNewDp.Result = "Added the new Dp next to the receiver of the request."
		f.ReceiveRequestToAddNewDPChannel <- informationToAddNewDp
	} else {
		informationToAddNewDp.Done = "false"
		informationToAddNewDp.FromRequestHander = "true"
		err := f.relayRequestToAddNewDpToLeftOfTargetDp(&informationToAddNewDp)
		if err != nil {
			f.LogMessage("3", "ReceiveRequestToAddNewDpToLeftOfTargetDp, failed", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	select {
	case informationToAddNewDp = <-f.ReceiveRequestToAddNewDPChannel:
	case <-time.After(60 * time.Second):
		informationToAddNewDp.Result = "ReceiveRequestToAddNewDpToLeftOfTargetDp did not receive a reply, timed out after 60 seconds"
	}

	jsonData, err := json.Marshal(informationToAddNewDp)
	if err != nil {
		f.LogMessage("3", "Error json.Marshal(informationToAddNewDp)", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	w.Write(jsonData)
}

// *
// AddNewDpToLeftOfTargetDp
// *
func (f *Facilitator) AddNewDpToLeftOfTargetDp(w http.ResponseWriter, r *http.Request) {
	var informationToAddNewDp InformationToAddNewDp

	err := json.NewDecoder(r.Body).Decode(&informationToAddNewDp)
	if err != nil {
		f.LogMessage("3", "AddNewDpToLeftOfTargetDp, NewDecoder Decode failure", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if informationToAddNewDp.Done == "true" {
		err = f.relayRequestToAddNewDpToLeftOfTargetDp(&informationToAddNewDp)
		if err != nil {
			f.LogMessage("3", "AddNewDpToLeftOfTargetDp, relayRequestToAddNewDpToLeftOfTargetDp failure", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if informationToAddNewDp.TargetDpNumber == f.DpNumber {
		// add new dp
		err := f.addNewDpToLeftOfTargetDp(&informationToAddNewDp)
		if err != nil {
			f.LogMessage("3", "AddNewDpToLeftOfTargetDp, addNewDpToLeftOfTargetDp failure", err)
			informationToAddNewDp.Result = "failure"
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		informationToAddNewDp.Done = "true"

		err = f.relayRequestToAddNewDpToLeftOfTargetDp(&informationToAddNewDp)
		if err != nil {
			f.LogMessage("3", "AddNewDpToLeftOfTargetDp, returnValueRelayRequest failure", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		informationToAddNewDp.Done = "false"
		err = f.relayRequestToAddNewDpToLeftOfTargetDp(&informationToAddNewDp)
		if err != nil {
			f.LogMessage("3", "AddNewDpToLeftOfTargetDp, relayRequestToAddNewDpToLeftOfTargetDp failure", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// *
// ChangeAddressOfAdjacentDp
// the dp receives this from the dp that received the request to do the actions needed to add a new dp
// to the circle.  it tells this dp to change the address it stores for the left or right dp
// *
func (f *Facilitator) ChangeAddressOfAdjacentDp(w http.ResponseWriter, r *http.Request) {
	responseMap := make(map[string]string)
	var newPhilosopher NewPhilosopher
	err := json.NewDecoder(r.Body).Decode(&newPhilosopher)
	if err != nil {
		f.LogMessage("3", "ChangeAddressOfAdjacentDp, NewDecoder Decode failuer", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newPhilosopher.Side == "Right" {
		f.RightAddress = newPhilosopher.Address
	} else {
		f.LeftAddress = newPhilosopher.Address
	}

	responseMap["newPhilosopher.Side"] = newPhilosopher.Side
	responseMap["newPhilosopher.Address"] = newPhilosopher.Address
	jsonData, err := json.Marshal(responseMap)
	if err != nil {
		f.LogMessage("3", "ChangeAddressOfAdjacentDp, Marshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// *
// DpAttributesReturn
// *
func (f *Facilitator) DpAttributesReturn(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		f.LogMessage("3", "DpAttributesReturn, failed to read body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var dpAttributesRelay DpAttributesRelay
	err = json.Unmarshal(body, &dpAttributesRelay)
	if err != nil {
		f.LogMessage("3", "DpAttributesReturn,failed to unmarshal body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if dpAttributesRelay.OriginationAddress == f.OwnAddress {
		f.DpAttributesResponseChannel <- dpAttributesRelay.DpAttributesMap
		return
	}

	err = f.RelayDpAttributesReturnRequest(dpAttributesRelay)
	if err != nil {
		f.LogMessage("3", "DpAttributesReturn, RelayDpAttributesReturnRequest failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// *
// DpAttributesRelay
// get attributes and call RelayDpAttributesRequest
// *
func (f *Facilitator) DpAttributesRelay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		f.LogMessage("2", "DpAttributesRelay, not a POST", nil)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		f.LogMessage("3", "DpAttributesRelay, failed to read body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var dpAttributesRelay DpAttributesRelay
	err = json.Unmarshal(body, &dpAttributesRelay)
	if err != nil {
		f.LogMessage("3", "DpAttributesRelay,failed to unmarshal body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for key, data := range dpAttributesRelay.DpAttributesMap {
		dpAttributesRelay.DpAttributesMap[key] = data
	}

	if dpAttributesRelay.OriginationAddress == f.OwnAddress {
		f.DpAttributesResponseChannel <- dpAttributesRelay.DpAttributesMap
		return
	}

	previousSequenceNumber := dpAttributesRelay.PreviousSequenceNumber
	previousSequenceNumberInt, err := strconv.Atoi(previousSequenceNumber)
	if err != nil {
		f.LogMessage("3", "DpAttributesRelay, failed to convert previousSequenceNumber to an integer", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// globalData.DpAttributesRelayHandlerMutex.Lock()
getattributes:
	if f.DpStates.State != Quiesce1 && f.DpStates.State != Quiesce2 && f.DpStates.State != Quiesce3 {
		currentSequenceNumberInt := previousSequenceNumberInt + 1
		currentSequenceNumber := strconv.Itoa(currentSequenceNumberInt)
		dpAttributesRelay.PreviousSequenceNumber = currentSequenceNumber
		dpAttributesRelay.DpAttributesMap[f.OwnAddress] = DpAttributes{
			Address:        f.OwnAddress,
			LeftAddress:    f.LeftAddress,
			RightAddress:   f.RightAddress,
			DpNumber:       f.DpNumber,
			SequenceNumber: currentSequenceNumber,
			Resource:       f.Resource,
			Iteration:      f.Iteration,
		}
	} else {
		goto getattributes
	}
	// globalData.DpAttributesRelayHandlerMutex.Unlock()

	err = f.RelayDpAttributesRequest(dpAttributesRelay)
	if err != nil {
		f.LogMessage("3", "DpAttributesRelay, RelayDpAttributesRequest failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// *
// StoreData
// *
func (f *Facilitator) StoreData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// do something
		f.LogMessage("2", "StoreData, not a POST", nil)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		f.LogMessage("3", "StoreData, ReadAll body failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var DpEnterResourceData DpEnterResourceData
	err = json.Unmarshal(body, &DpEnterResourceData)
	if err != nil {
		f.LogMessage("3", "StoreData, Unmarshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	globalData.DpEngineMutex.Lock()
	f.EnterData = "true"
	f.ResourceDpNumber = DpEnterResourceData.ResourceDpNumber
	f.Data = DpEnterResourceData.Data
	globalData.DpEngineMutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

// *
// RequestDpAttributesRelay
// *
func (f *Facilitator) RequestDpAttributesRelay(w http.ResponseWriter, r *http.Request) {
	currentAttributes := DpAttributesRelay{
		OriginationAddress:     f.OwnAddress,
		PreviousSequenceNumber: "0",
		DpAttributesMap: DpAttributesCurrent{f.OwnAddress: DpAttributes{
			Address:        f.OwnAddress,
			LeftAddress:    f.LeftAddress,
			RightAddress:   f.RightAddress,
			DpNumber:       f.DpNumber,
			SequenceNumber: "0",
			Resource:       f.Resource,
			Iteration:      f.Iteration,
		},
		},
	}

	// globalData.DpAttributesRelayHandlerMutex.Lock()
	err := f.RelayDpAttributesRequest(currentAttributes)
	// globalData.DpAttributesRelayHandlerMutex.Unlock()
	if err != nil {
		f.LogMessage("3", "RequestDpAttributesRelay,  requestDpAttributesRelay failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// wait for the response from each dp
	var dpAttributesResponse DpAttributesCurrent
	select {
	case dpAttributesResponse = <-f.DpAttributesResponseChannel:
	case <-time.After(60 * time.Second):
		dpAttributesResponse = nil
	}

	jsonData, err := json.Marshal(dpAttributesResponse)
	if err != nil {
		f.LogMessage("3", "RequestDpAttributesRelay, Marshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// *
// RequestAttributesFromSingleDp
// *
func (f *Facilitator) RequestAttributesFromSingleDp(w http.ResponseWriter, r *http.Request) {
	dpAttributes := DpAttributes{
		Address:        f.OwnAddress,
		LeftAddress:    f.LeftAddress,
		RightAddress:   f.RightAddress,
		DpNumber:       f.DpNumber,
		SequenceNumber: "0",
		Resource:       f.Resource,
		Iteration:      f.Iteration,
	}

	jsonData, err := json.Marshal(dpAttributes)
	if err != nil {
		f.LogMessage("3", "requestAttributesFromSingleDp, Marshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// *
// DpResourceRelay
// *
func (f *Facilitator) DpResourceRelay(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		f.LogMessage("3", "DpResourceRelay, ReadAll body failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var clientMessage messageServerStack.ClientMessage
	err = json.Unmarshal(body, &clientMessage)
	if err != nil {
		f.LogMessage("3", "DpResourceRelay, Unmarshal failure", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if clientMessage.ResourceDpNumber == f.DpNumber {
		globalData.DpEngineMutex.Lock()
		clientMessage.Done = "true"
		globalData.RequestsToProcess.Push(clientMessage)
		globalData.DpEngineMutex.Unlock()
		w.WriteHeader(http.StatusOK)
		return
	}

	if clientMessage.OriginatorAddress == f.OwnAddress {
		switch clientMessage.Done {
		case "true":
			f.LogMessage("1", "DpResourceRelay, the retrieve or store response has been completed", nil)
		case "false":
			msg := fmt.Sprintf("DpResourceRelay, the target dp, %s, does not exist", clientMessage.ResourceDpNumber)
			f.LogMessage("1", msg, nil)
			clientMessage.ResultMessage = msg
			clientMessage.Data = ""
		}
		f.ResourceResponseChannel <- clientMessage
		w.WriteHeader(http.StatusOK)
		return
	}

	err = f.RelayDpResourceInformation(clientMessage)
	if err != nil {
		f.LogMessage("3", "DpResourceRelay, relayDpResourceInformation failure", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// *
// RequestDpMessagesRelay
// *
func (f *Facilitator) RequestDpMessagesRelay(w http.ResponseWriter, r *http.Request) {

	var currentDpMessagesRelay DpMessagesRelay
	currentDpMessagesRelay.DpMessagesMap = make(map[string]*ringBuffer.RingBuffer)

	currentDpMessagesRelay.OriginationAddress = f.OwnAddress
	currentDpMessagesRelay.DpMessagesMap[f.DpNumber] = globalData.DpMessages

	globalData.DpMessagesRelayHandlerMutex.Lock()
	err := f.relayDpMessagesRequest(currentDpMessagesRelay)
	globalData.DpMessagesRelayHandlerMutex.Unlock()
	if err != nil {
		f.LogMessage("3", "RequestDpMessagesRelay, relayDpMessagesRequest failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// wait for the response from each dp
	var dpMessagesRelayResponse DpMessagesRelay
	dpMessagesRelayResponse = <-f.DpMessagesRelayResponseChannel

	jsonData, err := json.Marshal(dpMessagesRelayResponse)
	if err != nil {
		f.LogMessage("3", "RequestDpMessagesRelay, Marshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// *
// RequestDpMessagesFromSingleDp
// *
func (f *Facilitator) RequestDpMessagesFromSingleDp(w http.ResponseWriter, r *http.Request) {

	data := struct {
		DpNumber    string                `json:"dpnumber"`
		LogMessages ringBuffer.RingBuffer `json:"logmessages"`
	}{
		DpNumber:    f.DpNumber,
		LogMessages: *globalData.DpMessages,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		f.LogMessage("3", "requestDpMessagesFromSingleDp, Marshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// *
// DpMessagesRelayHandler
// *
func (f *Facilitator) DpMessagesRelayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		f.LogMessage("2", "dpMessagesRelayHandler, not a POST", nil)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		f.LogMessage("3", "dpMessagesRelayHandler, ReadAll body failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var dpMessagesRelay DpMessagesRelay
	err = json.Unmarshal(body, &dpMessagesRelay)
	if err != nil {
		f.LogMessage("3", "dpMessagesRelayHandler, Unmarshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for key, data := range dpMessagesRelay.DpMessagesMap {
		dpMessagesRelay.DpMessagesMap[key] = data
	}

	if dpMessagesRelay.OriginationAddress == f.OwnAddress {
		f.DpMessagesRelayResponseChannel <- dpMessagesRelay
		return
	}

	globalData.DpMessagesRelayHandlerMutex.Lock()

getdpmessages:
	if f.DpStates.State != Quiesce1 && f.DpStates.State != Quiesce2 && f.DpStates.State != Quiesce3 {
		dpMessagesRelay.DpMessagesMap[f.DpNumber] = globalData.DpMessages
	} else {
		goto getdpmessages
	}
	globalData.DpMessagesRelayHandlerMutex.Unlock()
	err = f.relayDpMessagesRequest(dpMessagesRelay)
	if err != nil {
		f.LogMessage("3", "dpMessagesRelayHandler, relayDpMessagesRequest failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// *
// StoreOrRequestDpResourceInformation
// *
func (f *Facilitator) StoreOrRequestDpResourceInformation(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		f.LogMessage("2", "StoreOrRequestDpResourceInformation, not a POST", nil)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		f.LogMessage("3", "StoreOrRequestDpResourceInformation, ReadAll body failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var dpResourceInformation ResourceInformation
	err = json.Unmarshal(body, &dpResourceInformation)
	if err != nil {
		f.LogMessage("3", "StoreOrRequestDpResourceInformation, Unmarshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	requestToProcess := messageServerStack.ClientMessage{
		OriginatorAddress: f.OwnAddress,
		Resource:          dpResourceInformation.Resource,
		ResourceDpNumber:  dpResourceInformation.ResourceDpNumber,
		ResultMessage:     "",
		Done:              "false",
		StoreOrRetrieve:   dpResourceInformation.StoreOrRetrieve,
		Data:              dpResourceInformation.Data,
	}

	if requestToProcess.ResourceDpNumber == f.DpNumber {
		globalData.DpEngineMutex.Lock()
		globalData.RequestsToProcess.Push(requestToProcess)
		globalData.DpEngineMutex.Unlock()
		requestToProcess.Done = "true"
		requestToProcess.ResultMessage = "receiver of request processed the data"
		f.ResourceResponseChannel <- requestToProcess
	} else {
		// if the forwarding of the request failed, DO return an error message to the client/requestor
		// since this is not the dp with the resource, send on the request to the dp on the left
		err = f.RelayDpResourceInformation(requestToProcess)
		if err != nil {
			f.LogMessage("3", "StoreOrRequestDpResourceInformation, relayDpResourceInformation failed", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	var response messageServerStack.ClientMessage
	response = <-f.ResourceResponseChannel

	jsonData, err := json.Marshal(response)
	if err != nil {
		f.LogMessage("3", "StoreOrRequestDpResourceInformation Marshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// *
// ChangeLeftOrRightAddress
// *
func (f *Facilitator) ChangeLeftOrRightAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		f.LogMessage("2", "ChangeRightOrLeftAddressHandler, not a POST", nil)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		f.LogMessage("3", "ChangeRightOrLeftAddressHandler, ReadAll body failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var ChangeLeftOrRightAddress ChangeLeftOrRightAddress
	err = json.Unmarshal(body, &ChangeLeftOrRightAddress)
	if err != nil {
		f.LogMessage("3", "ChangeRightOrLeftAddressHandler, Unmarshal failure", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch ChangeLeftOrRightAddress.LeftOrRightAddress {
	case "Left":
		f.LeftAddress = ChangeLeftOrRightAddress.NewAddress
	case "Right":
		f.RightAddress = ChangeLeftOrRightAddress.NewAddress
	default:
		f.LogMessage("3", "ChangeRightOrLeftAddressHandler, invalid data, side must be Left or Right", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
