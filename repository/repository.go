package repository

import (
	"database/sql"
	"time"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/microcosm-cc/bluemonday"
)


type TodoList struct {
	ID        int64
	Hash      string
	Todos     []Todo
	CreatedAt string
	UpdatedAt string
}

type Todo struct {
	ID           int64
	TodoListHash string
	TodoListID   int64
	Text         string
	Done         int
	CreatedAt    string
	UpdatedAt    string
}

var Database *sql.DB

func InitDB(databaseName string) {
	var err error
	Database, err = sql.Open("sqlite3", databaseName)

	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
}

func Migrate() {
	var err error

	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS todos_lists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hash TEXT NOT NULL,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS todos_todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		list_id INTEGER,
		text TEXT NOT NULL,
		is_done INTEGER NOT NULL,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		FOREIGN KEY(list_id) REFERENCES todos_lists(id)
	);
	`
	_, err = Database.Exec(createTablesSQL)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	log.Println("Migrations ran successfully!")
}

func FindTodoListByHash(hash string) (*TodoList, error) {
	todoList := &TodoList{}
	query := "SELECT id, hash, created_at, updated_at FROM todos_lists WHERE hash = ?"

	err := Database.QueryRow(query, hash).Scan(&todoList.ID, &todoList.Hash, &todoList.CreatedAt, &todoList.UpdatedAt)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	rows, err := Database.Query("SELECT id, text, is_done FROM todos_todos WHERE list_id = ?", todoList.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Text, &todo.Done); err != nil {
			return nil, err
		}
		todo.TodoListID = todoList.ID
		todo.TodoListHash = todoList.Hash
		todoList.Todos = append(todoList.Todos, todo)
	}

	return todoList, nil
}

func CreateList() (*TodoList, error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	todoList := &TodoList{
		Hash:      generateHash(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	result, err := Database.Exec("INSERT INTO todos_lists (hash, created_at, updated_at) VALUES (?, ?, ?)", todoList.Hash, todoList.CreatedAt, todoList.UpdatedAt)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	todoList.ID = id

	return todoList, nil
}

func AddTodoToListByHash(hash string, text string) (*Todo, error) {
	list, err := FindTodoListByHash(hash)
	if err != nil {
		return nil, err
	}

	policy := bluemonday.UGCPolicy()
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	todo := &Todo{
		Text:      policy.Sanitize(text),
		Done:      0,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}


	query := "INSERT INTO todos_todos (list_id, text, is_done, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	result, err := Database.Exec(query, list.ID, todo.Text, todo.Done, todo.CreatedAt, todo.UpdatedAt)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	todo.ID = id
	todo.TodoListHash = hash

	return todo, nil
}

func FindTodo(id int64) (*Todo, error) {
	todo := &Todo{}
	query := "SELECT id, list_id, text, is_done, created_at, updated_at FROM todos_todos WHERE id = ?"

	err := Database.QueryRow(query, id).Scan(&todo.ID, &todo.TodoListID, &todo.Text, &todo.Done, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return todo, nil
}

func UpdateTodo(id int64, done int, text string) (*Todo, error) {
	currentTime := string(time.Now().Format("2006-01-02 15:04:05"))
	query := "UPDATE todos_todos SET is_done = ?, text = ?, updated_at = ? WHERE id = ?"

	_, err := Database.Exec(query, done, text, currentTime, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	todo, _ := FindTodo(id)

	return todo, nil
}

func DeleteTodo(id int64) error {
	query := "DELETE FROM todos_todos WHERE id = ?"

	_, err := Database.Exec(query, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
