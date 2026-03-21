

The Evolving Philosophers problem makes it posible to add new dining philosophers to the ring without
affecting entry into the critical section of each dining philosoper.  Each dining
philosopher here is a process, and the communication is REST.

The design is based on the paper "The Evolving Philosophers Problem:
Dynamic Change Management" by Jeff Kramer and Jeff Magee.  4 dining philosophers are put in quiesent states, and a new diningPhilosopher is added with 2 diningPhilosophera on each
side. Similar with removal.  Dijkstra's algorithm is used.

dpClient binary is used to send commands to a diningPhilosopher.

Build the binaries:
go build -ldflags '-extldflags "-static"' diningPhilosopher.go

go build -ldflags '-extldflags "-static"' dpClient.go

Enter ./diningPhilosopher or ./dpClient to get a list of commands.

Enter ./diningPhilosopher <command> or ./dpClient <command> to get the options for the commands


Start 5 diningPhilosopher:

./diningPhilosopher initialRing --addressOfDp=localhost:8080 --addressOfDpOnLeft=localhost:8084 --addressOfDpOnRight=localhost:8081 --dpNumber=1 --debugStdout=true

./diningPhilosopher initialRing --addressOfDp=localhost:8081 --addressOfDpOnLeft=localhost:8080 --addressOfDpOnRight=localhost:8082 --dpNumber=2 --debugStdout=true

./diningPhilosopher initialRing --addressOfDp=localhost:8082 --addressOfDpOnLeft=localhost:8081 --addressOfDpOnRight=localhost:8083 --dpNumber=3 --debugStdout=true

./diningPhilosopher initialRing --addressOfDp=localhost:8083 --addressOfDpOnLeft=localhost:8082 --addressOfDpOnRight=localhost:8084 --dpNumber=4 --debugStdout=true

./diningPhilosopher initialRing --addressOfDp=localhost:8084 --addressOfDpOnLeft=localhost:8083 --addressOfDpOnRight=localhost:8080 --dpNumber=5 --debugStdout=true

Each diningPhilosopher will enter the critical section after all are started.


View each diningPhilosopher in the ring:

./dpClient relayAttributes --dpStartAddress=localhost:8080

dp address, dp number, sequence, left address, right address, iteration

localhost:8080     dpNumber: 1   S: 0   L: localhost:8084   R: localhost:8081   Iter: 2279

localhost:8084     dpNumber: 5   S: 1   L: localhost:8083   R: localhost:8080   Iter: 1331

localhost:8083     dpNumber: 4   S: 2   L: localhost:8082   R: localhost:8084   Iter: 1330

localhost:8082     dpNumber: 3   S: 3   L: localhost:8081   R: localhost:8083   Iter: 1330

localhost:8081     dpNumber: 2   S: 4   L: localhost:8080   R: localhost:8082   Iter: 2280



Start a new dp:
./diningPhilosopher addToRing --addressOfDp=localhost:8085 --dpNumber=6 --debugStdout=true

Add to ring:
./dpClient relayAddNewDpToLeftOfTarget --addressOfDpToForwardRequest=localhost:8080 --numberOfDpToForwardRequest=1 --numberOfTargetDp=4 --addressOfNewDp=localhost:8085

Return status for relayAddNewDpToLeftOfTarget: 200

OriginatorDpNumber: 1

TargetDpNumber: 4

NewDpAddress: localhost:8085

Done: true

Result: Added new DP successfully


View diningPhilosopher in ring:

./dpClient relayAttributes --dpStartAddress=localhost:8080

dp address, dp number, sequence, left address, right address, iteration

localhost:8080     dpNumber: 1   S: 0   L: localhost:8084   R: localhost:8081   Iter: 10066

localhost:8084     dpNumber: 5   S: 1   L: localhost:8083   R: localhost:8080   Iter: 9780

localhost:8083     dpNumber: 4   S: 2   L: localhost:8085   R: localhost:8084   Iter: 9778

localhost:8085     dpNumber: 6   S: 3   L: localhost:8082   R: localhost:8083   Iter: 1152

localhost:8082     dpNumber: 3   S: 4   L: localhost:8081   R: localhost:8085   Iter: 9778

localhost:8081     dpNumber: 2   S: 5   L: localhost:8080   R: localhost:8082   Iter: 10063



Remove a diningPhilosopher:

./dpClient relayRemoveDp --addressOfDpToForwardRequest=localhost:8080 --numberOfDpToForwardRequest=1 --numberOfDpToRemove=3

Return status for relayRemoveDp: 200

Forwarder DpNumber: 1

DpNumberToBeRemoved: 3

LeftAddress of Forwarder: localhost:8084

Done: true

Result: Removed dp successfully

View dps in ring:

./dpClient relayAttributes --dpStartAddress=localhost:8080 

dp address, dp number, sequence, left address, right address, iteration

localhost:8080     dpNumber: 1   S: 0   L: localhost:8084   R: localhost:8081   Iter: 11986

localhost:8084     dpNumber: 5   S: 1   L: localhost:8083   R: localhost:8080   Iter: 11703

localhost:8083     dpNumber: 4   S: 2   L: localhost:8085   R: localhost:8084   Iter: 11699

localhost:8085     dpNumber: 6   S: 3   L: localhost:8081   R: localhost:8083   Iter: 3073

localhost:8081     dpNumber: 2   S: 4   L: localhost:8080   R: localhost:8085   Iter: 11983



Each diningPhilosopher has a ring buffer to store log messages.

./dpClient relayRequestLogEntries --addressOfDpToForwardRequest=localhost:8080

dpNumber: 1

