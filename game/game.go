package game

var (
	rooms  *Rooms
	lobby  *Lobby
	engine *Engine
)

func MainLobby() *Lobby {
	if lobby == nil {
		lobby = NewLobby()
	}
	return lobby
}

func MainEngine() *Engine {
	if engine == nil {
		engine = NewEngine()
	}
	return engine
}

func MainRooms() *Rooms {
	if rooms == nil {
		rooms = NewRooms()
	}
	return rooms
}
