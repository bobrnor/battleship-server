package handlers

import (
	"fmt"
	"net/http"

	"git.nulana.com/bobrnor/battleship-server/db"
	"git.nulana.com/bobrnor/battleship-server/game"
	json "git.nulana.com/bobrnor/json.git"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type startParams struct {
	ClientUID string    `json:"client_uid"`
	RoomUID   string    `json:"room_uid"`
	Grid      [13]uint8 `json:"grid"`
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
	zap.S().Infof("handling start request %+v", i)
	h := startHandler{}
	return h.handleStart(i)
}

func (h *startHandler) handleStart(i interface{}) interface{} {
	h.fetchParams(i)
	h.fetchClient()
	h.fetchRoom()
	h.notifyEngine()
	return h.response()
}

func (h *startHandler) fetchParams(i interface{}) {
	zap.S().Infof("fetching startParams %+v", i)
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

	zap.S().Infof("fetching client %+v", h.p.ClientUID)

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

	zap.S().Infof("fetching room %+v", h.p.RoomUID)

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

func (h *startHandler) notifyEngine() {
	if h.err != nil {
		return
	}

	engine := game.DefaultEngine()
	if err := engine.SetGrid(h.r, h.c, h.p.Grid); err != nil {
		h.err = err
	}
}

func (h *startHandler) response() interface{} {
	status := 0
	if h.err != nil {
		zap.S().Errorf("Error %+v", h.err)
		status = -1
	}

	return map[string]interface{}{
		"status": status,
	}
}
