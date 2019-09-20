package main

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/graph-gophers/graphql-go"
	"github.com/juju/loggo"
	"net/http"
	"randomInsultService/handler"
	"randomInsultService/repository"
	"randomInsultService/resolver"
	"randomInsultService/router"
	"randomInsultService/server"
	"randomInsultService/service"
)

var (
	ctx          = context.Background()
	l            loggo.Logger
	client       = http.Client{}
	handlerFuncs *handler.Funcs
)

func init() {
	l.SetLogLevel(loggo.DEBUG)
	server := server.NewServer(l)
	schemaString, err := server.GetSchema("./server/graphql/", l)
	if err != nil {
		l.Criticalf("error occurred while fetching graphql schema: %v", err)
	}

	repos := struct {
		fire repository.FireBase
	}{
		fire: repository.NewFireBase(),
	}
	services := struct {
		Insult service.Insult
	}{
		Insult: service.NewInsult(repos.fire, l),
	}
	handlerFuncs = &handler.Funcs{
		Ctx: ctx,
		Schema: graphql.MustParseSchema(schemaString, &resolver.Resolver{
			Services: services,
		}),
	}
}

func main() {
	r := router.NewRouter(handlerFuncs)
	allowedHeaders := handlers.AllowedHeaders([]string{"content-type"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{http.MethodPost})
	l.Criticalf(http.ListenAndServe(":8080", handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)).Error())
}
