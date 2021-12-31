package post

import (
	"DBMS/handlers/user"
	"database/sql"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"log"
)

type Post struct {
	DB *sql.DB
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
	} `json:"post"`
}

//easyjson:json
type ResWithAuthor struct {
	Author struct {
		About    string `json:"about"`
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
		Nickname string `json:"nickname"`
	} `json:"author"`
	Post struct {
		Author   string `json:"author"`
		Created  string `json:"created"`
		Forum    string `json:"forum"`
		ID       int    `json:"id"`
		Message  string `json:"message"`
		Thread   int    `json:"thread"`
		IsEdited bool   `json:"isEdited"`
	} `json:"post"`
}

//easyjson:json
type ResWithThread struct {
	Post struct {
		Author   string `json:"author"`
		Created  string `json:"created"`
		Forum    string `json:"forum"`
		ID       int    `json:"id"`
		Message  string `json:"message"`
		Thread   int    `json:"thread"`
		IsEdited bool   `json:"isEdited"`
	} `json:"post"`
	Thread struct {
		Author  string `json:"author"`
		Created string `json:"created"`
		Forum   string `json:"forum"`
		ID      int    `json:"id"`
		Message string `json:"message"`
		Slug    string `json:"slug"`
		Title   string `json:"title"`
	} `json:"thread"`
}

//easyjson:json
type ResWithForum struct {
	Post struct {
		Author   string `json:"author"`
		Created  string `json:"created"`
		Forum    string `json:"forum"`
		ID       int    `json:"id"`
		Message  string `json:"message"`
		Thread   int    `json:"thread"`
		IsEdited bool   `json:"isEdited"`
	} `json:"post"`
	Forum struct {
		Posts   int    `json:"posts"`
		Slug    string `json:"slug"`
		Threads int    `json:"threads"`
		Title   string `json:"title"`
		User    string `json:"user"`
	} `json:"forum"`
}

//easyjson:json
type ResWithAuthorAndThread struct {
	Post struct {
		Author   string `json:"author"`
		Created  string `json:"created"`
		Forum    string `json:"forum"`
		ID       int    `json:"id"`
		Message  string `json:"message"`
		Thread   int    `json:"thread"`
		IsEdited bool   `json:"isEdited"`
	} `json:"post"`
	Author struct {
		About    string `json:"about"`
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
		Nickname string `json:"nickname"`
	} `json:"author"`
	Thread struct {
		Author  string `json:"author"`
		Created string `json:"created"`
		Forum   string `json:"forum"`
		ID      int    `json:"id"`
		Message string `json:"message"`
		Slug    string `json:"slug"`
		Title   string `json:"title"`
	} `json:"thread"`
}

//easyjson:json
type ResWithAuthorAndForum struct {
	Post struct {
		Author   string `json:"author"`
		Created  string `json:"created"`
		Forum    string `json:"forum"`
		ID       int    `json:"id"`
		Message  string `json:"message"`
		Thread   int    `json:"thread"`
		IsEdited bool   `json:"isEdited"`
	} `json:"post"`
	Author struct {
		About    string `json:"about"`
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
		Nickname string `json:"nickname"`
	} `json:"author"`
	Forum struct {
		Posts   int    `json:"posts"`
		Slug    string `json:"slug"`
		Threads int    `json:"threads"`
		Title   string `json:"title"`
		User    string `json:"user"`
	} `json:"forum"`
}

//easyjson:json
type ResPost struct {
	Author   string `json:"author,omitempty"`
	Created  string `json:"created,omitempty"`
	Forum    string `json:"forum,omitempty"`
	ID       int    `json:"id,omitempty"`
	Message  string `json:"message,omitempty"`
	Thread   int    `json:"thread,omitempty"`
	IsEdited bool   `json:"isEdited,omitempty"`
}

