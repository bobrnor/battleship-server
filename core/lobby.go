package core

import (
	"sync"

	"log"

	"git.nulana.com/bobrnor/battleship-server/db"
	"git.nulana.com/bobrnor/longpoll.git"
)

type Lobby struct {
	sync.Mutex
	clients map[interface{}]db.Client
}

func NewLobby() *Lobby {
	lobby := &Lobby{
		clients: map[interface{}]db.Client{},
	}
	longpoll.DefaultLongpoll().SetPurgeFunc(lobby.HandleLongpollPurge)
	return lobby
}

func (l *Lobby) StartWaitingForRoom(client *db.Client) {
	log.Printf("adding client %+v", client)

	l.Lock()
	longpoll.DefaultLongpoll().Register(client.UID)

	l.clients[client.UID] = *client
	if len(l.clients) > 1 {
		l.createRoom()
	}
	l.Unlock()
}

func (l *Lobby) StopWaitingForRoom(uid interface{}) {
	log.Printf("deleting client %+v", uid)

	l.Lock()
	delete(l.clients, uid)
	l.Unlock()
}

func (l *Lobby) createRoom() {
	log.Printf("trying to regiter room")

	clients := l.fetchClientsForRoom()
	rooms := MainRooms()

	roomUID, err := rooms.Register(clients)
	if err != nil {
		log.Printf("Can't register room %+v", err)
		return
	}

	msg := map[string]interface{}{
		"type":     "search_result",
		"room_uid": roomUID,
	}
	if err := l.notifyClients(clients, msg); err != nil {
		log.Printf("Can't send notify clients about founded room %+v", err)
	}
	l.removeClients(clients)
}

func (l *Lobby) fetchClientsForRoom() []db.Client {
	clients := []db.Client{}
	for _, c := range l.clients {
		clients = append(clients, c)
		if len(clients) == 2 {
			break
		}
	}
	return clients
}

func (l *Lobby) notifyClients(clients []db.Client, msg map[string]interface{}) error {
	for _, c := range clients {
		if err := longpoll.DefaultLongpoll().Send(c.UID, msg); err != nil {
			return err
		}
	}
	return nil
}

func (l *Lobby) removeClients(clients []db.Client) {
	for _, c := range clients {
		l.clients[c.UID] = c
	}
}

func (l *Lobby) HandleLongpollPurge(lp *longpoll.Longpoll, i interface{}) {
	l.StopWaitingForRoom(i)
}
