package globalData

import (
	"sync"
	"time"

	"evolvingPhilosophers.local/dataStorageHeap"
	"evolvingPhilosophers.local/messageServerStack"
	ringBuffer "evolvingPhilosophers.local/ringbuffer"
)

var PauseBetweenStates time.Duration

var SkipBetweenStateOutput int

var DebugToStdout string

var RequestsToProcess *messageServerStack.MessageServerStack
var CompletedRequests *messageServerStack.MessageServerStack
var ResponseToRequestor *messageServerStack.MessageServerStack

var DpEngineMutex sync.Mutex

var TestMutex sync.Mutex

var DpAttributesRelayHandlerMutex sync.Mutex

var DpMessagesRelayHandlerMutex sync.Mutex

var DataMessageHeap *dataStorageHeap.DataStorageHeap

var DpMessages *ringBuffer.RingBuffer
