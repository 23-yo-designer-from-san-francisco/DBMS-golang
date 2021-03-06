package thread

import (
	"DBMS/handlers/user"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
	"time"
)

//easyjson:json
type ResThread struct {
	ID       int       `json:"id,omitempty"`
	Parent   int64     `json:"parent,omitempty"`
	Author   string    `json:"author,omitempty"`
	Message  string    `json:"message,omitempty"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	Thread   int       `json:"thread,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Slug     string    `json:"slug,omitempty"`
	Title    string    `json:"title,omitempty"`
	Votes    int       `json:"votes,omitempty"`
}

type Thread struct {
	DB *pgxpool.Pool
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
	var row pgx.Row
	var forumTitle string
	var forumSlug string
	if err == nil {
		row = thread.DB.QueryRow(context.Background(), `SELECT title, forum from threads where id=$1`, id)
		err = row.Scan(&forumTitle, &forumSlug)
	} else {
		id = -1
		row = thread.DB.QueryRow(context.Background(), `SELECT title, forum, id from threads where slug=$1`, SLUG)
		err = row.Scan(&forumTitle, &forumSlug, &id)
	}

	if len(forumTitle) != 0 {
		threads := &ResThreads{}
		easyjson.Unmarshal(ctx.PostBody(), threads)
		if len(*threads) == 0 {
			ctx.SetBody([]byte("[]"))
			ctx.SetStatusCode(201)
			ctx.SetContentType("application/json")
			return
		}
		query := `INSERT INTO posts (parent, author, message, thread, forum) VALUES`
		var values []interface{}
		usersQuery := `INSERT INTO forum_users (user_nickname, forum_slug) VALUES`
		var forumUsers []interface{}
		for i, thr := range *threads {
			value := fmt.Sprintf(
				"(NULLIF($%d, 0), $%d, $%d, $%d, $%d),",
				i*5+1, i*5+2, i*5+3, i*5+4, i*5+5,
			)
			usersQuery += fmt.Sprintf("($%d, $%d),", i*2+1, i*2+2)
			forumUsers = append(forumUsers, thr.Author, forumSlug)
			query += value
			values = append(values, thr.Parent, thr.Author, thr.Message, id, forumSlug)
		}
		query = strings.TrimSuffix(query, ",")
		usersQuery = strings.TrimSuffix(usersQuery, ",")
		query += ` RETURNING id, parent, author, message, isedited, forum, thread, created;`
		tx, err := thread.DB.Begin(context.TODO())
		rows, err := thread.DB.Query(context.Background(), query, values...)
		tx.Commit(ctx)
		if rows != nil {
			defer rows.Close()
		}
		usersQuery += " ON CONFLICT DO NOTHING"
		tx, err = thread.DB.Begin(context.TODO())
		userRows, err := thread.DB.Query(context.Background(), usersQuery, forumUsers...)
		if userRows != nil {
			defer userRows.Close()
		}
		tx.Commit(ctx)

		if err != nil {
			result := user.ErrMsg{Message: "Can't find post author by nickname: "}
			res, _ := easyjson.Marshal(result)
			ctx.SetBody(res)
			ctx.SetStatusCode(404)
			ctx.SetContentType("application/json")
			return
		}

		resPosts := make(ResThreads, 0)
		for rows.Next() {
			post := &ResThread{}
			var parent sql.NullInt64

			rows.Scan(
				&post.ID,
				&parent,
				&post.Author,
				&post.Message,
				&post.IsEdited,
				&post.Forum,
				&post.Thread,
				&post.Created)

			if parent.Valid {
				post.Parent = parent.Int64
			} else {
				post.Parent = 0
			}
			resPosts = append(resPosts, *post)
		}
		if er, ok := rows.Err().(*pgconn.PgError); ok {
			switch er.Code {
			case "23503":
				result := user.ErrMsg{Message: "Can't find post author by nickname: "}
				res, _ := easyjson.Marshal(result)
				ctx.SetBody(res)
				ctx.SetStatusCode(404)
				ctx.SetContentType("application/json")
				return
			case "42704":
				result := user.ErrMsg{Message: "Parent post was created in another thread"}
				res, _ := easyjson.Marshal(result)
				ctx.SetBody(res)
				ctx.SetStatusCode(409)
				ctx.SetContentType("application/json")
				return
			}
		}
		res, _ := easyjson.Marshal(resPosts)
		ctx.SetBody(res)
		ctx.SetStatusCode(201)
		ctx.SetContentType("application/json")
	} else {
		result := user.ErrMsg{Message: "Can't find post author by nickname: "}
		res, _ := easyjson.Marshal(result)
		ctx.SetBody(res)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
	}
}

