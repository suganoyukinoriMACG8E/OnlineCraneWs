package operation

import (
	"OnlineCraneWs/crane"
	"encoding/json"
	"github.com/olahol/melody"
)

func Start(_ *melody.Session, c *crane.Crane, msg []byte) {
	message := map[string]interface{}{}
	json.Unmarshal(msg, &message)
	message["target"] = "broadcast"
	str, _ := json.Marshal(message)
	c.Logger.Print(string(str))
	c.Send(str)
}
