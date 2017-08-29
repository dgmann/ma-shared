package stats

import (
	"testing"
	"github.com/dgmann/ma-shared/queue"
	"time"
)

func TestPrint(t *testing.T) {
	msg := queue.NewMessage(nil, 1, time.Now())
	msg.EnterStage("Stage 1")
	time.Sleep(20 * time.Millisecond)
	msg.LeaveStage("Stage 1")
	stat := NewStat(msg)
	if int(stat.Stages[0].ProcessingTime.Seconds() * 1000) != 20 {
		t.Error("Expected 20ms, got", int(stat.Stages[0].ProcessingTime.Seconds() * 1000))
	}
	stat.Print()
}
