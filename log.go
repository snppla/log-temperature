package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

var db *string
var username *string
var password *string
var host *string

func main() {
	db = flag.String("db", "temperature", "Influxdb database")
	username = flag.String("user", "", "Username")
	password = flag.String("password", "", "Password")
	host = flag.String("host", "http://localhost", "Host to connect to")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	//reader := bufio.NewReader(strings.NewReader(testString))
	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     *host,
		Username: *username,
		Password: *password,
	})
	if err != nil {
		fmt.Println("Could not connect")
		fmt.Println(err)
		return
	}
	_ = reader
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			return
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
				logTemp(temp, client)

			}
		}
	}
}

func logTemp(temp float64, client influx.Client) {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  *db,
		Precision: "s",
	})
	if err != nil {
		fmt.Println("Could not log temperature")
		fmt.Println(err)
		return
	}
	tags := map[string]string{"location": "conference-room"}
	fields := map[string]interface{}{
		"temperature": temp,
	}

	pt, err := influx.NewPoint("temperature", tags, fields, time.Now())

	if err != nil {
		fmt.Println("Could not log temperature")
		fmt.Println(err)
		return
	}
	bp.AddPoint(pt)
	err = client.Write(bp)

	if err != nil {
		fmt.Println("Could not log temperature")
		fmt.Println(err)
		return
	}
}

const (
	testString = `
Registering protocol [1] "Acurite 606TX Temperature Sensor"
Registered 1 out of 101 device decoding protocols
Found 1 device(s)

trying device  0:  Realtek, RTL2838UHIDIR, SN: 00000001
Detached kernel driver
Found Rafael Micro R820T tuner
Using device 0: Generic RTL2832U OEM
Exact sample rate is: 250000.000414 Hz
[R82XX] PLL not locked!
Sample rate set to 250000.
Bit detection level set to 0 (Auto).
Tuner gain set to Auto.
Reading samples in async mode...
Tuned to 433920000 Hz.
2018-04-07 21:51:13 :   Acurite 606TX Sensor    :       -5
	Battery:         OK
	Temperature:     25.3 C
2018-04-07 21:51:18 :   Acurite 606TX Sensor    :       -5
	Battery:         OK
	Temperature:     25.3 C
2018-04-07 21:51:20 :   Acurite 606TX Sensor    :       -5
	Battery:         OK
	Temperature:     25.3 C
2018-04-07 21:51:21 :   Acurite 606TX Sensor    :       -5
	Battery:         OK
	Temperature:     25.3 C
2018-04-07 21:51:51 :   Acurite 606TX Sensor    :       -5
	Battery:         OK
	Temperature:     25.7 C
2018-04-07 21:52:21 :   Acurite 606TX Sensor    :       -5
	Battery:         OK
	Temperature:     26.2 C
2018-04-07 21:52:53 :   Acurite 606TX Sensor    :       -5
	Battery:         OK
	Temperature:     26.3 C
2018-04-07 21:53:24 :   Acurite 606TX Sensor    :       -5
	Battery:         OK
	Temperature:     26.4 C


`
)
