package shared

import (
	"time"
	"encoding/json"
)

type Stage struct {
	Index int `json:"index"`
	EnteredAt time.Time `json:"enteredAt"`
	LeftAt time.Time `json:"leftAt"`
}

func(message *Message) EnterStage(stageName string) {
	message.Stages[stageName] = Stage{
		Index: len(message.Stages),
		EnteredAt: time.Now(),
	}
}

func(message *Message) LeaveStage(stageName string) {
	message.Stages[stageName] = Stage{
		Index: message.Stages[stageName].Index,
		EnteredAt: message.Stages[stageName].EnteredAt,
		LeftAt: time.Now(),
	}
}

func(message *Message) AddStage(stageName string, enteredAt, leftAt time.Time) {
	message.Stages[stageName] = Stage{
		Index: len(message.Stages),
		EnteredAt: enteredAt,
		LeftAt: leftAt,
	}
}

type Message struct {
	Origin string `json:"origin"`
	Image []byte `json:"image"`
	FrameNumber int `json:"frameNumber"`
	CreatedAt time.Time `json:"createdAt"`
	Stages map[string]Stage `json:"stages"`
	Result Result `json:"result"`
}

type Result struct {
	OpenALPR OpenAlprResponse `json:"openalpr"`
	IsWanted bool `json:"isWanted"`
}

func NewMessage(image []byte, frameNumer int, readAt, createdAt time.Time) (*Message) {
	msg := Message{
		Origin: "",
		Image: image,
		FrameNumber: frameNumer,
		CreatedAt: createdAt,
		Stages: make(map[string]Stage),
		Result: Result{},
	}
	msg.AddStage("Decode", readAt, createdAt)
	return &msg
}

func NewMessageFromSample(sample VideoSample) (*Message, error) {
	msg := NewMessage(nil, sample.FrameNumber, sample.ReadPacketAt, sample.CreatedAt)
	msg.EnterStage("Encode")
	img, err := sample.ToJPEG()
	if err != nil {
		return nil, err
	}
	msg.Image = img.Bytes()
	msg.LeaveStage("Encode")
	return msg, err
}

func NewMessageFromJSON(b []byte) (*Message) {
	var msg Message
	json.Unmarshal(b, &msg)
	return &msg
}

func(message *Message) ToJSON() ([]byte, error) {
	return json.Marshal(message)
}