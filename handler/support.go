package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"evolvingPhilosophers.local/globalData"
	"evolvingPhilosophers.local/messageServerStack"
)

func AntiDeadlockEngine(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		time.Sleep(time.Second * 60)
	}
}

func (f *Facilitator) CollectAndForwardResourceData(wg *sync.WaitGroup) {
	var response messageServerStack.ClientMessage
	defer wg.Done()

	for {
		response = <-f.ResourceResponseChannel

		// send this data to the resource data relay handler
		err := f.RelayDpResourceInformation(response)
		if err != nil {
			f.LogMessage("2", "collectAndForwardResourceData, relayDpResourceInformation failed", err)
		}
	}
}

func (f *Facilitator) GetStateLeft() (int, error) {
	leftRequestUrl := fmt.Sprintf("http://%s/stateFromAdjacentDp", f.LeftAddress)
	count := 0
lrr:
	if count > 100 {
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "GetStateLeft", "count of leftRequestUrl exceeds 100", "count exceeds 100")
		f.LogMessage("3", "GetStateLeft count exceeds 100", functionErr)
		return -1, functionErr
	}

	request, err := http.NewRequest("GET", leftRequestUrl, nil)
	if err != nil {
		f.LogMessage("3", "GetStateLeft build request failure", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "GetStateLeft", "failed to build request", err.Error())
		return -1, functionErr
	}

	client := &http.Client{Timeout: 30 * time.Second}
	leftRequestResponse, err := client.Do(request)

	if err != nil {
		text := fmt.Sprintf("getStateLeft, Error sending request, leftRequestUrl: %s, count: %d", leftRequestUrl, count)
		f.LogMessage("2", text, err)
		time.Sleep(time.Second * 30)
		count++
		goto lrr
	}

	defer leftRequestResponse.Body.Close()

	leftRequestResponseBody, err := io.ReadAll(leftRequestResponse.Body)
	if err != nil {
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "GetStateLeft", "failed to read response body", err.Error())
		f.LogMessage("3", "getStateLeft, ReadAll body failure", err)
		return -1, functionErr
	}

	var leftRequestResponseContent RequestData

	err = json.Unmarshal(leftRequestResponseBody, &leftRequestResponseContent)
	if err != nil {
		f.LogMessage("3", "getStateLeft, Error unmarshalling response body", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "GetStateLeft", "failed to unmarshal leftRequestResponseContent", err.Error())
		return -1, functionErr
	}

	leftState := leftRequestResponseContent.State
	f.DpStates.LeftState = leftState

	return leftState, nil
}

func (f *Facilitator) GetStateRight() (int, error) {
	rightRequestUrl := fmt.Sprintf("http://%s/stateFromAdjacentDp", f.RightAddress)
	count := 0
rrr:
	if count > 100 {
		functionErr := NewForwardError(f.RightAddress, f.OwnAddress, "GetStateRight", "count of rightRequestUrl exceeds 100", "count exceeds 100")
		f.LogMessage("3", "GetStateRight count exceeds 100", functionErr)
		return -1, functionErr
	}

	request, err := http.NewRequest("GET", rightRequestUrl, nil)
	if err != nil {
		f.LogMessage("3", "GetStateRight build request failure", err)
		functionErr := NewForwardError(f.RightAddress, f.OwnAddress, "GetStateRight", "failed to build request", err.Error())
		return -1, functionErr
	}

	client := &http.Client{Timeout: 30 * time.Second}
	rightRequestResponse, err := client.Do(request)

	if err != nil {
		text := fmt.Sprintf("GetStateRight, Error sending request, rightRequestUrl: %s, count: %d", rightRequestUrl, count)
		f.LogMessage("2", text, err)
		time.Sleep(time.Second * 60)
		count++
		goto rrr
	}

	defer rightRequestResponse.Body.Close()

	rightRequestResponseBody, err := io.ReadAll(rightRequestResponse.Body)
	if err != nil {
		functionErr := NewForwardError(f.RightAddress, f.OwnAddress, "GetStateRight", "failed to read response body", err.Error())
		f.LogMessage("3", "getStateRight, ReadAll body failure", err)
		return -1, functionErr
	}

	var rightRequestResponseContent RequestData

	err = json.Unmarshal(rightRequestResponseBody, &rightRequestResponseContent)
	if err != nil {
		f.LogMessage("3", "getStateLeft, Error unmarshalling response body", err)
		functionErr := NewForwardError(f.RightAddress, f.OwnAddress, "GetStateLeft", "failed to unmarshal rightRequestResponseContent", err.Error())
		f.LogMessage("3", "getStateRight, Unmarshal failure", err)
		return -1, functionErr
	}

	rightState := rightRequestResponseContent.State
	f.DpStates.RightState = rightState
	return rightState, nil
}

func (f *Facilitator) AddEndOfLineData(unreachableAddress string, dpAttributesRelay DpAttributesRelay) (*DpAttributesRelay, error) {
	previousSequenceNumber := dpAttributesRelay.PreviousSequenceNumber
	previousSequenceNumberInt, err := strconv.Atoi(previousSequenceNumber)
	if err != nil {
		f.LogMessage("3", "AddEndOfLineData, failed to convert previousSequenceNumber to an integer", err)
		return nil, err
	}

	// globalData.DpAttributesRelayHandlerMutex.Lock()
	currentSequenceNumberInt := previousSequenceNumberInt + 1
	currentSequenceNumber := strconv.Itoa(currentSequenceNumberInt)
	dpAttributesRelay.PreviousSequenceNumber = currentSequenceNumber
	dpAttributesRelay.DpAttributesMap[unreachableAddress] = DpAttributes{
		Address:        unreachableAddress,
		LeftAddress:    "XXXXXXXXXX",
		RightAddress:   "XXXXXXXXXX",
		DpNumber:       "-1",
		SequenceNumber: currentSequenceNumber,
		Resource:       "",
		Iteration:      -1,
		Message:        "refused connection",
	}
	// globalData.DpAttributesRelayHandlerMutex.Unlock()

	return &dpAttributesRelay, nil
}

func (f *Facilitator) TestIfAbleToDine() {
	globalData.TestMutex.Lock()
	defer globalData.TestMutex.Unlock()
	if f.DpStates.State == Hungry &&
		f.DpStates.LeftState != Dining &&
		f.DpStates.RightState != Dining {
		f.DpStates.State = Dining
		f.BothForksAvailable.Release()
	}
}
