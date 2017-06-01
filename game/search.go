package game

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/hashicorp/packer/common/uuid"
	"github.com/pkg/errors"

	"git.nulana.com/bobrnor/battleship-server/db/client"
	"git.nulana.com/bobrnor/battleship-server/db/room"
	jsonep "git.nulana.com/bobrnor/json-ep.git"
)

type params struct {
	ClientUID string `json:"client_uid"`
}

type handler struct {
	p *params
	c *client.Client

	err error
}

func Handler() http.HandlerFunc {
	var p params
	return jsonep.Decorate(handle, &p)
}

func handle(i interface{}) interface{} {
	zap.S().Info("Received", i)
	h := handler{}
	return h.handler(i)
}

func (h *handler) handler(i interface{}) interface{} {
	h.fetchParams(i)
	h.fetchClient()
	return h.response()
}

func (h *handler) fetchParams(i interface{}) {
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

func (h *handler) doSearch() {
	if h.err != nil {
		return
	}

	if c := DefaultPlaylist().PopAny(); c != nil {
		newRoom := room.Room{
			UID:       uuid.TimeOrderedUUID(),
			ClientID1: c.ID,
			ClientID2: h.c.ID,
		}
		if err := newRoom.Save(nil); err != nil {
			h.err = err
			DefaultPlaylist().Push(c)
		}
	} else {
		DefaultPlaylist().Wait(h.c)
	}
}

func (h *handler) response() interface{} {
	status := 0
	if h.err != nil {
		status = -1
	}
	return map[string]interface{}{
		"status": status,
	}
}

//
// 	// TODO: find request or create new (!!!: concurrency)
// 	// if request found:
// 	// - delete it
// 	// - create room
// 	// - return room uid
//
// 	return map[string]interface{}{
// 		"status": 0,
// 	}
// }
