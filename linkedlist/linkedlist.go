package linkedlist

import (
	"github.com/rcrowley/go-metrics"
	"sync"
	"time"
)

type Node struct {
	data interface{}
	next *Node
}

type LinkedList struct {
	head   *Node
	length int
	mu     sync.Mutex
}

// Metrics to track operations
var (
	AddCounter    metrics.Counter
	RemoveCounter metrics.Counter
	FindCounter   metrics.Counter
	AddDuration   metrics.Timer
	RemoveDuration metrics.Timer
	FindDuration   metrics.Timer
)

// Initialize the metrics
func initMetrics() {
	AddCounter = metrics.NewCounter()
	RemoveCounter = metrics.NewCounter()
	FindCounter = metrics.NewCounter()

	AddDuration = metrics.NewTimer()
	RemoveDuration = metrics.NewTimer()
	FindDuration = metrics.NewTimer()

	metrics.Register("AddCounter", AddCounter)
	metrics.Register("RemoveCounter", RemoveCounter)
	metrics.Register("FindCounter", FindCounter)
	metrics.Register("AddDuration", AddDuration)
	metrics.Register("RemoveDuration", RemoveDuration)
	metrics.Register("FindDuration", FindDuration)
}

// NewLinkedList creates and returns a new linked list with initialized metrics
func NewLinkedList() *LinkedList {
	// Initialize the metrics only once
	initMetrics()

	// Create a new linked list and return its pointer
	return &LinkedList{}
}

// Add an element at the end of the list
func (list *LinkedList) Add(data interface{}) {
	start := time.Now()         // Track the start time for the Add operation

	list.mu.Lock()
	defer list.mu.Unlock()

	newNode := &Node{data: data}

	if list.head == nil {
		list.head = newNode
	} else {
		current := list.head
		for current.next != nil {
			current = current.next
		}
		current.next = newNode
	}

	list.length++
	AddCounter.Inc(1)
	AddDuration.UpdateSince(start) // Record the duration of the Add operation
}

// Remove an element from the list
func (list *LinkedList) Remove(data interface{}) {
	start := time.Now()         // Track the start time for the Remove operation

	list.mu.Lock()
	defer list.mu.Unlock()

	if list.head == nil {
		return
	}

	if list.head.data == data {
		list.head = list.head.next
		list.length--
		RemoveCounter.Inc(1)
		RemoveDuration.UpdateSince(start) // Record the duration of the Remove operation
		return
	}

	current := list.head
	for current.next != nil {
		if current.next.data == data {
			current.next = current.next.next
			list.length--
			RemoveCounter.Inc(1)
			RemoveDuration.UpdateSince(start) // Record the duration of the Remove operation
			return
		}
		current = current.next
	}
}

// Find an element in the list
func (list *LinkedList) Find(data interface{}) *Node {
	start := time.Now()         // Track the start time for the Find operation

	list.mu.Lock()
	defer list.mu.Unlock()

	current := list.head
	for current != nil {
		if current.data == data {
			FindCounter.Inc(1)
			FindDuration.UpdateSince(start) // Record the duration of the Find operation
			return current
		}
		current = current.next
	}
	return nil
}

