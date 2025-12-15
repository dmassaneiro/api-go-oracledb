package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var Logger *log.Logger

func Init() {
	Logger = log.New()
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&log.JSONFormatter{
		PrettyPrint: true, // opcional, deixa o JSON mais legível
	})
	Logger.SetLevel(log.InfoLevel)
}