func (thread *Thread) Details(ctx *fasthttp.RequestCtx) {
	SLUG := ctx.UserValue("slug_or_id").(string)
	id, err := strconv.Atoi(SLUG)
	var row pgx.Row
	thr := &ResThread{}
	if err == nil {
		slug := sql.NullString{}
		row = thread.DB.QueryRow(context.Background(), `SELECT author, created, forum, id, message, slug, title, votes
										from threads where id=$1`,
			id)
		err = row.Scan(&thr.Author, &thr.Created, &thr.Forum, &thr.ID, &thr.Message, &slug, &thr.Title, &thr.Votes)
		if slug.Valid {
			thr.Slug = slug.String
		}
		if thr.ID == 0 {
			result := user.ErrMsg{Message: "Can't find thread by ID: "}
			res, _ := easyjson.Marshal(result)
			ctx.SetBody(res)
			ctx.SetStatusCode(404)
			ctx.SetContentType("application/json")
			return
		}
	} else {
		id = -1
		slug := sql.NullString{}
		row = thread.DB.QueryRow(context.Background(), `SELECT author, created, forum, id, message, slug, title, votes
										from threads where slug=$1`,
			SLUG)
		err = row.Scan(&thr.Author, &thr.Created, &thr.Forum, &thr.ID, &thr.Message, &slug, &thr.Title, &thr.Votes)
		if slug.Valid {
			thr.Slug = slug.String
		}
		if err != nil {
			result := user.ErrMsg{Message: "Can't find thread by slug: "}
			res, _ := easyjson.Marshal(result)
			ctx.SetBody(res)
			ctx.SetStatusCode(404)
			ctx.SetContentType("application/json")
			return
		}
	}
	res, _ := easyjson.Marshal(thr)
	ctx.SetBody(res)
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
}

func (thread *Thread) Update(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	var thr ResThread
	easyjson.Unmarshal(ctx.PostBody(), &thr)
	argc := 1
	var args []interface{}
	ID, err := strconv.Atoi(slugOrID)
	query := "UPDATE threads SET title=COALESCE(NULLIF($1, ''), title), message=COALESCE(NULLIF($2, ''), message) "
	args = append(args, thr.Title)
	argc++
	args = append(args, thr.Message)
	argc++
	if err == nil {
		query += " WHERE id=$" + strconv.Itoa(argc)
		args = append(args, ID)
	} else {
		query += " WHERE slug=$" + strconv.Itoa(argc)
		args = append(args, slugOrID)
	}
	query += " RETURNING id, title, author, forum, message, votes, slug, created;"

	row := thread.DB.QueryRow(context.Background(), query, args...)
	slug := sql.NullString{}
	err = row.Scan(&thr.ID, &thr.Title, &thr.Author, &thr.Forum, &thr.Message, &thr.Votes, &slug, &thr.Created)
	if slug.Valid {
		thr.Slug = slug.String
	}
	if thr.ID == 0 {
		errMsg := &user.ErrMsg{Message: "Can't find thread by slug: "}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}
	res, _ := easyjson.Marshal(thr)
	ctx.SetBody(res)
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
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
	ID, err := strconv.Atoi(slugOrID)
	var thr pgx.Row
	if err == nil {
		thr = thread.DB.QueryRow(context.Background(), `SELECT COUNT(1) FROM threads WHERE id=$1`, ID)
	} else {
		thr = thread.DB.QueryRow(context.Background(), `SELECT COUNT(1) FROM threads WHERE slug=$1`, slugOrID)
	}
	var threadFound int
	thr.Scan(&threadFound)
	if threadFound == 0 {
		errMsg := &user.ErrMsg{Message: "Can't find thread by slug: "}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

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
				query += "AND p.ID < $" + strconv.Itoa(argc)
				argc++
			} else {
				query += "AND p.ID > $" + strconv.Itoa(argc)
				argc++
			}
			args = append(args, since)
		}
		if desc == "true" {
			query += " ORDER BY created DESC, id DESC "
		} else {
			query += " ORDER BY created, p.id "
		}
		if len(limit) != 0 {
			query += " LIMIT $" + strconv.Itoa(argc)
			argc++
			args = append(args, limit)
		}
	case "tree":
		var sinceQuery string
		var descQuery string
		var limitSQL string
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
			sinceQuery += " AND (path "
			if desc == "true" {
				sinceQuery += " < "
			} else {
				sinceQuery += " > "
			}
			sinceQuery += "(SELECT path FROM posts WHERE id=$" + strconv.Itoa(argc) + ")) "
			args = append(args, since)
			argc++
		} else {
			sinceQuery = ""
		}

		if desc == "true" {
			descQuery = " DESC "
		} else {
			descQuery = ""
		}

		if len(limit) != 0 {
			limitSQL = " LIMIT $" + strconv.Itoa(argc)
			lim, _ := strconv.Atoi(limit)
			args = append(args, lim)
			argc++
		} else {
			limitSQL = ""
		}
		query += sinceQuery + `ORDER BY path` + descQuery + limitSQL
	default:
		var IDSubquery string
		var descQuery string
		var sinceQuery string
		var limitQuery string

		ID, err := strconv.Atoi(slugOrID)
		if err == nil {
			IDSubquery += "thr.id = $1 "
			args = append(args, ID)
		} else {
			IDSubquery += "thr.slug = $1 "
			args = append(args, slugOrID)
		}
		argc := 2

		if desc == "true" {
			descQuery = " DESC "
		} else {
			descQuery = ""
		}

		if len(since) != 0 {
			sinceQuery = " AND p.id "
			if desc == "true" {
				sinceQuery += "<"
			} else {
				sinceQuery += ">"
			}
			sinceQuery += "(SELECT path[1] FROM posts WHERE id = $" + strconv.Itoa(argc) + ") "
			argc++
			args = append(args, since)
		} else {
			sinceQuery = ""
		}

		if len(limit) != 0 {
			limitQuery = "LIMIT $" + strconv.Itoa(argc)
			lim, _ := strconv.Atoi(limit)
			args = append(args, lim)
		} else {
			limitQuery = "LIMIT 100000"
		}

		query = `
    SELECT p.id, p.thread, p.created, p.message, COALESCE(p.parent, 0),
      p.author, p.forum
      FROM posts p
    JOIN threads thr ON p.thread = thr.id 
      WHERE path[1] IN (
        SELECT p.id FROM posts p JOIN threads thr ON p.thread = thr.id WHERE `
		query += "" + IDSubquery + ` AND parent IS NULL `
		query += sinceQuery
		query += "ORDER BY id "
		query += descQuery
		query += limitQuery
		query += `
      ) AND `
		query += IDSubquery
		query += `ORDER BY path[1]` + descQuery
		query += `, path;
    `

	}
	rows, _ := thread.DB.Query(context.Background(), query, args...)
	defer rows.Close()
	result := make(ResThreads, 0)
	for rows.Next() {
		var thr ResThread
		rows.Scan(&thr.ID, &thr.Thread, &thr.Created, &thr.Message, &thr.Parent, &thr.Author, &thr.Forum)
		result = append(result, thr)
	}
	res, _ := easyjson.Marshal(result)
	ctx.SetBody(res)
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
}