func (post *Post) Details(ctx *fasthttp.RequestCtx) {
	ID := ctx.UserValue("id")
	related := string(ctx.QueryArgs().Peek("related"))
	log.Println(related)
	var resultPost Res

	row := post.DB.QueryRow(`SELECT author, created, forum, id, message, thread, isedited FROM posts WHERE id=$1`, ID)

	err := row.Scan(&resultPost.Post.Author, &resultPost.Post.Created, &resultPost.Post.Forum,
		&resultPost.Post.ID, &resultPost.Post.Message, &resultPost.Post.Thread, &resultPost.Post.IsEdited)
	if err != nil {
		log.Println(err)
		errMsg := &user.ErrMsg{Message: "Can't find post with id: "}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

	if related == "user" {
		var result ResWithAuthor
		result.Post = resultPost.Post
		usr := post.DB.QueryRow(`SELECT about, email, fullname, nickname FROM users WHERE nickname=$1`, resultPost.Post.Author)
		err := usr.Scan(&result.Author.About, &result.Author.Email, &result.Author.Fullname, &result.Author.Nickname)
		if err != nil {
			log.Println(err)
		}
		res, _ := easyjson.Marshal(result)
		ctx.SetBody(res)
		ctx.SetStatusCode(200)
		ctx.SetContentType("application/json")
		return
	} else if related == "thread" {
		var result ResWithThread
		result.Post = resultPost.Post
		thread := post.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title FROM threads WHERE id=$1`, resultPost.Post.Thread)
		err := thread.Scan(&result.Thread.Author, &result.Thread.Created, &result.Thread.Forum, &result.Thread.ID, &result.Thread.Message,
			&result.Thread.Slug, &result.Thread.Title)
		if err != nil {
			log.Println(err)
		}
		res, _ := easyjson.Marshal(result)
		ctx.SetBody(res)
		ctx.SetStatusCode(200)
		ctx.SetContentType("application/json")
		return
	} else if related == "forum" {
		var result ResWithForum
		result.Post = resultPost.Post
		forum := post.DB.QueryRow(`SELECT posts, slug, threads, title, "user" FROM forums WHERE slug=$1`, resultPost.Post.Forum)
		err := forum.Scan(&result.Forum.Posts, &result.Forum.Slug, &result.Forum.Threads, &result.Forum.Title, &result.Forum.User)
		if err != nil {
			log.Println(err)
		}
		res, _ := easyjson.Marshal(result)
		ctx.SetBody(res)
		ctx.SetStatusCode(200)
		ctx.SetContentType("application/json")
		return
	} else if related == "user,thread" {
		var result ResWithAuthorAndThread
		result.Post = resultPost.Post
		thread := post.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title FROM threads WHERE id=$1`, resultPost.Post.Thread)
		err := thread.Scan(&result.Thread.Author, &result.Thread.Created, &result.Thread.Forum, &result.Thread.ID, &result.Thread.Message,
			&result.Thread.Slug, &result.Thread.Title)
		if err != nil {
			log.Println(err)
		}
		usr := post.DB.QueryRow(`SELECT about, email, fullname, nickname FROM users WHERE nickname=$1`, resultPost.Post.Author)
		err = usr.Scan(&result.Author.About, &result.Author.Email, &result.Author.Fullname, &result.Author.Nickname)
		if err != nil {
			log.Println(err)
		}
		res, _ := easyjson.Marshal(result)
		ctx.SetBody(res)
		ctx.SetStatusCode(200)
		ctx.SetContentType("application/json")
		return
	} else if related == "user,forum" {
		var result ResWithAuthorAndForum
		result.Post = resultPost.Post
		forum := post.DB.QueryRow(`SELECT posts, slug, threads, title, "user" FROM forums WHERE slug=$1`, resultPost.Post.Forum)
		err := forum.Scan(&result.Forum.Posts, &result.Forum.Slug, &result.Forum.Threads, &result.Forum.Title, &result.Forum.User)
		if err != nil {
			log.Println(err)
		}
		usr := post.DB.QueryRow(`SELECT about, email, fullname, nickname FROM users WHERE nickname=$1`, resultPost.Post.Author)
		err = usr.Scan(&result.Author.About, &result.Author.Email, &result.Author.Fullname, &result.Author.Nickname)
		if err != nil {
			log.Println(err)
		}
		res, _ := easyjson.Marshal(result)
		ctx.SetBody(res)
		ctx.SetStatusCode(200)
		ctx.SetContentType("application/json")
		return
	}

	res, _ := easyjson.Marshal(resultPost)
	ctx.SetBody(res)
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
}

func (post *Post) UpdateMessage(ctx *fasthttp.RequestCtx) {
	ID := ctx.UserValue("id")
	var resultPost ResPost
	easyjson.Unmarshal(ctx.PostBody(), &resultPost)
	row := post.DB.QueryRow(`UPDATE posts SET message=$1 WHERE id=$2 RETURNING author, created, forum, id, isedited, message, thread`, resultPost.Message, ID)
	err := row.Scan(&resultPost.Author, &resultPost.Created, &resultPost.Forum, &resultPost.ID, &resultPost.IsEdited, &resultPost.Message, &resultPost.Thread)
	if err != nil {
		log.Println(err)
	}
	resultPost.IsEdited = true
	res, _ := easyjson.Marshal(resultPost)
	ctx.SetBody(res)
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
}
