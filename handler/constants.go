package handler

// State Quiesce3 means that the new philosopher will not leave the select loop
// until the dp that received the add dp request sends the dequiesce request
const (
	Thinking = iota
	Hungry
	Dining
	Quiesce1
	Quiesce2
	Quiesce3
)

const RingBufferCapacity = 40
