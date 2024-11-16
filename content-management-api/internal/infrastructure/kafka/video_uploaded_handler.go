package kafka

import (
	"context"
	"fmt"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/event"
	"github.com/guilhermefaleiros/codestream/content-management-system/pkg"
	"sync"
)

type VideoUploadedHandler struct {
	producer *Producer
	topic    string
}

func (v *VideoUploadedHandler) Handle(ctx context.Context, ev pkg.Event, wg *sync.WaitGroup) {
	defer wg.Done()
	videoUploadedPayload := ev.GetPayload().(event.VideoUploadedEventPayload)
	err := v.producer.Publish(videoUploadedPayload, nil, v.topic)
	if err != nil {
		fmt.Printf("Error on dispatch event to kafka: %s\n", err.Error())
		return
	}
	fmt.Printf("Video uploaded: %s\n", videoUploadedPayload.ID)
}

func NewVideoUploadedHandler(producer *Producer, topic string) *VideoUploadedHandler {
	return &VideoUploadedHandler{
		producer: producer,
		topic:    topic,
	}
}
