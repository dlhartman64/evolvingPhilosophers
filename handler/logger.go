package handler

import (
	"fmt"
	"time"

	"evolvingPhilosophers.local/globalData"
)

// type LogMessageFacilitator struct {
// 	*handler.Facilitator
// }

// func NewEngineFacilitator(f *handler.Facilitator) *LogMessageFacilitator {
// 	return &LogMessageFacilitator{
// 		handler.GetFacilitator(),
// 	}
// }

type ForwardError struct {
	TargetAddress string
	OwnAddress    string
	FunctionName  string
	Reason        string
	ErrorMessage  string
}

func NewForwardError(targetAddress string, ownAddress string, functionName string, reason string, errorMessage string) error {
	return &ForwardError{
		TargetAddress: targetAddress,
		OwnAddress:    ownAddress,
		FunctionName:  functionName,
		Reason:        reason,
		ErrorMessage:  errorMessage,
	}
}

func (fe *ForwardError) Error() string {
	return fmt.Sprintf("TargetAddress: %s, OwnAddress: %s, FunctionName: %s, Reason: %s",
		fe.TargetAddress, fe.OwnAddress, fe.FunctionName, fe.Reason)
}

func (f *Facilitator) LogMessage(severity string, text string, err error) {
	currentTime := time.Now()

	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	errText := ""
	if err != nil {
		errText = err.Error()
	}

	logEntry := fmt.Sprintf("%s, %s, %s, %s, %s", f.DpNumber, formattedTime, severity, text, errText)

	severityRange := "all"

	switch severityRange {
	case "all":
		globalData.DpMessages.Add(logEntry)
	}
}
