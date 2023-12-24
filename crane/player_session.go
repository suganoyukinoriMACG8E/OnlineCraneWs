package crane

import (
	"github.com/olahol/melody"
	"sync"
)

type PlayerSession struct {
	sessions sync.Map
}

func (p *PlayerSession) Set(playerID string, session *melody.Session) {
	p.sessions.Store(playerID, session)
}

func (p *PlayerSession) Get(playerID string) (session *melody.Session) {
	s, _ := p.sessions.Load(playerID)
	if s == nil {
		return nil
	}
	session = s.(*melody.Session)
	return
}

func (p *PlayerSession) Delete(playerID string) {
	p.sessions.Delete(playerID)
}
