package search

import (
	"sync"

	"go.uber.org/zap"

	"git.nulana.com/bobrnor/battleship-server/db/client"
	"git.nulana.com/bobrnor/battleship-server/game/rooms"
	longpoll "git.nulana.com/bobrnor/longpoll.git"
)

type Playlist struct {
	sync.Mutex

	playersLookingForOpponent map[string]client.Client
}

func NewPlaylist(lp *longpoll.Longpoll) *Playlist {
	p := Playlist{
		playersLookingForOpponent: map[string]client.Client{},
	}
	lp.SetPurgeFunc(p.purgeFunc)
	return &p
}

func (p *Playlist) Push(c *client.Client) {
	zap.S().Infof("pushing client %+v", c)

	p.Lock()
	p.playersLookingForOpponent[c.UID] = *c
	p.tryToRegisterRoom()
	p.Unlock()
}

func (p *Playlist) purgeFunc(lp *longpoll.Longpoll, i interface{}) {
	zap.S().Infof("lp purged %+v", i)
	p.Lock()
	if uid, ok := i.(string); ok {
		delete(p.playersLookingForOpponent, uid)
	} else {
		zap.S().Errorf("Bad uid type %T %+v", i, i)
	}
	p.Unlock()
}

func (p *Playlist) tryToRegisterRoom() {
	zap.S().Infof("trying to regiter room %+v", p.playersLookingForOpponent)

	if len(p.playersLookingForOpponent) > 1 {
		clients := p.pop2()
		if len(clients) == 2 {
			r := rooms.DefaultRooms()
			if roomUID, err := r.Register(clients); err != nil {
				zap.S().Errorf("Can't register room, push back clients %+v", err)
				p.pushAll(clients)
			} else {
				lp.Send(clients[0].UID, map[string]interface{}{
					"room_uid": roomUID,
					"status":   0,
				})
				lp.Send(clients[1].UID, map[string]interface{}{
					"room_uid": roomUID,
					"status":   0,
				})
			}
		}
	}
}

func (p *Playlist) pop2() []client.Client {
	if len(p.playersLookingForOpponent) > 1 {
		clients := []client.Client{}
		for uid, c := range p.playersLookingForOpponent {
			clients = append(clients, c)
			delete(p.playersLookingForOpponent, uid)
			if len(clients) == 2 {
				break
			}
		}
		return clients
	}
	return []client.Client{}
}

func (p *Playlist) pushAll(clients []client.Client) {
	for _, c := range clients {
		p.playersLookingForOpponent[c.UID] = c
	}
}
