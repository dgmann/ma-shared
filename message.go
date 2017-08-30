package shared

import (
	"time"
	"encoding/json"
)

func(message *Message) EnterStage(stageName string) {
	message.Stages[stageName] = Stage{
		EnteredAt: time.Now(),
	}
}

func(message *Message) LeaveStage(stageName string) {
	message.Stages[stageName] = Stage{
		EnteredAt: message.Stages[stageName].EnteredAt,
		LeftAt: time.Now(),
	}
}

type Message struct {
	Image []byte `json:"image"`
	FrameNumber int `json:"frameNumber"`
	CreatedAt time.Time `json:"createdAt"`
	Stages map[string]Stage `json:"timeline"`
}

type Stage struct {
	EnteredAt time.Time `json:"enteredAt"`
	LeftAt time.Time `json:"leftAt"`
}

func NewMessage(image []byte, frameNumer int, createdAt time.Time) (*Message) {
	return & Message{
		Image: image,
		FrameNumber: frameNumer,
		CreatedAt: createdAt,
		Stages: make(map[string]Stage),
	}
}

func NewMessageFromJSON(b []byte) (*Message) {
	var msg Message
	json.Unmarshal(b, &msg)
	return &msg
}

func(message *Message) ToJSON() ([]byte, error) {
	return json.Marshal(message)
}
