package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

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

	graphqlServer := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolver(*createOrderUseCase, listOrdersUseCase),
	}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", graphqlServer)
	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	go http.ListenAndServe(":"+configs.GraphQLServerPort, nil)

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

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
