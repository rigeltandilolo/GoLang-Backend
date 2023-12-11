package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Note structure
type Note struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var noteStore map[string]Note

func init() {
	noteStore = make(map[string]Note)
}

// Handler to get all notes
func getAllNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result := make([]Note, 0, len(noteStore))
	for _, v := range noteStore {
		result = append(result, v)
	}
	json.NewEncoder(w).Encode(result)
}

// Handler to get note by ID
func getNoteByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	if note, ok := noteStore[id]; ok {
		json.NewEncoder(w).Encode(note)
	} else {
		http.Error(w, "Note not found", http.StatusNotFound)
	}
}

// Handler to add new note
func addNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newNote Note
	err := json.NewDecoder(r.Body).Decode(&newNote)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Generate a unique ID (in a real-world scenario, you might use a UUID library)
	newNote.ID = fmt.Sprintf("%d", len(noteStore)+1)
	noteStore[newNote.ID] = newNote
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newNote)
}

// Handler to update note by ID
func updateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	if _, ok := noteStore[id]; ok {
		var updatedNote Note
		err := json.NewDecoder(r.Body).Decode(&updatedNote)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		noteStore[id] = updatedNote
		json.NewEncoder(w).Encode(updatedNote)
	} else {
		http.Error(w, "Note not found", http.StatusNotFound)
	}
}

// Handler to delete note by ID
func deleteNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	if _, ok := noteStore[id]; ok {
		delete(noteStore, id)
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Note not found", http.StatusNotFound)
	}
}

func main() {
	r := mux.NewRouter()

	// Define endpoints
	r.HandleFunc("/api/notes", getAllNotes).Methods("GET")
	r.HandleFunc("/api/notes/{id}", getNoteByID).Methods("GET")
	r.HandleFunc("/api/notes", addNote).Methods("POST")
	r.HandleFunc("/api/notes/{id}", updateNote).Methods("PUT")
	r.HandleFunc("/api/notes/{id}", deleteNote).Methods("DELETE")

	// Start the server
	port := 8080
	log.Printf("Server started on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
