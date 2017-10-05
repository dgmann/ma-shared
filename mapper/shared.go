package mapper

import (
	"github.com/dgmann/ma-shared"
	"github.com/jeffail/tunny"
)

type RunnerMode interface {
	Listen(endpoint string, setResult func(message *shared.Message, result interface{}))
}

type Mode struct {
	StageName string
	pool *tunny.WorkPool
}
