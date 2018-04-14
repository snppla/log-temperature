package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	parser "log-temperature/parser"
	"strconv"
	"strings"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

type arrayFlags map[int]string

var db *string
var username *string
var password *string
var host *string
var locations = arrayFlags{}

func (i arrayFlags) String() string {
	text, _ := json.Marshal(i)
	return string(text)
}

func (i arrayFlags) Set(value string) error {

	values := strings.Split(value, ":")
	if len(values) != 2 {
		return errors.New("Invalid location")
	}
	id, _ := strconv.ParseInt(values[0], 10, 32)

	i[int(id)] = values[1]
	return nil
}

func main() {
	db = flag.String("db", "temperature", "Influxdb database")
	username = flag.String("user", "", "Username")
	password = flag.String("password", "", "Password")
	host = flag.String("host", "http://localhost", "Host to connect to")
	flag.Var(&locations, "location", "Id and location of the sensor (5:conference-room")
	flag.Parse()

	//reader := bufio.NewReader(os.Stdin)
	reader := bufio.NewReader(strings.NewReader(jasonTestString))
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
	parser := parser.JSONParser{reader}
	for {
		measurement, err := parser.Measure()
		if measurement != nil {
			logTemp(*measurement, client)
			_ = client
		}
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func logTemp(measurement parser.Measurement, client influx.Client) {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  *db,
		Precision: "s",
	})
	if err != nil {
		fmt.Println("Could not log temperature")
		fmt.Println(err)
		return
	}
	var tags = map[string]string{}
	if val, ok := locations[measurement.ID]; ok {
		tags["location"] = val

	} else {
		tags["location"] = "unknown"
	}
	fields := map[string]interface{}{
		"temperature": measurement.Temperature,
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
	simpleTestString = `
	2018-04-13 22:57:36 :   Acurite 606TX Sensor    :       -5
        Battery:         OK
        Temperature:     21.8 C
2018-04-13 22:58:01 :   Acurite 606TX Sensor    :       103
        Battery:         OK
        Temperature:     24.3 C
2018-04-13 22:58:07 :   Acurite 606TX Sensor    :       -5
        Battery:         OK
        Temperature:     21.8 C
2018-04-13 22:58:32 :   Acurite 606TX Sensor    :       103
        Battery:         OK
        Temperature:     24.6 C
2018-04-13 22:58:38 :   Acurite 606TX Sensor    :       -5
        Battery:         OK
        Temperature:     21.8 C
2018-04-13 22:59:03 :   Acurite 606TX Sensor    :       103
        Battery:         OK
        Temperature:     24.8 C
2018-04-13 22:59:09 :   Acurite 606TX Sensor    :       -5
        Battery:         OK
        Temperature:     21.8 C

	`
	jasonTestString = `
	Registering protocol [1] "Rubicson Temperature Sensor"
	Registering protocol [2] "Prologue Temperature Sensor"
	Registering protocol [3] "Waveman Switch Transmitter"
	Registering protocol [4] "LaCrosse TX Temperature / Humidity Sensor"
	Registering protocol [5] "Acurite 609TXC Temperature and Humidity Sensor"
	Registering protocol [6] "Oregon Scientific Weather Sensor"
	Registering protocol [7] "KlikAanKlikUit Wireless Switch"
	Registering protocol [8] "AlectoV1 Weather Sensor (Alecto WS3500 WS4500 Ventus W155/W044 Oregon)"
	Registering protocol [9] "Cardin S466-TX2"
	Registering protocol [10] "Fine Offset Electronics, WH2 Temperature/Humidity Sensor"
	Registering protocol [11] "Nexus Temperature & Humidity Sensor"
	Registering protocol [12] "Ambient Weather Temperature Sensor"
	Registering protocol [13] "Calibeur RF-104 Sensor"
	Registering protocol [14] "GT-WT-02 Sensor"
	Registering protocol [15] "Danfoss CFR Thermostat"
	Registering protocol [16] "Chuango Security Technology"
	Registering protocol [17] "Generic Remote SC226x EV1527"
	Registering protocol [18] "TFA-Twin-Plus-30.3049 and Ea2 BL999"
	Registering protocol [19] "Fine Offset Electronics WH1080/WH3080 Weather Station"
	Registering protocol [20] "WT450"
	Registering protocol [21] "LaCrosse WS-2310 Weather Station"
	Registering protocol [22] "Esperanza EWS"
	Registering protocol [23] "Efergy e2 classic"
	Registering protocol [24] "Generic temperature sensor 1"
	Registering protocol [25] "WG-PB12V1"
	Registering protocol [26] "HIDEKI TS04 Temperature, Humidity, Wind and Rain Sensor"
	Registering protocol [27] "Watchman Sonic / Apollo Ultrasonic / Beckett Rocket oil tank monitor"
	Registering protocol [28] "CurrentCost Current Sensor"
	Registering protocol [29] "emonTx OpenEnergyMonitor"
	Registering protocol [30] "HT680 Remote control"
	Registering protocol [31] "S3318P Temperature & Humidity Sensor"
	Registering protocol [32] "Akhan 100F14 remote keyless entry"
	Registering protocol [33] "Quhwa"
	Registering protocol [34] "OSv1 Temperature Sensor"
	Registering protocol [35] "Proove"
	Registering protocol [36] "Bresser Thermo-/Hygro-Sensor 3CH"
	Registering protocol [37] "Springfield Temperature and Soil Moisture"
	Registering protocol [38] "Oregon Scientific SL109H Remote Thermal Hygro Sensor"
	Registering protocol [39] "Acurite 606TX Temperature Sensor"
	Registering protocol [40] "TFA pool temperature sensor"
	Registering protocol [41] "Kedsum Temperature & Humidity Sensor"
	Registering protocol [42] "blyss DC5-UK-WH (433.92 MHz)"
	Registering protocol [43] "Steelmate TPMS"
	Registering protocol [44] "Schrader TPMS"
	Registering protocol [45] "Elro DB286A Doorbell"
	Registering protocol [46] "Efergy Optical"
	Registering protocol [47] "Honda Car Key"
	Registering protocol [48] "Fine Offset Electronics, XC0400"
	Registering protocol [49] "Radiohead ASK"
	Registering protocol [50] "Kerui PIR Sensor"
	Registering protocol [51] "Fine Offset WH1050 Weather Station"
	Registering protocol [52] "Honeywell Door/Window Sensor"
	Registering protocol [53] "Maverick ET-732/733 BBQ Sensor"
	Registering protocol [54] "LaCrosse TX141-Bv2/TX141TH-Bv2 sensor"
	Registering protocol [55] "Acurite 00275rm,00276rm Temp/Humidity with optional probe"
	Registering protocol [56] "LaCrosse TX35DTH-IT Temperature sensor"
	Registering protocol [57] "LaCrosse TX29IT Temperature sensor"
	Registering protocol [58] "Vaillant calorMatic 340f Central Heating Control"
	Registering protocol [59] "Fine Offset Electronics, WH25 Temperature/Humidity/Pressure Sensor"
	Registering protocol [60] "Fine Offset Electronics, WH0530 Temperature/Rain Sensor"
	Registering protocol [61] "IBIS beacon"
	Registering protocol [62] "Oil Ultrasonic STANDARD FSK"
	Registering protocol [63] "Citroen TPMS"
	Registering protocol [64] "Oil Ultrasonic STANDARD ASK"
	Registering protocol [65] "Thermopro TP11 Thermometer"
	Registering protocol [66] "Solight TE44"
	Registering protocol [67] "Wireless Smoke and Heat Detector GS 558"
	Registering protocol [68] "Generic wireless motion sensor"
	Registering protocol [69] "Toyota TPMS"
	Registering protocol [70] "Ford TPMS"
	Registering protocol [71] "Renault TPMS"
	Registering protocol [72] "FT-004-B Temperature Sensor"
	Registering protocol [73] "Ford Car Key"
	Registering protocol [74] "Philips outdoor temperature sensor"
	Registering protocol [75] "Schrader TPMS EG53MA4"
	Registering protocol [76] "Nexa"
	Registering protocol [77] "Thermopro TP12 Thermometer"
	Registering protocol [78] "GE Color Effects"
	Registering protocol [79] "X10 Security"
	Registering protocol [80] "Interlogix GE UTC Security Devices"
	Registered 80 out of 101 device decoding protocols
	Found 1 device(s)
	
	trying device  0:  Realtek, RTL2838UHIDIR, SN: 00000001
	Found Rafael Micro R820T tuner
	Using device 0: Generic RTL2832U OEM
	Exact sample rate is: 250000.000414 Hz
	[R82XX] PLL not locked!
	Sample rate set to 250000.
	Bit detection level set to 0 (Auto).
	Tuner gain set to Auto.
	Reading samples in async mode...
	Tuned to 433920000 Hz.
	^[ {"time" : "2018-04-13 23:00:05", "model" : "Acurite 606TX Sensor", "id" : 103, "battery" : "OK", "temperature_C" : 24.900}
	{"time" : "2018-04-13 23:00:11", "model" : "Acurite 606TX Sensor", "id" : -5, "battery" : "OK", "temperature_C" : 21.800}
	{"time" : "2018-04-13 23:00:36", "model" : "Acurite 606TX Sensor", "id" : 103, "battery" : "OK", "temperature_C" : 24.900}
	{"time" : "2018-04-13 23:00:42", "model" : "Acurite 606TX Sensor", "id" : -5, "battery" : "OK", "temperature_C" : 21.800}
	{"time" : "2018-04-13 23:01:07", "model" : "Acurite 606TX Sensor", "id" : 103, "battery" : "OK", "temperature_C" : 24.900}
	{"time" : "2018-04-13 23:01:13", "model" : "Acurite 606TX Sensor", "id" : -5, "battery" : "OK", "temperature_C" : 21.800}
	{"time" : "2018-04-13 23:01:38", "model" : "Acurite 606TX Sensor", "id" : 103, "battery" : "OK", "temperature_C" : 24.900}
	{"time" : "2018-04-13 23:01:44", "model" : "Acurite 606TX Sensor", "id" : -5, "battery" : "OK", "temperature_C" : 21.800}
	{"time" : "2018-04-13 23:02:09", "model" : "Acurite 606TX Sensor", "id" : 103, "battery" : "OK", "temperature_C" : 24.900}
	{"time" : "2018-04-13 23:02:15", "model" : "Acurite 606TX Sensor", "id" : -5, "battery" : "OK", "temperature_C" : 21.800}
	{"time" : "2018-04-13 23:02:40", "model" : "Acurite 606TX Sensor", "id" : 103, "battery" : "OK", "temperature_C" : 24.800}
	{"time" : "2018-04-13 23:02:46", "model" : "Acurite 606TX Sensor", "id" : -5, "battery" : "OK", "temperature_C" : 21.800}
	{"time" : "2018-04-13 23:03:11", "model" : "Acurite 606TX Sensor", "id" : 103, "battery" : "OK", "temperature_C" : 24.800}
	{"time" : "2018-04-13 23:03:17", "model" : "Acurite 606TX Sensor", "id" : -5, "battery" : "OK", "temperature_C" : 21.800}
	{"time" : "2018-04-13 23:03:42", "model" : "Acurite 606TX Sensor", "id" : 103, "battery" : "OK", "temperature_C" : 24.800}
	{"time" : "2018-04-13 23:03:48", "model" : "Acurite 606TX Sensor", "id" : -5, "battery" : "OK", "temperature_C" : 21.800}
	{"time" : "2018-04-13 23:04:13", "model" : "Acurite 606TX Sensor", "id" : 103, "battery" : "OK", "temperature_C" : 24.700}
	{"time" : "2018-04-13 23:04:19", "model" : "Acurite 606TX Sensor", "id" : -5, "battery" : "OK", "temperature_C" : 21.800}
	{"time" : "2018-04-13 23:04:44", "model" : "Acurite 606TX Sensor", "id" : 103, "battery" : "OK", "temperature_C" : 24.700}
	{"time" : "2018-04-13 23:04:50", "model" : "Acurite 606TX Sensor", "id" : -5, "battery" : "OK", "temperature_C" : 21.800}
	
	

`
)
