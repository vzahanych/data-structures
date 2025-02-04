package array

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppend(t *testing.T) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](5, config)

    // Append values to the array
    err := arr.Append(10)
    assert.NoError(t, err)
    err = arr.Append(20)
    assert.NoError(t, err)
    err = arr.Append(30)
    assert.NoError(t, err)

    // Ensure the array has the correct values
    value, err := arr.Get(0)
    assert.NoError(t, err)
    assert.Equal(t, 10, value)

    value, err = arr.Get(1)
    assert.NoError(t, err)
    assert.Equal(t, 20, value)

    value, err = arr.Get(2)
    assert.NoError(t, err)
    assert.Equal(t, 30, value)
}

func TestAppendFullArray(t *testing.T) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](3, config)

    // Append values to the array
    err := arr.Append(10)
    assert.NoError(t, err)
    err = arr.Append(20)
    assert.NoError(t, err)
    err = arr.Append(30)
    assert.NoError(t, err)

    // Attempt to append to a full array
    err = arr.Append(40)
    assert.Error(t, err)
    assert.Equal(t, "array is full", err.Error())
}

func TestGet(t *testing.T) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](5, config)

    // Append some values
    arr.Append(10)
    arr.Append(20)
    arr.Append(30)

    // Get the values
    value, err := arr.Get(0)
    assert.NoError(t, err)
    assert.Equal(t, 10, value)

    value, err = arr.Get(1)
    assert.NoError(t, err)
    assert.Equal(t, 20, value)

    value, err = arr.Get(2)
    assert.NoError(t, err)
    assert.Equal(t, 30, value)

    // Attempt to get an out-of-bounds index
    value, err = arr.Get(5)
    assert.Error(t, err)
    assert.Equal(t, "index out of bounds", err.Error())
}

func TestLength(t *testing.T) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](5, config)

    // Append some values
    arr.Append(10)
    arr.Append(20)
    arr.Append(30)

    // Check the length
    assert.Equal(t, 3, arr.Length())

    // Append more values
    arr.Append(40)
    arr.Append(50)

    // Check the updated length
    assert.Equal(t, 5, arr.Length())
}

func TestResize(t *testing.T) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](5, config)

    // Append values
    arr.Append(10)
    arr.Append(20)

    // Resize the array to a larger size
    err := arr.Resize(10)
    assert.NoError(t, err)

    // Append more values after resizing
    arr.Append(30)
    arr.Append(40)

    // Ensure the values are correct
    value, err := arr.Get(0)
    assert.NoError(t, err)
    assert.Equal(t, 10, value)

    value, err = arr.Get(1)
    assert.NoError(t, err)
    assert.Equal(t, 20, value)

    value, err = arr.Get(2)
    assert.NoError(t, err)
    assert.Equal(t, 30, value)

    value, err = arr.Get(3)
    assert.NoError(t, err)
    assert.Equal(t, 40, value)

    // Try resizing to a smaller capacity (should fail)
    err = arr.Resize(3)
    assert.Error(t, err)
    assert.Equal(t, "new capacity must be greater than or equal to the current size", err.Error())
}

func TestDelete(t *testing.T) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](5, config)

    // Append some values
    arr.Append(10)
    arr.Append(20)
    arr.Append(30)
    arr.Append(40)

    // Delete the second element (index 1)
    err := arr.Delete(1)
    assert.NoError(t, err)

    // Ensure the values are correct after deletion
    value, err := arr.Get(0)
    assert.NoError(t, err)
    assert.Equal(t, 10, value)

    value, err = arr.Get(1)
    assert.NoError(t, err)
    assert.Equal(t, 30, value)

    value, err = arr.Get(2)
    assert.NoError(t, err)
    assert.Equal(t, 40, value)

    // Try to delete an element at an invalid index
    err = arr.Delete(10)
    assert.Error(t, err)
    assert.Equal(t, "index out of bounds", err.Error())
}

func BenchmarkAppend(b *testing.B) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](1000, config)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        arr.Append(i)
    }
}

// BenchmarkAppendConcurrent benchmarks the concurrent version of Append.
func BenchmarkAppendConcurrent(b *testing.B) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](1000, config)

    var wg sync.WaitGroup
    b.ResetTimer()

    // Launch multiple goroutines to append concurrently
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            wg.Add(1)
            go func() {
                defer wg.Done()
                arr.Append(1)
            }()
        }
    })

    // Wait for all goroutines to complete
    wg.Wait()
}

// BenchmarkGet benchmarks the Get method for retrieving an element at a given index.
func BenchmarkGet(b *testing.B) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](1000, config)

    // Pre-fill the array
    for i := 0; i < 1000; i++ {
        arr.Append(i)
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        arr.Get(i % 1000) // Get elements at varying indices
    }
}

// BenchmarkGetConcurrent benchmarks the concurrent version of Get.
func BenchmarkGetConcurrent(b *testing.B) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](1000, config)

    // Pre-fill the array
    for i := 0; i < 1000; i++ {
        arr.Append(i)
    }

    var wg sync.WaitGroup
    b.ResetTimer()

    // Run multiple goroutines concurrently to get elements
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            wg.Add(1)
            go func() {
                defer wg.Done()
                arr.Get(1 % 1000) // Get an element at a varying index
            }()
        }
    })

    // Wait for all goroutines to complete
    wg.Wait()
}

// BenchmarkDelete benchmarks the Delete method for deleting an element at a specific index.
func BenchmarkDelete(b *testing.B) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](1000, config)

    // Pre-fill the array
    for i := 0; i < 1000; i++ {
        arr.Append(i)
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        arr.Delete(i % 1000) // Delete elements at varying indices
    }
}


// BenchmarkDeleteConcurrent benchmarks the concurrent version of Delete.
func BenchmarkDeleteConcurrent(b *testing.B) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](1000, config)

    // Pre-fill the array
    for i := 0; i < 1000; i++ {
        arr.Append(i)
    }

    var wg sync.WaitGroup
    b.ResetTimer()

    // Run multiple goroutines concurrently to delete elements
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            wg.Add(1)
            go func() {
                defer wg.Done()
                arr.Delete(1 % 1000) // Delete elements at varying indices
            }()
        }
    })

    // Wait for all goroutines to complete
    wg.Wait()
}


// BenchmarkResize benchmarks the Resize method to measure the cost of resizing the array.
func BenchmarkResize(b *testing.B) {
    config := ArrayConfig{MetricsEnabled: false}
    arr := NewArray[int](1000, config)

    // Pre-fill the array
    for i := 0; i < 1000; i++ {
        arr.Append(i)
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        arr.Resize(2000) // Resize to a larger size
    }
}