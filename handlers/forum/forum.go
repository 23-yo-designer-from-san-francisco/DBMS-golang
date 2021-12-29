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
type Reqs []Req

//easyjson:json
type ThreadReq struct {
	ID      int64     `json:"id,omitempty"`
	Author  string    `json:"author,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Forum   string    `json:"forum,omitempty"`
	Message string    `json:"message,omitempty"`
	Title   string    `json:"title,omitempty"`
	Slug    string    `json:"slug,omitempty"`
	Votes   int       `json:"votes,omitempty"`
}

//easyjson:json
type ThreadsReq []ThreadReq

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
	if len(thr.Slug) != 0 {
		SLUG = thr.Slug
	}
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
	if err, ok := err.(*pq.Error); ok {
		log.Println(err.Code)
		log.Println(err.Message)
		switch err.Code {
		case "23505":
			thread := &ThreadReq{}
			forum.DB.QueryRow("SELECT author, created, forum, id, message, slug, title "+
				"FROM threads "+
				"WHERE slug=$1", SLUG).Scan(&thread.Author,
				&thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title)
			res, _ := easyjson.Marshal(thread)
			ctx.SetBody(res)
			ctx.SetStatusCode(409)
			ctx.SetContentType("application/json")
			return
		}
	}

	res, _ := easyjson.Marshal(thr)
	ctx.SetBody(res)
	ctx.SetStatusCode(201)
	ctx.SetContentType("application/json")
}

func (forum *Forum) Users(ctx *fasthttp.RequestCtx) {

}

func (forum *Forum) GetThreads(ctx *fasthttp.RequestCtx) {
	SLUG := ctx.UserValue("slug").(string)
	desc := string(ctx.QueryArgs().Peek("desc"))
	limit := ctx.QueryArgs().Peek("limit")
	since := ctx.QueryArgs().Peek("since")
	var limitQueryArg string
	if len(limit) != 0 {
		limitQueryArg = " LIMIT " + string(limit)
	} else {
		limitQueryArg = ""
	}
	var rows *sql.Rows
	var err error
	if len(limit) != 0 {
		if len(since) == 0 {
			if desc == "true" {
				rows, err = forum.DB.Query("SELECT id, title, author, forum, message, votes, slug, created "+
					"FROM threads WHERE forum=$1 ORDER BY created DESC "+limitQueryArg, SLUG)
			} else {
				rows, err = forum.DB.Query("SELECT id, title, author, forum, message, votes, slug, created "+
					"FROM threads WHERE forum=$1 ORDER BY created ASC "+limitQueryArg, SLUG)
			}
		} else {
			if desc == "true" {
				rows, err = forum.DB.Query("SELECT id, title, author, forum, message, votes, slug, created "+
					"FROM threads WHERE forum=$1 AND created <= $2 ORDER BY created DESC "+limitQueryArg, SLUG, since)
			} else {
				rows, err = forum.DB.Query("SELECT id, title, author, forum, message, votes, slug, created "+
					"FROM threads WHERE forum=$1 AND created >= $2 ORDER BY created ASC "+limitQueryArg, SLUG, since)
			}
		}
		threads := make(ThreadsReq, 0)
		if err != nil {
			log.Fatalln(err)
		}
		defer rows.Close()
		found := false
		for rows.Next() {
			found = true
			thr := &ThreadReq{}
			rows.Scan(&thr.ID, &thr.Title, &thr.Author, &thr.Forum, &thr.Message, &thr.Votes, &thr.Slug, &thr.Created)
			threads = append(threads, *thr)
		}
		if !found {
			forum, _ := forum.DB.Query("SELECT id FROM forums WHERE slug=$1", SLUG)
			if !forum.Next() {
				errMsg := &user.ErrMsg{Message: fmt.Sprintf("Can't find forum by slug: %s", SLUG)}
				response, _ := easyjson.Marshal(errMsg)
				ctx.SetBody(response)
				ctx.SetStatusCode(404)
				ctx.SetContentType("application/json")
				return
			}
		}
		resp, err := easyjson.Marshal(threads)
		if err != nil {
			fmt.Println(err)
		}
		ctx.Response.SetBody(resp)
		ctx.SetContentType("application/json")
		ctx.Response.SetStatusCode(200)
		return
	}

	thr := &ThreadReq{}
	usr, _ := forum.DB.Query("SELECT id, title, author, forum, message, votes, slug, created "+
		"FROM threads WHERE forum=$1", SLUG)
	if usr.Next() {
		usr.Scan(&thr.ID, &thr.Title, &thr.Author, &thr.Forum, &thr.Message, &thr.Votes, &thr.Slug, &thr.Created)
	} else {
		errMsg := &user.ErrMsg{Message: fmt.Sprintf("Can't find forum by slug: %s", SLUG)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
	}
}
