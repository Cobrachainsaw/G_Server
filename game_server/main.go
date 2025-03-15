package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
)

type PlayerSession struct {
	sessionID int
	clientID  int
	username  string
	inLobby   bool
	conn      *websocket.Conn
}

func newPlayerSession(sid int, conn *websocket.Conn) actor.Producer {
	return func() actor.Receiver {
		return &PlayerSession{
			sessionID: sid,
			conn:      conn,
		}
	}
}

func (s *PlayerSession) Receive(c *actor.Context) {

}

type GameServer struct {
	ctx      *actor.Context
	sessions map[*actor.PID]struct{}
}

func newGameServer() actor.Receiver {
	return &GameServer{
		sessions: make(map[*actor.PID]struct{}),
	}
}

func (s *GameServer) Receive(c *actor.Context) {
	if s.ctx == nil {
		s.ctx = c // Assign the actor context when it first starts
	}

	switch msg := c.Message().(type) {
	case actor.Started:
		s.startHTTP()
		_ = msg
	}
}

func (s *GameServer) startHTTP() {
	fmt.Println("Starting HTTP server on port 40000")
	go func() {
		http.HandleFunc("/ws", s.handleWS)
		http.ListenAndServe(":40000", nil)
	}()
}

func (s *GameServer) handleWS(w http.ResponseWriter, r *http.Request) {
	if s.ctx == nil {
		log.Println("GameServer context is nil!")
		return
	}
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("ws upgrade error:", err)
		return
	}
	fmt.Println("New Client trying to connect")
	sid := rand.Intn(math.MaxInt)
	pid := s.ctx.SpawnChild(newPlayerSession(sid, conn), fmt.Sprintf("session_%d", sid))
	s.sessions[pid] = struct{}{}
	fmt.Printf("Client with sid: %d and pid: %s just connected\n", sid, pid)
}

func main() {
	config := actor.NewEngineConfig()
	e, err := actor.NewEngine(config)
	if err != nil {
		log.Fatal(err)
	}
	e.Spawn(newGameServer, "server")
	select {}
}

type HTTPServer struct {
}
