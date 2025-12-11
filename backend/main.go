package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nedpals/supabase-go"
)

const supabaseUrl = "https://ixreecjqnzroibwvuzxx.supabase.co"
const supabaseKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Iml4cmVlY2pxbnpyb2lid3Z1enh4Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjU0MzQ3NjAsImV4cCI6MjA4MTAxMDc2MH0.8_ULIKMHEOoOdNi_pNcHcg_rYEIqaKofAtuUA61lSJg"

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

		client := supabase.CreateClient(supabaseUrl, supabaseKey)

		var results []TimerLog
		err = client.DB.From("sessions").Insert(log).Execute(&results)

		if err != nil {
			fmt.Println("Error inserting to Supabase:", err)
			http.Error(w, "Failed to save to database", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Saved to Supabase: %+v\n", results)

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