package cli

// Decoder represents an interface which decodes string
type Decoder interface {
	Decode(s string) error
}

// SliceDecoder represents an interface which decodes string as slice
type SliceDecoder interface {
	Decoder
	DecodeSlice()
}

// Encoder represents an interface which encodes to string
type Encoder interface {
	Encode() string
}

// CounterDecoder represents an counter decoder
type CounterDecoder interface {
	Decoder
	IsCounter()
}

// Counter implements counter decoder
type Counter struct {
	value int
}

// Value returns value of counter
func (c Counter) Value() int { return c.value }

// Decode decodes counter from string
func (c *Counter) Decode(s string) error {
	c.value++
	return nil
}

// IsCounter implements method of interface CounterDecoder
func (c Counter) IsCounter() {}
