package post

import (
	"database/sql"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

type Post struct {
	DB *sql.DB
}

//easyjson:json
type Res struct {
	Post struct {
		Author  string    `json:"author"`
		Created time.Time `json:"created"`
		Forum   string    `json:"forum"`
		ID      int       `json:"id"`
		Message string    `json:"message"`
		Thread  int       `json:"thread"`
	} `json:"post"`
}

func (post *Post) Details(ctx *fasthttp.RequestCtx) {
	ID := ctx.UserValue("id")
	row := post.DB.QueryRow(`SELECT author, created, forum, id, message, thread FROM posts WHERE id=$1`, ID)
	var resultPost Res
	err := row.Scan(&resultPost.Post.Author, &resultPost.Post.Created, &resultPost.Post.Forum,
		&resultPost.Post.ID, &resultPost.Post.Message, &resultPost.Post.Thread)
	if err != nil {
		log.Println(err)
	}

	res, _ := easyjson.Marshal(resultPost)
	ctx.SetBody(res)
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
}

func (post *Post) UpdateMessage(ctx *fasthttp.RequestCtx) {

}
