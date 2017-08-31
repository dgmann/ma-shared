package stats

import (
	"testing"
	"time"
	"github.com/dgmann/ma-shared"
	"bytes"
)

func TestCSV(t *testing.T) {
	msg := shared.NewMessage(nil, 1, time.Now(), time.Now())
	msg.EnterStage("Stage 1")
	time.Sleep(20 * time.Millisecond)
	msg.LeaveStage("Stage 1")

	msg.EnterStage("Stage 2")
	time.Sleep(16 * time.Millisecond)
	msg.LeaveStage("Stage 2")

	stat := NewStat(msg)
	collector := NewCollector()
	collector.Add(stat)

	var b bytes.Buffer
	collector.ToCSV(&b)
	println("Output: " + b.String())
}