func (thread *Thread) Vote(ctx *fasthttp.RequestCtx) {
	threadSlug := ctx.UserValue("slug_or_id").(string)
	var row pgx.Row
	id, err := strconv.Atoi(threadSlug)
	thr := &ResThread{}
	if err == nil {
		slug := sql.NullString{}
		row = thread.DB.QueryRow(context.Background(), `SELECT author, created, forum, id, message, slug, title from threads where id=$1`, id)
		row.Scan(&thr.Author, &thr.Created, &thr.Forum, &thr.ID, &thr.Message, &slug, &thr.Title)
		if slug.Valid {
			thr.Slug = slug.String
		}
	} else {
		id = -1
		row = thread.DB.QueryRow(context.Background(), `SELECT author, created, forum, id, message, slug, title from threads where slug=$1`, threadSlug)
		row.Scan(&thr.Author, &thr.Created, &thr.Forum, &thr.ID, &thr.Message, &thr.Slug, &thr.Title)
	}
	if len(threadSlug) != 0 {
		vote := &Vote{}
		easyjson.Unmarshal(ctx.PostBody(), vote)
		row := thread.DB.QueryRow(context.Background(), `INSERT INTO votes as vote 
                (nickname, thread, voice)
                VALUES ($1, $2, $3) 
                ON CONFLICT ON CONSTRAINT votes_user_thread_unique DO
                UPDATE SET voice = $3 WHERE vote.voice <> $3`, vote.Nickname, thr.ID, vote.Voice)
		err := row.Scan()
		if err, ok := err.(*pgconn.PgError); ok {
			switch err.Code {
			case "23503":
				result := user.ErrMsg{Message: "Can't find user by nickname: "}
				res, _ := easyjson.Marshal(result)
				ctx.SetBody(res)
				ctx.SetStatusCode(404)
				ctx.SetContentType("application/json")
				return
			}
		}
		row = thread.DB.QueryRow(context.Background(), `SELECT votes FROM threads WHERE id=$1`, thr.ID)
		err = row.Scan(&thr.Votes)
		if err != nil {
			result := user.ErrMsg{Message: "Can't find thread by slug: " + threadSlug}
			res, _ := easyjson.Marshal(result)
			ctx.SetBody(res)
			ctx.SetStatusCode(404)
			ctx.SetContentType("application/json")
			return
		}

		res, _ := easyjson.Marshal(thr)
		ctx.SetBody(res)
		ctx.SetStatusCode(200)
		ctx.SetContentType("application/json")
	} else {
		result := user.ErrMsg{Message: "Can't find thread by slug: "}
		res, _ := easyjson.Marshal(result)
		ctx.SetBody(res)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
	}
}
