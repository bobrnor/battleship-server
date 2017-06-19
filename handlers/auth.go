package handlers

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"go.uber.org/zap"

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
	zap.S().Info("Received", i)
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

	c := h.fetchClient()
	if c == nil {
		c = h.createNewClient()
	}

	h.client = c
}

func (h *authHandler) fetchClient() *db.Client {
	if h.err != nil {
		return nil
	}

	c, err := db.FindClientByUID(h.p.ClientUID)
	if err != nil {
		h.err = err
		return nil
	}

	return c
}

func (h *authHandler) createNewClient() *db.Client {
	if h.err != nil {
		return nil
	}

	newClient := db.Client{
		UID: h.p.ClientUID,
	}
	if err := newClient.Save(nil); err != nil {
		h.err = err
		return nil
	}

	return &newClient
}

func (h *authHandler) response() interface{} {
	status := 0
	if h.err != nil {
		zap.S().Errorf("Error %+v", h.err)
		status = -1
	}
	return map[string]interface{}{
		"status": status,
	}
}
