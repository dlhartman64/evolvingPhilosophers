package engine

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"

	"evolvingPhilosophers.local/dataStorageHeap"
	"evolvingPhilosophers.local/globalData"
	"evolvingPhilosophers.local/handler"
)

type EngineFacilitator struct {
	*handler.Facilitator
}

func NewEngineFacilitator(f *handler.Facilitator) *EngineFacilitator {
	return &EngineFacilitator{
		handler.GetFacilitator(),
	}
}

// *
// Figure out if the dp can get both forks, ie, the fork of the dp on the left and the fork of the dp on the right
// It waits for both forks (semaphores) to be acquired.  The semaphores are released by a handler for a request for the
// dp on the right and the dp on the left
// *
func (ef *EngineFacilitator) takeForks() error {
	var stateLeft int
	var stateRight int
	ef.DpStates.State = handler.Hungry

	if ef.Iteration%globalData.SkipBetweenStateOutput == 0 {
		if globalData.DebugToStdout == "true" {
			printState(ef.DpStates.State, ef.OwnAddress, ef.DpNumber)
			printLeftState(ef.DpStates.LeftState, ef.LeftAddress, ef.DpNumber)
			printRightState(ef.DpStates.RightState, ef.RightAddress, ef.DpNumber)
			fmt.Printf("\n\n")
		}
	}

	stateLeft, err := ef.GetStateLeft()
	if err != nil {
		ef.LogMessage("C", "takeForks, failure of getStateLeft", nil)
		return err
	}
	ef.DpStates.LeftState = stateLeft

	// return state right instead of dpstates
	stateRight, err = ef.GetStateRight()
	if err != nil {
		ef.LogMessage("C", "takeForks, failure of getStateRight", nil)
		return err
	}
	// assign state right to f from return value for f.getStateRight
	ef.DpStates.RightState = stateRight

	ef.TestIfAbleToDine()

	// released by the handler getStateRightAndReceiveStateLeft
	// or by the handler getStateLeftAndReceiveStateRight
	ef.BothForksAvailable.Acquire()

	return nil
}

// This function relinquishes the forks after the dp has finished dining
// It sends the requests getStateLeftAndReceiveStateRight and getStateRightAndReceiveStateLeft
// which in turn can release the semaphores for the receivers
func (ef *EngineFacilitator) putForks() error {
	var err error
	ef.DpStates.State = handler.Thinking

	leftStateAndStateRightRequest := fmt.Sprintf("http://%s/getStateLeftAndReceiveStateRight", ef.LeftAddress)
	params := url.Values{}
	params.Add("state", strconv.Itoa(ef.DpStates.State))
	leftStateAndStateRightRequest = leftStateAndStateRightRequest + "?" + params.Encode()

	rightStateAndStateLeftRequest := fmt.Sprintf("http://%s/getStateRightAndReceiveStateLeft", ef.RightAddress)
	rightStateAndStateLeftRequest = rightStateAndStateLeftRequest + "?" + params.Encode()

	_, err = http.Get(leftStateAndStateRightRequest)
	if err != nil {
		text := fmt.Sprintf("putForks, error sending %s", leftStateAndStateRightRequest)
		ef.LogMessage("C", text, err)
		return err
	}

	_, err = http.Get(rightStateAndStateLeftRequest)
	if err != nil {
		text := fmt.Sprintf("putForks, error sending %s", rightStateAndStateLeftRequest)
		ef.LogMessage("C", text, err)
		return err
	}

	return nil
}