index: 6,  1, 2026-03-16 01:05:54, C, relayRequestToAddNewDpToLeftOfTargetDp, There is no dp with the dpNumber 3 in the ring, 

index: 5,  1, 2026-03-16 01:05:54, C, relayRequestToAddNewDpToLeftOfTargetDp, There is no dp with the dpNumber 3 in the ring, 

index: 4,  1, 2026-03-16 01:05:54, I, DeQuiesce2, address: localhost:8080, left: localhost:8084, right: localhost:8081, 

index: 3,  1, 2026-03-16 01:05:54, I, GetStateRightAndReceiveStateLeft,f.dpStates.state == Quiesce2, 

index: 2,  1, 2026-03-16 01:05:54, I, GetStateRightAndReceiveStateLeft,f.dpStates.state == Quiesce2, 

index: 1,  1, 2026-03-16 01:05:54, I, Quiesce2, address: localhost:8080, left: localhost:8084, right: localhost:8081, 

index: 0,  1, 2026-03-16 00:58:24, W, getStateLeft, Error sending request, leftRequestUrl: http://localhost:8084/stateFromAdjacentDp, count: 0, Get "http://localhost:8084/stateFromAdjacentDp": dial tcp 127.0.0.1:8084: connect: connection refused

dpNumber: 2

index: 6,  2, 2026-03-16 01:05:54, I, DeQuiesce1Left, address: localhost:8081, left: localhost:8080, right: localhost:8085, 

index: 5,  2, 2026-03-16 01:05:54, I, Quiesce1Left, address: localhost:8081, left: localhost:8080, right: localhost:8082, 

index: 4,  2, 2026-03-16 01:04:32, I, DeQuiesce2, address: localhost:8081, left: localhost:8080, right: localhost:8082, 

index: 3,  2, 2026-03-16 01:04:32, I, GetStateRightAndReceiveStateLeft,f.dpStates.state == Quiesce2, 

index: 2,  2, 2026-03-16 01:04:32, I, GetStateRightAndReceiveStateLeft,f.dpStates.state == Quiesce2, 

index: 1,  2, 2026-03-16 01:04:32, I, Quiesce2, address: localhost:8081, left: localhost:8080, right: localhost:8082, 

index: 0,  2, 2026-03-16 00:58:29, W, GetStateRight, Error sending request, rightRequestUrl: http://localhost:8082/stateFromAdjacentDp, count: 0, Get "http://localhost:8082/stateFromAdjacentDp": dial tcp 127.0.0.1:8082: connect: connection refused

dpNumber: 4

index: 4,  4, 2026-03-16 01:05:54, I, DeQuiesce2, address: localhost:8083, left: localhost:8085, right: localhost:8084, 

index: 3,  4, 2026-03-16 01:05:54, I, GetStateLeftAndReceiveStateRight,f.dpStates.state == Quiesce2, 

index: 2,  4, 2026-03-16 01:05:54, I, Quiesce2, address: localhost:8083, left: localhost:8085, right: localhost:8084, 

index: 1,  4, 2026-03-16 01:04:32, I, addNewDpToLeftOfTargetDp, address: localhost:8083, left: localhost:8085, right: localhost:8084, 

index: 0,  4, 2026-03-16 00:58:39, W, GetStateRight, Error sending request, rightRequestUrl: http://localhost:8084/stateFromAdjacentDp, count: 0, Get "http://localhost:8084/stateFromAdjacentDp": dial tcp 127.0.0.1:8084: connect: connection refused

dpNumber: 5

index: 2,  5, 2026-03-16 01:04:32, I, DeQuiesce2, address: localhost:8084, left: localhost:8083, right: localhost:8080, 

index: 1,  5, 2026-03-16 01:04:32, I, GetStateLeftAndReceiveStateRight,f.dpStates.state == Quiesce2, 

index: 0,  5, 2026-03-16 01:04:32, I, Quiesce2, address: localhost:8084, left: localhost:8083, right: localhost:8080, 

dpNumber: 6

index: 2,  6, 2026-03-16 01:05:54, I, DeQuiesce1Right, address: localhost:8085, left: localhost:8081, right: localhost:8083, 

index: 1,  6, 2026-03-16 01:05:54, I, Quiesce1Left, address: localhost:8085, left: localhost:8082, right: localhost:8083, 

index: 0,  6, 2026-03-16 01:04:32, I, DeQuiesce3, address: localhost:8085, left: localhost:8082, right: localhost:8083, 


One of the motivations behind a multi-process ring of dp is that one of the hosts might 
have a unique resource valuable  to dp running on other hosts.  
There is a data heap on each diningPhilosopher to illustrate this that stores strings:


./dpClient relayStoreDataOnDp --addressOfDpToForwardRequest=localhost:8080 --numberOfDpThatStoresData=4 --dataToStore="The quick brown fox jumped over the lazy dog's back."

Forwarder Address: localhost:8080

DpNumber where stored: 4

ResultMessage: Data stored no problem.

StoreOrRetrieve: store
Data: 


Retrieve data:

./dpClient relayRetrieveDataFromDp --addressOfDpToForwardRequest=localhost:8080 --numberOfDpToRetrieveDataFrom=4

Forwarder Address: localhost:8080

DpNumber where stored: 4

ResultMessage: Data retrieved

StoreOrRetrieve: retrieve

Data: The quick brown fox jumped over the lazy dog's back.


dpClient commands that begin with the word "relay" specify a dp that receives 
the request and forwards it through the ring.

dpClient commands that begin with the work "direct" can connect and send a request, to any diningPhilosopher in the ring, and receive the reply from that diningPhilosopher.



