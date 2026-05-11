package main

import (
	"bytes" // För att skicka JSON-data i POST-request
	"encoding/json"
	"fmt" // I/O
	"io"  // För att läsa svar från API:et
	"net/http"
	"os"
	"os/exec" // För att rensa konsolen
	"runtime" // För att kolla operativsystemet

	"github.com/spf13/cobra"
)

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
}

func main() {
	// Root-kommandot (det som körs om man bara skriver namnet på programmet)
	var rootCmd = &cobra.Command{
		Use:   "quiz",
		Short: "Ett quiz!",
	}

	// Kommandot 'start' (det som körs när man skriver 'quiz start')
	var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Startar quizet",
    Run: func(cmd *cobra.Command, args []string) {

    //Hämta frågorna
    resp, err := http.Get("http://localhost:8080/questions")

	if err != nil {
	fmt.Println("Kunde inte ansluta till servern. Kontrollera att API:et körs!", err)
		return
	}

    defer resp.Body.Close()
    var questions []Question
    json.NewDecoder(resp.Body).Decode(&questions)

	//Min lista för att spara användarens val (ID -> svar)
    userAnswers := make(map[int]string)

    fmt.Println("--- QUIZ START ---")

	for {
		for _, q := range questions {
			var badInput = false; 

			for {
				clearConsole();
				
				if badInput {
					fmt.Println("\n Ogiltigt val, försök igen.")
				}

				fmt.Printf("\n%s\n", q.Text) 			// Ställer frågan
				for i, opt := range q.Options { 		// Printar ut valen i en fin lista
					fmt.Printf("%d) %s\n", i+1, opt)    // först kommer formateringen, sedan vad som ska in i den. 
				}

				// Sparar vad anvädaren skriver
				var choice int
				fmt.Print("Ditt svar (nummer): ")
				_, err :=fmt.Scanln(&choice)

				//Om svaret är en string
				if err != nil {
					badInput = true;
					clearConsole();
					var dump string
					fmt.Scanln(&dump)
					continue
				}
				
				//Om svaret är giltigt
				if choice > 0 && choice <= len(q.Options) { 
					userAnswers[q.ID] = q.Options[choice-1]
					badInput = false;
					break
				} else {
					badInput = true;
					clearConsole();
				}
			}
		}

	// Skicka svaren till API:et
    jsonData, _ := json.Marshal(userAnswers)
    
    // Skicka en POST-request
    postResp, err := http.Post("http://localhost:8080/submit", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Kunde inte skicka svar:", err)
        return
    }
    defer postResp.Body.Close()

    // 4. Läs och visa resultatet från API:et
    body, _ := io.ReadAll(postResp.Body)
    fmt.Printf("\n--- RESULTAT ---\n%s\n", string(body))
		
		//Slutet av spelet
		var replay string
		fmt.Print("Vill du spela igen? (j/n): ")
		fmt.Scanln(&replay)

		if replay != "j" {
			fmt.Println("Okej, hejdå!")
			break;
		}
	}
},
}
	// Lägg till 'start' under huvudkommandot
	rootCmd.AddCommand(startCmd)

	// Kör programmet
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func clearConsole() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Windows kör 'cls' via kommandotolken
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		// Mac/Linux kör 'clear'
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
