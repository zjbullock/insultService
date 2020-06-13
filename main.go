package main

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/graph-gophers/graphql-go"
	"github.com/juju/loggo"
	"insultService/datasource"
	"insultService/handler"
	"insultService/repository"
	"insultService/resolver"
	"insultService/router"
	"insultService/server"
	"insultService/service"
	"math/rand"
	"net/http"
	"time"
)

var (
	ctx          = context.Background()
	l            loggo.Logger
	client       = http.Client{}
	handlerFuncs *handler.Funcs
)

func init() {
	l.SetLogLevel(loggo.DEBUG)
	serv := server.NewServer(l)
	schemaString, err := serv.GetSchema("./server/graphql/", l)
	if err != nil {
		l.Criticalf("error occurred while fetching graphql schema: %v", err)
	}
	projectId := "insult"
	dataSource := datasource.NewDataSource(l, ctx, projectId)
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	repos := struct {
		fire repository.FireStore
	}{
		fire: repository.NewFireStore(dataSource, l),
	}
	services := struct {
		Insult service.Insult
		Sms    service.SMS
	}{
		Insult: service.NewInsult(repos.fire, l, *r),
		Sms:    service.NewSMS(l, client, "https://api.twilio.com"),
	}
	handlerFuncs = &handler.Funcs{
		Ctx: ctx,
		Schema: graphql.MustParseSchema(schemaString, &resolver.Resolver{
			Services: services,
			Log:      l,
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
