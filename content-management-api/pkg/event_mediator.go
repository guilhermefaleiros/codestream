package pkg

import (
	"context"
	"errors"
	"sync"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

type EventMediator struct {
	handlers map[string][]EventHandler
}

func (em *EventMediator) Register(eventType string, handler EventHandler) error {
	if _, ok := em.handlers[eventType]; ok {
		for _, h := range em.handlers[eventType] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	em.handlers[eventType] = append(em.handlers[eventType], handler)
	return nil
}

func (em *EventMediator) Dispatch(ctx context.Context, event Event) error {
	if handlers, ok := em.handlers[event.GetType()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(ctx, event, wg)
		}
		wg.Wait()
	}
	return nil
}

func (em *EventMediator) Unregister(eventType string, handler EventHandler) error {
	if _, ok := em.handlers[eventType]; ok {
		for i, h := range em.handlers[eventType] {
			if h == handler {
				em.handlers[eventType] = append(em.handlers[eventType][:i], em.handlers[eventType][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func NewEventMediator() *EventMediator {
	return &EventMediator{
		handlers: make(map[string][]EventHandler),
	}
}
