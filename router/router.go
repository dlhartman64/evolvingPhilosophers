package router

import (
	"net/http"

	"evolvingPhilosophers.local/handler"
)

func RegisterRoutes(mux *http.ServeMux, f *handler.Facilitator) {
	mux.HandleFunc("/stateFromAdjacentDp", f.StateFromAdjacentDp)

	mux.HandleFunc("/getStateRightAndReceiveStateLeft", f.GetStateRightAndReceiveStateLeft)
	mux.HandleFunc("/getStateLeftAndReceiveStateRight", f.GetStateLeftAndReceiveStateRight)

	// quiesce handlers
	mux.HandleFunc("/quiesce1Left", f.Quiesce1Left)
	mux.HandleFunc("/quiesce1Right", f.Quiesce1Right)
	mux.HandleFunc("/quiesce2", f.Quiesce2)
	mux.HandleFunc("/deQuiesce1Left", f.DeQuiesce1Left)
	mux.HandleFunc("/deQuiesce1Right", f.DeQuiesce1Right)
	mux.HandleFunc("/deQuiesce2", f.DeQuiesce2)
	mux.HandleFunc("/deQuiesce3", f.DeQuiesce3)

	// add, remove
	mux.HandleFunc("/directAddDp", f.DirectAddDp)
	mux.HandleFunc("/directRemoveDp", f.DirectRemoveDp)
	mux.HandleFunc("/changeAddressOfAdjacentDp", f.ChangeAddressOfAdjacentDp)
	mux.HandleFunc("/setLeftAndRightAddressOfNewDp", f.SetLeftAndRightAddressOfNewDp)
	mux.HandleFunc("/receiveRequestToAddNewDpToLeftOfTargetDp", f.ReceiveRequestToAddNewDpToLeftOfTargetDp)
	mux.HandleFunc("/addNewDpToLeftOfTargetDp", f.AddNewDpToLeftOfTargetDp)
	mux.HandleFunc("/receiveRequestToRemoveDp", f.ReceiveRequestToRemoveDp)
	mux.HandleFunc("/relayRemoveDp", f.RelayRemoveDp)

	// collect attributes from each dp
	mux.HandleFunc("/dpAttributesRelay", f.DpAttributesRelay)
	mux.HandleFunc("/dpAttributesReturn", f.DpAttributesReturn)
	mux.HandleFunc("/requestDpAttributesRelay", f.RequestDpAttributesRelay)
	mux.HandleFunc("/requestAttributesFromSingleDp", f.RequestAttributesFromSingleDp)

	// data handling
	mux.HandleFunc("/storeData", f.StoreData)
	mux.HandleFunc("/storeOrRequestDpResourceInformation", f.StoreOrRequestDpResourceInformation)
	mux.HandleFunc("/dpResourceRelay", f.DpResourceRelay)
	mux.HandleFunc("/requestDpMessagesRelay", f.RequestDpMessagesRelay)
	mux.HandleFunc("/dpMessagesRelayHandler", f.DpMessagesRelayHandler)
	mux.HandleFunc("/requestDpMessagesFromSingleDp", f.RequestDpMessagesFromSingleDp)
	mux.HandleFunc("/dataStorageHeapResponseRelay", f.DataStorageHeapResponseRelay)

	// log message from another dp
	mux.HandleFunc("/requestToLogMessage", f.RequestToLogMessage)

	// change right or left address, experimental
	mux.HandleFunc("/changeLeftOrRightAddress", f.ChangeLeftOrRightAddress)
}
