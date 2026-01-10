package handlers

import (
	// "database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"backend/internal/db"
	// "backend/internal/model"
)

func (h *Handler) GetAllRoomsHandler(w http.ResponseWriter, r *http.Request) {
	// check if request is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// get all rooms
	allRooms, err := db.GetRooms(h.MsgDB)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(allRooms)

}

func (h *Handler) GetRoomData(w http.ResponseWriter, r *http.Request) {
	// check if request is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Page doesn't exist",http.StatusBadRequest)
		return
	}

	roomID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid room id", http.StatusBadRequest)
		return
	}
	
	allMessages,err:=db.GetMsgOfRoom(h.MsgDB,roomID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(allMessages)

}
