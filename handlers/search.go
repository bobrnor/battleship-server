package handlers

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"git.nulana.com/bobrnor/battleship-server/db"
	"git.nulana.com/bobrnor/battleship-server/game"
	json "git.nulana.com/bobrnor/json.git"
)

type searchParams struct {
	ClientUID string `json:"client_uid"`
	Seq       uint64 `json:"seq"`
	Reset     bool   `json:"reset"`
}

type searchHandler struct {
	p *searchParams
	c *db.Client

	err error
}

func SearchHandler() http.HandlerFunc {
	return json.Decorate(handleSearch, (*searchParams)(nil))
}

func handleSearch(i interface{}) interface{} {
	zap.S().Infof("handling search reuqest %+v", i)
	h := searchHandler{}
	return h.handleSearch(i)
}

func (h *searchHandler) handleSearch(i interface{}) interface{} {
	h.fetchParams(i)
	h.fetchClient()
	h.addClientToLobby()
	return h.response()
}

func (h *searchHandler) fetchParams(i interface{}) {
	zap.S().Infof("fetching searchParams %+v", i)
	p, ok := i.(*searchParams)
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

func (h *searchHandler) fetchClient() {
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

func (h *searchHandler) addClientToLobby() {
	if h.err != nil {
		return
	}

	lobby := game.MainLobby()
	lobby.StartWaitingForRoom(h.c)
}

func (h *searchHandler) response() interface{} {
	status := 0
	if h.err != nil {
		zap.S().Errorf("Error %+v", h.err)
		status = -1
	}

	msg := map[string]interface{}{
		"status": status,
	}
	return msg
}
