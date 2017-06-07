package auth

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"go.uber.org/zap"

	"git.nulana.com/bobrnor/battleship-server/db/client"
	json "git.nulana.com/bobrnor/json.git"
)

type params struct {
	ClientUID string `json:"client_uid"`
}

type handler struct {
	p      *params
	client *client.Client

	err error
}

func Handler() http.HandlerFunc {
	return json.Decorate(handle, (*params)(nil))
}

func handle(i interface{}) interface{} {
	zap.S().Info("Received", i)
	h := handler{}
	return h.handle(i)
}

func (h *handler) handle(i interface{}) interface{} {
	h.fetchParams(i)
	h.authClient()
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

func (h *handler) authClient() {
	if h.err != nil {
		return
	}

	c := h.fetchClient()
	if c == nil {
		c = h.createNewClient()
	}

	h.client = c
}

func (h *handler) fetchClient() *client.Client {
	if h.err != nil {
		return nil
	}

	c, err := client.FindByUID(h.p.ClientUID)
	if err != nil {
		h.err = err
		return nil
	}

	return c
}

func (h *handler) createNewClient() *client.Client {
	if h.err != nil {
		return nil
	}

	newClient := client.Client{
		UID: h.p.ClientUID,
	}
	if err := newClient.Save(); err != nil {
		h.err = err
		return nil
	}

	return &newClient
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
