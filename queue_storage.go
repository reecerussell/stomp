package stomp

import (
	"container/list"
	"github.com/jjeffery/stomp/message"
)

// Interface for queue storage. The intent is that
// different queue storage implementations can be
// used, depending on preference. Queue storage
// mechanisms could include in-memory, and various
// persistent storage mechanisms (eg file system, DB, etc)
type QueueStorage interface {
	// Pushes a MESSAGE frame to the end of the queue. Sets
	// the "message-id" header of the frame before adding to
	// the queue.
	Enqueue(queue string, frame *message.Frame) error

	// Pushes a MESSAGE frame to the head of the queue. Sets
	// the "message-id" header of the frame if it is not
	// already set.
	Requeue(queue string, frame *message.Frame) error

	// Removes a frame from the head of the queue.
	// Returns nil if no frame is available.
	Dequeue(queue string) (*message.Frame, error)

	// Called at server startup. Allows the queue storage
	// to perform any initialization.
	Start()

	// Called prior to server shutdown. Allows the queue storage
	// to perform any cleanup.
	Stop()
}

type MemoryQueueStorage struct {
	lists map[string]*list.List
}

func NewMemoryQueueStorage() QueueStorage {
	m := new(MemoryQueueStorage)
	return m
}

func (m *MemoryQueueStorage) Enqueue(queue string, frame *message.Frame) error {
	l, ok := m.lists[queue]
	if !ok {
		l = list.New()
		m.lists[queue] = l
	}
	l.PushBack(frame)

	return nil
}

// Pushes a frame to the head of the queue. Sets
// the "message-id" header of the frame if it is not
// already set.
func (m *MemoryQueueStorage) Requeue(queue string, frame *message.Frame) error {
	l, ok := m.lists[queue]
	if !ok {
		l = list.New()
		m.lists[queue] = l
	}
	l.PushFront(frame)

	return nil
}

// Removes a frame from the head of the queue.
// Returns nil if no frame is available.
func (m *MemoryQueueStorage) Dequeue(queue string) (*message.Frame, error) {
	l, ok := m.lists[queue]
	if !ok {
		return nil, nil
	}

	element := l.Front()
	if element == nil {
		return nil, nil
	}

	return l.Remove(element).(*message.Frame), nil
}

// Called at server startup. Allows the queue storage
// to perform any initialization.
func (m *MemoryQueueStorage) Start() {
	m.lists = make(map[string]*list.List)
}

// Called prior to server shutdown. Allows the queue storage
// to perform any cleanup.
func (m *MemoryQueueStorage) Stop() {
	m.lists = nil
}