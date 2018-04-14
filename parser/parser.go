package parser

// Measurement provides an interface for each reading
type Measurement struct {
	// ID of the sensor
	ID int
	// Temperature in celcius
	Temperature float64
}

// Parser provides an interface for parsing temperature strings
type Parser interface {
	Measure() (*Measurement, error)
}
