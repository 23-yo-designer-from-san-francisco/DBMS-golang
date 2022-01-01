package main

import (
	"DBMS/handlers/forum"
	"DBMS/handlers/post"
	"DBMS/handlers/service"
	"DBMS/handlers/thread"
	"DBMS/handlers/user"
	"database/sql"
	"fmt"
	"github.com/fasthttp/router"
	_ "github.com/lib/pq"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

func main() {
	connStr := fmt.Sprintf("port=%s dbname=%s username=%s password=%s sslmode=disable",
		os.Getenv("DBPORT"),
		os.Getenv("DBNAME"),
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
	)
	var err error
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	r := router.New()

	forum := &forum.Forum{DB: db}
	r.POST("/api/forum/create", forum.Create)
	r.GET("/api/forum/{slug}/details", forum.Details)
	r.POST("/api/forum/{slug}/create", forum.CreateThread)
	r.GET("/api/forum/{slug}/users", forum.Users)
	r.GET("/api/forum/{slug}/threads", forum.GetThreads)

	post := &post.Post{DB: db}
	r.GET("/api/post/{id}/details", post.Details)
	r.POST("/api/post/{id}/details", post.UpdateMessage)

	service := &service.Service{DB: db}
	r.POST("/api/service/clear", service.Clear)
	r.GET("/api/service/status", service.Status)

	thread := &thread.Thread{DB: db}
	r.POST("/api/thread/{slug_or_id}/create", thread.Create)
	r.GET("/api/thread/{slug_or_id}/details", thread.Details)
	r.POST("/api/thread/{slug_or_id}/details", thread.Update)
	r.GET("/api/thread/{slug_or_id}/posts", thread.GetPosts)
	r.POST("/api/thread/{slug_or_id}/vote", thread.Vote)

	user := &user.User{DB: db}
	r.POST("/api/user/{nickname}/create", user.Create)
	r.GET("/api/user/{nickname}/profile", user.Profile)
	r.POST("/api/user/{nickname}/profile", user.Update)

	log.Fatal(fasthttp.ListenAndServe(":5000", r.Handler))
}
