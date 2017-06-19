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

type turnParams struct {
	ClientUID string `json:"client_uid"`
	RoomUID   string `json:"room_uid"`
	X         uint   `json:"x"`
	Y         uint   `json:"y"`
}

type turnHandler struct {
	p *turnParams
	c *db.Client
	r *db.Room

	result game.TurnResult

	err error
}

func TurnHandler() http.HandlerFunc {
	return json.Decorate(handleTurn, (*turnParams)(nil))
}

func handleTurn(i interface{}) interface{} {
	zap.S().Infof("handling turn request %+v", i)
	h := turnHandler{}
	return h.handleTurn(i)
}

func (h *turnHandler) handleTurn(i interface{}) interface{} {
	h.fetchParams(i)
	h.fetchClient()
	h.fetchRoom()
	h.doTurn()
	return h.response()
}

func (h *turnHandler) fetchParams(i interface{}) {
	zap.S().Infof("fetching turnParams %+v", i)
	p, ok := i.(*turnParams)
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

func (h *turnHandler) fetchClient() {
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
		h.err = errors.Errorf("Client not found `%v`", h.p.ClientUID)
		return
	}

	h.c = c
}

func (h *turnHandler) fetchRoom() {
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

func (h *turnHandler) doTurn() {
	if h.err != nil {
		return
	}

	engine := game.MainEngine()
	if result, err := engine.Turn(h.r, h.c, h.p.X, h.p.Y); err != nil {
		h.err = err
	} else {
		h.result = result
	}
}

func (h *turnHandler) response() interface{} {
	resp := map[string]interface{}{}
	if h.err != nil {
		zap.S().Errorf("Error %+v", h.err)
		resp["status"] = -1
	} else {
		resp["result"] = h.result
		resp["status"] = 0
	}

	return resp
}
