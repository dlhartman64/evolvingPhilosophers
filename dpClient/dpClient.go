package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"syscall"
	"time"

	"evolvingPhilosophers.local/messageServerStack"
	ringBuffer "evolvingPhilosophers.local/ringbuffer"
)

type DpAttributesRelay struct {
	OriginationAddress     string              `json:"originationaddress"`
	PreviousSequenceNumber string              `json:"previoussequencenumber"`
	DpAttributesMap        DpAttributesCurrent `json:"dpattributescurrent"`
}

type DpAttributes struct {
	Address        string `json:"address"`
	LeftAddress    string `json:"leftaddress"`
	RightAddress   string `json:"rightaddress"`
	DpNumber       string `json:"dpnumber"`
	SequenceNumber string `json:"sequencenumber"`
	Resource       string `json:"resource"`
	Iteration      int    `json:"iteration"`
	Message        string `json:"message"`
}

type DpAttributesCurrent map[string]DpAttributes

type NewPhilosopher struct {
	Side    string `json:"side"`
	Address string `json:"address"`
}

type QuiesceResponse struct {
	Address string `json:"address"`
	State   string `json:"state"`
}

type QuiesceResponseMap map[string]QuiesceResponse

type StoreOrRetrieveData struct {
	ResourceDpNumber string `json:"resourcedpnumber"`
	Resource         string `json:"resource"`
	StoreOrRetrieve  string `json:"storeorretrieve"`
	Data             string `json:"data"`
}

type InformationToAddNewDp struct {
	OriginatorDpNumber           string `json:"originatordpnumber"`
	TargetDpNumber               string `json:"targetdpnumber"`
	AddNewDpToThisSideOfTargetDp string `json:"addnewdptothissideoftargetdp"`
	NewDpAddress                 string `json:"newdpaddress"`
	Result                       string `json:"result"`
	Done                         string `json:"done"`
	FromRequestHander            string `json:"fromrequesthandler"`
}

type InformationToRemoveDp struct {
	OriginatorDpNumber  string `json:"originatordpnumber"`
	DpNumberToBeRemoved string `json:"dpnumbertoberemoved"`
	OriginalLeftAddress string `json:"originalleftaddress"`
	TerminateSelf       string `json:"terminateself"`
	Result              string `json:"result"`
	Done                string `json:"done"`
	FromRequestHander   string `json:"fromrequesthandler"`
}

type InformationToDirectRemoveDp struct {
	TerminateSelf string `json:"terminateself"`
	Result        string `json:"result"`
	Done          string `json:"done"`
}

type DpMessagesRelay struct {
	OriginationAddress string                            `originatoraddress:"string"`
	DpMessagesMap      map[string]*ringBuffer.RingBuffer `dpmessagesmap:"string"`
}

type ChangeLeftOrRightAddress struct {
	LeftOrRightAddress string `json:"leftorrightaddress"`
	NewAddress         string `json:"newaddress"`
}

