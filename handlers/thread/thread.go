package thread

import (
	"database/sql"
	"fmt"
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
	Slug     string `json:"slug,omitempty"`
	Title    string `json:"title,omitempty"`
	Votes    int    `json:"votes,omitempty"`
}

type Thread struct {
	DB *sql.DB
}

//easyjson:json
type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
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
	threadSlug := ctx.UserValue("slug_or_id").(string)
	var row *sql.Row
	id, err := strconv.Atoi(threadSlug)
	thr := &ResThread{}
	if err == nil {
		row = thread.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title from threads where id=$1`, id)
		err = row.Scan(&thr.Author, &thr.Created, &thr.Forum, &thr.ID, &thr.Message, &thr.Slug, &thr.Title)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		id = -1
		row = thread.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title from threads where slug=$1`, threadSlug)
		err = row.Scan(&thr.Author, &thr.Created, &thr.Forum, &thr.ID, &thr.Message, &thr.Slug, &thr.Title)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if len(threadSlug) != 0 {
		vote := &Vote{}
		easyjson.Unmarshal(ctx.PostBody(), vote)
		//tx, _ := thread.DB.Begin()
		log.Println("Thread ID")
		log.Println(thr.ID)
		row := thread.DB.QueryRow(`INSERT INTO votes as vote 
                (nickname, thread, voice)
                VALUES ($1, $2, $3) 
                ON CONFLICT ON CONSTRAINT votes_user_thread_unique DO
                UPDATE SET voice = $3 WHERE vote.voice <> $3`, vote.Nickname, thr.ID, vote.Voice)
		log.Println(row.Scan().Error())
		row = thread.DB.QueryRow(`SELECT votes FROM threads WHERE id=$1`, thr.ID)
		err := row.Scan(&thr.Votes)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Votest")
		fmt.Println(thr.Votes)

		res, _ := easyjson.Marshal(thr)
		ctx.SetBody(res)
		ctx.SetStatusCode(200)
		ctx.SetContentType("application/json")
		return
	} else {

	}
}
