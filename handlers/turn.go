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

type turnParams struct {
	ClientUID  string `json:"client_uid"`
	RoomUID    string `json:"room_uid"`
	Coordinate struct {
		X uint `json:"x"`
		Y uint `json:"y"`
	} `json:"coord"`
}

type turnHandler struct {
	p *turnParams
	c *db.Client
	r *db.Room

	result core.TurnResult

	err error
}

func TurnHandler() http.HandlerFunc {
	return json.Decorate(handleTurn, (*turnParams)(nil))
}

func handleTurn(i interface{}) interface{} {
	log.Printf("handling turn request %+v", i)
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
	log.Printf("fetching turnParams %+v", i)
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

	log.Printf("fetching client %+v", h.p.ClientUID)

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

func (h *turnHandler) doTurn() {
	if h.err != nil {
		return
	}

	engine := core.MainEngine()
	if result, err := engine.Turn(h.r, h.c, h.p.Coordinate.X, h.p.Coordinate.Y); err != nil {
		h.err = err
	} else {
		h.result = result
	}
}

func (h *turnHandler) response() interface{} {
	msg := map[string]interface{}{
		"type": "turn",
	}
	if h.err != nil {
		log.Printf("Error %+v", h.err)
		msg["error"] = map[string]interface{}{
			"code": 1,
			"msg":  h.err.Error(),
		}
	} else {
		switch h.result {
		case core.TurnResultMiss:
			msg["result"] = "miss"
		case core.TurnResultHit:
			msg["result"] = "hit"
		}
	}
	return msg
}
