package main

import (
	"learn_gqlgen/graph"
	"learn_gqlgen/graph/generated"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func GraphqlHandler() gin.HandlerFunc {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver()}))
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})

	log.Println("called once")

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "content-type"},
		ExposeHeaders:    []string{"Content-Length", "content-type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	gqlHandler := GraphqlHandler()
	router.POST("/query", gqlHandler)
	router.GET("/subscriptions", gqlHandler)

	router.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL", "/query").ServeHTTP(c.Writer, c.Request)
	})

	router.GET("/sub", func(c *gin.Context) {
		playground.Handler("TestSub", "/subscriptions").ServeHTTP(c.Writer, c.Request)
	})

	router.Run(":8080")
}
