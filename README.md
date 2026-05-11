Mitt Go Quiz Project
Ett interaktivt quiz-spel byggt i Go. Projektet är uppdelat i ett REST API (backend) och ett CLI-verktyg (frontend) som kommunicerar via HTTP och JSON.

Kom igång
För att köra projektet behöver du ha Go installerat på din dator.

1. Starta API:et (Backend)
Öppna en terminal, navigera till /api-mappen och kör:

```
go run main.go
```
API:et startar nu på http://localhost:8080 och väntar på anrop.

2. Starta CLI:et (Frontend)
Öppna en ny terminal, navigera till /cli-mappen och kör:

```
go run main.go start
```
Teknisk översikt
Arkitektur
Projektet följer en klient-server-modell:

Backend: Ett REST API byggt med Gos inbyggda net/http-paket.

Frontend: Ett CLI byggt med Cobra-biblioteket.

API Endpoints
GET /questions - Hämtar alla quiz-frågor.

POST /submit - Skickar in användarens svar och returnerar poängstatistik.
