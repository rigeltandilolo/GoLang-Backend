package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Note struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var noteStore map[string]Note

func init() {
	noteStore = make(map[string]Note)
}

func getAllNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result := make([]Note, 0, len(noteStore))
	for _, v := range noteStore {
		result = append(result, v)
	}
	json.NewEncoder(w).Encode(result)
}

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

func addNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newNote Note
	err := json.NewDecoder(r.Body).Decode(&newNote)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newNote.ID = fmt.Sprintf("%d", len(noteStore)+1)
	noteStore[newNote.ID] = newNote
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newNote)
}

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

	r.HandleFunc("/api/notes", getAllNotes).Methods("GET")
	r.HandleFunc("/api/notes/{id}", getNoteByID).Methods("GET")
	r.HandleFunc("/api/notes", addNote).Methods("POST")
	r.HandleFunc("/api/notes/{id}", updateNote).Methods("PUT")
	r.HandleFunc("/api/notes/{id}", deleteNote).Methods("DELETE")

	port := 8080
	log.Printf("Server started on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
