package handlers

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"log"

	"git.nulana.com/bobrnor/battleship-server/core"
	"git.nulana.com/bobrnor/battleship-server/db"
	json "git.nulana.com/bobrnor/json.git"
)

type searchParams struct {
	ClientUID string `json:"client_uid"`
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
	log.Printf("handling search reuqest %+v", i)
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
	log.Printf("fetching searchParams %+v", i)
	p, ok := i.(*searchParams)
	if !ok {
		h.err = errors.WithStack(fmt.Errorf("wrong parameters type %T %v", i, i))
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

func (h *searchHandler) addClientToLobby() {
	if h.err != nil {
		return
	}

	lobby := core.MainLobby()
	lobby.StartWaitingForRoom(h.c)
}

func (h *searchHandler) response() interface{} {
	msg := map[string]interface{}{
		"type": "search",
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
