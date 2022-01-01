package post

import (
	"DBMS/handlers/user"
	"database/sql"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"log"
	"strings"
	"time"
)

type Post struct {
	DB *sql.DB
}

//easyjson:json
type Author struct {
	About    string `json:"about,omitempty"`
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	Nickname string `json:"nickname,omitempty"`
}

//easyjson:json
type Thread struct {
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Forum   string    `json:"forum"`
	ID      int       `json:"id"`
	Message string    `json:"message"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Votes   int       `json:"votes"`
}

type Forum struct {
	Posts   int    `json:"posts"`
	Slug    string `json:"slug"`
	Threads int    `json:"threads"`
	Title   string `json:"title"`
	User    string `json:"user"`
}

//easyjson:json
type Res struct {
	Post struct {
		Author   string `json:"author"`
		Created  string `json:"created"`
		Forum    string `json:"forum"`
		ID       int    `json:"id"`
		Message  string `json:"message"`
		Thread   int    `json:"thread"`
		IsEdited bool   `json:"isEdited"`
		Parent   int64  `json:"parent"`
	} `json:"post"`
	Author *Author `json:"author"`
	Thread *Thread `json:"thread"`
	Forum  *Forum  `json:"forum"`
}

//easyjson:json
type ResPost struct {
	Author   string `json:"author"`
	Created  string `json:"created"`
	Forum    string `json:"forum"`
	ID       int    `json:"id"`
	Message  string `json:"message"`
	Thread   int    `json:"thread"`
	IsEdited bool   `json:"isEdited"`
	Parent   int64  `json:"parent"`
}

func (post *Post) Details(ctx *fasthttp.RequestCtx) {
	log.Println("GET /post/{id}/details")
	ID := ctx.UserValue("id")
	related := string(ctx.QueryArgs().Peek("related"))
	log.Println(related)
	var resultPost Res

	row := post.DB.QueryRow(`SELECT author, created, forum, id, message, thread, isedited, parent FROM posts WHERE id=$1`, ID)
	par := sql.NullInt64{}
	err := row.Scan(&resultPost.Post.Author, &resultPost.Post.Created, &resultPost.Post.Forum,
		&resultPost.Post.ID, &resultPost.Post.Message, &resultPost.Post.Thread, &resultPost.Post.IsEdited, &par)
	if par.Valid {
		resultPost.Post.Parent = par.Int64
	}
	if err != nil {
		log.Println(err)
		errMsg := &user.ErrMsg{Message: "Can't find post with id: "}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

	if strings.Contains(related, "user") {
		resultPost.Author = &Author{}
		usr := post.DB.QueryRow(`SELECT about, email, fullname, nickname FROM users WHERE nickname=$1`, resultPost.Post.Author)
		err := usr.Scan(&resultPost.Author.About, &resultPost.Author.Email, &resultPost.Author.Fullname, &resultPost.Author.Nickname)
		if err != nil {
			log.Println(err)
		}
	}

	if strings.Contains(related, "thread") {
		resultPost.Thread = &Thread{}
		swag := sql.NullString{}
		thread := post.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title, votes FROM threads WHERE id=$1`, resultPost.Post.Thread)
		err := thread.Scan(&resultPost.Thread.Author, &resultPost.Thread.Created, &resultPost.Thread.Forum, &resultPost.Thread.ID, &resultPost.Thread.Message,
			&swag, &resultPost.Thread.Title, &resultPost.Thread.Votes)
		log.Println("thr")
		log.Println(swag.String)
		log.Println("thr")

		if err != nil {
			log.Println(err)
		}
		if swag.Valid {
			resultPost.Thread.Slug = swag.String
			log.Println(swag.String)
			log.Println(resultPost.Thread.Slug)
		}
	}

	if strings.Contains(related, "forum") {
		resultPost.Forum = &Forum{}
		swag := sql.NullString{}
		forum := post.DB.QueryRow(`SELECT posts, slug, threads, title, "user" FROM forums WHERE slug=$1`, resultPost.Post.Forum)
		err := forum.Scan(&resultPost.Forum.Posts, &swag, &resultPost.Forum.Threads, &resultPost.Forum.Title, &resultPost.Forum.User)
		if err != nil {
			log.Println(err)
		}
		if swag.Valid {
			resultPost.Forum.Slug = swag.String
		}
	}
	res, _ := easyjson.Marshal(resultPost)
	ctx.SetBody(res)
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
}

func (post *Post) UpdateMessage(ctx *fasthttp.RequestCtx) {
	log.Println("POST /post/{id}/details")
	ID := ctx.UserValue("id")
	var resultPost ResPost
	easyjson.Unmarshal(ctx.PostBody(), &resultPost)
	newMessage := resultPost.Message
	if len(resultPost.Message) == 0 {
		row := post.DB.QueryRow(`SELECT author, created, forum, id, isedited, message, thread from posts WHERE id=$1`, ID)
		err := row.Scan(&resultPost.Author, &resultPost.Created, &resultPost.Forum, &resultPost.ID, &resultPost.IsEdited, &resultPost.Message, &resultPost.Thread)
		if err != nil {
			log.Println(err)
		}
		res, _ := easyjson.Marshal(resultPost)
		ctx.SetBody(res)
		ctx.SetStatusCode(200)
		ctx.SetContentType("application/json")
		return
	}
	oldPost := post.DB.QueryRow(`SELECT author, created, forum, id, isedited, message, thread FROM posts WHERE id=$1`, ID)
	err := oldPost.Scan(&resultPost.Author, &resultPost.Created, &resultPost.Forum, &resultPost.ID,
		&resultPost.IsEdited, &resultPost.Message, &resultPost.Thread)
	if err != nil {
		log.Println(err)
	}

	if resultPost.ID == 0 {
		err := user.ErrMsg{Message: "Can't find post with id: "}
		res, _ := easyjson.Marshal(err)
		ctx.SetBody(res)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

	if resultPost.Message == newMessage {
		res, _ := easyjson.Marshal(resultPost)
		ctx.SetBody(res)
		ctx.SetStatusCode(200)
		ctx.SetContentType("application/json")
		log.Println("not edited")
		return
	}
	row := post.DB.QueryRow(`
		UPDATE
		posts
		SET
		message=$1
		WHERE
		id=$2
		RETURNING
		message
		`, newMessage, ID)
	err = row.Scan(&resultPost.Message)
	if err != nil {
		log.Println(err)
	}

	resultPost.IsEdited = true
	res, _ := easyjson.Marshal(resultPost)
	ctx.SetBody(res)
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
}
