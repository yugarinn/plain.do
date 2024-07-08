package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/yugarinn/plain.do/repository"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)


var (
	clientsMutex sync.Mutex
	clients      = make(map[*websocket.Conn]string)
	lists        = make(map[string]map[*websocket.Conn]bool)
	upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WsIngressMessage struct {
	ListID string `json:"listId"`
	TodoID string `json:"todoId"`
	Hash   string `json:"hash"`
	Text   string `json:"text"`
	Done   string `json:"done"`
	Action string `json:"action"`
}

func renderPage(w http.ResponseWriter, templateName string, context any) {
	tmpl, err := template.New("").ParseFiles(
		filepath.Join("templates", templateName),
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", "todo.html"),
		filepath.Join("templates", "todo-input.html"),
	)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.ExecuteTemplate(w, "base", context)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    hash := vars["listHash"]
	list, _ := repository.FindTodoListByHash(hash)

	if list == nil {
		list, _ = repository.CreateList()
	}

	renderPage(w, "index.html", list)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "about.html", nil)
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	listID := r.URL.Query().Get("listId")
	if listID == "" {
		log.Println("No list ID specified")
		conn.Close()
		return
	}

	clientsMutex.Lock()
	clients[conn] = listID
	if _, ok := lists[listID]; !ok {
		lists[listID] = make(map[*websocket.Conn]bool)
	}
	lists[listID][conn] = true
	clientsMutex.Unlock()

	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		delete(lists[listID], conn)
		if len(lists[listID]) == 0 {
			delete(lists, listID)
		}
		clientsMutex.Unlock()
		conn.Close()
	}()

	for {
		var message WsIngressMessage
		err := conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("WebSocket connection closed normally")
				return
			}
			log.Println("Error reading message:", err)
			return
		}

		if message.Action == "addTodo" {
			addTodoHandler(w, message)
		}

		if message.Action == "updateTodo" {
			updateTodoHandler(w, message)
		}
	}
}

func addTodoHandler(w http.ResponseWriter, message WsIngressMessage) {
	todo, err := repository.AddTodoToListByHash(message.ListID, message.Text)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "todo.html"),
		filepath.Join("templates", "todo-input.html"),
	)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "todo", todo)
	if err != nil {
		log.Println(err)
		return
	}

	html := buf.String()

	broadcastAddTodo(message.ListID, html)
}

func broadcastAddTodo(listID, html string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range lists[listID] {
		err := client.WriteMessage(websocket.TextMessage, []byte(html))
		if err != nil {
			log.Println("Error writing message to client:", err)
			client.Close()
			delete(clients, client)
			delete(lists[listID], client)
		}
	}
}

func updateTodoHandler(w http.ResponseWriter, message WsIngressMessage) {
	text := message.Text
	done := 0

	doneValue := message.Done
	if doneValue == "on" {
		done = 1
	}

    todoID, err := strconv.ParseInt(message.TodoID, 10, 64)
    if err != nil {
		log.Println(err)
        return
    }

	log.Println(message)

	if (text == "") {
		repository.DeleteTodo(todoID)
		broadcastDeleteTodo(message.ListID, todoID)
		return
	}

	todo, _ := repository.UpdateTodo(todoID, done, text)
	todo.TodoListHash = message.ListID

	tmpl, err := template.New("").ParseFiles(filepath.Join("templates", "todo-input.html"))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "todo-input", todo)
	if err != nil {
		log.Println(err)
		return
	}

	html := buf.String()

	broadcastUpdateTodo(message.ListID, html)
}

func broadcastUpdateTodo(listID, html string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range lists[listID] {
		err := client.WriteMessage(websocket.TextMessage, []byte(html))
		if err != nil {
			log.Println("Error writing message to client:", err)
			client.Close()
			delete(clients, client)
			delete(lists[listID], client)
		}
	}
}

func broadcastDeleteTodo(listID string, todoID int64) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	response := map[string]interface{}{
		"action": "deleteTodo",
		"todoID": todoID,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling JSON response:", err)
		return
	}

	for client := range lists[listID] {
		err := client.WriteMessage(websocket.TextMessage, jsonResponse)
		if err != nil {
			log.Println("Error writing message to client:", err)
			client.Close()
			delete(clients, client)
			delete(lists[listID], client)
		}
	}
}
