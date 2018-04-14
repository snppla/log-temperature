package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
)

// JSONParser reads the json output
type JSONParser struct {
	Reader *bufio.Reader
}

// {"time" : "2018-04-13 23:04:50", "model" : "Acurite 606TX Sensor", "id" : -5, "battery" : "OK", "temperature_C" : 21.800}
type jsonPacket struct {
	Time          string
	Model         string
	ID            int
	Battery       string
	Temperature_C float64
}

// Measure returns the temperature. Null if it can't
func (jp JSONParser) Measure() (*Measurement, error) {
	text, err := jp.Reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	buffer := []byte(text)
	if json.Valid(buffer) {
		var jsonPacket jsonPacket

		err = json.Unmarshal(buffer, &jsonPacket)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		fmt.Println("temp is " + strconv.FormatFloat(jsonPacket.Temperature_C, 'f', 6, 64))
		return &Measurement{ID: jsonPacket.ID, Temperature: jsonPacket.Temperature_C}, nil
	}
	return nil, nil
}
