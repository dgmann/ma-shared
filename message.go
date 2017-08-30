package shared

import (
	"time"
	"encoding/json"
	"github.com/dgmann/ma-shared/openalpr"
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
	Origin string `json:"origin"`
	Image []byte `json:"image"`
	FrameNumber int `json:"frameNumber"`
	CreatedAt time.Time `json:"createdAt"`
	Stages map[string]Stage `json:"timeline"`
	Results Results `json:"results"`
}

type Stage struct {
	EnteredAt time.Time `json:"enteredAt"`
	LeftAt time.Time `json:"leftAt"`
}

type Results struct {
	OpenALPR openalpr.OpenAlprResponse `json:"openalpr"`
}

func NewMessage(image []byte, frameNumer int, createdAt time.Time) (*Message) {
	return & Message{
		Origin: "",
		Image: image,
		FrameNumber: frameNumer,
		CreatedAt: createdAt,
		Stages: make(map[string]Stage),
		Results: Results{},
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
