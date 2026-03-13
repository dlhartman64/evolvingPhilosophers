# evolvingPhilosophers
Evolving Philosophers

The Evolving Philosophers problem makes it posible to add new dining philosophers to the ring without
affecting entry into the critical section of each dining philosoper.  Each dining
philosopher here is a process, and the communication is REST.

The design here is based on the paper "The Evolving Philosophers Problem:
Dynamic Change Management" by Jeff Kramer and Jeff Magee.  4 dining philosophers, 
referred to here as dps, are put in quiesent states, and a new dp is added with 2 on each
side. Similar with removal.  Dijkstra's algorithm is used.

The dp binary is the dp, and the dpClient binary is used to send commands to the dp.

Build the binaries:

go build -ldflags '-extldflags "-static"' dp.go

go build -ldflags '-extldflags "-static"' dpClient.go

Enter ./dp or ./dpClient to get a list of commands.
Enter ./do <command> or ./dpClient <command> to get the options for the commands

Start 5 dp:

./dp initialRing --addressOfDp=localhost:8080 --addressOfDpOnLeft=localhost:8084 --addressOfDpOnRight=localhost:8081 --dpNumber=1 --debugStdout=true

./dp initialRing --addressOfDp=localhost:8081 --addressOfDpOnLeft=localhost:8080 --addressOfDpOnRight=localhost:8082 --dpNumber=2 --debugStdout=true

./dp initialRing --addressOfDp=localhost:8082 --addressOfDpOnLeft=localhost:8081 --addressOfDpOnRight=localhost:8083 --dpNumber=3 --debugStdout=true

./dp initialRing --addressOfDp=localhost:8083 --addressOfDpOnLeft=localhost:8082 --addressOfDpOnRight=localhost:8084 --dpNumber=4 --debugStdout=true

./dp initialRing --addressOfDp=localhost:8084 --addressOfDpOnLeft=localhost:8083 --addressOfDpOnRight=localhost:8080 --dpNumber=5 --debugStdout=true

Each dp will enter the critical section after all are started.


View each dp in the ring:

./dpClient relayAttributes --dpStartAddress=localhost:8080

command relayAttributes
	dpStartAddress : localhost:8080

dp address, dp number, sequence, left address, right address, iteration

localhost:8080   1   0   L: localhost:8084   R: localhost:8081   M:     3932

localhost:8084   5   1   L: localhost:8083   R: localhost:8080   M:     3800

localhost:8083   4   2   L: localhost:8082   R: localhost:8084   M:     3796

localhost:8082   3   3   L: localhost:8081   R: localhost:8083   M:     3799

localhost:8081   2   4   L: localhost:8080   R: localhost:8082   M:     3931


Start a new dp:
./dp addToRing --addressOfDp=localhost:8085 --dpNumber=6 --debugStdout=true

Add to ring:
./dpClient relayAddNewDpToLeftOfTarget --addressOfDpToForwardRequest=localhost:8080 --numberOfDpToForwardRequest=1 --numberOfTargetDp=4 --addressOfNewDp=localhost:8085

command relayAddNewDpToLeftOfTarget

	addressOfDpToForwardRequest : localhost:8080
	
	numberOfDpToForwardRequest : 1
	
	numberOfTargetDp : 4
	
	addressOfNewDp : localhost:8085
	
Return status for relayAddNewDpToLeftOfTarget: 200

OriginatorDpNumber: 1

TargetDpNumber: 4

NewDpAddress: localhost:8085

Done: true

Result: Added new DP successfully


View dps in ring:

./dpClient relayAttributes --dpStartAddress=localhost:8080 
command relayAttributes
	dpStartAddress : localhost:8080

dp address, dp number, sequence, left address, right address, iteration

localhost:8080   1   0   L: localhost:8084   R: localhost:8081   M:     98858

localhost:8084   5   1   L: localhost:8083   R: localhost:8080   M:     98042

localhost:8083   4   2   L: localhost:8085   R: localhost:8084   M:     98041

