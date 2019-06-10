package api

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

func InitLog() (err error) {
	fileInfo, err := os.OpenFile("./logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = fileInfo
	} else {
		return err
		log.Error("Failed to open log file, using default stderr")
	}
	return nil
}
