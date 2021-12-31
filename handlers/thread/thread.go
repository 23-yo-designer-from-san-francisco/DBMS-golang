package thread

import (
	"database/sql"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"strings"
)

//easyjson:json
type ResThread struct {
	ID       int    `json:"id,omitempty"`
	Parent   int64  `json:"parent,omitempty"`
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
		if err := easyjson.Unmarshal(ctx.PostBody(), threads); err != nil {
			log.Println(err)
		}
		if len(*threads) == 0 {
			ctx.SetBody([]byte("[]"))
			ctx.SetStatusCode(201)
			ctx.SetContentType("application/json")
			return
		}
		query := `INSERT INTO posts (parent, author, message, thread, forum) VALUES `
		var values []interface{}
		for i, thr := range *threads {
			value := fmt.Sprintf(
				"(NULLIF($%d, 0), $%d, $%d, $%d, $%d),",
				i*5+1, i*5+2, i*5+3, i*5+4, i*5+5,
			)
			query += value
			values = append(values, thr.Parent, thr.Author, thr.Message, id, forumSwag)
		}
		query = strings.TrimSuffix(query, ",")
		query += ` RETURNING id, parent, author, message, isedited, forum, thread, created;`
		rows, err := thread.DB.Query(query, values...)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()

		resPosts := make(ResThreads, 0)
		for rows.Next() {
			post := &ResThread{}
			var parent sql.NullInt64

			err := rows.Scan(
				&post.ID,
				&parent,
				&post.Author,
				&post.Message,
				&post.IsEdited,
				&post.Forum,
				&post.Thread,
				&post.Created)
			if err != nil {
				log.Println(err)
			}

			if parent.Valid {
				post.Parent = parent.Int64
			} else {
				post.Parent = 0
			}
			resPosts = append(resPosts, *post)
		}
		res, _ := easyjson.Marshal(resPosts)
		ctx.SetBody(res)
		ctx.SetStatusCode(201)
		ctx.SetContentType("application/json")
		return
	} else {

	}
}

func (thread *Thread) Details(ctx *fasthttp.RequestCtx) {
	SLUG := ctx.UserValue("slug_or_id").(string)
	id, err := strconv.Atoi(SLUG)
	var row *sql.Row
	thr := &ResThread{}
	if err == nil {
		row = thread.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title
										from threads where id=$1`,
			id)
		err = row.Scan(&thr.Author, &thr.Created, &thr.Forum, &thr.ID, &thr.Message, &thr.Slug, &thr.Title)
		if err != nil {
			log.Println(err)
		}
	} else {
		id = -1
		row = thread.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title
										from threads where slug=$1`,
			SLUG)
		err = row.Scan(&thr.Author, &thr.Created, &thr.Forum, &thr.ID, &thr.Message, &thr.Slug, &thr.Title)
		if err != nil {
			log.Println(err)
		}
	}
	res, _ := easyjson.Marshal(thr)
	ctx.SetBody(res)
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
}

func (thread *Thread) Update(ctx *fasthttp.RequestCtx) {

}

func (thread *Thread) GetPosts(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	limit := string(ctx.QueryArgs().Peek("limit"))
	since := string(ctx.QueryArgs().Peek("since"))
	sort := string(ctx.QueryArgs().Peek("sort"))
	desc := string(ctx.QueryArgs().Peek("desc"))

	if len(sort) == 0 {
		sort = "flat"
	}

	var query string
	var args []interface{}
	switch sort {
	case "flat":
		query = `SELECT p.id, p.thread, p.created,
				p.message, COALESCE(p.parent, 0), p.author, p.forum FROM posts p JOIN threads thr ON p.thread = thr.id WHERE `
		ID, err := strconv.Atoi(slugOrID)
		if err == nil {
			query += "thr.ID = $1 "
			args = append(args, ID)
		} else {
			query += "thr.slug = $1 "
			args = append(args, slugOrID)
		}
		argc := 2
		if len(since) != 0 {
			if desc == "true" {
				query += "AND ID < $" + strconv.Itoa(argc)
				argc++
			} else {
				query += "AND ID > $" + strconv.Itoa(argc)
				argc++
			}
			args = append(args, since)
		}
		if desc == "true" {
			query += " ORDER BY created DESC, id DESC "
		} else {
			query += " ORDER BY created, id "
		}
		if len(limit) != 0 {
			query += " LIMIT $" + strconv.Itoa(argc)
			argc++
			args = append(args, limit)
		}
	case "tree":
		//var sinceQuery string
		//var descQuery string
		//var limitSQL string
		//argc := 2
		//
		//var args []interface{}
		//args = append(args, slugOrID)

	}
	log.Println(query, args)
	rows, err := thread.DB.Query(query, args...)
	if err != nil {
		log.Println(err)
	}
	result := make(ResThreads, 0)
	for rows.Next() {
		var thr ResThread
		err := rows.Scan(&thr.ID, &thr.Thread, &thr.Created, &thr.Message, &thr.Parent, &thr.Author, &thr.Forum)
		if err != nil {
			log.Println(err)
		}
		result = append(result, thr)
	}
	res, _ := easyjson.Marshal(result)
	ctx.SetBody(res)
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
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

		res, _ := easyjson.Marshal(thr)
		ctx.SetBody(res)
		ctx.SetStatusCode(200)
		ctx.SetContentType("application/json")
		return
	} else {

	}
}
