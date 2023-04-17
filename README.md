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
	log.Info("Hello", "world") // logged as "Hello world"
	log.Infof("Hello %s", "world") // logged as "Hello world"

	// logged as "Hello world" with a an attached json
	// structure called "struct": `{"Maximum","Effort"}`
	wantAsJson := struct{ Maximum string }{"Effort"}
	log.Info("Hello", "world", log.Field("struct", wantAsJson))

	// also works on the formated log methods (e.g. Infof(...))
	// the fields are just ignored for the formated string
	log.Infof("Hello %", log.Field("struct", wantAsJson), "world")
}
```

NB: The "anything" in `log.Field(<name>, anything)` is sent unchanged
to the marshal function, i.e. if you need to show a byte string, you
need to wrap it in `string()`.
I.e. `log.Field("lorum", string(someBytes))`.
