package confirm

import (
	"fmt"
	"net/http"

	"git.nulana.com/bobrnor/battleship-server/db/client"
	"git.nulana.com/bobrnor/battleship-server/db/room"
	json "git.nulana.com/bobrnor/json.git"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type params struct {
	ClientUID string `json:"client_uid"`
	RoomUID   string `json:"room_uid"`
}

type handler struct {
	p *params
	c *client.Client
	r *room.Room

	err error
}

func Handler() http.HandlerFunc {
	return json.Decorate(handle, (*params)(nil))
}

func handle(i interface{}) interface{} {
	zap.S().Infof("handling confirm request %+v", i)
	h := handler{}
	return h.handle(i)
}

func (h *handler) handle(i interface{}) interface{} {
	h.fetchParams(i)
	h.fetchClient()
	h.fetchRoom()
	h.confirm()
	return h.response()
}

func (h *handler) fetchParams(i interface{}) {
	zap.S().Infof("fetching params %+v", i)
	p, ok := i.(*params)
	if !ok {
		h.err = errors.WithStack(fmt.Errorf("Wrong parameters type %T %v", i, i))
		return
	}

	if len(p.ClientUID) == 0 {
		h.err = errors.WithStack(fmt.Errorf("`client_uid` expected but empty %v", p))
		return
	}

	h.p = p
}

func (h *handler) fetchClient() {
	if h.err != nil {
		return
	}

	zap.S().Infof("fetching client %+v", h.p.ClientUID)

	c, err := client.FindByUID(h.p.ClientUID)
	if err != nil {
		h.err = err
		return
	}

	if c == nil {
		h.err = errors.WithStack(fmt.Errorf("Client not found `%v`", h.p.ClientUID))
		return
	}

	h.c = c
}

func (h *handler) fetchRoom() {
	if h.err != nil {
		return
	}

	r, err := room.FindByUID(nil, h.p.RoomUID)
	if err != nil {
		h.err = err
		return
	}

	if r == nil {
		h.err = errors.New("Room not found")
		return
	}

	clients, err := r.Clients(nil)
	if err != nil {
		h.err = err
		return
	}

	for _, client := range clients {
		if client.ID == h.c.ID {
			h.r = r
			return
		}
	}

	h.err = errors.New("Founded room is not for that client")
}

func (h *handler) confirm() {
	if h.err != nil {
		return
	}

	if err := h.r.Confirm(nil, h.c); err != nil {
		h.err = err
	}
}

func (h *handler) response() interface{} {
	status := 0
	if h.err != nil {
		zap.S().Errorf("Error %+v", h.err)
		status = -1
	}
	return map[string]interface{}{
		"status": status,
	}
}
