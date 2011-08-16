/* conductor.go
*/

package main

import (
	"flag"
	"os"
	o	"orchestra"
)

var (
	ConfigFile = flag.String("config-file", "/etc/orchestra/conductor.conf", "File containing the conductor configuration")
	// this is used a lot for marshalling. I just can't stuff it
	// anywhere else.
	InvalidValueError = os.NewError("Invalid value")
)


func main() {
	o.SetLogName("conductor")

	// parse command line options.
	flag.Parse()

	// Start the client registry - configuration parsing will block indefinately
	// if the registry listener isn't working
	StartRegistry()
	// do an initial configuration load
	ConfigLoad()

	// Build the Spool Tree if necessary
	MakeSpoolDir()

	// start the master dispatch system
	InitDispatch()
	defer CleanDispatch()

	// start the status listener
	StartHTTP()
	// start the audience listener
	StartAudienceSock()
	// service TLS requests.
	ServiceRequests()
}
