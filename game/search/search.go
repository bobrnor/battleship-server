package search

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"git.nulana.com/bobrnor/battleship-server/db/client"
	jsonep "git.nulana.com/bobrnor/json-ep.git"
	longpoll "git.nulana.com/bobrnor/longpoll.git"
	seqqueue "git.nulana.com/bobrnor/seqqueue.git"
)

type params struct {
	ClientUID string `json:"client_uid"`
	Seq       uint64 `json:"seq"`
	Reset     bool   `json:"reset"`
}

type handler struct {
	p *params
	c *client.Client
	e *seqqueue.Entry

	err error
}

const (
	longpollTimeout = 1 * time.Minute
)

var (
	lp       = longpoll.NewLongpoll()
	playlist = NewPlaylist(lp)
)

func Handler() http.HandlerFunc {
	return jsonep.Decorate(handle, (*params)(nil))
}

func handle(i interface{}) interface{} {
	zap.S().Infof("handling search reuqest %+v", i)
	h := handler{}
	return h.handle(i)
}

func (h *handler) handle(i interface{}) interface{} {
	h.fetchParams(i)
	h.fetchClient()
	h.poll()
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

func (h *handler) poll() {
	if h.err != nil {
		return
	}

	zap.S().Infof("polling %+v", h.c)

	q := lp.Register(h.p.ClientUID)
	var c <-chan *seqqueue.Entry
	if h.p.Reset {
		c = q.OutWithoutSeq()
	} else {
		c = q.Out(h.p.Seq)
	}

	playlist.Push(h.c)

	select {
	case h.e = <-c:
	case <-time.After(longpollTimeout):
	}
}

func (h *handler) response() interface{} {
	status := 0
	if h.err != nil {
		zap.S().Errorf("Error %+v", h.err)
		status = -1
	}

	var msg map[string]interface{}
	if h.e != nil {
		if m, ok := h.e.Value.(map[string]interface{}); !ok {
			status = -1
			msg = map[string]interface{}{}
		} else {
			msg = m
			msg["seq"] = h.e.Seq
		}
	} else {
		msg = map[string]interface{}{}
	}
	msg["status"] = status
	return msg
}
