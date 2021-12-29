package thread

import (
	"database/sql"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
)

//easyjson:json
type ResThread struct {
	ID       int    `json:"id,omitempty"`
	Parent   int    `json:"parent,omitempty"`
	Author   string `json:"author,omitempty"`
	Message  string `json:"message,omitempty"`
	IsEdited bool   `json:"isEdited,omitempty"`
	Forum    string `json:"forum,omitempty"`
	Thread   int    `json:"thread,omitempty"`
	Created  string `json:"created,omitempty"`
}

type Thread struct {
	DB *sql.DB
}

//easyjson:json
type ResThreads []ResThread

func (thread *Thread) Create(ctx *fasthttp.RequestCtx) {
	SLUG := ctx.UserValue("slug_or_id").(string)
	id, err := strconv.Atoi(SLUG)
	var row *sql.Row
	var forumTitle string
	var forumSwag string
	if err == nil {
		row = thread.DB.QueryRow(`SELECT title, forum from threads where id=$1`, id)
		err = row.Scan(&forumTitle, &forumSwag)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		id = -1
		row = thread.DB.QueryRow(`SELECT title, forum, id from threads where slug=$1`, SLUG)
		err = row.Scan(&forumTitle, &forumSwag, &id)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if len(forumTitle) != 0 {
		threads := &ResThreads{}
		result := make(ResThreads, 0)
		if err := easyjson.Unmarshal(ctx.PostBody(), threads); err != nil {
			log.Println(err)
		}
		for _, thr := range *threads {
			row := thread.DB.QueryRow("INSERT INTO posts (author, message, forum, thread) "+
				"VALUES ($1, $2, $3, $4) RETURNING id, created",
				thr.Author,
				thr.Message,
				forumTitle,
				id,
			)
			if err != nil {
				log.Fatalln(err)
			}
			err := row.Scan(&thr.ID, &thr.Created)
			thr.Forum = forumSwag
			if err != nil {
				log.Fatalln(err)
			}
			thr.Thread = id
			result = append(result, thr)
		}
		res, _ := easyjson.Marshal(result)
		ctx.SetBody(res)
		ctx.SetStatusCode(201)
		ctx.SetContentType("application/json")
		return
	} else {

	}
}

func (thread *Thread) Details(ctx *fasthttp.RequestCtx) {

}

func (thread *Thread) Update(ctx *fasthttp.RequestCtx) {

}

func (thread *Thread) Messages(ctx *fasthttp.RequestCtx) {

}

func (thread *Thread) Vote(ctx *fasthttp.RequestCtx) {

}
