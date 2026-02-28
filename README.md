# Desafio Go: Client/Server Cotação USD-BRL

Este projeto contém dois programas:
- `server.go`: servidor HTTP na porta 8080 com endpoint `/cotacao`
- `client.go`: cliente que consome o servidor e salva a cotação em `cotacao.txt`

## Requisitos atendidos

### Server
- Endpoint: `GET /cotacao`
- Consome API externa: https://economia.awesomeapi.com.br/json/last/USD-BRL
- Timeout da API externa: **200ms** (context)
- Persiste em SQLite com timeout de **10ms** (context)
- Retorna JSON com o campo `bid`
- Loga erro quando timeout é excedido

### Client
- Requisição para `http://localhost:8080/cotacao`
- Timeout do client: **300ms** (context)
- Salva em `cotacao.txt` no formato: `Dólar: {valor}`
- Loga erro quando timeout é excedido

## Como rodar

### 1) Subir o servidor
```bash
go mod tidy
go run server.go
```
- saida experada: Servidor rodando em http://localhost:8080

### 2) Rodar o client
- Abra outro terminal (separado do Terminal 1) e rode:

```bash
go run client.go
```
- Se tudo ocorrer bem dentro dos limites de tempo:
```code
Arquivo cotacao.txt criado com sucesso!
````
- Será criado/atualizado o arquivo cotacao.txt com o formato exigido:
```code
Dólar: 5.1234
````

De acordo com a documentação do desafio: O timeout máximo do client é de 300ms

- Se exceder esse tempo, o client deve logar o erro no console
```code
timeout no client: context deadline exceeded
````

