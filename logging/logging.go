package logging

import (
	"fmt"
	"log"
	"time"
)

type customLogWriter struct {
}

func (writer customLogWriter) Write(bytes []byte) (int, error) {
	return fmt.Print("[DC-STATS] " + time.Now().UTC().Format(time.RFC3339) + " " + string(bytes))
}

func Start() {

	log.SetOutput(new(customLogWriter))
	log.SetFlags(0)
	log.SetFlags(log.Lshortfile)
}