localhost:8085   6   3   L: localhost:8082   R: localhost:8083   M:     1247

localhost:8082   3   4   L: localhost:8081   R: localhost:8085   M:     98038

localhost:8081   2   5   L: localhost:8080   R: localhost:8082   M:     98856


Remove a dp:

./dpClient relayRemoveDp --addressOfDpToForwardRequest=localhost:8080 --numberOfDpToForwardRequest=1 --numberOfDpToRemove=3

command relayRemoveDp

	addressOfDpToForwardRequest : localhost:8080
	
	numberOfDpToForwardRequest : 1
	
	relayNumberOfDpToRemove : 3
	
	terminate : true
	
Return status for relayRemoveDp: 200

Forwarder DpNumber: 1

DpNumberToBeRemoved: 3

LeftAddress of Forwarder: localhost:8084

Done: true

Result: Removed dp successfully


View dps in ring:

./dpClient relayAttributes --dpStartAddress=localhost:8080

command relayAttributes

	dpStartAddress : localhost:8080

dp address, dp number, sequence, left address, right address, iteration

localhost:8080   1   0   L: localhost:8084   R: localhost:8081   M:     103974

localhost:8084   5   1   L: localhost:8083   R: localhost:8080   M:     103158

localhost:8083   4   2   L: localhost:8085   R: localhost:8084   M:     103156

localhost:8085   6   3   L: localhost:8081   R: localhost:8083   M:     6364

localhost:8081   2   4   L: localhost:8080   R: localhost:8085   M:     103973



Each dp has a ring buffer to store log messages.

./dpClient relayRequestLogEntries --addressOfDpToForwardRequest=localhost:8080
command relayRequestLogEntries
	addressOfDpToForwardRequest : localhost:8080

dpNumber: 1

index: 7,  1, 2026-03-09 22:47:18, 3, relayRequestToAddNewDpToLeftOfTargetDp, There is no dp with the dpNumber 3 in the ring, 

index: 6,  1, 2026-03-09 22:47:18, 3, relayRequestToAddNewDpToLeftOfTargetDp, There is no dp with the dpNumber 3 in the ring, 

index: 5,  1, 2026-03-09 22:47:18, 1, DeQuiesce2, address: localhost:8080, left: localhost:8084, right: localhost:8081, 

index: 4,  1, 2026-03-09 22:47:18, 1, GetStateRightAndReceiveStateLeft,f.dpStates.state == Quiesce2, 

index: 3,  1, 2026-03-09 22:47:18, 1, Quiesce2, address: localhost:8080, left: localhost:8084, right: localhost:8081, 

index: 2,  1, 2026-03-09 21:46:42, 2, getStateLeft, Error sending request, leftRequestUrl: http://localhost:8084/stateFromAdjacentDp, count: 2, Get "http://localhost:8084/stateFromAdjacentDp": dial tcp 127.0.0.1:8084: connect: connection refused

index: 1,  1, 2026-03-09 21:46:12, 2, getStateLeft, Error sending request, leftRequestUrl: http://localhost:8084/stateFromAdjacentDp, count: 1, Get "http://localhost:8084/stateFromAdjacentDp": dial tcp 127.0.0.1:8084: connect: connection refused

index: 0,  1, 2026-03-09 21:45:42, 2, getStateLeft, Error sending request, leftRequestUrl: http://localhost:8084/stateFromAdjacentDp, count: 0, Get "http://localhost:8084/stateFromAdjacentDp": dial tcp 127.0.0.1:8084: connect: connection refused

dpNumber: 2

index: 6,  2, 2026-03-09 22:47:18, 1, DeQuiesce1Left, address: localhost:8081, left: localhost:8080, right: localhost:8085, 

index: 5,  2, 2026-03-09 22:47:18, 1, Quiesce1Left, address: localhost:8081, left: localhost:8080, right: localhost:8082, 

index: 4,  2, 2026-03-09 22:43:56, 1, DeQuiesce2, address: localhost:8081, left: localhost:8080, right: localhost:8082, 