// *
// This function is the heart of the dp. It checks to see
//   - if the dp can get both forks by calling the function takeForks
//   - takeForks sets the state to hungry and returns if it is able to change state to dining
//   - then then the dp is in the critical section with state dining
//   - when the dp leaves the critical section, it calls the function putForks
//   - putForks changes the state to thinking, and sends requests to the left and right dp that could
//     result in those dp being able to change state to dining
//   - at the top of the endless loop, DpEngine checks to see if a channel can be read that would change the state
//     to either Quiesce1, Quiesce2, or Thinking
//     if no channels can be read it continues on to call takeForks if the state is Thinking
//     If the state is not Thinking, it pauses for a set interval, and then checks to see if any channels can be read
//     If any of the exit Quiesce state channels can be read, the state is set to Thinking, and takeForks is called
//
// *
func (ef *EngineFacilitator) DpEngine(wg *sync.WaitGroup) {
	globalData.DataMessageHeap = &dataStorageHeap.DataStorageHeap{}
	var st syscall.Stat_t
	defer wg.Done()
	i := 0
	for {

		if globalData.DebugToStdout == "true" {
			if i%globalData.SkipBetweenStateOutput == 0 {
				printState(ef.DpStates.State, ef.OwnAddress, ef.DpNumber)
				printLeftState(ef.DpStates.LeftState, ef.LeftAddress, ef.DpNumber)
				printRightState(ef.DpStates.RightState, ef.RightAddress, ef.DpNumber)
				fmt.Printf("\n\n")
			}
		}

		think()

		for {
			select {
			case <-ef.Quiesce1Chan:
				ef.DpStates.State = handler.Quiesce1
			case <-ef.Quiesce2Chan:
				ef.DpStates.State = handler.Quiesce2
			case <-ef.Quiesce3Chan:
				ef.DpStates.State = handler.Quiesce3
			case <-ef.ExitQuiesce1Chan:
				ef.DpStates.State = handler.Thinking
			case <-ef.ExitQuiesce2Chan:
				ef.DpStates.State = handler.Thinking
			case <-ef.ExitQuiesce3Chan:
				ef.DpStates.State = handler.Thinking
			default:
			}

			if ef.DpStates.State == handler.Thinking {
				break
			} else {
				if globalData.DebugToStdout == "true" {
					if i%globalData.SkipBetweenStateOutput == 0 {
						fmt.Printf("\nfor in DpEngine\n")
						printState(ef.DpStates.State, ef.OwnAddress, ef.DpNumber)
					}
				}
				time.Sleep(globalData.PauseBetweenStates)
			}
		}

		err := ef.takeForks()
		if err != nil {
			ef.RequestAdjacentDpToLogMessage("takeForks() returned and error, terminating", "F", err)
			os.Exit(0)
		}

		if globalData.DebugToStdout == "true" {
			if i%globalData.SkipBetweenStateOutput == 0 {
				printState(ef.DpStates.State, ef.OwnAddress, ef.DpNumber)
				printLeftState(ef.DpStates.LeftState, ef.LeftAddress, ef.DpNumber)
				printRightState(ef.DpStates.RightState, ef.RightAddress, ef.DpNumber)
				fmt.Printf("\n\n")
			}
		}

		dine()

		if globalData.DpEngineMutex.TryLock() {
			requestToProcess, _ := globalData.RequestsToProcess.Pop()
			if requestToProcess != nil {
				switch requestToProcess.StoreOrRetrieve {
				case "store":
					// use dataStorageHeap
					ctimeSpec := st.Ctimespec
					globalData.DataMessageHeap.Push(&dataStorageHeap.DataStorage{Ctime: ctimeSpec.Sec, Data: requestToProcess.Data})
					// send response about successful storage of the data to a continuously
					// running go function via a channel
					requestToProcess.ResultMessage = "Data stored no problem."
					requestToProcess.Data = ""
					requestToProcess.Done = "true"
					ef.DataHeapChannel <- *requestToProcess
				case "retrieve":
					var output *dataStorageHeap.DataStorage
					var valueFromHeap any
					valueFromHeap = globalData.DataMessageHeap.Pop()
					if valueFromHeap == nil {
						requestToProcess.ResultMessage = fmt.Sprintf("DP %s's dataStorageHeap is empty", ef.DpNumber)
						requestToProcess.Data = ""
					} else {
						output = valueFromHeap.(*dataStorageHeap.DataStorage)
						requestToProcess.ResultMessage = "Data retrieved"
						requestToProcess.Data = output.Data
					}
					requestToProcess.Done = "true"
					ef.DataHeapChannel <- *requestToProcess
				}
			}
			globalData.DpEngineMutex.Unlock()
		}

		err = ef.putForks()
		if err != nil {
			ef.RequestAdjacentDpToLogMessage("putForks() returned and error, terminating", "F", err)
			os.Exit(0)
		}

		if globalData.DebugToStdout == "true" {
			if i%globalData.SkipBetweenStateOutput == 0 {
				printState(ef.DpStates.State, ef.OwnAddress, ef.DpNumber)
				printLeftState(ef.DpStates.LeftState, ef.LeftAddress, ef.DpNumber)
				printRightState(ef.DpStates.RightState, ef.RightAddress, ef.DpNumber)
				fmt.Printf("\n\n")
			}
		}

		i++
		ef.Iteration = i
		time.Sleep(globalData.PauseBetweenStates)

		if globalData.DebugToStdout == "true" {
			if i%globalData.SkipBetweenStateOutput == 0 {
				fmt.Printf("\n\n")
				fmt.Printf("Dp %s, address %s,  iteration: %d\n\n", ef.DpNumber, ef.OwnAddress, i)
			}
		}
	}
}

