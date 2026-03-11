package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"evolvingPhilosophers.local/messageServerStack"
)

// *
// RelayDpResourceInformation
// *
func (f *Facilitator) RelayDpResourceInformation(requestToProcess messageServerStack.ClientMessage) error {

	jsonData, err := json.Marshal(requestToProcess)
	if err != nil {
		f.LogMessage("3", "RelayDpResourceInformation, Marshal failure", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "RelayDpResourceInformation", "Failed to marshal request", err.Error())
		return functionErr
	}

	requestUrl := fmt.Sprintf("http://%s/dpResourceRelay", f.LeftAddress)
	request, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		f.LogMessage("3", "RelayDpResourceInformation, build POST request failure", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "RelayDpResourceInformation", "failed to build POST request", err.Error())
		return functionErr
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		f.LogMessage("3", "RelayDpResourceInformation, failed to send  POST request", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "RelayDpResourceInformation", "failed to send  POST request", err.Error())
		return functionErr
	}

	defer response.Body.Close()
	return nil
}

func (f *Facilitator) RelayDpAttributesRequest(currentAttributes DpAttributesRelay) error {

	jsonData, err := json.Marshal(currentAttributes)
	if err != nil {
		f.LogMessage("3", "RelayDpAttributesRequest, Marshal failure", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "RelayDpAttributesRequest", "failed to marshal currentAttributes", err.Error())
		return functionErr
	}

	requestUrl := fmt.Sprintf("http://%s/dpAttributesRelay", f.LeftAddress)
	request, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		f.LogMessage("3", "RelayDpAttributesRequest, build request failure", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "RelayDpAttributesRequest", "failed to build POST request", err.Error())
		return functionErr
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		if os.IsTimeout(err) {
			f.LogMessage("3", "RelayDpAttributesRequest, dpAttributesRelay request timed out after 30 seconds", err)
		} else {
			f.LogMessage("3", "RelayDpAttributesRequest, send request failure", err)
		}

		newAttributesPtr := &currentAttributes
		newAttributesPtr, err := f.AddEndOfLineData(f.LeftAddress, *newAttributesPtr)
		if err != nil {
			f.LogMessage("3", "RelayDpAttributesRequest, failed to complete AddEndOfLineData call", err)
			functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "RelayDpAttributesRequest", "failed to complete AddEndOfLineData call", err.Error())
			return functionErr
		}

		msg := fmt.Sprintf("RelayDpAttributesRequest, dpNumber %s is calling RelayDpAttributesReturnRequest", f.DpNumber)
		f.LogMessage("1", msg, err)
		err = f.RelayDpAttributesReturnRequest(*newAttributesPtr)
		if err != nil {
			f.LogMessage("3", "RelayDpAttributesRequest, RelayDpAttributesReturnRequest failed", err)
			functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "RelayDpAttributesRequest", "failed to complete RelayDpAttributesReturnRequest call", err.Error())
			return functionErr
		}
		return nil
	}
	defer response.Body.Close()
	return nil
}

// *
// RelayDpAttributesReturnRequest
// *
func (f *Facilitator) RelayDpAttributesReturnRequest(currentAttributes DpAttributesRelay) error {

	jsonData, err := json.Marshal(currentAttributes)
	if err != nil {
		f.LogMessage("3", "RelayDpAttributesReturnRequest, failed to marshal currentAttributes", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "RelayDpAttributesReturnRequest", "failed to marshal currentAttributes", err.Error())
		return functionErr
	}

	var address string
	if currentAttributes.OriginationAddress == f.OwnAddress {
		address = f.OwnAddress
	} else {
		address = f.RightAddress
	}

	requestUrl := fmt.Sprintf("http://%s/dpAttributesReturn", address)
	request, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		f.LogMessage("3", "RelayDpAttributesReturnRequest, build request failure", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayDpMessagesRequest", "failed to build POST request", err.Error())
		return functionErr
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		if os.IsTimeout(err) {
			f.LogMessage("2", "RelayDpAttributesReturnRequest, dpAttributesReturn request timed out after 30 seconds", err)
			functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "RelayDpAttributesReturnRequest", "failed to send POST request", err.Error())
			return functionErr
		}
		f.LogMessage("3", "RelayDpAttributesReturnRequest, send request failure", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "RelayDpAttributesReturnRequest", "failed to send POST request", err.Error())
		return functionErr
	}

	defer response.Body.Close()
	return nil
}

