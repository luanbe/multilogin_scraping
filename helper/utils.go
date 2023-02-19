package helper

import (
	"bytes"
	"encoding/json"
	"log"
)

type Message map[string]interface {
}

type Utils interface {
	Serialize(msg Message) ([]byte, error)
	Deserialize(b []byte) (Message, error)
	failOnError(err error, msg string)
}
type UtilHelper struct {
}

func NewUtils() Utils {
	return &UtilHelper{}
}

func (u *UtilHelper) failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (u *UtilHelper) Serialize(msg Message) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(msg)
	return b.Bytes(), err
}

func (u *UtilHelper) Deserialize(b []byte) (Message, error) {
	var msg Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
