package crane

import (
	"github.com/olahol/melody"
	"sync"
)

type Handler struct {
	handlers sync.Map
}

func (h *Handler) Set(name string, hand func(*melody.Session, *Crane, []byte)) {
	h.handlers.Store(name, hand)
}

func (h *Handler) Get(name string) (hand func(*melody.Session, *Crane, []byte)) {
	base, _ := h.handlers.Load(name)
	hand = base.(func(*melody.Session, *Crane, []byte))
	return
}
