package server

type lobby struct {
	clients map[string]*clientStr
}

func newLobby() lobby {
	return lobby{
		clients: make(map[string]*clientStr),
	}
}
