package api

type (
	Message struct {
		Ope      string `json:"ope"`
		Target   string `json:"target"`
		UUID     string `json:"uuid"`
		CraneID  int64  `json:"crane_id"`
		RewardID int64  `json:"reward_id"`
	}

	MonitorMessage struct {
		Ope     string `json:"ope"`
		Target  string `json:"target"`
		CraneID int64  `json:"crane_id"`
	}
)
