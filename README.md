# payment-processor

API HTTP em Go usando arquitetura hexagonal.

## Estrutura

```
payment-processor/
├── cmd/payment_processor/
│   └── main.go               ← entry point: wiring + HTTP server
├── internal/boleto/
│   ├── adapters/
│   │   ├── inbound/
│   │   │   ├── dto/           ← request/response structs
│   │   │   └── handler/       ← HTTP handlers (net/http)
│   │   └── outbound/
│   │       ├── repository/    ← MySQL
│   │       └── service/       ← serviços externos
│   ├── domain/                ← entidades e regras de negócio
│   ├── ports/                 ← interfaces (contratos)
│   └── usecase/               ← casos de uso
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## Rodar localmente

```bash
# Subir API + MySQL
docker compose up --build

# Apenas a API (MySQL já no ar)
go run ./cmd/payment_processor
```

## Endpoints

| Método | Path | Descrição |
|--------|------|-----------|
| GET | /health | Health check |
| POST | /api/v1/boletos | Processar boleto |

## Variáveis de ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| PORT | 8080 | Porta do servidor |
| DB_HOST | localhost | Host MySQL |
| DB_PORT | 3306 | Porta MySQL |
| DB_NAME | payment_processor | Database |
| DB_USER | app | Usuário |
| DB_PASS | secret | Senha |

## Deploy na VPS

```bash
# Build da imagem
docker build -t payment-processor .

# Rodar
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=seu-host \
  -e DB_PASS=sua-senha \
  --name payment-processor \
  payment-processor
```
