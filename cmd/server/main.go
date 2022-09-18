package main

import (
	"log"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/ws-proxy/server"
)

// TODO (#2)
func main() {
	if err := server.Run(server.Config{
		Debug: true,
		Obs: []server.ConfigObs{{
			DkTrackId: 1,
			Host:      "127.0.0.1:4455",
			Password:  "",
		}},
		BindAddr: ":20080",
	}); err != nil {
		log.Fatal(err)
	}
}
