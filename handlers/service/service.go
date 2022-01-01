package service

import (
	"database/sql"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

type Service struct {
	DB *sql.DB
}

//easyjson:json
type Res struct {
	Forum  int `json:"forum"`
	Post   int `json:"post"`
	Thread int `json:"thread"`
	User   int `json:"user"`
}

func (service *Service) Clear(ctx *fasthttp.RequestCtx) {
	service.DB.Exec(`TRUNCATE users, forums, threads, posts, votes, forum_users;`)
}

func (service *Service) Status(ctx *fasthttp.RequestCtx) {
	var status Res
	row := service.DB.QueryRow(`SELECT * FROM
		(SELECT COUNT(*) FROM users) as user_count,
 		(SELECT COUNT(*) FROM forums) as forum_count,
		(SELECT COUNT(*) FROM threads) as thread_count,
		(SELECT COUNT(*) FROM posts) as post_count;`)
	row.Scan(&status.User, &status.Forum, &status.Thread, &status.Post)
	res, _ := easyjson.Marshal(status)
	ctx.SetBody(res)
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
}
