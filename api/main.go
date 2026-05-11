package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
}


// Vi håller de rätta svaren hemliga här i API:et
var correctAnswers = map[int]string{
	1: "En gopher (kindpåsråtta)",
	2: "Goroutines",
	3: "Italien",
	4: "Venus",
	5: "3",
	6: "Pascal",
}

// Här sparar vi alla tidigare poäng för att kunna räkna ut statistik
var allScores []int

func getQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var questions = []Question{
    {
        ID:      1,
        Text:    "Vilket djur är Gos officiella maskot?",
        Options: []string{"En hamster", "En gopher (kindpåsråtta)", "En bäver", "En ekorre"},
    },
    {
        ID:      2,
        Text:    "Vad kallas det när man kör flera funktioner samtidigt i Go?",
        Options: []string{"Threads", "Async-await", "Goroutines", "Parallel-loops"},
    },
    {
        ID:      3,
        Text:    "Vilket land uppfann tekniskt sett pizzan?",
        Options: []string{"Italien", "Grekland", "USA", "Kina"},
    },
    {
        ID:      4,
        Text:    "Vilken planet i vårt solsystem är faktiskt den varmaste?",
        Options: []string{"Merkurius", "Venus", "Mars", "Jupiter"},
    },
    {
        ID:      5,
        Text:    "Hur många hjärtan har en bläckfisk?",
        Options: []string{"1", "2", "3", "8"},
    },
    {
        ID:      6,
        Text:    "Vilket av dessa språk är Go INTE inspirerat av?",
        Options: []string{"C", "Pascal", "Python", "Oberon"},
    },
}
	json.NewEncoder(w).Encode(questions)
}

func submitAnswers(w http.ResponseWriter, r *http.Request) {
	var userAnswers map[int]string // Format: {1: "Stockholm", 2: "Go"}
	json.NewDecoder(r.Body).Decode(&userAnswers)

	score := 0
	for id, ans := range userAnswers {
		if correctAnswers[id] == ans {
			score++
		}
	}

	// Räkna ut statistik
	betterThan := 0
	for _, s := range allScores {
		if score > s {
			betterThan++
		}
	}
	
	percentile := 0.0
	if len(allScores) > 0 {
		percentile = (float64(betterThan) / float64(len(allScores))) * 100
	}

	allScores = append(allScores, score)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Du fick % d rätt! Du var bättre än %.0f%% av alla som svarat.", score, percentile)
}

func main() {
	http.HandleFunc("/questions", getQuestions)
	http.HandleFunc("/submit", submitAnswers)

	fmt.Println("API:et körs på http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}