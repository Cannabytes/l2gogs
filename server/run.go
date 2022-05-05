package server

import (
	"fmt"
	"l2gogameserver/data/logger"
	"l2gogameserver/gameserver"
	"l2gogameserver/gameserver/handlers"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/models/stats/regeneration"
	"net"
	"time"
)

type GameServer struct {
	clientsListener net.Listener
	//OnlineCharacters *models.OnlineCharacters
	//clients          sync.Map
}

func New() *GameServer {
	return &GameServer{}
}

var onlineChars models.OnlineCharacters

func (g *GameServer) Start() {
	var err error
	/* #nosec */
	g.clientsListener, err = net.Listen("tcp4", ":7777")
	if err != nil {
		logger.Error.Panicln(err.Error())
	}
	onlineChars.Char = make(map[int32]*models.Character)
	gameserver.OnlineCharacters = &onlineChars

	//go g.Tick()
	defer g.clientsListener.Close()
	for {
		client := models.NewClient()
		client.Socket, err = g.clientsListener.Accept()

		if err != nil {
			fmt.Println("Couldn't accept the incoming connection.", err)
			continue
		}

		//g.AddClient(client) //todo надо ли добавлять клиентов в отдельную мапу или массив?
		go handlers.Handler(client)
	}
}

//События
func Events() {
	for {
		for _, character := range onlineChars.Char {
			regeneration.RenerationHpMpCp(character.Conn)
		}
		time.Sleep(time.Second)
	}
}
