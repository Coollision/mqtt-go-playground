package serverMgmt

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"reflect"
	"sort"
	"strings"
	"syscall"
	"time"
)

var GracefulStop chan os.Signal

func init() {
	GracefulStop = make(chan os.Signal, 1)
	signal.Notify(GracefulStop, syscall.SIGTERM) // Kubernetes shutdown code
	signal.Notify(GracefulStop, syscall.SIGINT)  // CTRL + C
}

func CustomLogging(logLevel string) {
	// Setup logger
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalf("Failed to parse log level. %v", err)
	}

	logrus.SetLevel(lvl)

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = time.RFC3339Nano
	// customFormatter.DisableColors = true
	customFormatter.FullTimestamp = true
	customFormatter.FieldMap = logrus.FieldMap{
		logrus.FieldKeyTime:  "time",
		logrus.FieldKeyLevel: "lvl",
		logrus.FieldKeyMsg:   "msg",
	}
	// get message as last value
	msgIsLastValue := func(s []string) {
		sort.Slice(s, func(i, j int) bool { return s[j] == "msg" })
	}
	customFormatter.SortingFunc = msgIsLastValue
	logrus.SetFormatter(customFormatter)
}

// logReflectValue recursively logs the given struct value
func logReflectValue(v reflect.Value, level int, name string) {
	t := v.Type()

	logrus.Infof("%v%v:", strings.Repeat("  ", level), name)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		// log nested config
		if field.Type.Kind() == reflect.Struct {
			logReflectValue(reflect.ValueOf(value), level+1, field.Name)
			continue
		}

		// Hide sensitive data
		secret := field.Tag.Get("secret")
		if secret == "true" {
			value = "[REDACTED]"
		}
		logrus.Infof("%v%v (%v): %v", strings.Repeat("  ", level+1), field.Name, field.Type, value)
	}
}

// log a config without exposing secrets.
// Takes a pointer to a config struct
// Secrets are marked with the tag `secret:"true"`
func Log(cfg interface{}) {
	rootElem := reflect.ValueOf(cfg).Elem()
	logReflectValue(rootElem, 0, "configuration")
	logrus.Infof("---------------------------------")
}
