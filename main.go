package main

import (
	"UTS/controllers"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Get All Rooms
	router.HandleFunc("/rooms", controllers.GetAllRooms).Methods("GET")
	// Get Detail Rooms
	router.HandleFunc("/detailrooms", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetDetailRooms(w, r, 1)
	}).Methods("GET")
	// Insert Room
	router.HandleFunc("/insertroom", func(w http.ResponseWriter, r *http.Request) {
		controllers.InsertRoom(w, r, "Ruang Datuk", 3)
	}).Methods("POST")
	// Leave Room
	router.HandleFunc("/leaveroom", func(w http.ResponseWriter, r *http.Request) {
		controllers.LeaveRoom(w, r, "Ruang Delima")
	}).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
