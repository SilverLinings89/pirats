package main

import (
    "fmt"
    "math/rand"
    "net/http"
    "sync"
    "time"

    "github.com/gorilla/websocket"
    "github.com/google/uuid"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Player struct {
    Id     string
    Conn   *websocket.Conn
    Name   string
    InGame bool
}

type GameState struct {
    Id  string
    BallX, BallY        float64
    Paddle1Y, Paddle2Y  float64
    PaddleWidth, PaddleHeight float64
    BallSpeedX, BallSpeedY    float64
    PlayerLeft, PlayerRight string
    GameOver bool
}

var players = make(map[string]*Player)
var games = make(map[string]*GameState)
var broadcast = make(chan Message)
var mutex = &sync.Mutex{}

func getPlayerForConnection(conn *websocket.Conn) (Player, bool) {
    for _, player := range players {
        if player.Conn == conn {
            return *player, true
        }
    }
    return Player{}, false
}

func getGameForPlayer(id string) (GameState, bool) {
    for _, game:= range games {
        if game.PlayerLeft == id || game.PlayerRight == id {
            return *game, true
        }
    }
    return GameState{}, false
}

type Message struct {
    Type    string `json:"type"`
    Content interface{} `json:"content"`
}

func main() {
    fs := http.FileServer(http.Dir("./site/build"))
    http.Handle("/", fs)
    http.HandleFunc("/ws", handleWebSocket)

    port := ":8080"
    fmt.Printf("Server running at http://localhost%s\n", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("WebSocket upgrade failed:", err)
        return
    }
    defer conn.Close()

    player := &Player{
        Id: uuid.NewString(),
        Conn:  conn,
        Name:  fmt.Sprintf("Player%d", rand.Intn(1000)), // Generate a random name
        InGame: false,
    }

    mutex.Lock()
    players[player.Id] = player
    fmt.Printf("Player %s connected\n", player.Name)
    mutex.Unlock()

    msg := Message{Type: "welcome", Content: player.Name}
    conn.WriteJSON(msg) 

    broadcastPlayerList()

    for {
        var incoming Message
        err := conn.ReadJSON(&incoming)
        if err != nil {
            fmt.Println("Player disconnected:", player.Name)
            removePlayer(conn)
            broadcastPlayerList()
            break
        }

        switch incoming.Type {
        case "challenge":
            handleChallenge(conn, incoming.Content)
        case "move":
            handleMove(conn, incoming.Content)
        }
    }
}

func broadcastPlayerList() {
    mutex.Lock()
    defer mutex.Unlock()

    var playerNames []string
    for _, player := range players {
        playerNames = append(playerNames, player.Name)
    }

    msg := Message{Type: "players", Content: playerNames}

    for _, player := range players {
        if err := player.Conn.WriteJSON(msg); err != nil {
            fmt.Println("Error sending player list:", err)
        }
    }
}

func removePlayer(conn *websocket.Conn) {
    mutex.Lock()
    p, success := getPlayerForConnection(conn)
    if success {
        delete(players, p.Id)
        g, success := getGameForPlayer(p.Id)
        if success {
            delete(games, g.Id)
        }
    }
    mutex.Unlock()
}

func handleChallenge(conn *websocket.Conn, t interface{}) {
    mutex.Lock()
    defer mutex.Unlock()
    targetName, ok := t.(string)

    if !ok {
        fmt.Println("Error: challenge target is not a string")
        return
    }

    for _, target := range players {
        if target.Name == targetName && !target.InGame {
            p1, successP1 := getPlayerForConnection(conn)
            p2, successP2 := getPlayerForConnection(target.Conn)
            if successP1 && successP2 {
                id := uuid.NewString()

                games[id] = &GameState{
                    Id: id,
                    BallX: 50, BallY: 50, Paddle1Y: 40, Paddle2Y: 40,
                    BallSpeedX: 1.5, BallSpeedY: 1.2, PaddleWidth: 10, PaddleHeight: 20,
                    PlayerLeft: p1.Id, 
                    PlayerRight: p2.Id,
                }

                players[target.Id].InGame = true
                players[p1.Id].InGame = true

                conn.WriteJSON(Message{Type: "start", Content: "Game started!"})
                target.Conn.WriteJSON(Message{Type: "start", Content: "Game started!"})

                go runGame(id)
            }
            return
        }
    }

    conn.WriteJSON(Message{Type: "error", Content: "Target not available"})
}

func runGame(id string) {

    ticker := time.NewTicker(16 * time.Millisecond)
    defer ticker.Stop()

    state := games[id]
    
    for range ticker.C {
        state.BallX += state.BallSpeedX
        state.BallY += state.BallSpeedY

        if state.BallY <= 0 || state.BallY >= 100 {
            state.BallSpeedY *= -1
        }
	if checkGameOver(state) {
    		msg:= Message{
    			Type: "gameOver", 
    			Content: "The game has ended",
    		}
    		mutex.Lock()
    		players[state.PlayerRight].Conn.WriteJSON(msg)
    		players[state.PlayerLeft].Conn.WriteJSON(msg)
    		mutex.Unlock()
    		return
    	}

        mutex.Lock()
		msg:= Message{
			Type: "gameState",
			Content: state,
		}
    		players[state.PlayerRight].Conn.WriteJSON(msg)
    		players[state.PlayerLeft].Conn.WriteJSON(msg)
        mutex.Unlock()
    }
}

func checkGameOver(gameState *GameState) bool {
    if gameState.BallX < 0 || gameState.BallX > 800 {
        return true 
    }
    return false
}

func handleMove(conn *websocket.Conn, moveDirection interface{}) {
    player, exists := getPlayerForConnection(conn)
    if !exists {
        fmt.Println("Player not found")
        return
    }
    
    gameState, exists := getGameForPlayer(player.Id)
    if !exists {
        fmt.Println("Game not found")
        return
    }

    direction, ok := moveDirection.(string)
    if !ok {
        fmt.Println("Invalid move direction")
        return
    }

    if player.Name == "Player1" {
        if direction == "up" {
            gameState.Paddle1Y -= 5
        } else if direction == "down" {
            gameState.Paddle1Y += 5
        }
    } else if player.Name == "Player2" {
        if direction == "up" {
            gameState.Paddle2Y -= 5
        } else if direction == "down" {
            gameState.Paddle2Y += 5
        }
    }

    if gameState.Paddle1Y < 0 {
        gameState.Paddle1Y = 0
    } else if gameState.Paddle1Y > 380 {
        gameState.Paddle1Y = 380
    }

    if gameState.Paddle2Y < 0 {
        gameState.Paddle2Y = 0
    } else if gameState.Paddle2Y > 380 {
        gameState.Paddle2Y = 380
    }

    if gameState.BallX < 0 || gameState.BallX > 800 {
        gameState.GameOver = true
        broadcastGameOver(gameState.Id)
    }

    broadcastGameState(gameState.Id)
}

func broadcastGameState(gameID string) {
    gameState := games[gameID]

    mutex.Lock()
    msg := Message{
         Type:    "gameState",
         Content: gameState,
    }
    players[gameState.PlayerLeft].Conn.WriteJSON(msg)
    players[gameState.PlayerRight].Conn.WriteJSON(msg)

    mutex.Unlock()
}

func broadcastGameOver(gameID string) {
    mutex.Lock()

    gameState := games[gameID]
    msg := Message{
        Type:    "gameOver",
        Content: "Game Over! The ball left the field.",
    }
    if err := players[gameState.PlayerLeft].Conn.WriteJSON(msg); err != nil {
        fmt.Println("Error sending game over message:", err)
    }
    if err := players[gameState.PlayerRight].Conn.WriteJSON(msg); err != nil {
        fmt.Println("Error sending game over message:", err)
    }

    mutex.Unlock()
}

