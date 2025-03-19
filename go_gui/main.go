package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Store active WebSocket connections
var clients = make(map[*websocket.Conn]struct{})
var clientsLock sync.Mutex

// List of Python scripts
var scripts = []string{"script1.py", "script2.py"}

// Serve the HTML page
func serveHome(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.New("index").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Script Runner</title>
			<script>
				let socket = new WebSocket("ws://" + window.location.host + "/ws");
				socket.onmessage = function(event) {
					document.getElementById("output").innerText += event.data + "\n";
				};
				
				function startScript(scriptName) {
					socket.send(scriptName);
				}
			</script>
		</head>
		<body>
			<h1>Script Runner</h1>
			{{range .}}
				<button onclick="startScript('{{.}}')">{{.}}</button>
			{{end}}
			<pre id="output" style="border:1px solid #000; padding:10px; width:80%; height:300px; overflow:auto;"></pre>
		</body>
		</html>
	`)
	tmpl.Execute(w, scripts)
}

// Run a Python script and stream output to all connected WebSockets
func runPythonScript(script string) {
	cmd := exec.Command("python", script)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		sendToAllClients("[Error]: Failed to start " + script)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			sendToAllClients(script + ": " + scanner.Text())
		}
	}()

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			sendToAllClients("[Error] " + script + ": " + scanner.Text())
		}
	}()

	wg.Wait()
	cmd.Wait()
}

// Send a message to all connected clients
func sendToAllClients(msg string) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	for client := range clients {
		client.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

// WebSocket handler
func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[Error]: WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	// Register client
	clientsLock.Lock()
	clients[conn] = struct{}{}
	clientsLock.Unlock()

	defer func() {
		clientsLock.Lock()
		delete(clients, conn)
		clientsLock.Unlock()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		script := string(msg)
		go runPythonScript(script)
	}
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleWS)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
