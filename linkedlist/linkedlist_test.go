package linkedlist

import (
	"sync"
	"testing"
)

// Test NewLinkedList initializes a new linked list correctly
func TestNewLinkedList(t *testing.T) {
	list := NewLinkedList()

	if list == nil {
		t.Fatalf("Expected NewLinkedList() to return a non-nil linked list")
	}

	if list.head != nil {
		t.Fatalf("Expected linked list head to be nil initially")
	}

	if list.length != 0 {
		t.Fatalf("Expected linked list length to be 0 initially, got %d", list.length)
	}
}

// Test Add function adds elements to the linked list
func TestAdd(t *testing.T) {
	list := NewLinkedList()

	list.Add(10)
	list.Add(20)
	list.Add(30)

	// Test the size of the list after adding elements
	if list.length != 3 {
		t.Fatalf("Expected linked list length to be 3, got %d", list.length)
	}

	// Test the head of the list
	if list.head.data != 10 {
		t.Fatalf("Expected first element to be 10, got %v", list.head.data)
	}

	// Test the second element in the list
	if list.head.next.data != 20 {
		t.Fatalf("Expected second element to be 20, got %v", list.head.next.data)
	}
}

// Test Remove function removes elements from the linked list
func TestRemove(t *testing.T) {
	list := NewLinkedList()

	list.Add(10)
	list.Add(20)
	list.Add(30)

	// Remove an element from the list
	list.Remove(20)

	// Test the length after removal
	if list.length != 2 {
		t.Fatalf("Expected linked list length to be 2 after removal, got %d", list.length)
	}

	// Test the second element in the list
	if list.head.next.data != 30 {
		t.Fatalf("Expected second element to be 30 after removal, got %v", list.head.next.data)
	}

	// Try removing a non-existing element (no-op)
	list.Remove(40)

	// Length should still be 2
	if list.length != 2 {
		t.Fatalf("Expected linked list length to remain 2, got %d", list.length)
	}
}

// Test Find function finds an element in the linked list
func TestFind(t *testing.T) {
	list := NewLinkedList()

	list.Add(10)
	list.Add(20)
	list.Add(30)

	// Test finding an existing element
	node := list.Find(20)
	if node == nil || node.data != 20 {
		t.Fatalf("Expected to find node with data 20, but got %v", node)
	}

	// Test finding a non-existing element
	node = list.Find(40)
	if node != nil {
		t.Fatalf("Expected to not find node with data 40, but got %v", node)
	}
}

// Test metrics are updated correctly
func TestMetrics(t *testing.T) {
	// Create a new linked list and add some elements
	list := NewLinkedList()

	// Add elements and check the metrics
	list.Add(10)
	list.Add(20)
	list.Add(30)

	// Check AddCounter
	if AddCounter.Count() != 3 {
		t.Fatalf("Expected AddCounter to be 3, but got %d", AddCounter.Count())
	}

	// Check AddDuration
	if AddDuration.Count() != 3 {
		t.Fatalf("Expected AddDuration to be 3, but got %d", AddDuration.Count())
	}

	// Remove elements and check the metrics
	list.Remove(20)

	// Check RemoveCounter
	if RemoveCounter.Count() != 1 {
		t.Fatalf("Expected RemoveCounter to be 1, but got %d", RemoveCounter.Count())
	}

	// Check RemoveDuration
	if RemoveDuration.Count() != 1 {
		t.Fatalf("Expected RemoveDuration to be 1, but got %d", RemoveDuration.Count())
	}

	// Find elements and check the metrics
	list.Find(10)
	list.Find(30)

	// Check FindCounter
	if FindCounter.Count() != 2 {
		t.Fatalf("Expected FindCounter to be 2, but got %d", FindCounter.Count())
	}

	// Check FindDuration
	if FindDuration.Count() != 2 {
		t.Fatalf("Expected FindDuration to be 2, but got %d", FindDuration.Count())
	}
}

func BenchmarkAdd(b *testing.B) {
	list := NewLinkedList()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Add(i)
	}
}

func BenchmarkRemove(b *testing.B) {
	list := NewLinkedList()

	// Pre-populate the list with 1000 elements
	for i := 0; i < 1000; i++ {
		list.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Remove(i % 1000)
	}
}

func BenchmarkFind(b *testing.B) {
	list := NewLinkedList()

	// Pre-populate the list with 1000 elements
	for i := 0; i < 1000; i++ {
		list.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Find(i % 1000)
	}
}

// Concurrent benchmark for Add operation
func BenchmarkConcurrentAdd(b *testing.B) {
	list := NewLinkedList()
	var wg sync.WaitGroup

	// Number of concurrent goroutines to add elements
	concurrency := 10
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		wg.Add(concurrency)
		for j := 0; j < concurrency; j++ {
			go func(j int) {
				defer wg.Done()
				list.Add(j)
			}(j)
		}
		wg.Wait()
	}
}

// Concurrent benchmark for Remove operation
func BenchmarkConcurrentRemove(b *testing.B) {
	list := NewLinkedList()
	var wg sync.WaitGroup

	// Pre-populate the list with 1000 elements
	for i := 0; i < 1000; i++ {
		list.Add(i)
	}

	// Number of concurrent goroutines to remove elements
	concurrency := 10
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		wg.Add(concurrency)
		for j := 0; j < concurrency; j++ {
			go func(j int) {
				defer wg.Done()
				list.Remove(j)
			}(j)
		}
		wg.Wait()
	}
}

// Concurrent benchmark for Find operation
func BenchmarkConcurrentFind(b *testing.B) {
	list := NewLinkedList()
	var wg sync.WaitGroup

	// Pre-populate the list with 1000 elements
	for i := 0; i < 1000; i++ {
		list.Add(i)
	}

	// Number of concurrent goroutines to find elements
	concurrency := 10
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		wg.Add(concurrency)
		for j := 0; j < concurrency; j++ {
			go func(j int) {
				defer wg.Done()
				list.Find(j)
			}(j)
		}
		wg.Wait()
	}
}
