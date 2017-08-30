package shared

import (
	"time"
	"encoding/json"
	"github.com/dgmann/ma-shared/sampler"
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

func(message *Message) AddStage(stageName string, enteredAt, leftAt time.Time) {
	message.Stages[stageName] = Stage{
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
	Results Results `json:"results"`
}

type Stage struct {
	EnteredAt time.Time `json:"enteredAt"`
	LeftAt time.Time `json:"leftAt"`
}

type Results struct {
	OpenALPR OpenAlprResponse `json:"openalpr"`
}

func NewMessage(image []byte, frameNumer int, readAt, createdAt time.Time) (*Message) {
	msg := Message{
		Origin: "",
		Image: image,
		FrameNumber: frameNumer,
		CreatedAt: createdAt,
		Stages: make(map[string]Stage),
		Results: Results{},
	}
	msg.AddStage("Recorded", readAt, createdAt)
	return &msg
}

func NewMessageFromSample(sample sampler.VideoSample) (*Message) {
	msg := NewMessage(nil, sample.FrameNumber, sample.ReadPacketAt, sample.CreatedAt)
	msg.EnterStage("Encode")
	img := sample.ToJPEG()
	msg.Image = img.Bytes()
	msg.LeaveStage("Encode")
	return msg
}

func NewMessageFromJSON(b []byte) (*Message) {
	var msg Message
	json.Unmarshal(b, &msg)
	return &msg
}

func(message *Message) ToJSON() ([]byte, error) {
	return json.Marshal(message)
}
