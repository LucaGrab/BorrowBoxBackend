# Verwende das offizielle Golang-Base-Image
FROM golang:latest

# Setze das Arbeitsverzeichnis innerhalb des Containers
WORKDIR /app

# Kopiere die lokale Go-Anwendung in das Arbeitsverzeichnis im Container
COPY . .

# Installiere erforderliche Abhängigkeiten (wenn benötigt)
# Hier könntest du beispielsweise zusätzliche Abhängigkeiten installieren, die dein Projekt benötigt.

# Baue die Go-Anwendung
RUN go build -o myapp main.go

# Exponiere den Port, auf dem die Anwendung lauscht
EXPOSE 8088

# Starte die Anwendung beim Ausführen des Containers
CMD ["./myapp"]
