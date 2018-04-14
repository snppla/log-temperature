package parser

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// SimpleParser reads the default output and just scans for temperature
type SimpleParser struct {
	Reader *bufio.Reader
}

// Measure returns the temperature. Null if it can't
func (sp SimpleParser) Measure() (*Measurement, error) {
	text, err := sp.Reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(text)
	for i := 0; i < len(fields); i++ {
		fields[i] = strings.ToLower(fields[i])
	}
	if len(fields) == 3 {
		if strings.Compare(fields[0], "temperature:") == 0 {
			temp, _ := strconv.ParseFloat(fields[1], 64)
			unit := fields[2]
			if strings.Compare(unit, "f") == 0 {
				temp = (temp - 32) * (5.0 / 9.0)
			}
			fmt.Println("temp is " + strconv.FormatFloat(temp, 'f', 6, 64))
			return &Measurement{ID: 0, Temperature: temp}, nil
		}
	}
	return nil, nil
}
