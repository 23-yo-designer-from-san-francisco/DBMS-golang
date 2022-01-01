package main

import (
	"DBMS/handlers/forum"
	"DBMS/handlers/post"
	"DBMS/handlers/service"
	"DBMS/handlers/thread"
	"DBMS/handlers/user"
	"database/sql"
	"github.com/fasthttp/router"
	_ "github.com/lib/pq"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	connStr := "port=54321 dbname=postgrid sslmode=disable"
	var err error
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	r := router.New()

	forum := &forum.Forum{DB: db}
	r.POST("/forum/create", forum.Create)
	r.GET("/forum/{slug}/details", forum.Details)
	r.POST("/forum/{slug}/create", forum.CreateThread)
	r.GET("/forum/{slug}/users", forum.Users)
	r.GET("/forum/{slug}/threads", forum.GetThreads)

	post := &post.Post{DB: db}
	r.GET("/post/{id}/details", post.Details)
	r.POST("/post/{id}/details", post.UpdateMessage)

	service := &service.Service{DB: db}
	r.POST("/service/clear", service.Clear)
	r.GET("/service/status", service.Status)

	thread := &thread.Thread{DB: db}
	r.POST("/thread/{slug_or_id}/create", thread.Create)
	r.GET("/thread/{slug_or_id}/details", thread.Details)
	r.POST("/thread/{slug_or_id}/details", thread.Update)
	r.GET("/thread/{slug_or_id}/posts", thread.GetPosts)
	r.POST("/thread/{slug_or_id}/vote", thread.Vote)

	user := &user.User{DB: db}
	r.POST("/user/{nickname}/create", user.Create)
	r.GET("/user/{nickname}/profile", user.Profile)
	r.POST("/user/{nickname}/profile", user.Update)

	log.Fatal(fasthttp.ListenAndServe(":5000", r.Handler))
}
