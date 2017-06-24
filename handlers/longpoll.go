package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"log"

	"git.nulana.com/bobrnor/battleship-server/db"
	json "git.nulana.com/bobrnor/json.git"
	"git.nulana.com/bobrnor/longpoll.git"
	seqqueue "git.nulana.com/bobrnor/seqqueue.git"
)

type longpollParams struct {
	ClientUID string `json:"client_uid"`
	Seq       uint64 `json:"seq"`
	Reset     bool   `json:"reset"`
}

type longpollHandler struct {
	p *longpollParams
	c *db.Client
	e *seqqueue.Entry

	err error
}

const (
	longpollTimeout = 1 * time.Minute
)

func LongpollHandler() http.HandlerFunc {
	return json.Decorate(handleLongpoll, (*longpollParams)(nil))
}

func handleLongpoll(i interface{}) interface{} {
	log.Printf("handling longpoll request %+v", i)
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
	log.Printf("fetching longpoll params %+v", i)
	p, ok := i.(*longpollParams)
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

func (h *longpollHandler) fetchClient() {
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

func (h *longpollHandler) poll() {
	if h.err != nil {
		return
	}

	log.Printf("polling %+v", h.c)

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
	msg := map[string]interface{}{
		"type": "longpoll",
	}
	if h.err != nil {
		log.Printf("Error %+v", h.err)
		msg["error"] = map[string]interface{}{
			"code": 1,
			"msg":  h.err.Error(),
		}
	} else {
		if h.e != nil {
			msg["seq"] = h.e.Seq
			msg["content"] = h.e.Value
		}
	}
	return msg
}
