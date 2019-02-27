package model

import "fmt"

//MessageOK represents a successful http call
type MessageOK struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

func (message *MessageOK) String() string {
	if message.Ok == true {
		return "true"
	}
	return fmt.Sprintf("false, %s", message.Message)
}

//Ok returns a message object with the set state
func Ok(state bool) MessageOK {
	return MessageOK{Ok: state}
}

//OkMessage returns a message object with the set state
func OkMessage(state bool, message string) MessageOK {
	return MessageOK{Ok: state, Message: message}
}
