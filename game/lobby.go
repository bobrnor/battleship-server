package game

import (
	"sync"

	"git.nulana.com/bobrnor/battleship-server/db"
	"git.nulana.com/bobrnor/longpoll.git"
	"go.uber.org/zap"
)

type Lobby struct {
	sync.Mutex
	clients map[interface{}]db.Client
	rooms   *Rooms
}

const (
	RoomFoundMessage = "room_found"
)

var (
	lobby *Lobby
)

func MainLobby() *Lobby {
	if lobby == nil {
		lobby = &Lobby{
			clients: map[interface{}]db.Client{},
			rooms:   &Rooms{},
		}
	}
	return lobby
}

func (l *Lobby) StartWaitingForRoom(client *db.Client) {
	zap.S().Infof("adding client %+v", client)

	l.Lock()
	l.clients[client.UID] = *client
	if len(l.clients) > 1 {
		l.createRoom()
	}
	l.Unlock()
}

func (l *Lobby) StopWaitingForRoom(uid interface{}) {
	zap.S().Infof("deleting client %+v", uid)

	l.Lock()
	delete(l.clients, uid)
	l.Unlock()
}

func (l *Lobby) createRoom() {
	zap.S().Infof("trying to regiter room")

	clients := l.fetchClientsForRoom()
	if roomUID, err := l.rooms.Register(clients); err != nil {
		zap.S().Errorf("Can't register room %+v", err)
	} else {
		msg := map[string]interface{}{
			"msg":      RoomFoundMessage,
			"room_uid": roomUID,
			"status":   0,
		}
		if err := l.notifyClients(clients, msg); err != nil {
			zap.S().Errorf("Can't send notify clients about founded room %+v", err)
		}
		l.removeClients(clients)
	}
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
