Purpose
=======
A _minimal_ zero dependency logger for GCP Cloud Run

Usage
=====
```
package main

import "github.com/karl-gustav/runlogger"

var log *runlogger.Logger
func init() {
	serviceName := os.Getenv("K_SERVICE")
	revision := os.Getenv("K_REVISION")
	configuration := os.Getenv("K_CONFIGURATION")
	var err error
	if serviceName != "" { // Check if running on localhost
        log, err = runlogger.CloudRunLogger(serviceName, revision, configuration)
        if err != nil {
                panic("not able to generate logger because of " + err.Error())
        }
	} else {
        log = runlogger.LocalLogger()
	}
}

func main() {
	log.Info("Hello", "world")
}
```
