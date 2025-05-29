# Clean Architecture - Order System

Este é um sistema de pedidos implementado com Clean Architecture em Go.

## Serviços

- REST API: http://localhost:8080
- GraphQL API: http://localhost:8081
- gRPC: http://localhost:50051
- MySQL: http://localhost:3306
- RabbitMQ: http://localhost:15672 (interface web)

## Requisitos

- Docker e Docker Compose
- Go 1.21+

## Inicialização

1. Clone o repositório
2. Execute:
```bash
docker compose up -d
```

Isso iniciará:
- MySQL (porta 3306)
- RabbitMQ (porta 5672 e interface web na porta 15672)
- Aplicação (portas 8080, 8081 e 50051)

## Endpoints

### REST API
- POST /orders - Criar pedido
- GET /orders - Listar pedidos

### GraphQL
- POST /graphql - Consultar pedidos

### gRPC
- CreateOrder - Criar pedido
- ListOrders - Listar pedidos

## Testando a API

Use o arquivo `api.http` com o HTTP Client do VS Code para testar os endpoints.

## Estrutura do Projeto

```
internal/
├── entity/        # Entidades do domínio
├── usecase/       # Casos de uso
├── infra/         # Infraestrutura
│   ├── database/  # Repositórios
│   ├── grpc/      # Serviço gRPC
│   ├── web/       # Handlers HTTP
│   └── graph/     # Schema GraphQL
└── events/        # Eventos e dispatcher
```

## Tecnologias

- Go 1.21+
- MySQL 5.7
- RabbitMQ
- gRPC
- GraphQL
- Clean Architecture
