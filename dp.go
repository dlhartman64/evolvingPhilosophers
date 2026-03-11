package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"syscall"
	"time"

	"evolvingPhilosophers.local/engine"
	"evolvingPhilosophers.local/globalData"
	"evolvingPhilosophers.local/handler"
	"evolvingPhilosophers.local/messageServerStack"
	ringBuffer "evolvingPhilosophers.local/ringbuffer"
	"evolvingPhilosophers.local/router"
)

func printValidCommands() {
	fmt.Println("valid commands are:")
	fmt.Println("help")
	fmt.Println("initialRing, start dining philosopher create ring")
	fmt.Println("directAddToRing, start dining philosopher for directAddNewDpToLeftOfTarget client command")
	fmt.Println("addToRing, start dining philosopher to add to existing ring")
}

func main() {

	syscall.Mlockall(syscall.MCL_CURRENT | syscall.MCL_FUTURE)

	args := os.Args

	if len(args) < 2 {
		printValidCommands()
		os.Exit(0)
	}

	command := args[1]
	switch command {
	case "help":
	case "initialRing":
	case "directAddToRing":
	case "addToRing":
	default:
		fmt.Printf("%s is not a valid command.  Must be one of:\n", command)
		fmt.Println("help")
		fmt.Println("initialRing")
		fmt.Println("directAddToRing")
		fmt.Println("addToRing")
		os.Exit(0)
	}

	if command == "help" {
		printValidCommands()
		os.Exit(0)
	}

	var initialStateStr string

	initialRingCmd := flag.NewFlagSet("initialRing", flag.ExitOnError)
	initialRingCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s initialRing\n", os.Args[0])
		initialRingCmd.PrintDefaults()
	}
	addressOfDpInitialRingPtr := initialRingCmd.String("addressOfDp", "", "address of new dp (required)")
	addressOfDpOnLeftInitialRingPtr := initialRingCmd.String("addressOfDpOnLeft", "", "address of dp on left (required)")
	addressOfDpOnRightInitialRingPtr := initialRingCmd.String("addressOfDpOnRight", "", "address of dp on right (required)")
	dpNumberInitialRingPtr := initialRingCmd.String("dpNumber", "", "number of new dp (required)")
	dpDeleteTextInitialRingPtr := initialRingCmd.String("deleteExecutable", "false", "Optional delete executable from disk")
	pauseBetweenStatesInitialRingPtr := initialRingCmd.Int("pauseBetweenStates", 10, "Optional pause between states in milliseconds")
	iterationsToSkipBetweenStateDebugOutputInitialRingPtr := initialRingCmd.Int("iterationsToSkipBetweenStateDebugOutput", 1000, "Optional iterations to skip between state debug ouput")
	debugStdoutInitialRingPtr := initialRingCmd.String("debugStdout", "false", "Optional debug output to stdout")

	addToRingCmd := flag.NewFlagSet("addToRing", flag.ExitOnError)
	addToRingCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s addToRing\n", os.Args[0])
		addToRingCmd.PrintDefaults()
	}
	addressOfDpAddToRingPtr := addToRingCmd.String("addressOfDp", "", "address of new dp (required)")
	dpNumberAddToRingPtr := addToRingCmd.String("dpNumber", "", "number of new dp (required)")
	dpDeleteTextAddToRingPtr := addToRingCmd.String("deleteExecutable", "false", "Optional delete executable from disk")
	pauseBetweenStatesAddToRingPtr := addToRingCmd.Int("pauseBetweenStates", 10, "Optional pause between states in milliseconds")
	iterationsToSkipBetweenStateDebugOutputAddToRingPtr := addToRingCmd.Int("iterationsToSkipBetweenStateDebugOutput", 1000, "Optional iterations to skip between state debug ouput")
	debugStdoutAddToRingPtr := addToRingCmd.String("debugStdout", "false", "Optional debug output to stdout")

	var address string
	var leftAddress string
	var rightAddress string
	var dpNumber string
	resource := ""
	var deleteExecutable string

	if command == "initialRing" || command == "directAddToRing" {
		initialRingCmd.Parse(os.Args[2:])

		if *addressOfDpInitialRingPtr == "" {
			fmt.Println("Error: --addressOfDp flag is required")
			initialRingCmd.Usage()
			os.Exit(0)
		}

		if *addressOfDpOnLeftInitialRingPtr == "" {
			fmt.Println("Error: --addressOfDpOnLeft flag is required")
			initialRingCmd.Usage()
			os.Exit(0)
		}

		if *addressOfDpOnRightInitialRingPtr == "" {
			fmt.Println("Error: --addressOfDpOnRight flag is required")
			initialRingCmd.Usage()
			os.Exit(0)
		}

		if *dpNumberInitialRingPtr == "" {
			fmt.Println("Error: --dpNumber flag is required")
			initialRingCmd.Usage()
			os.Exit(0)
		}

		switch command {
		case "initialRing":
			initialStateStr = "Thinking"
		case "directAddToRing":
			initialStateStr = "Quiesce3"
		default:
			initialStateStr = "Thinking"
		}

		address = *addressOfDpInitialRingPtr
		leftAddress = *addressOfDpOnLeftInitialRingPtr
		rightAddress = *addressOfDpOnRightInitialRingPtr
		dpNumber = *dpNumberInitialRingPtr
		deleteExecutable = *dpDeleteTextInitialRingPtr
		globalData.PauseBetweenStates = time.Millisecond * time.Duration(*pauseBetweenStatesInitialRingPtr)
		globalData.SkipBetweenStateOutput = *iterationsToSkipBetweenStateDebugOutputInitialRingPtr
		globalData.DebugToStdout = *debugStdoutInitialRingPtr

		if globalData.DebugToStdout == "true" {
			fmt.Println("command initialRing")
			fmt.Println("	addressOfDp:", *addressOfDpInitialRingPtr)
			fmt.Println("	addressOfDpOnLeft:", *addressOfDpOnLeftInitialRingPtr)
			fmt.Println("	addressOfDpOnRight:", *addressOfDpOnRightInitialRingPtr)
			fmt.Println("	dpNumber:", *dpNumberInitialRingPtr)
			fmt.Println("	deleteExecutable:", *dpDeleteTextInitialRingPtr)
			fmt.Println("	pauseBetweenStates:", *pauseBetweenStatesInitialRingPtr)
			fmt.Println("	iterationsToSkipBetweenStateDebugOutput:", *iterationsToSkipBetweenStateDebugOutputInitialRingPtr)
			fmt.Println("	debugStdout:", *debugStdoutInitialRingPtr)
		}
	}

	if command == "addToRing" {
		addToRingCmd.Parse(os.Args[2:])
		if *addressOfDpAddToRingPtr == "" {
			if globalData.DebugToStdout == "true" {
				fmt.Println("Error: --addressOfDp flag is required")
			}
			addToRingCmd.Usage()
			os.Exit(0)
		}

		if *dpNumberAddToRingPtr == "" {
			if globalData.DebugToStdout == "true" {
				fmt.Println("Error: --dpNumber flag is required")
			}
			addToRingCmd.Usage()
			os.Exit(0)
		}

		initialStateStr = "Quiesce3"
		address = *addressOfDpAddToRingPtr
		leftAddress = ""
		rightAddress = ""
		dpNumber = *dpNumberAddToRingPtr
		deleteExecutable = *dpDeleteTextAddToRingPtr
		globalData.PauseBetweenStates = time.Millisecond * time.Duration(*pauseBetweenStatesAddToRingPtr)
		globalData.SkipBetweenStateOutput = *iterationsToSkipBetweenStateDebugOutputAddToRingPtr
		globalData.DebugToStdout = *debugStdoutAddToRingPtr

		if globalData.DebugToStdout == "true" {
			fmt.Println("command addToRing")
			fmt.Println("	addressOfDp:", *addressOfDpAddToRingPtr)
			fmt.Println("	dpNumber:", *dpNumberAddToRingPtr)
			fmt.Println("	deleteExecutable:", *dpDeleteTextInitialRingPtr)
			fmt.Println("	pauseBetweenStates:", *pauseBetweenStatesInitialRingPtr)
			fmt.Println("	iterationsToSkipBetweenStateDebugOutput:", *iterationsToSkipBetweenStateDebugOutputInitialRingPtr)
			fmt.Println("	debugStdout:", *debugStdoutInitialRingPtr)
		}
	}

	var initialState int
	switch initialStateStr {
	case "Quiesce3":
		initialState = handler.Quiesce3
	case "Thinking":
		initialState = handler.Thinking
	default:
		initialState = handler.Thinking
	}

	leftState := handler.Thinking
	rightState := handler.Thinking

	var wg sync.WaitGroup

	globalData.RequestsToProcess = messageServerStack.NewMessageServerStack()
	globalData.CompletedRequests = messageServerStack.NewMessageServerStack()
	globalData.ResponseToRequestor = messageServerStack.NewMessageServerStack()

	globalData.DpMessages = ringBuffer.NewRingBuffer(handler.RingBufferCapacity)

	f := handler.NewFacilitator(address, leftAddress, rightAddress, dpNumber, initialState, leftState, rightState, resource)

	f.BothForksAvailable.Acquire()

	engineFacilitator := engine.NewEngineFacilitator(f)

	wg.Add(1)
	go engineFacilitator.DpEngine(&wg)

	mux := http.NewServeMux()
	router.RegisterRoutes(mux, f)

	wg.Add(1)
	go f.CollectAndForwardResourceData(&wg)

	wg.Add(1)
	go handler.AntiDeadlockEngine(&wg)

	if deleteExecutable == "true" {
		onDiskBinary, err := os.Executable()
		if err != nil {
			msg := fmt.Sprintf("Error error getting executable path: %s", err)
			f.LogMessage("2", msg, err)
			if globalData.DebugToStdout == "true" {
				fmt.Printf("%s\n", msg)
			}
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Second * 60)
			err = os.Remove(onDiskBinary)
			if err != nil {
				msg := fmt.Sprintf("Error deleting file: %s", onDiskBinary)
				f.LogMessage("2", msg, err)
			}
		}()
	}

	if globalData.DebugToStdout == "true" {
		fmt.Println("dp server listening on address " + address)
	}

	err := http.ListenAndServe(address, mux)

	if err != nil {
		if globalData.DebugToStdout == "true" {
			fmt.Printf("Error starting the dp server: %s\n\n", err)
		}
	}

	// Will never get to this line
	wg.Wait()
}
