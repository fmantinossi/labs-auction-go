
# Labs Auction Go

Este projeto é um sistema de leilão desenvolvido em Go, utilizando MongoDB como persistência e Docker para orquestração. Agora, com suporte a **fechamento automático de leilões** via goroutine!

---

## Funcionalidade nova

- Leilões criados são **fechados automaticamente** após um tempo definido (padrão: 60 segundos).
- O tempo pode ser configurado via `.env`.

---

## Requisitos

- Docker e Docker Compose instalados  
- Go (opcional, para rodar testes fora do container)

---

## Como rodar o projeto

### 1. Clone o repositório

```bash
git clone https://github.com/fmantinossi/labs-auction-go.git
cd labs-auction-go
```

### 2. Suba os containers

```bash
docker-compose up --build
```

---

## Como executar os testes

### Dentro do container (recomendado)

```bash
docker exec -it <nome_do_container_app> go test ./internal/infra/database/auction -v
```

### Localmente (se tiver Go instalado)

```bash
go test ./internal/infra/database/auction -v
```

> O teste `TestAuctionAutoClose` cria um leilão e valida se ele é fechado após 5 segundos automaticamente.

---

## Estrutura relevante

- `internal/infra/database/auction/create_auction.go`: lógica de criação + goroutine de expiração  
- `internal/infra/database/auction/create_auction_test.go`: teste automatizado de expiração  
- `.env.example`: variáveis de configuração

---

## 🛠 Variáveis de ambiente

```env
AUCTION_DURATION_SECONDS=60
AUCTION_CHECK_INTERVAL_SECONDS=10
```

---
