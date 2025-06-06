package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devfullcycle/20-CleanArch/configs"
	"github.com/devfullcycle/20-CleanArch/internal/event"
	"github.com/devfullcycle/20-CleanArch/internal/event/handler"
	"github.com/devfullcycle/20-CleanArch/internal/infra/database"
	"github.com/devfullcycle/20-CleanArch/internal/infra/graph"
	"github.com/devfullcycle/20-CleanArch/internal/infra/grpc/pb"
	"github.com/devfullcycle/20-CleanArch/internal/infra/grpc/service"
	"github.com/devfullcycle/20-CleanArch/internal/infra/web"
	"github.com/devfullcycle/20-CleanArch/internal/infra/web/webserver"
	"github.com/devfullcycle/20-CleanArch/internal/usecase"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// Depurar as variáveis de ambiente
	fmt.Println("DB_DRIVER:", configs.DBDriver)
	fmt.Println("DB_HOST:", configs.DBHost)
	fmt.Println("DB_PORT:", configs.DBPort)
	fmt.Println("DB_USER:", configs.DBUser)
	fmt.Println("DB_PASSWORD:", configs.DBPassword)
	fmt.Println("DB_NAME:", configs.DBName)

	// Garantir que o driver mysql seja usado
	if configs.DBDriver == "" {
		configs.DBDriver = "mysql"
		fmt.Println("Usando driver mysql padrão")
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	orderCreatedEvent := event.NewOrderCreated()
	orderCreatedHandler := handler.NewOrderCreatedHandler(rabbitMQChannel)
	if err := eventDispatcher.Register("OrderCreated", orderCreatedHandler); err != nil {
		fmt.Println("Error registering event handler:", err)
		return
	}
	orderRepository := database.NewOrderRepository(db)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, orderCreatedEvent, eventDispatcher)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := web.NewWebOrderHandler(eventDispatcher, orderRepository, orderCreatedEvent)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webListHandler := web.NewOrderListHandler(listOrdersUseCase)
	webserver.AddHandler("/orders", webListHandler.ListOrders)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	// Configurando o servidor gRPC
	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	// Configurando o servidor GraphQL
	graphqlServer := graphql_handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolver(*createOrderUseCase, listOrdersUseCase),
	}))

	// Configurando as rotas HTTP para o GraphQL
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", graphqlServer)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	// Usar variáveis de ambiente para conectar ao RabbitMQ
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "rabbitmq" // Valor padrão para o ambiente Docker
	}
	port := os.Getenv("RABBITMQ_PORT")
	if port == "" {
		port = "5672" // Porta padrão do RabbitMQ
	}
	user := os.Getenv("RABBITMQ_USER")
	if user == "" {
		user = "guest" // Usuário padrão do RabbitMQ
	}
	pass := os.Getenv("RABBITMQ_PASS")
	if pass == "" {
		pass = "guest" // Senha padrão do RabbitMQ
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
	fmt.Println("Tentando conectar ao RabbitMQ em:", url)

	// Adicionar lógica de retry para a conexão com o RabbitMQ
	var conn *amqp.Connection
	var err error
	for i := 0; i < 30; i++ { // Tentar por até 30 segundos
		conn, err = amqp.Dial(url)
		if err == nil {
			break // Conexão bem-sucedida
		}
		fmt.Printf("Tentativa %d de conectar ao RabbitMQ falhou: %s. Tentando novamente em 1 segundo...\n", i+1, err)
		time.Sleep(1 * time.Second) // Esperar 1 segundo antes de tentar novamente
	}

	if err != nil {
		panic(fmt.Sprintf("Não foi possível conectar ao RabbitMQ após várias tentativas: %s", err))
	}

	fmt.Println("Conectado ao RabbitMQ com sucesso!")
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
