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

type authParams struct {
	ClientUID string `json:"client_uid"`
}

type authHandler struct {
	p      *authParams
	client *db.Client

	err error
}

func AuthHandler() http.HandlerFunc {
	return json.Decorate(handleAuth, (*authParams)(nil))
}

func handleAuth(i interface{}) interface{} {
	log.Printf("received %+v", i)
	h := authHandler{}
	return h.handleAuth(i)
}

func (h *authHandler) handleAuth(i interface{}) interface{} {
	h.fetchParams(i)
	h.authClient()
	return h.response()
}

func (h *authHandler) fetchParams(i interface{}) {
	p, ok := i.(*authParams)
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

func (h *authHandler) authClient() {
	if h.err != nil {
		return
	}

	c, err := core.AuthClient(h.p.ClientUID)
	if err != nil {
		h.err = err
		return
	}

	h.client = c
}

func (h *authHandler) response() interface{} {
	msg := map[string]interface{}{
		"type": "auth",
	}
	if h.err != nil {
		msg["error"] = map[string]interface{}{
			"code": 1,
			"msg":  h.err.Error(),
		}
	} else {
		msg["client_uid"] = h.client.UID
	}
	return msg
}
