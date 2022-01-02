package forum

import (
	"DBMS/handlers/user"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"time"
)

type Forum struct {
	DB *pgxpool.Pool
}

//easyjson:json
type Req struct {
	Slug    string `json:"slug,omitempty"`
	Title   string `json:"title,omitempty"`
	User    string `json:"user,omitempty"`
	Posts   int64  `json:"posts,omitempty"`
	Threads int64  `json:"threads,omitempty"`
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
	usr, err := forum.DB.Query(context.Background(), `SELECT nickname FROM users WHERE nickname=$1`, request.User)
	defer usr.Close()
	if usr.Next() {
		usr.Scan(&request.User)
	}
	_, err = forum.DB.Exec(context.Background(), `INSERT INTO forums (title, "user", slug) VALUES($1, $2, $3)`,
		request.Title,
		request.User,
		request.Slug,
	)

	if err, ok := err.(*pq.Error); ok {
		switch err.Code {
		case "23505":
			rows, _ := forum.DB.Query(context.Background(), `SELECT slug, title, "user"`+
				"FROM forums "+
				"WHERE slug=$1",
				request.Slug)
			defer rows.Close()
			if rows.Next() {
				rows.Scan(&request.Slug, &request.Title, &request.User)
			}
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
	rows, _ := forum.DB.Query(context.Background(), `SELECT slug, title, "user", posts, threads `+
		"FROM forums "+
		"WHERE slug=$1",
		request.Slug)
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&request.Slug, &request.Title, &request.User, &request.Posts, &request.Threads)
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
	SLUG := ctx.UserValue("slug").(string) // SWAG форума, который точно должен быть
	thr := &ThreadReq{}
	easyjson.Unmarshal(ctx.PostBody(), thr)

	var SWAG string
	err := forum.DB.QueryRow(context.Background(), `SELECT slug FROM forums WHERE slug=$1`, SLUG).Scan(&SWAG)
	if err == sql.ErrNoRows {
		errMsg := &user.ErrMsg{Message: fmt.Sprintf("Can't find thread forum by slug: %s", SWAG)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

	row := forum.DB.QueryRow(context.Background(), `INSERT INTO threads (title,author, forum, message, created, slug) 
		VALUES($1, $2, $3, $4, $5, CASE WHEN $6 <> '' THEN $6 END) RETURNING id`,
		thr.Title,
		thr.Author,
		SWAG,
		thr.Message,
		thr.Created,
		thr.Slug,
	)
	err = row.Scan(&thr.ID)
	if err, ok := err.(*pq.Error); ok {
		switch err.Code {
		case "23505":
			thread := &ThreadReq{}
			swag := sql.NullString{}
			forum.DB.QueryRow(context.Background(), "SELECT author, created, forum, id, message, slug, title "+
				"FROM threads "+
				"WHERE slug=$1", thr.Slug).Scan(&thread.Author,
				&thread.Created, &thread.Forum, &thread.ID, &thread.Message, &swag, &thread.Title)
			if swag.Valid {
				thread.Slug = swag.String
			}
			res, _ := easyjson.Marshal(thread)
			ctx.SetBody(res)
			ctx.SetStatusCode(409)
			ctx.SetContentType("application/json")
			return
		case "23503":
			errMsg := &user.ErrMsg{Message: fmt.Sprintf("Can't find thread author by nickname: %s", thr.Author)}
			response, _ := easyjson.Marshal(errMsg)
			ctx.SetBody(response)
			ctx.SetStatusCode(404)
			ctx.SetContentType("application/json")
			return
		}
	}
	thr.Forum = SWAG
	userQuery := forum.DB.QueryRow(context.Background(), `INSERT INTO forum_users (user_nickname, forum_swag) VALUES ($1, $2)`, thr.Author, thr.Forum)
	userQuery.Scan()
	res, _ := easyjson.Marshal(thr)
	ctx.SetBody(res)
	ctx.SetStatusCode(201)
	ctx.SetContentType("application/json")
}

func (forum *Forum) Users(ctx *fasthttp.RequestCtx) {
	SLUG := ctx.UserValue("slug").(string)
	desc := string(ctx.QueryArgs().Peek("desc"))
	limit := string(ctx.QueryArgs().Peek("limit"))
	since := string(ctx.QueryArgs().Peek("since"))
	users := make(user.Reqs, 0)

	query := forum.DB.QueryRow(context.Background(), `SELECT id from forums WHERE slug=$1`, SLUG)
	var forumID int
	query.Scan(&forumID)
	var queryOpts string
	if len(since) != 0 {
		queryOpts += " AND nickname "
		if desc == "true" {
			queryOpts += "<"
		} else {
			queryOpts += ">"
		}
		queryOpts += "'" + since + "'"
	}
	queryOpts += " ORDER BY nickname "
	if desc == "true" {
		queryOpts += " DESC "
	}
	if len(limit) != 0 {
		queryOpts += " LIMIT " + limit //TODO=Сделать через $
	}

	rows, _ := forum.DB.Query(context.Background(), `SELECT about, email, fullname, nickname FROM users u
										JOIN forum_users fu ON fu.user_nickname = u.nickname WHERE forum_swag = $1`+queryOpts, SLUG)
	defer rows.Close()
	found := false
	for rows.Next() {
		found = true
		var user user.Req
		rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		users = append(users, user)
	}

	if !found {
		forum := forum.DB.QueryRow(context.Background(), `SELECT id FROM forums where slug=$1`, SLUG)
		var forumID int
		forum.Scan(&forumID)
		if forumID == 0 {
			errMsg := &user.ErrMsg{}
			errMsg.Message = "Can't find forum by slug: " + SLUG
			res, _ := easyjson.Marshal(errMsg)
			ctx.SetBody(res)
			ctx.SetStatusCode(404)
			ctx.SetContentType("application/json")
			return
		}
	}

	response, _ := easyjson.Marshal(users)
	ctx.SetBody(response)
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
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
	var rows pgx.Rows

	if len(limit) != 0 {
		if len(since) == 0 {
			if desc == "true" {
				rows, _ = forum.DB.Query(context.Background(), "SELECT id, title, author, forum, message, votes, slug, created "+
					"FROM threads WHERE forum=$1 ORDER BY created DESC "+limitQueryArg, SLUG)
			} else {
				rows, _ = forum.DB.Query(context.Background(), "SELECT id, title, author, forum, message, votes, slug, created "+
					"FROM threads WHERE forum=$1 ORDER BY created ASC "+limitQueryArg, SLUG)
			}
		} else {
			if desc == "true" {
				rows, _ = forum.DB.Query(context.Background(), "SELECT id, title, author, forum, message, votes, slug, created "+
					"FROM threads WHERE forum=$1 AND created <= $2 ORDER BY created DESC "+limitQueryArg, SLUG, since)
			} else {
				rows, _ = forum.DB.Query(context.Background(), "SELECT id, title, author, forum, message, votes, slug, created "+
					"FROM threads WHERE forum=$1 AND created >= $2 ORDER BY created ASC "+limitQueryArg, SLUG, since)
			}
		}
		threads := make(ThreadsReq, 0)
		defer rows.Close()
		thr := &ThreadReq{}
		for rows.Next() {
			thr := &ThreadReq{}
			swagga := sql.NullString{}
			rows.Scan(&thr.ID, &thr.Title, &thr.Author, &thr.Forum, &thr.Message, &thr.Votes, &swagga, &thr.Created)
			if swagga.Valid {
				thr.Slug = swagga.String
			}
			threads = append(threads, *thr)
		}
		if thr.ID == 0 {
			forum, _ := forum.DB.Query(context.Background(), "SELECT id FROM forums WHERE slug=$1", SLUG)
			defer forum.Close()
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
	swag := sql.NullString{}
	usr, _ := forum.DB.Query(context.Background(), "SELECT id, title, author, forum, message, votes, slug, created "+
		"FROM threads WHERE forum=$1", SLUG)
	defer usr.Close()
	usr.Next()
	usr.Scan(&thr.ID, &thr.Title, &thr.Author, &thr.Forum, &thr.Message, &thr.Votes, &swag, &thr.Created)
	if swag.Valid {
		thr.Slug = swag.String
	}

	if thr.ID == 0 {
		errMsg := &user.ErrMsg{Message: fmt.Sprintf("Can't find forum by slug: %s", SLUG)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
	}
	errMsg := &user.ErrMsg{Message: fmt.Sprintf("Can't find forum by slug: %s", SLUG)}
	response, _ := easyjson.Marshal(errMsg)
	ctx.SetBody(response)
	ctx.SetStatusCode(404)
	ctx.SetContentType("application/json")
}
