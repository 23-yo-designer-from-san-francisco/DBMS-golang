package forum

import (
	"DBMS/handlers/user"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

type Forum struct {
	DB *sql.DB
}

//easyjson:json
type Req struct {
	Slug  string `json:"slug,omitempty"`
	Title string `json:"title,omitempty"`
	User  string `json:"user,omitempty"`
}

//easyjson:json
type ThreadReq struct {
	ID      int64     `json:"id,omitempty"`
	Author  string    `json:"author,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Forum   string    `json:"forum,omitempty"`
	Message string    `json:"message,omitempty"`
	Title   string    `json:"title,omitempty"`
	Slug    string    `json:"slug,omitempty"`
}

//easyjson:json
type Reqs []Req

func (forum *Forum) Create(ctx *fasthttp.RequestCtx) {
	request := &Req{}
	easyjson.Unmarshal(ctx.PostBody(), request)
	usr, err := forum.DB.Query(`SELECT nickname FROM users WHERE nickname=$1`, request.User)
	if usr.Next() {
		usr.Scan(&request.User)
	}
	_, err = forum.DB.Exec(`INSERT INTO forums (title, "user", slug) VALUES($1, $2, $3)`,
		request.Title,
		request.User,
		request.Slug,
	)

	if err, ok := err.(*pq.Error); ok {
		fmt.Println(err.Code)
		switch err.Code {
		case "23505":
			rows, _ := forum.DB.Query(`SELECT slug, title, "user"`+
				"FROM forums "+
				"WHERE slug=$1",
				request.Slug)
			if rows.Next() {
				rows.Scan(&request.Slug, &request.Title, &request.User)
			}
			rows.Close()
			result, _ := easyjson.Marshal(request)
			ctx.SetBody(result)
			ctx.SetStatusCode(409)
			ctx.SetContentType("application/json")
			return
		case "23503":
			errMsg := &user.ErrMsg{Message: fmt.Sprintf("Can't find user with nickname %s", request.User)}
			response, _ := easyjson.Marshal(errMsg)
			ctx.SetBody(response)
			ctx.SetStatusCode(404)
			ctx.SetContentType("application/json")
			return
		}
	}
	res, _ := easyjson.Marshal(request)
	ctx.SetBody(res)
	ctx.SetStatusCode(201)
	ctx.SetContentType("application/json")
}

func (forum *Forum) Details(ctx *fasthttp.RequestCtx) {
	request := &Req{}
	request.Slug = ctx.UserValue("slug").(string)
	rows, _ := forum.DB.Query(`SELECT slug, title, "user" `+
		"FROM forums "+
		"WHERE slug=$1",
		request.Slug)
	if rows.Next() {
		rows.Scan(&request.Slug, &request.Title, &request.User)
		resp, _ := easyjson.Marshal(request)
		ctx.Response.SetBody(resp)
		ctx.SetContentType("application/json")
		ctx.Response.SetStatusCode(200)
		return
	} else {
		errMsg := &user.ErrMsg{Message: fmt.Sprintf("Can't find forum with slug:  %s", request.Slug)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}
}

func (forum *Forum) CreateThread(ctx *fasthttp.RequestCtx) {
	SLUG := ctx.UserValue("slug").(string)
	thr := &ThreadReq{}
	easyjson.Unmarshal(ctx.PostBody(), thr)
	row := forum.DB.QueryRow(`INSERT INTO threads (title,author, forum, message, created, slug) 
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id`,
		thr.Title,
		thr.Author,
		thr.Forum,
		thr.Message,
		thr.Created,
		SLUG,
	)
	err := row.Scan(&thr.ID)
	if err != nil {
		log.Println(err)
	}
	log.Println(thr.ID)

	res, _ := easyjson.Marshal(thr)
	ctx.SetBody(res)
	ctx.SetStatusCode(201)
	ctx.SetContentType("application/json")
}

func (forum *Forum) Users(ctx *fasthttp.RequestCtx) {

}

func (forum *Forum) Threads(ctx *fasthttp.RequestCtx) {

}
