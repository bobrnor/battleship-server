package turn

import (
	"fmt"
	"net/http"

	"git.nulana.com/bobrnor/battleship-server/db/client"
	"git.nulana.com/bobrnor/battleship-server/db/room"
	"git.nulana.com/bobrnor/battleship-server/game/longpoll"
	json "git.nulana.com/bobrnor/json.git"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type params struct {
	ClientUID string `json:"client_uid"`
	RoomUID   string `json:"room_uid"`
	X         uint   `json:"x"`
	Y         uint   `json:"y"`
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
	zap.S().Infof("handling turn request %+v", i)
	h := handler{}
	return h.handle(i)
}

func (h *handler) handle(i interface{}) interface{} {
	h.fetchParams(i)
	h.fetchClient()
	h.fetchRoom()
	h.performTurn()
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
		h.err = errors.Errorf("Client not found `%v`", h.p.ClientUID)
		return
	}

	h.c = c
}

func (h *handler) fetchRoom() {
	if h.err != nil {
		return
	}

	zap.S().Infof("fetching room %+v", h.p.RoomUID)

	r, err := room.FindByUID(nil, h.p.RoomUID)
	if err != nil {
		h.err = err
		return
	}

	if r == nil {
		h.err = errors.Errorf("room not found %+v", h.p.RoomUID)
		return
	}

	h.r = r
}

func (h *handler) performTurn() {
	if h.err != nil {
		return
	}

	// use game engine

	clients, err := h.r.Clients(nil)
	if err != nil {
		h.err = err
		return
	}

	var otherClient *client.Client
	for _, c := range clients {
		if h.c.ID != c.ID {
			otherClient = h.c
			break
		}
	}

	if otherClient == nil {
		h.err = errors.Errorf("opponent client not found %+v", otherClient)
		return
	}

	msg := map[string]interface{}{
		"room_uid": h.p.RoomUID,
		"x":        h.p.X,
		"y":        h.p.Y,
	}
	longpoll.Longpoll.Send(otherClient.UID, msg)
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