// *
// relayDpMessagesRequest
// *
func (f *Facilitator) relayDpMessagesRequest(dpMessagesRelay DpMessagesRelay) error {
	jsonData, err := json.Marshal(dpMessagesRelay)
	if err != nil {
		f.LogMessage("3", "relayDpMessagesRequest, failed to marshal dpMessagesRelay", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayDpMessagesRequest", "failed to marshal dpMessagesRelay", err.Error())
		return functionErr
	}

	requestUrl := fmt.Sprintf("http://%s/dpMessagesRelayHandler", f.LeftAddress)
	request, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		f.LogMessage("3", "relayDpMessagesRequest, build request failure", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayDpMessagesRequest", "failed to build POST request", err.Error())
		return functionErr
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 180 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		f.LogMessage("3", "relayDpMessagesRequest, send request failure", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayDpMessagesRequest", "failed to send POST request", err.Error())
		return functionErr
	}

	defer response.Body.Close()
	return nil
}

// *
// addNewDpToLeftOfTargetDp
// *
func (f *Facilitator) addNewDpToLeftOfTargetDp(informationToAddNewDp *InformationToAddNewDp) error {

	// send quiesce1 to the dp on the left
	leftQuiesce1RequestUrl := fmt.Sprintf("http://%s/quiesce1Left", f.LeftAddress)

	_, err := http.Get(leftQuiesce1RequestUrl)
	if err != nil {
		text := fmt.Sprintf("addNewDpToLeftOfTargetDp, Error sending request, GET leftQuiesce1RequestUrl: %s", leftQuiesce1RequestUrl)
		f.LogMessage("3", text, err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed to send POST Quiesce1Left", err.Error())
		return functionErr
	}

	f.Quiesce1Chan <- Quiesce1

	rightQuiesce2Request := fmt.Sprintf("http://%s/quiesce2", f.RightAddress)

	responseRightQuiesce2Request, err := http.Get(rightQuiesce2Request)
	if err != nil {
		text := fmt.Sprintf("addNewDpToLeftOfTargetDp, error sending request: %s", rightQuiesce2Request)
		f.LogMessage("3", text, err)
		functionErr := NewForwardError(f.RightAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed to send RightQuiesce2", err.Error())
		return functionErr
	}

	defer responseRightQuiesce2Request.Body.Close()

	var quiesce2Response QuiesceResponse

	err = json.NewDecoder(responseRightQuiesce2Request.Body).Decode(&quiesce2Response)
	if err != nil {
		f.LogMessage("3", "addNewDpToLeftOfTargetDp, NewDecoder Decode failure", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed to decode quiesce2Response", err.Error())
		return functionErr
	}

	originalLeftAddress := f.LeftAddress

	handlerSetLeftAndRightAddressOfNewDpUrl := fmt.Sprintf("http://%s/setLeftAndRightAddressOfNewDp", informationToAddNewDp.NewDpAddress)

	newDpLeftAndRightAddress := NewDpLeftAndRightAddress{
		LeftAddress:  f.LeftAddress,
		RightAddress: f.OwnAddress,
	}

	jsonDataNewDpLeftAndRightAddressOfNewDp, err := json.Marshal(newDpLeftAndRightAddress)
	if err != nil {
		f.LogMessage("3", "addNewDpToLeftOfTargetDp, Marshal failure", err)
		functionErr := NewForwardError(informationToAddNewDp.NewDpAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed marshal newDpLeftAndRightAddress", err.Error())
		return functionErr
	}

	setLeftAndRightAddressOfNewDpRequest, err := http.NewRequest("POST", handlerSetLeftAndRightAddressOfNewDpUrl, bytes.NewBuffer(jsonDataNewDpLeftAndRightAddressOfNewDp))
	if err != nil {
		f.LogMessage("3", "addNewDpToLeftOfTargetDp, failure to create setLeftAndRightAddressOfNewDpRequest", err)
		functionErr := NewForwardError(informationToAddNewDp.NewDpAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed build setLeftAndRightAddressOfNewDpRequest POST", err.Error())
		return functionErr
	}

	setLeftAndRightAddressOfNewDpRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	responseSetLeftAndRightAddressOfNewDp, err := client.Do(setLeftAndRightAddressOfNewDpRequest)
	if err != nil {
		f.LogMessage("3", "addNewDpToLeftOfTargetDp failure to send setLeftAndRightAddressOfNewDpRequest", err)
		functionErr := NewForwardError(informationToAddNewDp.NewDpAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed to send setLeftAndRightAddressOfNewDpRequest POST", err.Error())
		return functionErr
	}
	defer responseSetLeftAndRightAddressOfNewDp.Body.Close()

	// tell the (left or right) adjacent partner to change the address it stores for the left or right dp
	requestUrlChangeOfAddress := fmt.Sprintf("http://%s/changeAddressOfAdjacentDp", f.LeftAddress)

	requestDataChangeOfAddress := map[string]string{
		"side":    "Right",
		"address": informationToAddNewDp.NewDpAddress,
	}

	jsonDataHandlerChangeOfAddress, err := json.Marshal(requestDataChangeOfAddress)
	if err != nil {
		f.LogMessage("3", "addNewDpToLeftOfTargetDp, failed to marshal data", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed to marshal requestDataChangeOfAddress", err.Error())
		return functionErr
	}

	requestChangeOfAddress, err := http.NewRequest("POST", requestUrlChangeOfAddress, bytes.NewBuffer(jsonDataHandlerChangeOfAddress))
	if err != nil {
		f.LogMessage("3", "addNewDpToLeftOfTargetDp, failed to create requestChangeOfAddress", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed to marshal requestDataChangeOfAddress", err.Error())
		return functionErr
	}

	requestChangeOfAddress.Header.Set("Content-Type", "application/json")

	// change the adjacent dp's address from the sender to the new dp
	client = &http.Client{Timeout: 120 * time.Second}
	responseChangeOfAddressRequest, err := client.Do(requestChangeOfAddress)
	if err != nil {
		f.LogMessage("3", "addNewDpToLeftOfTargetDp, failure to send requestChangeOfAddress", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed to send requestChangeOfAddress", err.Error())
		return functionErr
	}

	defer responseChangeOfAddressRequest.Body.Close()

	f.LeftAddress = informationToAddNewDp.NewDpAddress

	deQuiesce2RightRequestUrl := fmt.Sprintf("http://%s/deQuiesce2", f.RightAddress)

	responsDeQuiesce2Right, err := http.Get(deQuiesce2RightRequestUrl)
	if err != nil {
		text := fmt.Sprintf("addNewDpToLeftOfTargetDp, error sending: %s", deQuiesce2RightRequestUrl)
		f.LogMessage("3", text, err)
		functionErr := NewForwardError(f.RightAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed to send DeQuiesce2", err.Error())
		return functionErr
	}

	defer responsDeQuiesce2Right.Body.Close()

	f.ExitQuiesce1Chan <- Thinking
	text := fmt.Sprintf("addNewDpToLeftOfTargetDp, address: %s, left: %s, right: %s", f.OwnAddress, f.LeftAddress, f.RightAddress)
	f.LogMessage("1", text, nil)

	requestUrlToDeQuiesce1LeftOfTarget := fmt.Sprintf("http://%s/deQuiesce1Left", originalLeftAddress)

	requestDeQuiesce1LeftOfTarget, err := http.NewRequest("GET", requestUrlToDeQuiesce1LeftOfTarget, nil)
	if err != nil {
		f.LogMessage("3", "addNewDpToLeftOfTargetDp, build requestUrlToDeQuiesce1LeftOfTarget failure", err)
		functionErr := NewForwardError(originalLeftAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed to build request for DeQuiesce1LeftOfTarget", err.Error())
		return functionErr
	}

	requestDeQuiesce1LeftOfTarget.Header.Add("Accept", "application/json")

	client = &http.Client{Timeout: 120 * time.Second}
	responseDeQuiesce1Left, err := client.Do(requestDeQuiesce1LeftOfTarget)
	if err != nil {
		f.LogMessage("3", "addNewDpToLeftOfTargetDp of target, send request failure", err)
		functionErr := NewForwardError(originalLeftAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed to send DeQuiesce1LeftOfTarget", err.Error())
		return functionErr
	}

	defer responseDeQuiesce1Left.Body.Close()

	if responseDeQuiesce1Left.StatusCode != http.StatusOK {
		text := fmt.Sprintf("addNewDpToLeftOfTargetDp, responseDeQuiesce1Left failed, response code: %s", responseDeQuiesce1Left.Status)
		f.LogMessage("3", text, err)
		functionErr := NewForwardError(originalLeftAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "not OK status returned by send of DeQuiesce1LeftOfTarget", err.Error())
		return functionErr
	}

	requestUrlToDeQuiesce3NewDP := fmt.Sprintf("http://%s/deQuiesce3", informationToAddNewDp.NewDpAddress)

	requestToDeQuiesce3NewDP, err := http.NewRequest("GET", requestUrlToDeQuiesce3NewDP, nil)
	if err != nil {
		f.LogMessage("3", "addNewDpToLeftOfTargetDp, build requestUrlToDeQuiesce3NewDP failure", err)
		functionErr := NewForwardError(informationToAddNewDp.NewDpAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed to build request for DeQuiesce3NewDP", err.Error())
		return functionErr
	}

	requestToDeQuiesce3NewDP.Header.Add("Accept", "application/json")

	client = &http.Client{Timeout: 60 * time.Second}
	responseDeQuiesce3NewDp, err := client.Do(requestToDeQuiesce3NewDP)
	if err != nil {
		f.LogMessage("3", "addNewDpToLeftOfTargetDp, send requestToDeQuiesce3NewDP failure", err)
		functionErr := NewForwardError(informationToAddNewDp.NewDpAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "failed to send request for DeQuiesce3NewDP", err.Error())
		return functionErr
	}

	defer responseDeQuiesce3NewDp.Body.Close()

	if responseDeQuiesce3NewDp.StatusCode != http.StatusOK {
		text := fmt.Sprintf("addNewDpToLeftOfTargetDp,  requestToDeQuiesce3NewDP failed, response code: %s", responseDeQuiesce3NewDp.Status)
		f.LogMessage("3", text, nil)
		functionErr := NewForwardError(informationToAddNewDp.NewDpAddress, f.OwnAddress, "addNewDpToLeftOfTargetDp", "not OK status returned by send of DeQuiesce3NewDP", err.Error())
		return functionErr
	}

	return nil
}

// *
// relayRequestToAddNewDpToLeftOfTargetDp
// *
func (f *Facilitator) relayRequestToAddNewDpToLeftOfTargetDp(informationToAddNewDp *InformationToAddNewDp) error {

	switch informationToAddNewDp.Done {
	case "false":
		if informationToAddNewDp.OriginatorDpNumber == f.DpNumber && informationToAddNewDp.FromRequestHander == "false" {
			msg := fmt.Sprintf("There is no dp with the dpNumber %s in the ring", informationToAddNewDp.TargetDpNumber)
			informationToAddNewDp.Result = msg
			logMessage := fmt.Sprintf("relayRequestToAddNewDpToLeftOfTargetDp, %s", msg)
			f.LogMessage("3", logMessage, nil)
			f.ReceiveRequestToAddNewDPChannel <- *informationToAddNewDp
		} else {
			informationToAddNewDp.FromRequestHander = "false"
			handlerAddNewDpToLeftOfTargetDpUrl := fmt.Sprintf("http://%s/addNewDpToLeftOfTargetDp", f.LeftAddress)

			jsonDataInformationToAddNewDp, err := json.Marshal(informationToAddNewDp)
			if err != nil {
				f.LogMessage("3", "relayRequestToAddNewDpToLeftOfTargetDp, Marshal failure", err)
				functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayRequestToAddNewDpToLeftOfTargetDp", "failed to marshal informationToAddNewDp", err.Error())
				return functionErr
			}

			handlerAddNewDpToLeftOfTargetDpRequest, err := http.NewRequest("POST", handlerAddNewDpToLeftOfTargetDpUrl, bytes.NewBuffer(jsonDataInformationToAddNewDp))
			if err != nil {
				f.LogMessage("3", "relayRequestToAddNewDpToLeftOfTargetDp failure to create request", err)
				functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayRequestToAddNewDpToLeftOfTargetDp", "failed to build POST request for AddNewDpToLeftOfTarget", err.Error())
				return functionErr
			}

			handlerAddNewDpToLeftOfTargetDpRequest.Header.Set("Content-Type", "application/json")

			client := &http.Client{Timeout: 120 * time.Second}
			handlerAddNewDpToLeftOfTargetDpResponse, err := client.Do(handlerAddNewDpToLeftOfTargetDpRequest)
			if err != nil {
				f.LogMessage("3", "relayRequestToAddNewDpToLeftOfTargetDp failure to send request", err)
				functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayRequestToAddNewDpToLeftOfTargetDp", "failed to build POST request for AddNewDpToLeftOfTarget", err.Error())
				return functionErr
			}
			defer handlerAddNewDpToLeftOfTargetDpResponse.Body.Close()
		}
	case "true":
		if informationToAddNewDp.OriginatorDpNumber == f.DpNumber {
			informationToAddNewDp.Result = "Added new DP successfully"
			f.ReceiveRequestToAddNewDPChannel <- *informationToAddNewDp
		} else {
			handlerAddNewDpToLeftOfTargetDpUrl := fmt.Sprintf("http://%s/addNewDpToLeftOfTargetDp", f.LeftAddress)

			jsonDataInformationToAddNewDp, err := json.Marshal(informationToAddNewDp)
			if err != nil {
				f.LogMessage("3", "relayRequestToAddNewDpToLeftOfTargetDp marshal request data failure", err)
				functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayRequestToAddNewDpToLeftOfTargetDp", "failed to marshal informationToAddNewDp", err.Error())
				return functionErr
			}

			handlerAddNewDpToLeftOfTargetDpRequest, err := http.NewRequest("POST", handlerAddNewDpToLeftOfTargetDpUrl, bytes.NewBuffer(jsonDataInformationToAddNewDp))
			if err != nil {
				f.LogMessage("3", "relayRequestToAddNewDpToLeftOfTargetDp failure to create request, %v\n", err)
				functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayRequestToAddNewDpToLeftOfTargetDp", "failed to build POST request for AddNewDpToLeftOfTargetDp", err.Error())
				return functionErr
			}

			handlerAddNewDpToLeftOfTargetDpRequest.Header.Set("Content-Type", "application/json")

			client := &http.Client{Timeout: 120 * time.Second}
			handlerAddNewDpToLeftOfTargetDpResponse, err := client.Do(handlerAddNewDpToLeftOfTargetDpRequest)
			if err != nil {
				f.LogMessage("3", "relayRequestToAddNewDpToLeftOfTargetDp failure to send request", err)
				functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayRequestToAddNewDpToLeftOfTargetDp", "failed to send POST request for AddNewDpToLeftOfTargetDp", err.Error())
				return functionErr
			}
			defer handlerAddNewDpToLeftOfTargetDpResponse.Body.Close()
		}
		return nil
	default:
		// unknown Done value
		f.LogMessage("3", "relayRequestToAddNewDpToLeftOfTargetDp unknown informationToAddNewDp.Done value", nil)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayRequestToAddNewDpToLeftOfTargetDp", "unexpected value for informationToAddNewDp.Done", informationToAddNewDp.Done)
		return functionErr
	}

	return nil
}

// *
// removeDp
// *
func (f *Facilitator) removeDp() error {

	// send quiesce1 to the dp on the left
	leftQuiesce1RequestUrl := fmt.Sprintf("http://%s/quiesce1Left", f.LeftAddress)

	_, err := http.Get(leftQuiesce1RequestUrl)
	if err != nil {
		text := fmt.Sprintf("removeDp, Error sending request, http.Get(leftQuiesce1RequestUrl), leftQuiesce1RequestUrl: %s", leftQuiesce1RequestUrl)
		f.LogMessage("3", text, err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "removeDp", "failed to send leftQuiesce1Request", err.Error())
		return functionErr
	}

	rightQuiesce1RequestUrl := fmt.Sprintf("http://%s/quiesce1Right", f.RightAddress)

	_, err = http.Get(rightQuiesce1RequestUrl)
	if err != nil {
		text := fmt.Sprintf("removeDp, Error sending request, http.Get(rightQuiesce1RequestUrl), rightQuiesce1RequestUrl: %s", rightQuiesce1RequestUrl)
		f.LogMessage("3", text, err)
		functionErr := NewForwardError(f.RightAddress, f.OwnAddress, "removeDp", "failed to send rightQuiesce1Request", err.Error())
		return functionErr
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
		f.LogMessage("3", "removeDp marshal failure", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "removeDp", "failed to marshal leftChangeAddressRequestData", err.Error())
		return functionErr
	}

	leftChangeAddresstRequest, err := http.NewRequest("POST", leftChangeAddressRequestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		f.LogMessage("3", "removeDp failure to create request", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "removeDp", "failed to build leftChangeAddresstRequest", err.Error())
		return functionErr
	}

	leftChangeAddresstRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	_, err = client.Do(leftChangeAddresstRequest)
	if err != nil {
		// error may require rollback here
		f.LogMessage("3", "removeDp failure to send request", err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "removeDp", "failed to send leftChangeAddresstRequest", err.Error())
		return functionErr
	}

	rightChangeAddresstRequestUrl := fmt.Sprintf("http://%s/changeAddressOfAdjacentDp", f.RightAddress)

	rightChangeAddresstRequestData := map[string]string{
		"side":    "Left",
		"address": f.LeftAddress,
	}

	jsonData, err = json.Marshal(rightChangeAddresstRequestData)
	if err != nil {
		f.LogMessage("3", "removeDp, Marshal failure", err)
		functionErr := NewForwardError(f.RightAddress, f.OwnAddress, "removeDp", "failed to marshal rightChangeAddresstRequestData", err.Error())
		return functionErr
	}

	rightChangeAddresstRequest, err := http.NewRequest("POST", rightChangeAddresstRequestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		f.LogMessage("3", "removeDp failure to create request", err)
		functionErr := NewForwardError(f.RightAddress, f.OwnAddress, "removeDp", "failed to build rightChangeAddresstRequest", err.Error())
		return functionErr
	}

	rightChangeAddresstRequest.Header.Set("Content-Type", "application/json")

	client = &http.Client{Timeout: 120 * time.Second}
	_, err = client.Do(rightChangeAddresstRequest)
	if err != nil {
		f.LogMessage("3", "removeDp failure to send request", err)
		functionErr := NewForwardError(f.RightAddress, f.OwnAddress, "removeDp", "failed to send rightChangeAddresstRequest", err.Error())
		return functionErr
	}

	// send dequiesce1 to the dp on its left
	leftDeQuiesce1RequestUrl := fmt.Sprintf("http://%s/deQuiesce1Left", f.LeftAddress)

	_, err = http.Get(leftDeQuiesce1RequestUrl)
	if err != nil {
		text := fmt.Sprintf("removeDp, Error sending request, http.Get(leftDeQuiesce1RequestUrl), leftDeQuiesce1RequestUrl: %s", leftDeQuiesce1RequestUrl)
		f.LogMessage("3", text, err)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "removeDp", "failed to send leftDeQuiesce1Request", err.Error())
		return functionErr
	}

	// send dequiesce1 to the dp on its right
	rightDeQuiesce1RequestUrl := fmt.Sprintf("http://%s/deQuiesce1Right", f.RightAddress)

	_, err = http.Get(rightDeQuiesce1RequestUrl)
	if err != nil {
		text := fmt.Sprintf("removeDp, Error sending request, http.Get(leftDeQuiesce1RequestUrl), leftDeQuiesce1RequestUrl: %s", leftDeQuiesce1RequestUrl)
		f.LogMessage("3", text, err)
		functionErr := NewForwardError(f.RightAddress, f.OwnAddress, "removeDp", "failed to send rightDeQuiesce1Request", err.Error())
		return functionErr
	}

	f.RightAddress = ""
	f.LeftAddress = ""
	return nil
}

// *
// relayRequestToRemoveDp
// *
func (f *Facilitator) relayRequestToRemoveDp(informationToRemoveDp *InformationToRemoveDp) error {

	switch informationToRemoveDp.Done {
	case "false":
		if informationToRemoveDp.OriginatorDpNumber == f.DpNumber && informationToRemoveDp.FromRequestHander == "false" {
			msg := fmt.Sprintf("There is no dp with the dpNumber %s in the ring", informationToRemoveDp.DpNumberToBeRemoved)
			informationToRemoveDp.Result = msg
			logMessage := fmt.Sprintf("relayRequestToAddNewDpToLeftOfTargetDp, %s", msg)
			f.LogMessage("3", logMessage, nil)
			f.ReceiveRequestToRemoveDPChannel <- *informationToRemoveDp
		} else {
			informationToRemoveDp.FromRequestHander = "false"
			// send to handlerRelayRemoveDp
			handlerRelayRemoveDpUrl := fmt.Sprintf("http://%s/relayRemoveDp", f.LeftAddress)

			jsonDataInformationToRemoveDp, err := json.Marshal(informationToRemoveDp)
			if err != nil {
				f.LogMessage("3", "relayRequestToRemoveDp, Marshal failure", err)
				functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayRequestToRemoveDp", "failed to marshal informationToRemoveDp", err.Error())
				return functionErr
			}

			handlerRelayRemoveDpRequest, err := http.NewRequest("POST", handlerRelayRemoveDpUrl, bytes.NewBuffer(jsonDataInformationToRemoveDp))
			if err != nil {
				f.LogMessage("3", "relayRequestToRemoveDp failure to create POST request", err)
				functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayRequestToRemoveDp", "failed to build POST for RelayRemoveDpRequest", err.Error())
				return functionErr
			}

			handlerRelayRemoveDpRequest.Header.Set("Content-Type", "application/json")

			client := &http.Client{Timeout: 120 * time.Second}
			handlerRelayRemoveDpResponse, err := client.Do(handlerRelayRemoveDpRequest)
			if err != nil {
				f.LogMessage("3", "relayRequestToRemoveDp failure to send request", err)
				functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayRequestToRemoveDp", "failed to send POST for RelayRemoveDpRequest", err.Error())
				return functionErr
			}
			defer handlerRelayRemoveDpResponse.Body.Close()
		}
	case "true":
		if informationToRemoveDp.OriginatorDpNumber == f.DpNumber {
			informationToRemoveDp.Result = "Removed dp successfully"
			f.ReceiveRequestToRemoveDPChannel <- *informationToRemoveDp
		} else {

			// send to handlerRelayRemoveDp
			handlerRelayRemoveDpUrl := fmt.Sprintf("http://%s/relayRemoveDp", informationToRemoveDp.OriginalLeftAddress)

			jsonDataInformationToRemoveDp, err := json.Marshal(informationToRemoveDp)
			if err != nil {
				f.LogMessage("3", "relayRequestToRemoveDp Marshal failure", err)
				functionErr := NewForwardError(informationToRemoveDp.OriginalLeftAddress, f.OwnAddress, "relayRequestToRemoveDp", "failed to marshal informationToRemoveDp", err.Error())
				return functionErr
			}

			handlerRelayRemoveDpRequest, err := http.NewRequest("POST", handlerRelayRemoveDpUrl, bytes.NewBuffer(jsonDataInformationToRemoveDp))
			if err != nil {
				f.LogMessage("3", "relayRequestToRemoveDp failure to create POST request", err)
				functionErr := NewForwardError(informationToRemoveDp.OriginalLeftAddress, f.OwnAddress, "relayRequestToRemoveDp", "failed to build request for RelayRemoveDpRequest", err.Error())
				return functionErr
			}

			handlerRelayRemoveDpRequest.Header.Set("Content-Type", "application/json")

			client := &http.Client{Timeout: 120 * time.Second}
			handlerRelayRemoveDpResponse, err := client.Do(handlerRelayRemoveDpRequest)
			if err != nil {
				f.LogMessage("3", "relayRequestToRemoveDp failure to send request", err)
				functionErr := NewForwardError(informationToRemoveDp.OriginalLeftAddress, f.OwnAddress, "relayRequestToRemoveDp", "failed to send POST request for RelayRemoveDpRequest", err.Error())
				return functionErr
			}
			defer handlerRelayRemoveDpResponse.Body.Close()
		}
		return nil
	default:
		// unknown Done value
		f.LogMessage("3", "relayRequestToRemoveDp unknown informationToAddNewDp.Done value", nil)
		functionErr := NewForwardError(f.LeftAddress, f.OwnAddress, "relayRequestToAddNewDpToLeftOfTargetDp", "unexpected value for informationToRemoveDp.Done", informationToRemoveDp.Done)
		return functionErr
	}

	return nil
}