func main() {

	syscall.Mlockall(syscall.MCL_CURRENT | syscall.MCL_FUTURE)

	args := os.Args

	if len(args) < 2 || args[1] == "help" {
		fmt.Printf("Commands:\n")
		fmt.Printf("dpClient help\n")
		fmt.Printf("dpClient relayAttributes: list attributes of all dp, relayed\n")
		fmt.Printf("dpClient directRequestAttributesFromSingleDP, list attributes of single dp\n")
		fmt.Printf("dpClient directAddNewDpToLeftOfTarget: add an existing unconnected dp into the ring\n")
		fmt.Printf("dpClient relayAddNewDpToLeftOfTarget: add an existing unconnected dp to the ring\n")
		fmt.Printf("dpClient relayStoreDataOnDp: store a string on a dp\n")
		fmt.Printf("dpClient relayRetrieveDataFromDp: retrieve one of the strings stored by a dp\n")
		fmt.Printf("dpClient directRemoveDp:  remove a dp from ring\n")
		fmt.Printf("dpClient relayRemoveDp: remove a dp from the ring, relayed\n")
		fmt.Printf("dpClient relayRequestLogEntries: get log messages from each dp\n")
		fmt.Printf("dpClient directRequestLogEntriesFromSingleDp: get log messaged from a single dp\n")
		fmt.Printf("dpClient directChangeLeftOrRightAddress: for use when an adjacent dp has terminated\n")
		os.Exit(0)
	}

	command := args[1]
	switch command {
	case "relayAttributes":
	case "directRequestAttributesFromSingleDP":
	case "directAddNewDpToLeftOfTarget":
	case "relayAddNewDpToLeftOfTarget":
	case "relayStoreDataOnDp":
	case "relayRetrieveDataFromDp":
	case "directRemoveDp":
	case "relayRemoveDp":
	case "relayRequestLogEntries":
	case "directRequestLogEntriesFromSingleDp":
	case "directChangeLeftOrRightAddress":
	default:
		fmt.Printf("%s is not a command.  Must be one of:\n", command)
		fmt.Println("help")
		fmt.Println("relayAttributes")
		fmt.Println("directRequestAttributesFromSingleDP")
		fmt.Println("directAddNewDpToLeftOfTarget")
		fmt.Println("relayAddNewDpToLeftOfTarget")
		fmt.Println("relayStoreDataOnDp")
		fmt.Println("relayRetrieveDataFromDp")
		fmt.Println("directRemoveDp")
		fmt.Println("relayRemoveDp")
		fmt.Println("relayRequestLogEntries")
		fmt.Println("directRequestLogEntriesFromSingleDp")
		fmt.Println("directChangeLeftOrRightAddress")
		os.Exit(0)
	}

	//
	// relayAttributes
	//
	relayAttributesCmd := flag.NewFlagSet("relayAttributes", flag.ExitOnError)
	relayAttributesCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage of %s relayAttributes\n", os.Args[0])
		relayAttributesCmd.PrintDefaults()
	}
	dpRelayAttributesStartAddressPtr := relayAttributesCmd.String("dpStartAddress", "", "address:port to receive and forward the request (required)")

	//
	// directRequestAttributesFromSingleDP
	//
	directRequestAttributesFromSingleDpCmd := flag.NewFlagSet("directRequestAttributesFromSingleDP", flag.ExitOnError)
	directRequestAttributesFromSingleDpCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage of %s directRequestAttributesFromSingleDP\n", os.Args[0])
		directRequestAttributesFromSingleDpCmd.PrintDefaults()
	}
	dpAddressDirectRequestAttributesFromSingleDpPtr := directRequestAttributesFromSingleDpCmd.String("dpAddress", "", "address:port of dp to contact directly (required)")

	//
	// directAddNewDpToLeftOfTarget
	//
	directAddNewDpToLeftOfTargetCmd := flag.NewFlagSet("directAddNewDpToLeftOfTarget", flag.ExitOnError)
	directAddNewDpToLeftOfTargetCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage of %s directAddNewDpToLeftOfTarget\n", os.Args[0])
		directAddNewDpToLeftOfTargetCmd.PrintDefaults()
	}
	directTargetDpAddressPtr := directAddNewDpToLeftOfTargetCmd.String("targetAddress", "", "new dp will be added to the left of this dp, address:port (required)")
	directAddressOfDpToLeftOfTargetPtr := directAddNewDpToLeftOfTargetCmd.String("addressOfDPToLeftOfTarget", "", "address of dp which will be to the left of the new dp (required)")
	directAddressOfNewDpPtr := directAddNewDpToLeftOfTargetCmd.String("addressOfNewDp", "", "address:port of new dp (required)")

	//
	// relayStoreDataOnDp
	//
	relayStoreDataOnDpCmd := flag.NewFlagSet("relayStoreDataOnDp", flag.ExitOnError)
	relayStoreDataOnDpCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s relayStoreDataOnDp:\n", os.Args[0])
		relayStoreDataOnDpCmd.PrintDefaults()
	}
	addressOfDpToRelayStoreRequestPtr := relayStoreDataOnDpCmd.String("addressOfDpToForwardRequest", "", "address of dp to receive and forward request (required)")
	numberOfDpThatStoresDataPtr := relayStoreDataOnDpCmd.String("numberOfDpThatStoresData", "", "number of dp that stores the data (required)")
	dataToStorePtr := relayStoreDataOnDpCmd.String("dataToStore", "", "data to store as a string (required)")

	//
	// relayRetrieveDataFromDp
	// relayRetrieveDataFromDp targetAddress DpNumberOfResource resourceName
	//
	relayRetrieveDataFromDpCmd := flag.NewFlagSet("retrieveDataFromDp", flag.ExitOnError)
	relayRetrieveDataFromDpCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s relayRetrieveDataFromDp:\n", os.Args[0])
		relayRetrieveDataFromDpCmd.PrintDefaults()
	}
	addressOfDpToRelayRetrieveRequestFromPtr := relayRetrieveDataFromDpCmd.String("addressOfDpToForwardRequest", "", "address:port of dp to receive and forward the request (required)")
	numberOfDpThatDataIsRetrievedFromPtr := relayRetrieveDataFromDpCmd.String("numberOfDpToRetrieveDataFrom", "", "number of dp that the data is retrieved from (required)")

	//
	// directRemoveDp
	// directRemoveDp targetAddress
	//
	directRemoveDpCmd := flag.NewFlagSet("directRemoveDp", flag.ExitOnError)
	directRemoveDpCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s directRemoveDp:\n", os.Args[0])
		relayRetrieveDataFromDpCmd.PrintDefaults()
	}
	addressOfDpToBeRemovedDirectPtr := directRemoveDpCmd.String("addressOfDpToRemove", "", "address:port of dp to remove (required)")
	directTerminatePtr := directRemoveDpCmd.String("terminate", "true", "Optional dp exits after disconnection from ring")

	//
	// relayRemoveDp
	// originatorDpAddress originatorDpNumber dpNumberToBeRemoved
	//
	relayRemoveDpCmd := flag.NewFlagSet("relayRemoveDp", flag.ExitOnError)
	relayRemoveDpCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s relayRemoveDp:\n", os.Args[0])
		relayRemoveDpCmd.PrintDefaults()
	}
	relayAddressOfDpToForwardRequestPtr := relayRemoveDpCmd.String("addressOfDpToForwardRequest", "", "address:port to receive and forward request (required)")
	relayNumberOfDpToForwardRequestPtr := relayRemoveDpCmd.String("numberOfDpToForwardRequest", "", "number of Dp to receive and forward request (required)")
	relayNumberOfDpToRemovePtr := relayRemoveDpCmd.String("numberOfDpToRemove", "", "number of Dp to remove, (required)")
	relayTerminatePtr := relayRemoveDpCmd.String("terminate", "true", "Optional dp exits after disconnection from ring")

	//
	// relayAddNewDpToLeftOfTarget originatorDpAddress originatorDpNumber targetDpNumber newDpAddress
	//
	relayAddNewDpToLeftOfTargetCmd := flag.NewFlagSet("relayAddNewDpToLeftOfTarget", flag.ExitOnError)
	relayAddNewDpToLeftOfTargetCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s relayAddNewDpToLeftOfTarget:\n", os.Args[0])
		relayAddNewDpToLeftOfTargetCmd.PrintDefaults()
	}
	relayAddressOfDpToForwardAddRequestPtr := relayAddNewDpToLeftOfTargetCmd.String("addressOfDpToForwardRequest", "", "address:port to receive and forward request (required)")
	relayNumberOfDpToForwardAddRequestPtr := relayAddNewDpToLeftOfTargetCmd.String("numberOfDpToForwardRequest", "", "number of Dp to receive and forward request (required)")
	relayNumberOfTargetDpPtr := relayAddNewDpToLeftOfTargetCmd.String("numberOfTargetDp", "", "number of Dp to left of new dp (required)")
	relayAddressOfNewDpPtr := relayAddNewDpToLeftOfTargetCmd.String("addressOfNewDp", "", "address:port of new dp (required)")

	//
	// relayRequestLogEntries
	// originatorDpAddress
	//
	relayRequestLogEntriesCmd := flag.NewFlagSet("relayRequestLogEntries", flag.ExitOnError)
	relayRequestLogEntriesCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s relayRequestLogEntries:\n", os.Args[0])
		relayRequestLogEntriesCmd.PrintDefaults()
	}
	relayAddressOfDpToForwardLogEntriesRequestPtr := relayRequestLogEntriesCmd.String("addressOfDpToForwardRequest", "", "address:port to receive and forward request (required)")

	//
	// directRequestLogEntriesFromSingeleDp
	//
	directRequestLogEntriesFromSingleDpCmd := flag.NewFlagSet("directRequestLogEntriesFromSingleDp", flag.ExitOnError)
	directRequestLogEntriesFromSingleDpCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s directRequestLogEntriesFromSingleDp:\n", os.Args[0])
		directRequestLogEntriesFromSingleDpCmd.PrintDefaults()
	}
	directLogEntriesAddressPtr := directRequestLogEntriesFromSingleDpCmd.String("addressOfDp", "", "address:port of dp to request log entries from (required)")

	//
	// directChangeLeftOrRightAddress
	//
	directChangeLeftOrRightAddressCmd := flag.NewFlagSet("directChangeLeftOrRightAddress", flag.ExitOnError)
	directChangeLeftOrRightAddressCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s directChangeLeftOrRightAddress:\n", os.Args[0])
		directChangeLeftOrRightAddressCmd.PrintDefaults()
	}
	addressOTargetfDpPtr := directChangeLeftOrRightAddressCmd.String("addressOfTargetDp", "", "change the left or right address of this dp (required)")
	leftOrRightAddressPtr := directChangeLeftOrRightAddressCmd.String("leftOrRightAddress", "", "leftOrRightAddress which address to change (required)")
	newAddressPtr := directChangeLeftOrRightAddressCmd.String("newAddress", "", "address:port to change address to (required)")

	//
	// relayAttributes
	// get attributes from each dp around the ring
	//
	if args[1] == "relayAttributes" {

		relayAttributesCmd.Parse(os.Args[2:])
		// fmt.Println("command relayAttributes")
		// fmt.Println("	dpStartAddress :", *dpRelayAttributesStartAddressPtr)

		if *dpRelayAttributesStartAddressPtr == "" {
			fmt.Println("Error: --dpStartAddress flag is required")
			os.Exit(0)
		}

		requestUrl := fmt.Sprintf("http://%s/requestDpAttributesRelay", *dpRelayAttributesStartAddressPtr)
		request, err := http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			fmt.Printf("attributes request build request failure, %v\n", err)
			os.Exit(0)
		}

		request.Header.Add("Accept", "application/json")

		client := &http.Client{Timeout: 60 * time.Second}
		response, err := client.Do(request)
		if err != nil {
			if os.IsTimeout(err) {
				fmt.Printf("attributes request timed out after 120 seconds, %v\n", err)
				os.Exit(0)
			}
			fmt.Printf("attributes request, send request failure, %v\n", err)
			os.Exit(0)
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Printf("Request failed with error code: %s\n", response.Status)
			os.Exit(0)
		}

		var attributes DpAttributesCurrent

		err = json.NewDecoder(response.Body).Decode(&attributes)
		if err != nil {
			fmt.Printf("attributes request failed to decode json response body: %v\n", err)
			os.Exit(0)
		}

		attributesSlice := []DpAttributes{}

		for _, data := range attributes {
			attributesSlice = append(attributesSlice, data)
		}

		sort.Slice(attributesSlice, func(i, j int) bool {
			valueLeft, err := strconv.Atoi(attributesSlice[i].SequenceNumber)
			if err != nil {
				fmt.Printf("attributes request, could not convert %s to an integer\n", (attributesSlice[i].SequenceNumber))
				os.Exit(0)
			}

			valueRight, err := strconv.Atoi(attributesSlice[j].SequenceNumber)
			if err != nil {
				fmt.Printf("attributes request, could not convert %s to an intger\n", attributesSlice[j].SequenceNumber)
				os.Exit(0)
			}
			return valueLeft < valueRight
		})

		fmt.Printf("\ndp address, dp number, sequence, left address, right address, iteration\n")
		for _, data := range attributesSlice {
			fmt.Printf("%s     dpNumber: %s   S: %s   L: %s   R: %s   Iter: %d\n",
				data.Address, data.DpNumber, data.SequenceNumber, data.LeftAddress, data.RightAddress, data.Iteration)
			// fmt.Printf("%s   %s   %s   L: %s   R: %s   M: %s    %d\n",
			// 	data.Address, data.DpNumber, data.SequenceNumber, data.LeftAddress, data.RightAddress, data.Message, data.Iteration)
		}
		fmt.Printf("\n")

		os.Exit(0)
	}

	//
	// directAddNewDpToLeftOfTarget, contact the target with direct connection
	//
	if args[1] == "directAddNewDpToLeftOfTarget" {

		directAddNewDpToLeftOfTargetCmd.Parse(os.Args[2:])
		// fmt.Println("command directAddNewDpToLeftOfTarget")
		// fmt.Println("	targetAddress :", *directTargetDpAddressPtr)
		// fmt.Println("	addressOfDPToLeftOfTarget :", *directAddressOfDpToLeftOfTargetPtr)
		// fmt.Println("	addressOfNewDp :", *directAddressOfNewDpPtr)

		if *directTargetDpAddressPtr == "" {
			fmt.Println("Error: --targetAddress flag is required")
			os.Exit(0)
		}

		if *directAddressOfDpToLeftOfTargetPtr == "" {
			fmt.Println("Error: --addressOfDPToLeftOfTarget flag is required")
			os.Exit(0)
		}

		if *directAddressOfNewDpPtr == "" {
			fmt.Println("Error: --addressOfNewDp flag is required")
			os.Exit(0)
		}

		targetAddress := *directTargetDpAddressPtr
		addressOfDPToLeftOfTarget := *directAddressOfDpToLeftOfTargetPtr
		addressOfNewDP := *directAddressOfNewDpPtr

		//
		// Quiesce1 for target
		//
		requestUrlToQuiesce1ForTarget := fmt.Sprintf("http://%s/quiesce1Right", targetAddress)

		request, err := http.NewRequest("GET", requestUrlToQuiesce1ForTarget, nil)
		if err != nil {
			// do something
			fmt.Printf("addNewDPToLeftOfTarget, build request failure, %v\n", err)
			os.Exit(0)
		}

		request.Header.Add("Accept", "application/json")

		client := &http.Client{Timeout: 120 * time.Second}
		responseQuiesce1Right, err := client.Do(request)
		if err != nil {
			// do something
			fmt.Printf("relayDpAttributesRequest, send request failure, %v\n", err)
			os.Exit(0)
		}

		defer responseQuiesce1Right.Body.Close()

		if responseQuiesce1Right.StatusCode != http.StatusOK {
			fmt.Printf("Request failed with error code: %s\n", responseQuiesce1Right.Status)
			os.Exit(0)
		}

		var quiesce1RightMap QuiesceResponseMap

		err = json.NewDecoder(responseQuiesce1Right.Body).Decode(&quiesce1RightMap)
		if err != nil {
			fmt.Printf("Failed to decode json response body: %v\n", err)
			os.Exit(0)
		}
		fmt.Printf("Return status for Quiesce1Right: %d\n", responseQuiesce1Right.StatusCode)
		for key, data := range quiesce1RightMap {
			// fmt.Printf("key: %s, data: %+v\n", key, data)
			fmt.Printf("key: %s, Address: %s, State: %s\n",
				key, data.Address, data.State)
		}
		fmt.Printf("\n")

		//
		// Quiesce1 for DP to left of target
		//
		requestUrlToQuiesce1ForLeftDPOfTarget := fmt.Sprintf("http://%s/quiesce1Left", addressOfDPToLeftOfTarget)

		request, err = http.NewRequest("GET", requestUrlToQuiesce1ForLeftDPOfTarget, nil)
		if err != nil {
			// do something
			fmt.Printf("addNewDPToLeftOfTarget, build request failure, %v\n", err)
			os.Exit(0)
		}

		request.Header.Add("Accept", "application/json")

		client = &http.Client{Timeout: 120 * time.Second}
		responseQuiesce1Left, err := client.Do(request)
		if err != nil {
			// do something
			fmt.Printf("relayDpAttributesRequest, send request failure, %v\n", err)
			os.Exit(0)
		}

		defer responseQuiesce1Left.Body.Close()

		if responseQuiesce1Left.StatusCode != http.StatusOK {
			fmt.Printf("Request failed with error code: %s\n", responseQuiesce1Left.Status)
			os.Exit(0)
		}

		var quiesce1LeftMap QuiesceResponseMap

		err = json.NewDecoder(responseQuiesce1Left.Body).Decode(&quiesce1LeftMap)
		if err != nil {
			fmt.Printf("responseQuiesce1Right, failed to decode json response body: %v\n", err)
			os.Exit(0)
		}

		fmt.Printf("Return status for Quiesce1Left: %d\n", responseQuiesce1Left.StatusCode)
		for key, data := range quiesce1LeftMap {
			// fmt.Printf("key: %s, data: %+v\n", key, data)
			fmt.Printf("key: %s, Address: %s, State: %s\n",
				key, data.Address, data.State)
		}
		fmt.Printf("\n")

		//
		// Add new DP
		//
		newDP := NewPhilosopher{Side: "Left", Address: addressOfNewDP}

		jsonData, err := json.Marshal(newDP)
		if err != nil {
			fmt.Printf("directAddDp marshal request data failure, %v\n", err.Error())
			os.Exit(0)
		}

		requestUrl := fmt.Sprintf("http://%s/directAddDp", targetAddress)

		request, err = http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("directAddDp failure to create request, %v\n", err.Error())
			os.Exit(0)
		}

		request.Header.Set("Content-Type", "application/json")

		client = &http.Client{Timeout: 120 * time.Second}
		responseAddDp, err := client.Do(request)
		if err != nil {
			// error may require rollback here
			fmt.Printf("directAddDp failure to send request, %v\n", err.Error())
			os.Exit(0)
		}
		defer responseAddDp.Body.Close()

		if responseAddDp.StatusCode != http.StatusOK {
			fmt.Printf("Request %s failed with error code: %s\n", requestUrl, responseAddDp.Status)
			os.Exit(0)
		}

		fmt.Printf("Return status for responseAddDp: %d\n", responseAddDp.StatusCode)

		var responseMapAddNewDp map[string]string

		err = json.NewDecoder(responseAddDp.Body).Decode(&responseMapAddNewDp)
		if err != nil {
			fmt.Printf("directAddDp, failed to decode json response body: %v\n", err)
			os.Exit(0)
		}

		for key, data := range responseMapAddNewDp {
			fmt.Printf("key: %s, data: %s\n", key, data)
		}
		fmt.Printf("\n")

		//
		// DeQuiesce1Right for target
		//
		requestUrlToDeQuiesce1RightForTarget := fmt.Sprintf("http://%s/deQuiesce1Right", targetAddress)

		request, err = http.NewRequest("GET", requestUrlToDeQuiesce1RightForTarget, nil)
		if err != nil {
			// do something
			fmt.Printf("deQuiesce1Right for target, build request failure, %v\n", err)
			os.Exit(0)
		}

		request.Header.Add("Accept", "application/json")

		client = &http.Client{Timeout: 60 * time.Second}
		responseDeQuiesce1Right, err := client.Do(request)
		if err != nil {
			// do something
			fmt.Printf("deQuiesce1Right for target, send request failure, %v\n", err)
			os.Exit(0)
		}

		defer responseDeQuiesce1Right.Body.Close()

		if responseDeQuiesce1Right.StatusCode != http.StatusOK {
			fmt.Printf("deQuiesce1Right for target failed with error code: %s\n", responseDeQuiesce1Right.Status)
			os.Exit(0)
		}

		fmt.Printf("Return status for deQuiesce1Right: %d\n", responseDeQuiesce1Right.StatusCode)

		var responseMapDeQuiesce1Right map[string]string

		err = json.NewDecoder(responseDeQuiesce1Right.Body).Decode(&responseMapDeQuiesce1Right)
		if err != nil {
			fmt.Printf("deQuiesce1Right, failed to decode json response body: %v\n", err)
			os.Exit(0)
		}

		for key, data := range responseMapDeQuiesce1Right {
			fmt.Printf("key: %s, data: %s\n", key, data)
		}
		fmt.Printf("\n")

		//
		// DeQuiesce1Left of target
		//
		requestUrlToDeQuiesce1LeftOfTarget := fmt.Sprintf("http://%s/deQuiesce1Left", addressOfDPToLeftOfTarget)

		request, err = http.NewRequest("GET", requestUrlToDeQuiesce1LeftOfTarget, nil)
		if err != nil {
			// do something
			fmt.Printf("DeQuiesce1Left of target, build request failure, %v\n", err)
			os.Exit(0)
		}

		request.Header.Add("Accept", "application/json")

		client = &http.Client{Timeout: 60 * time.Second}
		responseDeQuiesce1Left, err := client.Do(request)
		if err != nil {
			// do something
			fmt.Printf("DeQuiesce1Left of target, send request failure, %v\n", err)
			os.Exit(0)
		}

		defer responseDeQuiesce1Left.Body.Close()

		if responseDeQuiesce1Left.StatusCode != http.StatusOK {
			fmt.Printf("DeQuiesce1Left of target failed with error code: %s\n", responseDeQuiesce1Left.Status)
			os.Exit(0)
		}

		fmt.Printf("Return status for deQuiesce1Left: %d\n", responseDeQuiesce1Left.StatusCode)

		var responseMapDeQuiesce1Left map[string]string

		err = json.NewDecoder(responseDeQuiesce1Left.Body).Decode(&responseMapDeQuiesce1Left)
		if err != nil {
			fmt.Printf("deQuiesce1Left, failed to decode json response body: %v\n", err)
			os.Exit(0)
		}

		for key, data := range responseMapDeQuiesce1Left {
			fmt.Printf("key: %s, data: %s\n", key, data)
		}
		fmt.Printf("\n")

		//
		// Dequiesce3 for new DP
		//
		requestUrlToDeQuiesce3NewDP := fmt.Sprintf("http://%s/deQuiesce3", addressOfNewDP)

		request, err = http.NewRequest("GET", requestUrlToDeQuiesce3NewDP, nil)
		if err != nil {
			// do something
			fmt.Printf("DeQuiesce3 new DP, build request failure, %v\n", err)
			os.Exit(0)
		}

		request.Header.Add("Accept", "application/json")

		client = &http.Client{Timeout: 60 * time.Second}
		responseDeQuiesce3NewDp, err := client.Do(request)
		if err != nil {
			// do something
			fmt.Printf("DeQuiesce3 new DP, send request failure, %v\n", err)
			os.Exit(0)
		}

		defer responseDeQuiesce3NewDp.Body.Close()

		if responseDeQuiesce3NewDp.StatusCode != http.StatusOK {
			fmt.Printf("DeQuiesce3 new DP failed with error code: %s\n", responseDeQuiesce3NewDp.Status)
			os.Exit(0)
		}

		fmt.Printf("Return status for DeQuiesce3: %d\n", responseDeQuiesce3NewDp.StatusCode)

		var responseMapDeQuiesce3 map[string]string

		err = json.NewDecoder(responseDeQuiesce3NewDp.Body).Decode(&responseMapDeQuiesce3)
		if err != nil {
			fmt.Printf("DeQuiesce3, failed to decode json response body: %v\n", err)
			os.Exit(0)
		}

		for key, data := range responseMapDeQuiesce3 {
			fmt.Printf("key: %s, data: %s\n", key, data)
		}
		fmt.Printf("\n")
	}

	//
	// Send data, in the form of a string, to a dpNumber to store in the dataStorageHeap
	//
	if args[1] == "relayStoreDataOnDp" {

		relayStoreDataOnDpCmd.Parse(os.Args[2:])
		// fmt.Println("command relayStoreDataOnDp")
		// fmt.Println("	addressOfDpToForwardRequest :", *addressOfDpToRelayStoreRequestPtr)
		// fmt.Println("	numberOfDpThatStoresData :", *numberOfDpThatStoresDataPtr)
		// fmt.Println("	dataToStore :", *dataToStorePtr)

		if *addressOfDpToRelayStoreRequestPtr == "" {
			fmt.Println("Error: --addressOfDpToForwardRequest flag is required")
			os.Exit(0)
		}

		if *numberOfDpThatStoresDataPtr == "" {
			fmt.Println("Error: --numberOfDpThatStoresData flag is required")
			os.Exit(0)
		}

		if *dataToStorePtr == "" {
			fmt.Println("Error: --dataToStore flag is required")
			os.Exit(0)
		}

		resourceName := ""

		requestUrl := fmt.Sprintf("http://%s/storeOrRequestDpResourceInformation", *addressOfDpToRelayStoreRequestPtr)

		storeOrRetrieveData := StoreOrRetrieveData{
			ResourceDpNumber: *numberOfDpThatStoresDataPtr,
			Resource:         resourceName,
			StoreOrRetrieve:  "store",
			Data:             *dataToStorePtr,
		}

		jsonData, err := json.Marshal(storeOrRetrieveData)
		if err != nil {
			fmt.Printf("storeOrRequestDpResourceInformation, marshal request data failure, %v\n", err.Error())
			os.Exit(0)
		}

		storeDataOnDpRequest, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("storeOrRequestDpResourceInformation failure to create request, %v\n", err.Error())
			os.Exit(0)
		}

		storeDataOnDpRequest.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 60 * time.Second}
		response, err := client.Do(storeDataOnDpRequest)
		if err != nil {
			fmt.Printf("storeOrRequestDpResourceInformation failure to send request, %v\n", err.Error())
			os.Exit(0)
		}

		var clientMessage messageServerStack.ClientMessage

		err = json.NewDecoder(response.Body).Decode(&clientMessage)
		if err != nil {
			fmt.Printf("storeOrRequestDpResourceInformation, failed to decode json response body: %v\n", err)
			os.Exit(0)
		}

		fmt.Printf("Forwarder Address: %s\n", clientMessage.OriginatorAddress)
		fmt.Printf("DpNumber where stored: %s\n", clientMessage.ResourceDpNumber)
		fmt.Printf("ResultMessage: %s\n", clientMessage.ResultMessage)
		fmt.Printf("StoreOrRetrieve: %s\n", clientMessage.StoreOrRetrieve)
		fmt.Printf("Data: %s\n", clientMessage.Data)
		fmt.Printf("\n")
	}

	// *
	// Request data store in the data storage heap for the DP
	// The request is passed around the ring until it reaches the target DP
	// *
	if args[1] == "relayRetrieveDataFromDp" {

		relayRetrieveDataFromDpCmd.Parse(os.Args[2:])
		// fmt.Println("command relayRetrieveDataFromDp")
		// fmt.Println("	addressOfDpToForwardRequest :", *addressOfDpToRelayRetrieveRequestFromPtr)
		// fmt.Println("	numberOfDpToRetrieveDataFrom :", *numberOfDpThatDataIsRetrievedFromPtr)

		if *addressOfDpToRelayRetrieveRequestFromPtr == "" {
			fmt.Println("Error: --addressOfDpToForwardRequest flag is required")
			os.Exit(0)
		}

		if *numberOfDpThatDataIsRetrievedFromPtr == "" {
			fmt.Println("Error: --numberOfDpToRetrieveDataFrom flag is required")
			os.Exit(0)
		}

		resourceName := ""

		requestUrl := fmt.Sprintf("http://%s/storeOrRequestDpResourceInformation", *addressOfDpToRelayRetrieveRequestFromPtr)

		storeOrRetrieveData := StoreOrRetrieveData{
			ResourceDpNumber: *numberOfDpThatDataIsRetrievedFromPtr,
			Resource:         resourceName,
			StoreOrRetrieve:  "retrieve",
			Data:             "",
		}

		jsonData, err := json.Marshal(storeOrRetrieveData)
		if err != nil {
			fmt.Printf("storeOrRequestDpResourceInformation, marshal request data failure, %v\n", err.Error())
			os.Exit(0)
		}

		storeDataOnDpRequest, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("storeOrRequestDpResourceInformation failure to create request, %v\n", err.Error())
			os.Exit(0)
		}

		storeDataOnDpRequest.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 60 * time.Second}
		response, err := client.Do(storeDataOnDpRequest)
		if err != nil {
			fmt.Printf("storeOrRequestDpResourceInformation failure to send request, %v\n", err.Error())
			os.Exit(0)
		}

		var clientMessage messageServerStack.ClientMessage

		err = json.NewDecoder(response.Body).Decode(&clientMessage)
		if err != nil {
			fmt.Printf("storeOrRequestDpResourceInformation, failed to decode json response body: %v\n", err)
			os.Exit(0)
		}

		fmt.Printf("Forwarder Address: %s\n", clientMessage.OriginatorAddress)
		fmt.Printf("DpNumber where stored: %s\n", clientMessage.ResourceDpNumber)
		fmt.Printf("ResultMessage: %s\n", clientMessage.ResultMessage)
		fmt.Printf("StoreOrRetrieve: %s\n", clientMessage.StoreOrRetrieve)
		fmt.Printf("Data: %s\n", clientMessage.Data)
		fmt.Printf("\n")
	}

	// *
	// directRemoveDp
	// Send request to remove a DP, to the DP being removed.
	// *
	if args[1] == "directRemoveDp" {

		directRemoveDpCmd.Parse(os.Args[2:])
		// fmt.Println("command directRemoveDp")
		// fmt.Println("	addressOfDpToRemove :", *addressOfDpToBeRemovedDirectPtr)
		// fmt.Println("	terminate :", *directTerminatePtr)

		if *addressOfDpToBeRemovedDirectPtr == "" {
			fmt.Println("Error: --addressOfDpToRemove flag is required")
			os.Exit(0)
		}

		addressOfDpToBeRemoved := *addressOfDpToBeRemovedDirectPtr
		terminateSelf := *directTerminatePtr

		if terminateSelf != "true" && terminateSelf != "false" {
			fmt.Printf("terminateSelf must be either true or false\n\n")
			return
		}

		informationToDirectRemoveDp := InformationToDirectRemoveDp{
			TerminateSelf: terminateSelf,
			Result:        "",
			Done:          "false",
		}

		jsonData, err := json.Marshal(informationToDirectRemoveDp)
		if err != nil {
			fmt.Printf("directRemoveDp marshal request data failure, %v\n", err.Error())
			os.Exit(0)
		}

		requestUrl := fmt.Sprintf("http://%s/directRemoveDp", addressOfDpToBeRemoved)
		fmt.Printf("directRemoveDp, requestUrl: %s\n", requestUrl)

		request, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("directRemoveDp failure to create request, %v\n", err.Error())
			os.Exit(0)
		}

		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 60 * time.Second}
		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("directRemoveDp, send request failure, %v\n", err)
			os.Exit(0)
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Printf("directRemoveDp request failed with error code: %s\n", response.Status)
			os.Exit(0)
		}

		fmt.Printf("directRemoveDp return status: %s\n", response.Status)
		os.Exit(0)
	}

	// *
	// relayAddNewDpToLeftOfTarget
	// Send the request to add a new DP and add the new dp to the left of the target dp.
	// Send the request to any DP and the request is passed around the circle until it reaches the target DP
	// *
	if args[1] == "relayAddNewDpToLeftOfTarget" {

		relayAddNewDpToLeftOfTargetCmd.Parse(os.Args[2:])
		// fmt.Println("command relayAddNewDpToLeftOfTarget")
		// fmt.Println("	addressOfDpToForwardRequest :", *relayAddressOfDpToForwardAddRequestPtr)
		// fmt.Println("	numberOfDpToForwardRequest :", *relayNumberOfDpToForwardAddRequestPtr)
		// fmt.Println("	numberOfTargetDp :", *relayNumberOfTargetDpPtr)
		// fmt.Println("	addressOfNewDp :", *relayAddressOfNewDpPtr)

		if *relayAddressOfDpToForwardAddRequestPtr == "" {
			fmt.Println("Error: --addressOfDpToForwardRequest flag is required")
			os.Exit(0)
		}

		if *relayNumberOfDpToForwardAddRequestPtr == "" {
			fmt.Println("Error: --numberOfDpToForwardRequest flag is required")
			os.Exit(0)
		}

		if *relayNumberOfTargetDpPtr == "" {
			fmt.Println("Error: --numberOfTargetDp flag is required")
			os.Exit(0)
		}

		if *relayAddressOfNewDpPtr == "" {
			fmt.Println("Error: --addressOfNewDp flag is required")
			os.Exit(0)
		}

		originatorDpAddress := *relayAddressOfDpToForwardAddRequestPtr
		originatorDpNumber := *relayNumberOfDpToForwardAddRequestPtr
		targetDpNumber := *relayNumberOfTargetDpPtr
		sideOfTargetDp := "Left"
		newDpAddress := *relayAddressOfNewDpPtr

		//
		// Add new DP
		//

		informationToAddNewDp := InformationToAddNewDp{
			OriginatorDpNumber:           originatorDpNumber,
			TargetDpNumber:               targetDpNumber,
			AddNewDpToThisSideOfTargetDp: sideOfTargetDp,
			NewDpAddress:                 newDpAddress,
			Result:                       "",
			Done:                         "false",
		}

		jsonData, err := json.Marshal(informationToAddNewDp)
		if err != nil {
			fmt.Printf("relayAddNewDpToLeftOfTarget marshal request data failure, %v\n", err.Error())
			os.Exit(0)
		}

		requestUrl := fmt.Sprintf("http://%s/receiveRequestToAddNewDpToLeftOfTargetDp", originatorDpAddress)

		request, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("relayAddNewDpToLeftOfTarget failure to create request, %v\n", err.Error())
			os.Exit(0)
		}

		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 240 * time.Second}
		relayAddNewDpResponse, err := client.Do(request)
		if err != nil {
			// error may require rollback here
			if os.IsTimeout(err) {
				fmt.Printf("relayAddNewDpToLeftOfTarget request timed out, %v\n", err.Error())
			} else {
				fmt.Printf("relayAddNewDpToLeftOfTarget failure to send request, %v\n", err.Error())
			}
			os.Exit(0)
		}
		defer relayAddNewDpResponse.Body.Close()

		if relayAddNewDpResponse.StatusCode != http.StatusOK {
			fmt.Printf("Request %s failed with error code: %s\n", requestUrl, relayAddNewDpResponse.Status)
			os.Exit(0)
		}

		fmt.Printf("Return status for relayAddNewDpToLeftOfTarget: %d\n", relayAddNewDpResponse.StatusCode)

		var informationToAddNewDpResponse InformationToAddNewDp

		err = json.NewDecoder(relayAddNewDpResponse.Body).Decode(&informationToAddNewDpResponse)
		if err != nil {
			fmt.Printf("DeQuiesce3, failed to decode json response body: %v\n", err)
			os.Exit(0)
		}

		fmt.Printf("OriginatorDpNumber: %s\n", informationToAddNewDpResponse.OriginatorDpNumber)
		fmt.Printf("TargetDpNumber: %s\n", informationToAddNewDpResponse.TargetDpNumber)
		fmt.Printf("NewDpAddress: %s\n", informationToAddNewDpResponse.NewDpAddress)
		fmt.Printf("Done: %s\n", informationToAddNewDpResponse.Done)
		fmt.Printf("Result: %s\n", informationToAddNewDpResponse.Result)

		fmt.Printf("\n")

		os.Exit(0)
	}

	// *
	// relayRemoveDp
	// Send the request to remove a DP.
	// Send the request to any DP and the request is passed around the circle until it reaches the target DP
	// *
	if args[1] == "relayRemoveDp" {

		relayRemoveDpCmd.Parse(os.Args[2:])
		// fmt.Println("command relayRemoveDp")
		// fmt.Println("	addressOfDpToForwardRequest :", *relayAddressOfDpToForwardRequestPtr)
		// fmt.Println("	numberOfDpToForwardRequest :", *relayNumberOfDpToForwardRequestPtr)
		// fmt.Println("	relayNumberOfDpToRemove :", *relayNumberOfDpToRemovePtr)
		// fmt.Println("	terminate :", *relayTerminatePtr)

		if *relayAddressOfDpToForwardRequestPtr == "" {
			fmt.Println("Error: --addressOfDpToForwardRequest flag is required")
			os.Exit(0)
		}

		if *relayNumberOfDpToForwardRequestPtr == "" {
			fmt.Println("Error: --numberOfDpToForwardRequest flag is required")
			os.Exit(0)
		}

		if *relayNumberOfDpToRemovePtr == "" {
			fmt.Println("Error: --numberOfDpToRemove flag is required")
			os.Exit(0)
		}

		originatorDpAddress := *relayAddressOfDpToForwardRequestPtr
		originatorDpNumber := *relayNumberOfDpToForwardRequestPtr
		dpNumberToRemoved := *relayNumberOfDpToRemovePtr
		terminateSelf := *relayTerminatePtr

		if terminateSelf != "true" && terminateSelf != "false" {
			fmt.Printf("terminateSelf must be either true or false\n\n")
			return
		}

		informationToRemoveDp := InformationToRemoveDp{
			OriginatorDpNumber:  originatorDpNumber,
			DpNumberToBeRemoved: dpNumberToRemoved,
			OriginalLeftAddress: "",
			TerminateSelf:       terminateSelf,
			Result:              "",
			Done:                "false",
		}

		jsonData, err := json.Marshal(informationToRemoveDp)
		if err != nil {
			fmt.Printf("relayRemoveDp marshal request data failure, %v\n", err.Error())
			os.Exit(0)
		}

		requestUrl := fmt.Sprintf("http://%s/receiveRequestToRemoveDp", originatorDpAddress)

		request, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("relayRemoveDp failure to create request, %v\n", err.Error())
			os.Exit(0)
		}

		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 400 * time.Second}
		removeDpResponse, err := client.Do(request)
		if err != nil {
			if os.IsTimeout(err) {
				fmt.Printf("relayRemoveDp request timed out, %v\n", err.Error())
			} else {
				fmt.Printf("relayRemoveDp failure to send request, %v\n", err.Error())
			}
			os.Exit(0)
		}
		defer removeDpResponse.Body.Close()

		if removeDpResponse.StatusCode != http.StatusOK {
			fmt.Printf("Request %s failed with error code: %s\n", requestUrl, removeDpResponse.Status)
			os.Exit(0)
		}

		fmt.Printf("Return status for relayRemoveDp: %d\n", removeDpResponse.StatusCode)

		var informationToRemoveDpResponse InformationToRemoveDp

		err = json.NewDecoder(removeDpResponse.Body).Decode(&informationToRemoveDpResponse)
		if err != nil {
			fmt.Printf("relayRemoveDp, failed to decode json response body: %v\n", err)
			os.Exit(0)
		}

		fmt.Printf("Forwarder DpNumber: %s\n", informationToRemoveDpResponse.OriginatorDpNumber)
		fmt.Printf("DpNumberToBeRemoved: %s\n", informationToRemoveDpResponse.DpNumberToBeRemoved)
		fmt.Printf("LeftAddress of Forwarder: %s\n", informationToRemoveDpResponse.OriginalLeftAddress)
		fmt.Printf("Done: %s\n", informationToRemoveDpResponse.Done)
		fmt.Printf("Result: %s\n", informationToRemoveDpResponse.Result)

		fmt.Printf("\n")

		os.Exit(0)
	}

	//
	// relayRequestLogEntries
	//
	if args[1] == "relayRequestLogEntries" {

		relayRequestLogEntriesCmd.Parse(os.Args[2:])
		// fmt.Println("command relayRequestLogEntries")
		// fmt.Println("	addressOfDpToForwardRequest :", *relayAddressOfDpToForwardLogEntriesRequestPtr)

		if *relayAddressOfDpToForwardLogEntriesRequestPtr == "" {
			fmt.Println("Error: --addressOfDpToForwardRequest flag is required")
			os.Exit(0)
		}

		originatorDpAddress := *relayAddressOfDpToForwardLogEntriesRequestPtr

		requestUrl := fmt.Sprintf("http://%s/requestDpMessagesRelay", originatorDpAddress)

		request, err := http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			fmt.Printf("requestDpMessagesRelay, failure to create request, error: %v\n", err.Error())
			os.Exit(0)
		}

		request.Header.Add("Accept", "application/json")

		client := &http.Client{Timeout: 20 * time.Second}
		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("requestDpMessagesRelay, send request failure, error: %v\n", err)
			os.Exit(0)
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Printf("requestDpMessagesRelay, request failed with response code: %s\n", response.Status)
			os.Exit(0)
		}

		var dpMessagesRelay DpMessagesRelay

		err = json.NewDecoder(response.Body).Decode(&dpMessagesRelay)
		if err != nil {
			fmt.Printf("requestDpMessagesRelay, request failed to decode json response body, error: %v\n", err)
			os.Exit(0)
		}

		var dpMessages map[string]*ringBuffer.RingBuffer

		dpMessages = dpMessagesRelay.DpMessagesMap

		fmt.Printf("\n")
		for item, rb := range dpMessages {
			fmt.Printf("dpNumber: %s\n", item)
			rb.PrintMostRecentElementsFirst()
			fmt.Printf("\n")
		}
		os.Exit(0)
	}

	if args[1] == "directRequestAttributesFromSingleDP" {

		directRequestAttributesFromSingleDpCmd.Parse(os.Args[2:])
		// fmt.Println("command directRequestAttributesFromSingleDP")
		// fmt.Println("	dpAddress:", *dpAddressDirectRequestAttributesFromSingleDpPtr)

		if *dpAddressDirectRequestAttributesFromSingleDpPtr == "" {
			fmt.Println("Error: --dpAddress flag is required")
			os.Exit(0)
		}

		dpAddress := *dpAddressDirectRequestAttributesFromSingleDpPtr

		requestUrl := fmt.Sprintf("http://%s/requestAttributesFromSingleDp", dpAddress)

		request, err := http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			fmt.Printf("attributes request build request failure, %v\n", err)
			os.Exit(0)
		}

		request.Header.Add("Accept", "application/json")

		client := &http.Client{Timeout: 20 * time.Second}
		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("attributes request, send request failure, %v\n", err)
			os.Exit(0)
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Printf("Request failed with error code: %s\n", response.Status)
			os.Exit(0)
		}

		var attributes DpAttributes

		err = json.NewDecoder(response.Body).Decode(&attributes)
		if err != nil {
			fmt.Printf("attributes request failed to decode json response body: %v\n", err)
			os.Exit(0)
		}

		fmt.Printf("\ndp address, dp number, left address, right address, iteration\n")
		fmt.Printf("%s   %s  L: %s   R: %s   %d\n",
			attributes.Address, attributes.DpNumber, attributes.LeftAddress, attributes.RightAddress, attributes.Iteration)
		fmt.Printf("\n")
		os.Exit(0)
	}

	//
	// directRequestLogEntriesFromSingleDp
	//
	if args[1] == "directRequestLogEntriesFromSingleDp" {

		directRequestLogEntriesFromSingleDpCmd.Parse(os.Args[2:])
		// fmt.Println("command directRequestLogEntriesFromSingleDp")
		// fmt.Println("	addressOfDp:", *directLogEntriesAddressPtr)

		if *directLogEntriesAddressPtr == "" {
			fmt.Println("Error: --addressOfDp flag is required")
			os.Exit(0)
		}

		var data struct {
			DpNumber    string                `json:"dpnumber"`
			LogMessages ringBuffer.RingBuffer `json:"logmessages"`
		}

		dpAddress := *directLogEntriesAddressPtr

		requestUrl := fmt.Sprintf("http://%s/requestDpMessagesFromSingleDp", dpAddress)

		request, err := http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			fmt.Printf("attributes request build request failure, %v\n", err)
			os.Exit(0)
		}

		request.Header.Add("Accept", "application/json")

		client := &http.Client{Timeout: 20 * time.Second}
		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("requestDpMessagesFromSingleDp request, send request failure, %v\n", err)
			os.Exit(0)
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Printf("Request failed with error code: %s\n", response.Status)
			os.Exit(0)
		}

		err = json.NewDecoder(response.Body).Decode(&data)
		if err != nil {
			fmt.Printf("requestDpMessagesFromSingleDp request failed to decode json response body: %v\n", err)
			os.Exit(0)
		}

		fmt.Printf("\n")
		fmt.Printf("dpNumber: %s\n", data.DpNumber)
		data.LogMessages.PrintMostRecentElementsFirst()
		fmt.Printf("\n")

		os.Exit(0)
	}

	//
	// directChangeLeftOrRightAddress
	//
	if args[1] == "directChangeLeftOrRightAddress" {

		directChangeLeftOrRightAddressCmd.Parse(os.Args[2:])
		// fmt.Println("command directChangeLeftOrRightAddress")
		// fmt.Println("	addressOfTargetDp:", *addressOTargetfDpPtr)
		// fmt.Println("	leftOrRightAddress:", *leftOrRightAddressPtr)
		// fmt.Println("	newAddress:", *newAddressPtr)

		if *addressOTargetfDpPtr == "" {
			fmt.Println("Error: --leftOrRightAddress flag is required")
			os.Exit(0)
		}

		if *leftOrRightAddressPtr == "" {
			fmt.Println("Error: --leftOrRightAddress flag is required")
			os.Exit(0)
		}

		if *newAddressPtr == "" {
			fmt.Println("Error: --newAddress flag is required")
			os.Exit(0)
		}

		addressOfTargetDp := *addressOTargetfDpPtr

		requestUrl := fmt.Sprintf("http://%s/changeLeftOrRightAddress", addressOfTargetDp)

		changeLeftOrRightAddress := ChangeLeftOrRightAddress{
			LeftOrRightAddress: *leftOrRightAddressPtr,
			NewAddress:         *newAddressPtr,
		}

		jsonData, err := json.Marshal(changeLeftOrRightAddress)
		if err != nil {
			fmt.Printf("directChangeLeftOrRightAddress, marshal request data failure, %v\n", err.Error())
			os.Exit(0)
		}

		request, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("directChangeLeftOrRightAddress request build request failure, %v\n", err)
			os.Exit(0)
		}

		request.Header.Add("Content-Type", "application/json")

		client := &http.Client{Timeout: 30 * time.Second}
		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("directChangeLeftOrRightAddress request, send request failure, %v\n", err)
			os.Exit(0)
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Printf("directChangeLeftOrRightAddress request failed with error code: %s\n", response.Status)
			os.Exit(0)
		}

		fmt.Printf("directRemoveDp return status: %s\n", response.Status)

		os.Exit(0)
	}
}
