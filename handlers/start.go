package handlers

import (
	"fmt"
	"net/http"

	"log"

	"git.nulana.com/bobrnor/battleship-server/core"
	"git.nulana.com/bobrnor/battleship-server/db"
	json "git.nulana.com/bobrnor/json.git"
	"github.com/pkg/errors"
)

type startParams struct {
	ClientUID string   `json:"client_uid"`
	RoomUID   string   `json:"room_uid"`
	Grid      [13]byte `json:"grid"`
}

type startHandler struct {
	p *startParams
	c *db.Client
	r *db.Room

	err error
}

func StartHandler() http.HandlerFunc {
	return json.Decorate(handleStart, (*startParams)(nil))
}

func handleStart(i interface{}) interface{} {
	log.Printf("handling start request %+v", i)
	h := startHandler{}
	return h.handleStart(i)
}

func (h *startHandler) handleStart(i interface{}) interface{} {
	h.fetchParams(i)
	h.fetchClient()
	h.fetchRoom()
	h.setGrid()
	return h.response()
}

func (h *startHandler) fetchParams(i interface{}) {
	log.Printf("fetching startParams %+v", i)
	p, ok := i.(*startParams)
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

func (h *startHandler) fetchClient() {
	if h.err != nil {
		return
	}

	log.Printf("fetching client %+v", h.p.ClientUID)

	c, err := db.FindClientByUID(h.p.ClientUID)
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

func (h *startHandler) fetchRoom() {
	if h.err != nil {
		return
	}

	log.Printf("fetching room %+v", h.p.RoomUID)

	r, err := db.FindRoomByUID(nil, h.p.RoomUID)
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

func (h *startHandler) setGrid() {
	if h.err != nil {
		return
	}

	engine := core.MainEngine()
	if err := engine.SetGrid(h.r, h.c, h.p.Grid); err != nil {
		h.err = err
	}
}

func (h *startHandler) response() interface{} {
	msg := map[string]interface{}{
		"type": "start",
	}
	if h.err != nil {
		log.Printf("Error %+v", h.err)
		msg["error"] = map[string]interface{}{
			"code": 1,
			"msg":  h.err.Error(),
		}
	}
	return msg
}
