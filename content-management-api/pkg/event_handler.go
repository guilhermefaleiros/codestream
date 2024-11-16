package pkg

import (
	"context"
	"sync"
)

type EventHandler interface {
	Handle(ctx context.Context, event Event, wg *sync.WaitGroup)
}
