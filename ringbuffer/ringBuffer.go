package ringBuffer

import (
	"fmt"
)

// RingBuffer a fixed-size buffer that overwrites the oldest element when full.
type RingBuffer struct {
	Buffer []string
	Size   int
	Head   int // Index of the oldest element
	Tail   int // Index where the next element will be added
	Count  int // Current number of elements
}

// NewRingBuffer creates a new RingBuffer with the specified size.
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		Buffer: make([]string, size),
		Size:   size,
		Head:   0,
		Tail:   0,
		Count:  0,
	}
}

// Add an element to the buffer.
func (b *RingBuffer) Add(data string) {
	if b.Count < b.Size {
		b.Buffer[b.Tail] = data
		b.Tail = (b.Tail + 1) % b.Size
		b.Count++
	} else {
		// Overwrite the oldest element (at the head)
		b.Buffer[b.Head] = data
		b.Head = (b.Head + 1) % b.Size
		b.Tail = (b.Tail + 1) % b.Size
	}
}

// GetElements returns all elements in order from oldest to newest.
func (b *RingBuffer) GetElements() []string {
	elements := make([]string, 0, b.Count)
	if b.Count == 0 {
		return elements
	}

	for i := 0; i < b.Count; i++ {
		index := (b.Head + i) % b.Size
		elements = append(elements, b.Buffer[index])
	}
	return elements
}

// Get elements newest first
func (b *RingBuffer) GetMostRecentElementsFirst() []string {
	elements := make([]string, 0, b.Count)
	if b.Count == 0 {
		return elements
	}

	for i := b.Count - 1; i >= 0; i-- {
		index := (b.Head + i) % b.Size
		elements = append(elements, b.Buffer[index])
	}
	return elements
}

func (b *RingBuffer) PrintMostRecentElementsFirst() {
	if b.Count == 0 {
		fmt.Printf("no messages\n")
		return
	}

	for i := b.Count - 1; i >= 0; i-- {
		fmt.Printf("index: %d,  %s\n", i, b.Buffer[i])
		if i == 0 {
			break
		}
	}

	for i := b.Count - 1; i >= b.Tail; i-- {
		fmt.Printf("index %d, %s\n", i, b.Buffer[i])
	}
}
