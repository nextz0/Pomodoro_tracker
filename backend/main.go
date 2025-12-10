package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TimerLog struct {
	Duration int    `json:"duration"`
	Status   string `json:"status"`
}

func saveSessionHandler(w http.ResponseWriter, r *http.Request) {
	// ตั้งค่า CORS (สำคัญมาก ไม่งั้น Frontend จะยิงไม่เข้า)
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var log TimerLog
		// รับข้อมูล JSON จาก Frontend
		err := json.NewDecoder(r.Body).Decode(&log)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// (ในอนาคต) บันทึกลง Database ตรงนี้
		fmt.Printf("บันทึกเซสชัน: %d วินาที สถานะ: %s\n", log.Duration, log.Status)

		// ตอบกลับ Frontend
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Success"})
	}
}

func main() {
	http.HandleFunc("/api/save-session", saveSessionHandler)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}