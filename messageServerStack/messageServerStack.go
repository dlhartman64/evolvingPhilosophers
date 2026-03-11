package messageServerStack

import (
	"fmt"
)

type ClientMessage struct {
	OriginatorAddress string `json:"originatoraddress"`
	Resource          string `json:"resource"`
	ResourceDpNumber  string `json:"dpresourcenumber"`
	ResultMessage     string `json:"resultmessage,omitempty"`
	Done              string `json:"done"`
	StoreOrRetrieve   string `json:"storeorretreive"`
	Data              string `json:"data"`
}

type MessageServerStack struct {
	elements []ClientMessage
}

func NewMessageServerStack() *MessageServerStack {
	e := make([]ClientMessage, 10)
	return &MessageServerStack{elements: e}
}

func (s *MessageServerStack) Push(value ClientMessage) {
	s.elements = append(s.elements, value)
}

func (s *MessageServerStack) Pop() (*ClientMessage, error) {
	if s.IsEmpty() {
		return nil, fmt.Errorf("stack is empty")
	}
	topElement := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return &topElement, nil
}

func (s *MessageServerStack) IsEmpty() bool {
	return len(s.elements) == 0
}
