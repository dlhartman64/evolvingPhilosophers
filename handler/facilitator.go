package handler

import (
	"evolvingPhilosophers.local/messageServerStack"
	"evolvingPhilosophers.local/semaphore"
)

type Facilitator struct {
	BothForksAvailable              *semaphore.Semaphore
	DpStates                        *DpStates
	OwnAddress                      string
	LeftAddress                     string
	RightAddress                    string
	OriginalLeftAddress             string
	OriginalRightAddress            string
	NewDpAddress                    string
	DpNumber                        string
	Quiesce1Chan                    chan int
	Quiesce2Chan                    chan int
	Quiesce3Chan                    chan int
	ExitQuiesce1Chan                chan int
	ExitQuiesce2Chan                chan int
	ExitQuiesce3Chan                chan int
	Resource                        string
	Data                            string
	ResourceDpNumber                string
	EnterData                       string
	GetData                         string
	ResourceType                    string
	DataHeapChannel                 DataHeapChannel
	DataHeapRequestChannel          DataHeapRequestChannel
	DpAttributesResponseChannel     DpAttriburesResponseChannel
	ReceiveRequestToAddNewDPChannel ReceiveRequestToAddNewDPChannel
	ReceiveRequestToRemoveDPChannel ReceiveRequestToRemoveDPChannel
	DpMessagesRelayResponseChannel  DpMessagesRelayResponseChannel
	Iteration                       int
}

var F *Facilitator
var Count int = 0

func GetFacilitator() *Facilitator {
	return F
}

func NewFacilitator(ownAddress string, leftAddress string, rightAddress string,
	dpNumber string, initialState int, leftState int, rightState int, resource string) *Facilitator {

	if Count == 1 {
		return F
	}
	Count = 1

	dpStates := &DpStates{
		State:      initialState,
		LeftState:  leftState,
		RightState: rightState,
	}
	bothForksAvailable := semaphore.NewSempaphore(1)
	newDpAddress := ""

	quiesce1Chan := make(chan int)
	quiesce2Chan := make(chan int)
	quiesce3Chan := make(chan int)
	exitQuiesce1Chan := make(chan int)
	exitQuiesce2Chan := make(chan int)
	exitQuiesce3Chan := make(chan int)

	data := "default"
	resourceDpNumber := "default"
	enterData := "false"
	getData := "false"
	resourceType := ""

	dataHeapChannel := make(chan messageServerStack.ClientMessage, 100)
	dataHeapRequestChannel := make(chan messageServerStack.ClientMessage, 100)
	dpAttributesResponseChannel := make(chan DpAttributesCurrent, 100)
	receiveRequestToAddNewDpChannel := make(chan InformationToAddNewDp, 100)
	receiveRequestToRemoveDPChannel := make(chan InformationToRemoveDp, 100)
	dpMessagesRelayResponseChannel := make(chan DpMessagesRelay, 100)

	f := &Facilitator{
		BothForksAvailable:              &bothForksAvailable,
		DpStates:                        dpStates,
		OwnAddress:                      ownAddress,
		LeftAddress:                     leftAddress,
		RightAddress:                    rightAddress,
		OriginalLeftAddress:             leftAddress,
		OriginalRightAddress:            rightAddress,
		NewDpAddress:                    newDpAddress,
		DpNumber:                        dpNumber,
		Quiesce1Chan:                    quiesce1Chan,
		Quiesce2Chan:                    quiesce2Chan,
		Quiesce3Chan:                    quiesce3Chan,
		ExitQuiesce1Chan:                exitQuiesce1Chan,
		ExitQuiesce2Chan:                exitQuiesce2Chan,
		ExitQuiesce3Chan:                exitQuiesce3Chan,
		Resource:                        resource,
		Data:                            data,
		ResourceDpNumber:                resourceDpNumber,
		EnterData:                       enterData,
		GetData:                         getData,
		ResourceType:                    resourceType,
		DataHeapChannel:                 dataHeapChannel,
		DataHeapRequestChannel:          dataHeapRequestChannel,
		DpAttributesResponseChannel:     dpAttributesResponseChannel,
		ReceiveRequestToAddNewDPChannel: receiveRequestToAddNewDpChannel,
		ReceiveRequestToRemoveDPChannel: receiveRequestToRemoveDPChannel,
		DpMessagesRelayResponseChannel:  dpMessagesRelayResponseChannel,
		Iteration:                       0,
	}

	F = f

	return f
}
