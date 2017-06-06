package longpoll

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"go.uber.org/zap"

	jsonep "git.nulana.com/bobrnor/json-ep.git"
	longpoll "git.nulana.com/bobrnor/longpoll.git"
	seqqueue "git.nulana.com/bobrnor/seqqueue.git"
)

type params struct {
	ClientUID string `json:"client_uid"`
	Seq       uint64 `json:"seq"`
	Initial   bool   `json:"initial"`
}

type handler struct {
	p     *params
	entry *seqqueue.Entry

	err error
}

const (
	// Timeout = 1 * time.Minute
	Timeout = 5 * time.Second
)

func Handler() http.HandlerFunc {
	return jsonep.Decorate(handle, (*params)(nil))
}

func handle(i interface{}) interface{} {
	zap.S().Info("Received", i)
	h := handler{}
	return h.handle(i)
}

func (h *handler) handle(i interface{}) interface{} {
	h.fetchParams(i)
	h.poll()
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

func (h *handler) poll() {
	if h.err != nil {
		return
	}

	lp := longpoll.DefaultLongpoll()
	q := lp.Register(h.p.ClientUID)
	var c <-chan *seqqueue.Entry
	if h.p.Initial {
		c = q.OutWithoutSeq()
	} else {
		c = q.Out(h.p.Seq)
	}

	select {
	case h.entry = <-c:
	case <-time.After(Timeout):
	}
}

func (h *handler) response() interface{} {
	if h.err != nil {
		return map[string]interface{}{
			"status": -1,
		}
	}

	if h.entry != nil {
		return map[string]interface{}{
			"seq":     h.entry.Seq,
			"message": h.entry.Value,
			"status":  0,
		}
	}

	return map[string]interface{}{
		"status": 0,
	}
}
