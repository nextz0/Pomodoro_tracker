// src/app/page.tsx
"use client"; // จำเป็นต้องใส่เพราะมีการใช้ State

import { useState, useEffect } from "react";
import { Play, Pause, RotateCcw } from "lucide-react"; // ไอคอนที่มีอยู่แล้วในโปรเจกต์

export default function Home() {
  // สถานะเวลา (25 นาที = 1500 วินาที)
  const [timeLeft, setTimeLeft] = useState(5);
  const [isActive, setIsActive] = useState(false);

  // ฟังก์ชันยิง API ไปหา Go Backend
  const saveToBackend = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/save-session", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          duration: 25 * 60,
          status: "completed",
        }),
      });
      const data = await response.json();
      console.log("Saved to backend:", data);
    } catch (error) {
      console.error("Failed to save:", error);
    }
  };

  // แปลงวินาทีเป็น นาที:วินาที
  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  };

  // Logic การนับถอยหลัง
  useEffect(() => {
    let interval: NodeJS.Timeout | null = null;

    if (isActive && timeLeft > 0) {
      interval = setInterval(() => {
        setTimeLeft((prev) => prev - 1);
      }, 1000);
    } else if (timeLeft === 0 && isActive) {
      saveToBackend(); // บันทึกเซสชันเมื่อหมดเวลา
      setIsActive(false);
      alert("Time is up!");
    }

    return () => {
      if (interval) clearInterval(interval);
    };
  }, [isActive, timeLeft]);

  // ฟังก์ชันควบคุม
  const toggleTimer = () => setIsActive(!isActive);
  const resetTimer = () => {
    setIsActive(false);
    setTimeLeft(25 * 60);
  };

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background text-foreground p-8 font-sans">
      <main className="flex flex-col items-center gap-8 text-center">
        
        <h1 className="text-4xl font-bold tracking-tight">Pomodoro Tracker</h1>
        
        {/* ส่วนแสดงเวลา */}
        <div className="text-9xl font-mono font-bold tracking-widest my-8">
          {formatTime(timeLeft)}
        </div>

        {/* ปุ่มควบคุม */}
        <div className="flex gap-4">
          <button
            onClick={toggleTimer}
            className={`flex items-center gap-2 px-8 py-4 rounded-full text-xl font-semibold transition-colors ${
              isActive 
                ? "bg-destructive text-destructive-foreground hover:bg-destructive/90" 
                : "bg-primary text-primary-foreground hover:bg-primary/90"
            }`}
          >
            {isActive ? <Pause /> : <Play />}
            {isActive ? "Pause" : "Start"}
          </button>

          <button
            onClick={resetTimer}
            className="flex items-center gap-2 px-8 py-4 rounded-full text-xl font-semibold bg-secondary text-secondary-foreground hover:bg-secondary/80 transition-colors"
          >
            <RotateCcw />
            Reset
          </button>
        </div>

        <p className="text-muted-foreground mt-8">
          Focus for 25 minutes, then take a break.
        </p>
      </main>
    </div>
  );
}