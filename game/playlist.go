package game

import (
	"sync"

	"git.nulana.com/bobrnor/battleship-server/db/client"
)

type Playlist struct {
	sync.Mutex

	playersLookingForOpponent map[int64]*client.Client
}

var (
	defaultPlaylist Playlist
)

func init() {
	defaultPlaylist = Playlist{
		playersLookingForOpponent: map[int64]*client.Client{},
	}
}

func DefaultPlaylist() *Playlist {
	return &defaultPlaylist
}

func (p *Playlist) PopAny() *client.Client {
	p.Lock()
	defer p.Unlock()

	var c *client.Client
	if len(p.playersLookingForOpponent) > 0 {
		for k, v := range p.playersLookingForOpponent {
			c = v
			delete(p.playersLookingForOpponent, k)
		}
	}
	return c
}

func (p *Playlist) Push(c *client.Client) {
	p.Lock()
	p.playersLookingForOpponent[c.ID] = c
	p.Unlock()
}

func (p *Playlist) Wait(c *client.Client) {
	p.Push(c)
}
