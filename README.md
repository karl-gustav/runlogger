Purpose
=======
A _minimal_ zero dependency logger for Google Cloud Platform/Stackdriver

Usage
=====
```
package main

import (
	"os"
	"github.com/karl-gustav/runlogger"
)

var log *runlogger.Logger
func init() {
	if os.Getenv("K_SERVICE") != "" { // Check if running in cloud run
		log= runlogger.StructuredLogger()
	} else {
		log = runlogger.PlainLogger()
	}
}

func main() {
	log.Info("Hello", "world")
}
```
