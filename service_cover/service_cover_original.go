package service_cover

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/sirupsen/logrus"
//	"math"
//	"mqtt-go-playground/dataObjects"
//	"mqtt-go-playground/mqtt"
//	"time"
//)
//
//type response struct {
//	Position int64  `json:"position"`
//	Action   string `json:"action"`
//}
//
//type z2mCover struct {
//	BacklightMode   string  `json:"backlight_mode,omitempty"`
//	CalibrationTime float64 `json:"calibration_time,omitempty"`
//	Moving          string  `json:"moving,omitempty"`
//	Position        int64   `json:"position,omitempty"`
//	State           string  `json:"state,omitempty"`
//}
//
//func Start() {
//	dc := BasicDeviceConfig()
//
//	var currPos int64
//	var currCalibTime float64
//
//	mqtt.PublishStruct("homeassistant/cover/testing/cover/config", dc)
//
//	mqtt.Subscribe(`zigbee2mqtt/cover.slaapkamer.klein-raam`, func(topic string, message []byte) {
//		var z2mCover z2mCover
//		if err := json.Unmarshal(message, &z2mCover); err != nil {
//			logrus.Errorf("Error: %s\n", err)
//			return
//		}
//
//		if z2mCover.Position != currPos {
//			currPos = z2mCover.Position
//		}
//		if z2mCover.CalibrationTime != currCalibTime {
//			currCalibTime = z2mCover.CalibrationTime
//		}
//
//		switch {
//		case z2mCover.Moving == "STOP" && z2mCover.State == "OPEN":
//			mqtt.Publish(dc.StateTopic, dc.StateOpen)
//		case z2mCover.Moving == "STOP" && z2mCover.State == "CLOSED":
//			mqtt.Publish(dc.StateTopic, dc.StateClosed)
//		case z2mCover.Moving == "UP":
//			mqtt.Publish(dc.StateTopic, dc.StateOpening)
//		case z2mCover.Moving == "DOWN":
//			mqtt.Publish(dc.StateTopic, dc.StateClosing)
//		case z2mCover.Moving == "STOP":
//			mqtt.Publish(dc.StateTopic, dc.StateStopped)
//		}
//
//		mqtt.Publish(dc.PositionTopic, fmt.Sprintf(`{"position": %d}`, z2mCover.Position))
//
//		logrus.Infof("%++v", z2mCover)
//	})
//
//	mqtt.Subscribe(dc.SetPositionTopic, func(topic string, message []byte) {
//		var res response
//		if err := json.Unmarshal(message, &res); err != nil {
//			logrus.Errorf("Error: %s\n", err)
//			return
//		}
//		logrus.Tracef("Position Requested: %d\n", res.Position)
//
//		if res.Position > 95 {
//			mqtt.Publish(`zigbee2mqtt/cover.slaapkamer.klein-raam/set`, `OPEN`)
//			return
//		}
//		if res.Position < 5 {
//			mqtt.Publish(`zigbee2mqtt/cover.slaapkamer.klein-raam/set`, `CLOSE`)
//			return
//		}
//
//		var neededAdjustment = res.Position - currPos
//		var timeToMove = math.Abs(float64(neededAdjustment)) * (currCalibTime / 100)
//
//		if neededAdjustment > 0 {
//			mqtt.Publish(`zigbee2mqtt/cover.slaapkamer.klein-raam/set`, `OPEN`)
//			time.Sleep(time.Duration(timeToMove) * time.Second)
//			mqtt.Publish(`zigbee2mqtt/cover.slaapkamer.klein-raam/set`, `STOP`)
//		} else if neededAdjustment < 0 {
//			mqtt.Publish(`zigbee2mqtt/cover.slaapkamer.klein-raam/set`, `CLOSE`)
//			time.Sleep(time.Duration(timeToMove) * time.Second)
//			mqtt.Publish(`zigbee2mqtt/cover.slaapkamer.klein-raam/set`, `STOP`)
//		}
//	})
//
//	mqtt.Subscribe(dc.CommandTopic, func(topic string, message []byte) {
//		var res response
//		if err := json.Unmarshal(message, &res); err != nil {
//			logrus.Errorf("Error: %s\n", err)
//			return
//		}
//		if res.Action == "open" {
//			mqtt.Publish(`zigbee2mqtt/cover.slaapkamer.klein-raam/set`, `OPEN`)
//
//		} else if res.Action == "close" {
//			mqtt.Publish(`zigbee2mqtt/cover.slaapkamer.klein-raam/set`, `CLOSE`)
//		} else if res.Action == "stop" {
//			mqtt.Publish(`zigbee2mqtt/cover.slaapkamer.klein-raam/set`, `STOP`)
//		}
//	})
//}
//
//func BasicDeviceConfig() *dataObjects.DeviceConfigCover {
//	var data dataObjects.DeviceConfigCover
//	dataS := `{"availability":[{"topic":"zigbee2mqtt/bridge/state"},{"topic":"zigbee2mqtt/cover.slaapkamer.klein-raam/availability"}],"availability_mode":"all","command_topic":"zigbee2mqtt/cover.slaapkamer.klein-raam/set","device":{"configuration_url":"http://zigbee2mqtt.pi//#/device/0xa4c138b5bc064c64/info","identifiers":["zigbee2mqtt_0xa4c138b5bc064c64"],"manufacturer":"_TZ3000_qa8s8vca","model":"Automatically generated definition (TS130F)","name":"cover.slaapkamer.klein-raam","via_device":"zigbee2mqtt_bridge_0x00124b00258d3f4d"},"name":null,"object_id":"cover.slaapkamer.klein-raam","origin":{"name":"Zigbee2MQTT","sw":"1.39.1","url":"https://www.zigbee2mqtt.io"},"position_template":"{{ value_json.position }}","position_topic":"zigbee2mqtt/cover.slaapkamer.klein-raam","set_position_template":"{ \"position\": {{ position }} }","set_position_topic":"zigbee2mqtt/cover.slaapkamer.klein-raam/set","state_closing":"DOWN","state_opening":"UP","state_stopped":"STOP","state_topic":"zigbee2mqtt/cover.slaapkamer.klein-raam","unique_id":"0xa4c138b5bc064c64_cover_zigbee2mqtt","value_template":"{% if \"moving\" in value_json and value_json.moving %} {{ value_json.moving }} {% else %} STOP {% endif %}"}`
//	err := json.Unmarshal([]byte(dataS), &data)
//	if err != nil {
//		fmt.Printf("Error: %s\n", err)
//		panic(err)
//	}
//
//	data.SetPositionTopic = "playground/cover.slaapkamer.klein.raam/position/set"
//	data.PositionTopic = "playground/cover.slaapkamer.klein.raam/position"
//	data.CommandTopic = "playground/cover.slaapkamer.klein.raam/set"
//	data.StateTopic = "playground/cover.slaapkamer.klein.raam"
//
//	data.StateClosed = "closed"
//	data.StateClosing = "closing"
//	data.StateOpen = "open"
//	data.StateOpening = "opening"
//	data.StateStopped = "stopped"
//	data.ValueTemplate = ""
//
//	data.UniqueID = "playground_cover_slaapkamer_klein_raam"
//	data.Name = "cover.slaapkamer.klein.raam"
//	data.PayloadOpen = `{"action":"open"}`
//	data.PayloadClose = `{"action":"close"}`
//	data.PayloadStop = `{"action":"stop"}`
//
//	return &data
//
//}
