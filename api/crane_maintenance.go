package api

import (
	"OnlineCraneWs/crane"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CraneMaintenance(c *crane.Crane, ctx echo.Context) error {
	message := &MonitorMessage{}
	message.Target = "broadcast"
	if err := ctx.Bind(message); err != nil {
		return err
	}
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	c.Send(msg)
	ctx.String(http.StatusOK, "")
	return nil
}
