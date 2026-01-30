package main

import "github.com/youssef28m/LockIn/internal/storage"

func main() {
	
	db := storage.Connect()

	defer db.Close()
	storage.CreateDB()


}
