package main

import (
	"OnlineCraneWs/api"
	"OnlineCraneWs/crane"
	"OnlineCraneWs/operation"
)

func main() {
	if err := crane.Initialize(); err != nil {
		panic(err)
	}

	crane.HandleAPI("/start", api.Start)
	crane.HandleAPI("/end", api.End)
	crane.HandleAPI("/get", api.Get)
	crane.HandleAPI("/crane_reset", api.CraneReset)
	crane.HandleAPI("/crane_maintenance", api.CraneMaintenance)

	crane.HandleOperation("start", operation.Start)
	crane.HandleOperation("end", operation.End)
	crane.HandleOperation("get", operation.Get)

	if err := crane.Start(); err != nil {
		panic(err)
	}
}
