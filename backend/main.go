package main

import (
	"encoding/json"
	"fmt"
	"log" 
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
)

type TimerLog struct {
	Duration int    `json:"duration"`
	Status   string `json:"status"`
}

func saveSessionHandler(w http.ResponseWriter, r *http.Request) {
	// ตั้งค่า CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var logData TimerLog
		
		// รับข้อมูล JSON จาก Frontend
		err := json.NewDecoder(r.Body).Decode(&logData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		supabaseUrl := os.Getenv("SUPABASE_URL")
		supabaseKey := os.Getenv("SUPABASE_KEY")

		if supabaseUrl == "" || supabaseKey == "" {
			fmt.Println("Error: SUPABASE_URL or SUPABASE_KEY not found in .env")
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}

		// สร้าง Client
		client := supabase.CreateClient(supabaseUrl, supabaseKey)

		// บันทึกลง Database
		var results []TimerLog
		err = client.DB.From("sessions").Insert(logData).Execute(&results)

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
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file (อาจจะรันบน Server หรือหาไฟล์ไม่เจอ)")
	}

	http.HandleFunc("/api/save-session", saveSessionHandler)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}