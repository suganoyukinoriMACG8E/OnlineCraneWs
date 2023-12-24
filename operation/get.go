package operation

import (
	"OnlineCraneWs/crane"
	"encoding/json"
	"github.com/olahol/melody"
)

type (
	Message struct {
		UUID     string `json:"uuid"`
		CraneID  int64  `json:"crane_id"`
		RewardID int64  `json:"reward_id"`
		Target   string `json:"target"`
	}
)

func Get(_ *melody.Session, c *crane.Crane, msg []byte) {
	message := &Message{}
	json.Unmarshal(msg, message)
	message.Target = message.UUID
	c.Send(msg)
}
