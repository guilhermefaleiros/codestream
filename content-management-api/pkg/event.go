package pkg

import "time"

type Event interface {
	GetPayload() interface{}
	GetType() string
	GetID() string
	GetCreatedAt() time.Time
}
