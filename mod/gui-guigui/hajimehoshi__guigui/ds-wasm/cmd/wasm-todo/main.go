//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"syscall/js"
	"time"
)

// Todo represents a single todo item
type Todo struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
}

// TodoStore manages the todo list state
type TodoStore struct {
	Todos       []Todo `json:"todos"`
	Filter      string `json:"filter"` // "all", "active", "completed"
	NewTodoText string `json:"newTodoText"`
	NextID      int    `json:"nextId"`
}

var store = &TodoStore{
	Todos:       []Todo{},
	Filter:      "all",
	NewTodoText: "",
	NextID:      1,
}

func main() {
	fmt.Println("🌟 DataStar Todo WASM - Starting...")

	// Load todos from localStorage on startup
	loadTodos()

	// Register WASM functions that can be called from JavaScript
	js.Global().Set("addTodo", js.FuncOf(addTodo))
	js.Global().Set("toggleTodo", js.FuncOf(toggleTodo))
	js.Global().Set("deleteTodo", js.FuncOf(deleteTodo))
	js.Global().Set("setFilter", js.FuncOf(setFilter))
	js.Global().Set("updateNewTodoText", js.FuncOf(updateNewTodoText))
	js.Global().Set("clearCompleted", js.FuncOf(clearCompleted))
	js.Global().Set("getTodos", js.FuncOf(getTodos))

	fmt.Println("✅ Todo WASM functions registered")
	fmt.Println("📝 Todo app ready!")

	// Keep the program running
	select {}
}

// addTodo adds a new todo item
func addTodo(this js.Value, args []js.Value) interface{} {
	if store.NewTodoText == "" {
		return nil
	}

	newTodo := Todo{
		ID:        store.NextID,
		Text:      store.NewTodoText,
		Completed: false,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	store.Todos = append(store.Todos, newTodo)
	store.NextID++
	store.NewTodoText = ""

	saveTodos()
	updateUI()
	return nil
}

// toggleTodo toggles the completed status of a todo
func toggleTodo(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return nil
	}

	id, err := strconv.Atoi(args[0].String())
	if err != nil {
		return nil
	}

	for i := range store.Todos {
		if store.Todos[i].ID == id {
			store.Todos[i].Completed = !store.Todos[i].Completed
			break
		}
	}

	saveTodos()
	updateUI()
	return nil
}

// deleteTodo removes a todo item
func deleteTodo(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return nil
	}

	id, err := strconv.Atoi(args[0].String())
	if err != nil {
		return nil
	}

	for i, todo := range store.Todos {
		if todo.ID == id {
			store.Todos = append(store.Todos[:i], store.Todos[i+1:]...)
			break
		}
	}

	saveTodos()
	updateUI()
	return nil
}

// setFilter changes the current filter
func setFilter(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return nil
	}

	filter := args[0].String()
	if filter == "all" || filter == "active" || filter == "completed" {
		store.Filter = filter
		updateUI()
	}
	return nil
}

// updateNewTodoText updates the new todo input text
func updateNewTodoText(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return nil
	}

	store.NewTodoText = args[0].String()
	return nil
}

// clearCompleted removes all completed todos
func clearCompleted(this js.Value, args []js.Value) interface{} {
	var activeTodos []Todo
	for _, todo := range store.Todos {
		if !todo.Completed {
			activeTodos = append(activeTodos, todo)
		}
	}
	store.Todos = activeTodos

	saveTodos()
	updateUI()
	return nil
}

// getTodos returns the current todos (for debugging)
func getTodos(this js.Value, args []js.Value) interface{} {
	data, _ := json.Marshal(store)
	return string(data)
}

// updateUI sends the current state to DataStar for UI updates
func updateUI() {
	// Filter todos based on current filter
	var filteredTodos []Todo
	for _, todo := range store.Todos {
		switch store.Filter {
		case "active":
			if !todo.Completed {
				filteredTodos = append(filteredTodos, todo)
			}
		case "completed":
			if todo.Completed {
				filteredTodos = append(filteredTodos, todo)
			}
		default: // "all"
			filteredTodos = append(filteredTodos, todo)
		}
	}

	// Count stats
	totalCount := len(store.Todos)
	activeCount := 0
	completedCount := 0
	for _, todo := range store.Todos {
		if todo.Completed {
			completedCount++
		} else {
			activeCount++
		}
	}

	// Prepare data for DataStar
	uiData := map[string]interface{}{
		"todos":          filteredTodos,
		"filter":         store.Filter,
		"newTodoText":    store.NewTodoText,
		"totalCount":     totalCount,
		"activeCount":    activeCount,
		"completedCount": completedCount,
		"hasCompleted":   completedCount > 0,
	}

	// Convert to JSON and send to DataStar
	jsonData, err := json.Marshal(uiData)
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return
	}

	// Call DataStar merge signals function via JavaScript
	js.Global().Call("updateDataStarStore", string(jsonData))
}

// saveTodos saves todos to localStorage
func saveTodos() {
	data, err := json.Marshal(store)
	if err != nil {
		fmt.Printf("Error saving todos: %v\n", err)
		return
	}

	localStorage := js.Global().Get("localStorage")
	localStorage.Call("setItem", "datastar-todos", string(data))
}

// loadTodos loads todos from localStorage
func loadTodos() {
	localStorage := js.Global().Get("localStorage")
	data := localStorage.Call("getItem", "datastar-todos")

	if data.IsNull() || data.IsUndefined() {
		// Initialize with sample todos
		store.Todos = []Todo{
			{
				ID:        1,
				Text:      "Learn DataStar",
				Completed: true,
				CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			},
			{
				ID:        2,
				Text:      "Build WASM Todo App",
				Completed: false,
				CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			},
			{
				ID:        3,
				Text:      "Deploy to production",
				Completed: false,
				CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			},
		}
		store.NextID = 4
		saveTodos()
		return
	}

	err := json.Unmarshal([]byte(data.String()), store)
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		// Reset to default if corrupted
		store.Todos = []Todo{}
		store.NextID = 1
	}

	fmt.Printf("📝 Loaded %d todos from localStorage\n", len(store.Todos))
}
