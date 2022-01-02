package main

import (
	"DBMS/handlers/forum"
	"DBMS/handlers/post"
	"DBMS/handlers/service"
	"DBMS/handlers/thread"
	"DBMS/handlers/user"
	"context"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	connStr := fmt.Sprintf("dbname=pq sslmode=disable pool_max_conns=30")
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Println(err)
	}
	dbPool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer dbPool.Close()

	r := router.New()

	forum := &forum.Forum{DB: dbPool}
	r.POST("/api/forum/create", forum.Create)
	r.GET("/api/forum/{slug}/details", forum.Details)
	r.POST("/api/forum/{slug}/create", forum.CreateThread)
	r.GET("/api/forum/{slug}/users", forum.Users)
	r.GET("/api/forum/{slug}/threads", forum.GetThreads)

	post := &post.Post{DB: dbPool}
	r.GET("/api/post/{id}/details", post.Details)
	r.POST("/api/post/{id}/details", post.UpdateMessage)

	service := &service.Service{DB: dbPool}
	r.POST("/api/service/clear", service.Clear)
	r.GET("/api/service/status", service.Status)

	thread := &thread.Thread{DB: dbPool}
	r.POST("/api/thread/{slug_or_id}/create", thread.Create)
	r.GET("/api/thread/{slug_or_id}/details", thread.Details)
	r.POST("/api/thread/{slug_or_id}/details", thread.Update)
	r.GET("/api/thread/{slug_or_id}/posts", thread.GetPosts)
	r.POST("/api/thread/{slug_or_id}/vote", thread.Vote)

	user := &user.User{DB: dbPool}
	r.POST("/api/user/{nickname}/create", user.Create)
	r.GET("/api/user/{nickname}/profile", user.Profile)
	r.POST("/api/user/{nickname}/profile", user.Update)

	log.Fatal(fasthttp.ListenAndServe(":5000", r.Handler))
}
