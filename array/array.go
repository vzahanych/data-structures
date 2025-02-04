package array

import (
    "errors"
    "sync"
    "github.com/rcrowley/go-metrics"
)

type ArrayConfig struct {
    MetricsEnabled bool
}

// Array is a generic array structure that can hold elements of any type
type Array[T any] struct {
    data           []T
    size           int
    mu             sync.RWMutex
    config         ArrayConfig
    appendCounter  metrics.Counter
    getCounter     metrics.Counter
    resizeCounter  metrics.Counter
}

// NewArray creates a new generic array with a fixed capacity
// The config parameter is used to enable or disable metrics collection
func NewArray[T any](capacity int, config ArrayConfig) *Array[T] {
    arr := &Array[T]{
        data: make([]T, capacity),
        size: 0,
        config: config,
    }

    // Initialize the metrics only if enabled in the config
    if arr.config.MetricsEnabled {
        arr.appendCounter = metrics.NewCounter()
        arr.getCounter = metrics.NewCounter()
        arr.resizeCounter = metrics.NewCounter()
        metrics.DefaultRegistry.Register("array.append", arr.appendCounter)
        metrics.DefaultRegistry.Register("array.get", arr.getCounter)
        metrics.DefaultRegistry.Register("array.resize", arr.resizeCounter)
    }

    return arr
}

// Append adds a new element to the array
func (a *Array[T]) Append(value T) error {
    a.mu.Lock() // Lock for writing
    defer a.mu.Unlock()

    if a.size >= len(a.data) {
        return errors.New("array is full")
    }
    a.data[a.size] = value
    a.size++

    // Track metrics if enabled
    if a.config.MetricsEnabled {
        a.appendCounter.Inc(1)
    }
    return nil
}

// Get retrieves an element at the specified index
func (a *Array[T]) Get(index int) (T, error) {
    a.mu.RLock() // Lock for reading
    defer a.mu.RUnlock()

    if index < 0 || index >= a.size {
        var zero T // Return a zero value of type T
        return zero, errors.New("index out of bounds")
    }

    // Track metrics if enabled
    if a.config.MetricsEnabled {
        a.getCounter.Inc(1)
    }
    return a.data[index], nil
}

// Length returns the current number of elements in the array
func (a *Array[T]) Length() int {
    a.mu.RLock() // Lock for reading
    defer a.mu.RUnlock()

    return a.size
}

// Resize resizes the array to a new capacity
func (a *Array[T]) Resize(newCapacity int) error {
    a.mu.Lock() // Lock for writing
    defer a.mu.Unlock()

    if newCapacity < a.size {
        return errors.New("new capacity must be greater than or equal to the current size")
    }

    newData := make([]T, newCapacity)
    copy(newData, a.data)
    a.data = newData

    // Track metrics if enabled
    if a.config.MetricsEnabled {
        a.resizeCounter.Inc(1)
    }

    return nil
}

// Delete removes the element at the specified index and shifts subsequent elements
func (a *Array[T]) Delete(index int) error {
    a.mu.Lock() // Lock for writing
    defer a.mu.Unlock()

    if index < 0 || index >= a.size {
        return errors.New("index out of bounds")
    }

    // Shift elements to the left to fill the gap
    copy(a.data[index:], a.data[index+1:a.size])
    
    // Set the last element to zero and reduce size
    a.size--
    var zero T // Zero value for type T
    a.data[a.size] = zero // Optional: clear the last element

    // Track metrics if enabled
    if a.config.MetricsEnabled {
        a.resizeCounter.Inc(1)
    }

    return nil
}
