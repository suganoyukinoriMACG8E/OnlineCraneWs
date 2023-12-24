package api

import (
	"OnlineCraneWs/crane"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Get(c *crane.Crane, ctx echo.Context) error {
	c.Logger.Print("api get")
	message := &Message{}
	if err := ctx.Bind(message); err != nil {
		return err
	}
	message.Target = message.UUID
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	if err = c.Send(msg); err != nil {
		return err
	}
	ctx.String(http.StatusOK, "")
	return nil
}
