package thread

import (
	"DBMS/handlers/user"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"log"
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
	//log.Println("POST /thread/{slug_or_id}/create")
	SLUG := ctx.UserValue("slug_or_id").(string)
	id, err := strconv.Atoi(SLUG)
	var row *sql.Row
	var forumTitle string
	var forumSwag string
	if err == nil {
		row = thread.DB.QueryRow(`SELECT title, forum from threads where id=$1`, id)
		err = row.Scan(&forumTitle, &forumSwag)
		//if err != nil {
		//log.Println(err)
		//}
	} else {
		id = -1
		row = thread.DB.QueryRow(`SELECT title, forum, id from threads where slug=$1`, SLUG)
		err = row.Scan(&forumTitle, &forumSwag, &id)
		//if err != nil {
		//	log.Println(err)
		//}
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
		query := `INSERT INTO posts (parent, author, message, thread, forum) VALUES`
		var values []interface{}
		usersQuery := `INSERT INTO forum_users (user_nickname, forum_swag) VALUES`
		var forumUsers []interface{}
		for i, thr := range *threads {
			value := fmt.Sprintf(
				"(NULLIF($%d, 0), $%d, $%d, $%d, $%d),",
				i*5+1, i*5+2, i*5+3, i*5+4, i*5+5,
			)
			usersQuery += fmt.Sprintf("($%d, $%d),", i*2+1, i*2+2)
			forumUsers = append(forumUsers, thr.Author, forumSwag)
			query += value
			values = append(values, thr.Parent, thr.Author, thr.Message, id, forumSwag)
		}
		query = strings.TrimSuffix(query, ",")
		usersQuery = strings.TrimSuffix(usersQuery, ",")
		query += ` RETURNING id, parent, author, message, isedited, forum, thread, created;`
		rows, err := thread.DB.Query(query, values...)
		defer rows.Close()
		if err != nil {
			log.Println(err)
			result := user.ErrMsg{Message: "Parent post was created in another thread"}
			res, _ := easyjson.Marshal(result)
			ctx.SetBody(res)
			ctx.SetStatusCode(409)
			ctx.SetContentType("application/json")
			return
		}

		usersQuery += " ON CONFLICT DO NOTHING"
		userRows, err := thread.DB.Query(usersQuery, forumUsers...)
		defer userRows.Close()
		if err != nil {
			log.Println(err)
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
	} else {
		result := user.ErrMsg{Message: "Can't find post author by nickname: "}
		res, _ := easyjson.Marshal(result)
		ctx.SetBody(res)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
	}
}

func (thread *Thread) Details(ctx *fasthttp.RequestCtx) {
	log.Println("GET /thread/{slug_or_id}/details")
	SLUG := ctx.UserValue("slug_or_id").(string)
	id, err := strconv.Atoi(SLUG)
	var row *sql.Row
	thr := &ResThread{}
	if err == nil {
		swag := sql.NullString{}
		row = thread.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title, votes
										from threads where id=$1`,
			id)
		err = row.Scan(&thr.Author, &thr.Created, &thr.Forum, &thr.ID, &thr.Message, &swag, &thr.Title, &thr.Votes)
		if swag.Valid {
			thr.Slug = swag.String
		}
		if thr.ID == 0 {
			log.Println(err)
			result := user.ErrMsg{Message: "Can't find thread by ID: "}
			res, _ := easyjson.Marshal(result)
			ctx.SetBody(res)
			ctx.SetStatusCode(404)
			ctx.SetContentType("application/json")
			return
		}
	} else {
		id = -1
		row = thread.DB.QueryRow(`SELECT author, created, forum, id, message, COALESCE(slug, ''), title, votes
										from threads where slug=$1`,
			SLUG)
		err = row.Scan(&thr.Author, &thr.Created, &thr.Forum, &thr.ID, &thr.Message, &thr.Slug, &thr.Title, &thr.Votes)
		if err != nil {
			log.Println(err)
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
	log.Println("POST /thread/{slug_or_id}/details")
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

	row := thread.DB.QueryRow(query, args...)
	err = row.Scan()
	if err != nil {
		log.Println(err)
	}
	swag := sql.NullString{}
	err = row.Scan(&thr.ID, &thr.Title, &thr.Author, &thr.Forum, &thr.Message, &thr.Votes, &swag, &thr.Created)
	if swag.Valid {
		thr.Slug = swag.String
	}
	if err != nil {
		log.Println(err)
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
	log.Println("GET /thread/{slug_or_id}/posts")
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
	var thr *sql.Row
	if err == nil {
		thr = thread.DB.QueryRow(`SELECT count(*) FROM threads WHERE id=$1`, ID)
	} else {
		thr = thread.DB.QueryRow(`SELECT count(*) FROM threads WHERE slug=$1`, slugOrID)
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
	rows, err := thread.DB.Query(query, args...)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	result := make(ResThreads, 0)
	//found := false
	for rows.Next() {
		//found = true
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
	log.Println("POST /thread/{slug_or_id}/vote")
	threadSlug := ctx.UserValue("slug_or_id").(string)
	var row *sql.Row
	id, err := strconv.Atoi(threadSlug)
	thr := &ResThread{}
	if err == nil {
		swag := sql.NullString{}
		row = thread.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title from threads where id=$1`, id)
		err = row.Scan(&thr.Author, &thr.Created, &thr.Forum, &thr.ID, &thr.Message, &swag, &thr.Title)
		if err != nil {
			log.Println(err)
		}
		if swag.Valid {
			thr.Slug = swag.String
		}
	} else {
		id = -1
		row = thread.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title from threads where slug=$1`, threadSlug)
		err = row.Scan(&thr.Author, &thr.Created, &thr.Forum, &thr.ID, &thr.Message, &thr.Slug, &thr.Title)
		if err != nil {
			log.Println(err)
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
		err := row.Scan()
		if err, ok := err.(*pq.Error); ok {
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
		row = thread.DB.QueryRow(`SELECT votes FROM threads WHERE id=$1`, thr.ID)
		err = row.Scan(&thr.Votes)
		if err != nil {
			log.Println(err)
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
		log.Println("Here")
		result := user.ErrMsg{Message: "Can't find thread by slug: "}
		res, _ := easyjson.Marshal(result)
		ctx.SetBody(res)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
	}
}
