package controllers

import (
	m "UTS/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetAllRooms(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM rooms"

	roomName := r.URL.Query()["room_name"]
	idGame := r.URL.Query()["id_game"]
	if roomName != nil {
		fmt.Print(roomName[0])
		query += " WHERE room_name='" + roomName[0] + "'"
	}

	if idGame != nil {
		if roomName[0] != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " id_game='" + idGame[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var room m.Room
	var rooms []m.Room
	for rows.Next() {
		if err := rows.Scan(&room.ID, &room.RoomName, &room.IDGame); err != nil {
			log.Println(err)
			return
		} else {
			rooms = append(rooms, room)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	var response m.RoomsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = rooms
	json.NewEncoder(w).Encode(response)
}

func GetDetailRooms(w http.ResponseWriter, r *http.Request, roomID int) {
	db := connect()
	defer db.Close()

	var query string
	if roomID != 0 {
		query = "SELECT r.id, r.room_name, r.id_game, p.id, p.id_room, p.id_account FROM rooms r JOIN participants p ON r.id = p.id WHERE r.id = " + strconv.Itoa(roomID)
	} else {
		query = "SELECT r.id, r.room_name, r.id_game, p.id, p.id_room, p.id_account FROM rooms r JOIN participants p ON r.id = p.id"
	}

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Gagal eksekusi query!", http.StatusInternalServerError)
		return
	}

	var participants []m.Participant

	for rows.Next() {
		var room m.Room
		var participant m.Participant

		err := rows.Scan(&room.ID, &room.IDGame, &room.RoomName, &participant.ID, &participant.IDAccount, &participant.IDRoom)
		if err != nil {
			http.Error(w, "Gagal scan baris", http.StatusInternalServerError)
			return
		}
		participant.IDRoom = room
		participants = append(participants, participant)
	}

	w.Header().Set("Content-Type", "application/json")
	var response = m.ParticipantsResponse{}
	response.Status = 200
	response.Message = "Success"
	response.Data = participants
	json.NewEncoder(w).Encode(response)
}

func InsertRoom(w http.ResponseWriter, r *http.Request, roomName string, idGame int) {
	db := connect()
	defer db.Close()

	query := "INSERT INTO rooms (room_name, id_game) VALUES (?,?)"

	var err error
	_, err = db.Exec(query, roomName, idGame)
	if err != nil {
		log.Println(err)
		http.Error(w, "Insert Gagal", http.StatusInternalServerError)
		return
	} else {
		fmt.Println("Berhasil insert ke room.")
	}
}

func LeaveRoom(w http.ResponseWriter, r *http.Request, roomName string) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}
	query := "DELETE FROM rooms WHERE room_name = ?;"
	var response m.RoomResponse
	_, err = db.Exec(query, roomName)
	if err != nil {
		response.Status = 400
		response.Message = "Insert Gagal"
		json.NewEncoder(w).Encode(response)
	} else {
		response.Status = 200
		response.Message = "Success"
		json.NewEncoder(w).Encode(response)
	}
}
