package forum

import (
	"database/sql"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"log"
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
type Reqs []Req

func (forum *Forum) Create(ctx *fasthttp.RequestCtx) {
	request := &Req{}
	easyjson.Unmarshal(ctx.PostBody(), request)
	_, err := forum.DB.Exec("INSERT INTO forums (title, user, slug) "+
		"VALUES($1, $2, $3)",
		request.Title,
		request.User,
		request.Slug,
	)
	if err != nil {
		log.Println(err)
	}

	ctx.SetBody(ctx.PostBody())
	ctx.SetStatusCode(201)
	ctx.SetContentType("application/json")
}

func (forum *Forum) Details(ctx *fasthttp.RequestCtx) {

}

func (forum *Forum) CreateThread(ctx *fasthttp.RequestCtx) {

}

func (forum *Forum) Users(ctx *fasthttp.RequestCtx) {

}

func (forum *Forum) Threads(ctx *fasthttp.RequestCtx) {

}
