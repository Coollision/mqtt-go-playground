package service_cover

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"mqtt-go-playground/dataObjects"
	"mqtt-go-playground/mqtt"
	"strings"
	"sync"
	"time"
)

type store struct {
	devices []string
	sync.RWMutex
}

var s store

type response struct {
	Position int64  `json:"position"`
	Action   string `json:"action"`
}

type z2mCover struct {
	BacklightMode   string  `json:"backlight_mode,omitempty"`
	CalibrationTime float64 `json:"calibration_time,omitempty"`
	Moving          string  `json:"moving,omitempty"`
	Position        int64   `json:"position,omitempty"`
	State           string  `json:"state,omitempty"`
}

type CoverService struct {
	OriginalConfig   *dataObjects.DeviceConfigCover
	MyConfig         *dataObjects.DeviceConfigCover
	CurrentPos       int64
	CurrentCalibTime float64
}

func NewCoverService(dc *dataObjects.DeviceConfigCover, name string) *CoverService {
	// Return the original device configuration without modifying it
	return &CoverService{
		OriginalConfig: dc,
		MyConfig:       OurDeviceConfig(*dc, name),
	}
}

func (cs *CoverService) handleZ2MCoverMessage(topic string, message []byte) {
	logrus.Tracef("Received Z2M Cover: %s\n", message)
	var z2mCover z2mCover
	if err := json.Unmarshal(message, &z2mCover); err != nil {
		logrus.Errorf("Error: %s\n", err)
		return
	}

	if z2mCover.Position != cs.CurrentPos {
		cs.CurrentPos = z2mCover.Position
	}
	if z2mCover.CalibrationTime != cs.CurrentCalibTime {
		cs.CurrentCalibTime = z2mCover.CalibrationTime
	}

	switch {
	case z2mCover.Moving == "STOP" && z2mCover.State == "OPEN":
		mqtt.Publish(cs.MyConfig.StateTopic, cs.MyConfig.StateOpen)
	case z2mCover.Moving == "STOP" && z2mCover.State == "CLOSE":
		mqtt.Publish(cs.MyConfig.StateTopic, cs.MyConfig.StateClosed)
	case z2mCover.Moving == "UP":
		mqtt.Publish(cs.MyConfig.StateTopic, cs.MyConfig.StateOpening)
	case z2mCover.Moving == "DOWN":
		mqtt.Publish(cs.MyConfig.StateTopic, cs.MyConfig.StateClosing)
	case z2mCover.Moving == "STOP":
		mqtt.Publish(cs.MyConfig.StateTopic, cs.MyConfig.StateStopped)
	}

	mqtt.Publish(cs.MyConfig.PositionTopic, fmt.Sprintf(`{"position": %d}`, z2mCover.Position))
}

func (cs *CoverService) handleSetPositionMessage(topic string, message []byte) {
	logrus.Tracef("Received Set Position: %s\n", message)
	var res response
	if err := json.Unmarshal(message, &res); err != nil {
		logrus.Errorf("Error: %s\n", err)
		return
	}

	if res.Position > 95 {
		mqtt.Publish(cs.OriginalConfig.CommandTopic, `OPEN`)
		return
	}
	if res.Position < 5 {
		mqtt.Publish(cs.OriginalConfig.CommandTopic, `CLOSE`)
		return
	}

	var neededAdjustment = res.Position - cs.CurrentPos
	var timeToMove = math.Abs(float64(neededAdjustment)) * (cs.CurrentCalibTime / 100)

	if neededAdjustment > 0 {
		mqtt.Publish(cs.OriginalConfig.CommandTopic, `OPEN`)
		time.Sleep(time.Duration(timeToMove) * time.Second)
		mqtt.Publish(cs.OriginalConfig.CommandTopic, `STOP`)
	} else if neededAdjustment < 0 {
		mqtt.Publish(cs.OriginalConfig.CommandTopic, `CLOSE`)
		time.Sleep(time.Duration(timeToMove) * time.Second)
		mqtt.Publish(cs.OriginalConfig.CommandTopic, `STOP`)
	}
}

