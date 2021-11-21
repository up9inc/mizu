package kubernetes

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	eventsv1 "k8s.io/api/events/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type WatchEvent watch.Event

func (we *WatchEvent) ToPod() (*corev1.Pod, error) {
	pod, ok := we.Object.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("Invalid object type on pod event stream")
	}

	return pod, nil
}

func (we *WatchEvent) ToEvent() (*eventsv1.Event, error) {
	event, ok := we.Object.(*eventsv1.Event)
	if !ok {
		return nil, fmt.Errorf("Invalid object type on pod event stream")
	}

	return event, nil
}
