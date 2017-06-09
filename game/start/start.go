package start

import (
	"fmt"
	"net/http"

	"git.nulana.com/bobrnor/battleship-server/db/client"
	json "git.nulana.com/bobrnor/json.git"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type params struct {
	ClientUID string    `json:"client_uid"`
	RoomUID   string    `json:"room_uid"`
	Grid      [13]uint8 `json:"grid"`
}

type handler struct {
	p *params
	c *client.Client

	err error
}

func Handler() http.HandlerFunc {
	return json.Decorate(handle, (*params)(nil))
}

func handle(i interface{}) interface{} {
	zap.S().Infof("handling start request %+v", i)
	h := handler{}
	return h.handle(i)
}

func (h *handler) handle(i interface{}) interface{} {
	h.fetchParams(i)
	h.fetchClient()
	return h.response()
}

func (h *handler) fetchParams(i interface{}) {
	zap.S().Infof("fetching params %+v", i)
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

	zap.S().Infof("fetching client %+v", h.p.ClientUID)

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