func think() {
	time.Sleep(globalData.PauseBetweenStates)
}

func dine() {
	time.Sleep(globalData.PauseBetweenStates)
}

func printState(state int, address string, dpNumber string) {
	switch state {
	case handler.Thinking:
		fmt.Printf("DP %s, address %s, is ............ THINKING\n", dpNumber, address)
	case handler.Hungry:
		fmt.Printf("DP %s, address %s, is ............ HUNGRY\n", dpNumber, address)
	case handler.Dining:
		fmt.Printf("DP %s, address %s, is ............ DINING\n", dpNumber, address)
	case handler.Quiesce1:
		fmt.Printf("DP %s, address %s, is ............ QUIESCE1\n", dpNumber, address)
	case handler.Quiesce2:
		fmt.Printf("DP %s, address %s, is ............ QUIESCE2\n", dpNumber, address)
	case handler.Quiesce3:
		fmt.Printf("DP %s, address %s, is ............ QUIESCE3\n", dpNumber, address)
	default:
		fmt.Printf("DP %s, address %s, is in ............ UNKNOWN state: %d\n", dpNumber, address, state)
	}
}

func printLeftState(state int, address string, dpNumber string) {
	switch state {
	case handler.Thinking:
		fmt.Printf("Left DP of DP %s, address %s, is Thinking\n", dpNumber, address)
	case handler.Hungry:
		fmt.Printf("Left DP of DP %s, address %s, is Hungry\n", dpNumber, address)
	case handler.Dining:
		fmt.Printf("Left DP of DP %s, address %s, is Dining\n", dpNumber, address)
	case handler.Quiesce1:
		fmt.Printf("Left DP of DP %s, address %s, is Quiesce1\n", dpNumber, address)
	case handler.Quiesce2:
		fmt.Printf("Left DP of DP %s, address %s, is Quiesce2\n", dpNumber, address)
	case handler.Quiesce3:
		fmt.Printf("Left DP of DP %s, address %s, is Quiesce3\n", dpNumber, address)
	default:
		fmt.Printf("Left DP of DP %s, address %s, is in Unknown state: %d\n", dpNumber, address, state)
	}
}

func printRightState(state int, address string, dpNumber string) {
	switch state {
	case handler.Thinking:
		fmt.Printf("Right DP of DP %s, address %s, is Thinking\n", dpNumber, address)
	case handler.Hungry:
		fmt.Printf("Right DP of DP %s, address %s, is Hungry\n", dpNumber, address)
	case handler.Dining:
		fmt.Printf("Right DP of DP %s, address %s, is Dining\n", dpNumber, address)
	case handler.Quiesce1:
		fmt.Printf("Right DP of DP %s, address %s, is Quiesce1\n", dpNumber, address)
	case handler.Quiesce2:
		fmt.Printf("Right DP of DP %s, address %s, is Quiesce2\n", dpNumber, address)
	case handler.Quiesce3:
		fmt.Printf("Right DP of DP %s, address %s, is Quiesce3\n", dpNumber, address)
	default:
		fmt.Printf("Right DP of DP %s, address %s, is in Unknown state: %d\n", dpNumber, address, state)
	}
}
