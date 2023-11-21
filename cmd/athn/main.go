package main

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn"
import "github.com/sirgallo/athn/common/connpool"


const NAME = "Main"
var Log = logger.NewCustomLog(NAME)


func main() {
	// hostname, hostErr := os.Hostname()
	// if hostErr != nil { log.Fatal("unable to get hostname") }

	/*
	systemsList := []*system.System{
		{ Host: "athnsrv1" },
		{ Host: "athnsrv2" },
		{ Host: "athnsrv3" },
		{ Host: "athnsrv4" },
		{ Host: "athnsrv5" },
	}
	*/

	athnOpts := athn.AthnServiceOpts{
		Protocol: "tcp",
		Ports: athn.AthnPortOpts{
			Request: 8080,
			Liveness: 2345,
			Proposal: 3456,
		},
		ConnPoolOpts: connpool.ConnectionPoolOpts{ MaxConn: 10 },
	}

	athn, createErr := athn.NewAthn(athnOpts)
	if createErr != nil { panic("unable to create athn instance") }

	go athn.StartAthn()
	select{}
}