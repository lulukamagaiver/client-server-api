package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"desafio.client.serve/cotacao"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "cotacao.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cotacoes (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			bid TEXT NOT NULL
		);
	`)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		bid, err := buscarCotacoes(r.Context())

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Println("timeout na API externa: ", err)
			} else {
				log.Println("erro na API externa: ", err)
			}
			http.Error(w, "Erro na busca da cotação!", http.StatusGatewayTimeout)
			return
		}

		err = salvarNoBanco(r.Context(), db, bid)

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Println("timeout ao salvar no Banco: ", err)
			} else {
				log.Println("erro ao salvar no Banco: ", err)
			}
			http.Error(w, "Erro ao salvar cotação!", http.StatusGatewayTimeout)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cotacao.CotacaoResponse{Bid: bid})
	})

	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func buscarCotacoes(parent context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(parent, 200*time.Millisecond)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var api cotacao.ApiResponse

	err = json.NewDecoder(resp.Body).Decode(&api)

	if err != nil {
		return "", err
	}

	return api.UsdBrl.Bid, nil
}

func salvarNoBanco(parent context.Context, db *sql.DB, bid string) error {
	ctx, cancel := context.WithTimeout(parent, 10*time.Millisecond)

	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO cotacoes (bid) VALUES (?)", bid)

	return err
}
