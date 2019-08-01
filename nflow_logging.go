package main

import (
	"github.com/Sirupsen/logrus"
	"os"
)

var log = logrus.New()

func InitLog() {
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel
}

func SetLogger(l *logrus.Logger) {
	log = l
}
