package handler

import (
	"evolvingPhilosophers.local/messageServerStack"
	ringBuffer "evolvingPhilosophers.local/ringbuffer"
)

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

type ResourceInformation struct {
	ResourceDpNumber string `json:"resourcedpnumber"`
	Resource         string `json:"resource"`
	StoreOrRetrieve  string `json:"storeorretrieve"`
	Data             string `json:"data"`
}

type DpEnterResourceData struct {
	ResourceDpNumber string `json:"resourcedpnumber"`
	Data             string `json:"data"`
}

type DpStates struct {
	State      int
	LeftState  int
	RightState int
}

type RequestData struct {
	State int `json:"state"`
}

// simplify
type DpAttributesCurrent map[string]DpAttributes

// simplify this overcomlicated declaration of DpAttriburesMap
type DpAttributesRelay struct {
	OriginationAddress     string              `json:"originationaddress"`
	PreviousSequenceNumber string              `json:"previoussequencenumber"`
	DpAttributesMap        DpAttributesCurrent `json:"dpattributescurrent"`
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

type DataHeapChannel chan messageServerStack.ClientMessage
type DataHeapRequestChannel chan messageServerStack.ClientMessage
type DpAttriburesResponseChannel chan DpAttributesCurrent
type ReceiveRequestToAddNewDPChannel chan InformationToAddNewDp
type ReceiveRequestToRemoveDPChannel chan InformationToRemoveDp
type DpMessagesRelayResponseChannel chan DpMessagesRelay

type NewPhilosopher struct {
	Side    string `json:"side"`
	Address string `json:"address"`
}

type QuiesceResponse struct {
	Address string `json:"address"`
	State   string `json:"state"`
}

type NewDpLeftAndRightAddress struct {
	LeftAddress  string `json:"leftaddress"`
	RightAddress string `json:"rightaddress"`
}

type ChangeLeftOrRightAddress struct {
	LeftOrRightAddress string `json:"leftorrightaddress"`
	NewAddress         string `json:"newaddress"`
}

type ForeignLogMessage struct {
	SenderAddress  string `json:"senderAddress"`
	SenderDpNumber string `json:"senderDpNumber"`
	Message        string `json:"message"`
	Severity       string `json:"severity"`
	ForeignErr     error  `json:"foreignerr"`
}
