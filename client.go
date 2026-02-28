package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"desafio.client.serve/cotacao"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	bid, err := pegarBidDoServidor(ctx)

	if err != nil {
		log.Println("Erro no cliente: ", err)
		return
	}

	texto := fmt.Sprintf("Dolar: %s", bid)

	err = os.WriteFile("cotacao.txt", []byte(texto), 0644)

	if err != nil {
		log.Println("Erro ao escrever arquivo: ", err)
		return
	}

	log.Println("Arquivo cotacao criado com sucesso!")
}

func pegarBidDoServidor(ctx context.Context) (any, any) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("servidor retornou %s: %s", resp.Status, string(b))
	}

	var cotacao cotacao.CotacaoResponse

	err = json.NewDecoder(resp.Body).Decode(&cotacao)

	if err != nil {
		return "", err
	}

	return cotacao.Bid, nil
}
