package confirm

import (
	"fmt"
	"net/http"

	"git.nulana.com/bobrnor/battleship-server/db/client"
	"git.nulana.com/bobrnor/battleship-server/db/room"
	jsonep "git.nulana.com/bobrnor/json-ep.git"
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
	return jsonep.Decorate(handle, (*params)(nil))
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

	if r.ClientID1 != h.c.ID && r.ClientID2 != h.c.ID {
		h.err = errors.New("Founded room is not for that client")
		return
	}

	h.r = r
}

func (h *handler) confirm() {
	if h.err != nil {
		return
	}

	if h.r.ClientID1 == h.c.ID {
		h.r.Client1State = room.ConfirmedState
	} else if h.r.ClientID2 == h.c.ID {
		h.r.Client2State = room.ConfirmedState
	}

	if err := h.r.Save(nil); err != nil {
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