index: 3,  2, 2026-03-09 22:43:56, 1, GetStateRightAndReceiveStateLeft,f.dpStates.state == Quiesce2, 

index: 2,  2, 2026-03-09 22:43:56, 1, GetStateRightAndReceiveStateLeft,f.dpStates.state == Quiesce2, 

index: 1,  2, 2026-03-09 22:43:56, 1, Quiesce2, address: localhost:8081, left: localhost:8080, right: localhost:8082, 

index: 0,  2, 2026-03-09 21:46:30, 2, GetStateRight, Error sending request, rightRequestUrl: http://localhost:8082/stateFromAdjacentDp, count: 0, Get "http://localhost:8082/stateFromAdjacentDp": dial tcp 127.0.0.1:8082: connect: connection refused

dpNumber: 4

index: 4,  4, 2026-03-09 22:47:18, 1, DeQuiesce2, address: localhost:8083, left: localhost:8085, right: localhost:8084, 

index: 3,  4, 2026-03-09 22:47:18, 1, GetStateLeftAndReceiveStateRight,f.dpStates.state == Quiesce2, 

index: 2,  4, 2026-03-09 22:47:18, 1, Quiesce2, address: localhost:8083, left: localhost:8085, right: localhost:8084, 

index: 1,  4, 2026-03-09 22:43:56, 1, addNewDpToLeftOfTargetDp, address: localhost:8083, left: localhost:8085, right: localhost:8084, 

index: 0,  4, 2026-03-09 21:46:57, 2, GetStateRight, Error sending request, rightRequestUrl: http://localhost:8084/stateFromAdjacentDp, count: 0, Get "http://localhost:8084/stateFromAdjacentDp": dial tcp 127.0.0.1:8084: connect: connection refused

dpNumber: 5

index: 1,  5, 2026-03-09 22:43:56, 1, DeQuiesce2, address: localhost:8084, left: localhost:8083, right: localhost:8080, 

index: 0,  5, 2026-03-09 22:43:56, 1, Quiesce2, address: localhost:8084, left: localhost:8083, right: localhost:8080, 

dpNumber: 6

index: 2,  6, 2026-03-09 22:47:18, 1, DeQuiesce1Right, address: localhost:8085, left: localhost:8081, right: localhost:8083, 

index: 1,  6, 2026-03-09 22:47:18, 1, Quiesce1Left, address: localhost:8085, left: localhost:8082, right: localhost:8083, 

index: 0,  6, 2026-03-09 22:44:17, 1, DeQuiesce3, address: localhost:8085, left: localhost:8082, right: localhost:8083, 


One of the motivations behind a multi-process ring of dp is that one of the hosts might 

have a unique resource valuable  to dp running on other hosts.  


Here in this code there is a data heap on each dp that stores strings:

./dpClient relayStoreDataOnDp --addressOfDpToForwardRequest=localhost:8080 --numberOfDpThatStoresData=4 --dataToStore="The quick brown fox jumped over the lazy dog's back."

command relayStoreDataOnDp

	addressOfDpToForwardRequest : localhost:8080
	
	numberOfDpThatStoresData : 4
	
	dataToStore : The quick brown fox jumped over the lazy dog's back.

Forwarder Address: localhost:8080

DpNumber where stored: 4

ResultMessage: Data stored no problem.

StoreOrRetrieve: store

Data: 


Retrieve data:

./dpClient relayRetrieveDataFromDp --addressOfDpToForwardRequest=localhost:8080 --numberOfDpToRetrieveDataFrom=4

command relayRetrieveDataFromDp

	addressOfDpToForwardRequest : localhost:8080
	
	numberOfDpToRetrieveDataFrom : 4

Forwarder Address: localhost:8080

DpNumber where stored: 4
ResultMessage: Data retrieved

StoreOrRetrieve: retrieve

Data: The quick brown fox jumped over the lazy dog's back.


dpClient commands that begin with the word "relay" specify a dp that receives 
the request and forwards it through the ring.

dpClient commands that begin with the work "direct" contact a single dp.



