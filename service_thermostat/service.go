package thermostatservice

//
//import (
//	devicestore "mqtt-go-playground/device_config_store"
//	"mqtt-go-playground/device_config_store/device"
//	"mqtt-go-playground/mqtt"
//	"time"
//
//	"github.com/sirupsen/logrus"
//)
//
//type service struct {
//	ds devicestore.Devicestore
//}
//
//func CreateService() service {
//	ds := devicestore.CreateDeviceStore()
//	return service{ds}
//}
//
//func (s service) Run() {
//	s.listenForDevConfigs()
//}
//
//func (s service) listenForDevConfigs() {
//	topic := "homeassistant/climate/+/climate/config"
//
//	thing := make(chan ([2]string))
//
//	mqtt.Subscribe(topic, thing)
//	time.AfterFunc(time.Second, func() {
//		mqtt.Unsubscribe(topic)
//	})
//
//	for t := range thing {
//		d := device.ConfigCreator([]byte(t[1]))
//		s.ds.UpdateDeviceConfig(d)
//		logrus.Debugf("RECEIVED Device Config: %s MESSAGE: %+v\n", t[0], d.Name)
//	}
//
//	panic("will never get here")
//}
