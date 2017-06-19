package handlers

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"git.nulana.com/bobrnor/battleship-server/db"
	json "git.nulana.com/bobrnor/json.git"
	"git.nulana.com/bobrnor/longpoll.git"
	seqqueue "git.nulana.com/bobrnor/seqqueue.git"
)

type lonpollPparams struct {
	ClientUID string `json:"client_uid"`
	Seq       uint64 `json:"seq"`
	Reset     bool   `json:"reset"`
}

type longpollHandler struct {
	p *lonpollPparams
	c *db.Client
	e *seqqueue.Entry

	err error
}

const (
	longpollTimeout = 1 * time.Minute
)

func LongpollHandler() http.HandlerFunc {
	return json.Decorate(handleLongpoll, (*lonpollPparams)(nil))
}

func handleLongpoll(i interface{}) interface{} {
	zap.S().Infof("handling longpoll request %+v", i)
	h := longpollHandler{}
	return h.handleLongpoll(i)
}

func (h *longpollHandler) handleLongpoll(i interface{}) interface{} {
	h.fetchParams(i)
	h.fetchClient()
	h.poll()
	return h.response()
}

func (h *longpollHandler) fetchParams(i interface{}) {
	zap.S().Infof("fetching lonpollPparams %+v", i)
	p, ok := i.(*lonpollPparams)
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

func (h *longpollHandler) fetchClient() {
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

func (h *longpollHandler) poll() {
	if h.err != nil {
		return
	}

	zap.S().Infof("polling %+v", h.c)

	q := longpoll.DefaultLongpoll().Register(h.p.ClientUID)
	var c <-chan *seqqueue.Entry
	if h.p.Reset {
		c = q.OutWithoutSeq()
	} else {
		c = q.Out(h.p.Seq)
	}

	select {
	case h.e = <-c:
	case <-time.After(longpollTimeout):
	}
}

func (h *longpollHandler) response() interface{} {
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