func (cs *CoverService) handleCommandMessage(topic string, message []byte) {
	logrus.Tracef("Received Command: %s\n", message)
	var res response
	if err := json.Unmarshal(message, &res); err != nil {
		logrus.Errorf("Error: %s\n", err)
		return
	}
	switch res.Action {
	case "open":
		mqtt.Publish(cs.OriginalConfig.CommandTopic, `OPEN`)
	case "close":
		mqtt.Publish(cs.OriginalConfig.CommandTopic, `CLOSE`)
	case "stop":
		mqtt.Publish(cs.OriginalConfig.CommandTopic, `STOP`)
	}
}

func (cs *CoverService) Start() {
	t := fmt.Sprintf("homeassistant/cover/%s/cover/config", cs.MyConfig.Name)
	t = camelCaseOn(t, ".")
	t = camelCaseOn(t, "-")

	if err := mqtt.PublishStruct(t, cs.MyConfig); err != nil {
		logrus.Errorf("Error: %s\n", err)
	}

	mqtt.Subscribe(cs.OriginalConfig.PositionTopic, chanSubscriber(cs.handleZ2MCoverMessage))
	mqtt.Subscribe(cs.MyConfig.CommandTopic, chanSubscriber(cs.handleCommandMessage))
	mqtt.Subscribe(cs.MyConfig.SetPositionTopic, chanSubscriber(cs.handleSetPositionMessage))
	logrus.Infof("Cover Service Started for %s\n", cs.MyConfig.Name)
}

func chanSubscriber(function func(topic string, message []byte)) func(string, []byte) {
	c := make(chan ([2]string), 10)
	f := func(topic string, message []byte) {
		c <- [2]string{topic, string(message)}
	}

	go func() {
		for cd := range c {
			function(cd[0], []byte(cd[1]))
		}
	}()

	return f
}

func OurDeviceConfig(dc dataObjects.DeviceConfigCover, uniqueID string) *dataObjects.DeviceConfigCover {
	// Adjusting the DeviceConfig with unique values
	dc.SetPositionTopic = fmt.Sprintf("playground/%s/position/set", uniqueID)
	dc.PositionTopic = fmt.Sprintf("playground/%s/position", uniqueID)
	dc.CommandTopic = fmt.Sprintf("playground/%s/set", uniqueID)
	dc.StateTopic = fmt.Sprintf("playground/%s", uniqueID)

	dc.DeviceClass = "shutter"
	dc.UniqueID = fmt.Sprintf("playground_%s", uniqueID)
	dc.Name = uniqueID

	dc.StateClosed = "closed"
	dc.StateClosing = "closing"
	dc.StateOpen = "open"
	dc.StateOpening = "opening"
	dc.StateStopped = "stopped"
	dc.ValueTemplate = ""

	dc.PayloadOpen = `{"action":"open"}`
	dc.PayloadClose = `{"action":"close"}`
	dc.PayloadStop = `{"action":"stop"}`

	return &dc
}

func DiscoverAndStartCoverServices() {
	c := make(chan ([2]string), 10)
	mqtt.Subscribe("homeassistant/cover/+/cover/config", func(topic string, message []byte) {
		c <- [2]string{topic, string(message)}
	})

	go func() {
		for t := range c {
			var dc dataObjects.DeviceConfigCover
			if err := json.Unmarshal([]byte(t[1]), &dc); err != nil {
				logrus.Errorf("Error: %s\n", err)
				return
			}

			skip := false

			for _, d := range s.devices {
				if d == dc.Device.Name {
					logrus.Infof("Cover Service for %s already started\n", dc.Device.Name)
					skip = true
				}
			}
			if dc.Device.Manufacturer == "Tuya" && dc.Device.Model == "Curtain/blind switch" && dc.Name == nil && !skip {
				// Create a new CoverService for this cover
				cs := NewCoverService(&dc, dc.Device.Name)
				cs.Start()
				s.devices = append(s.devices, dc.Device.Name)
			}

		}
		panic("will never get here")
	}()
}

func Start() {
	s = store{}
	DiscoverAndStartCoverServices()
}

func camelCaseOn(text, char string) string {
	t := strings.Split(text, char)
	for i := range t {
		if i == 0 {
			continue
		}
		t[i] = strings.ToUpper(string(t[i][0])) + t[i][1:]
	}
	return strings.Join(t, "")
}
