package user

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"log"
)

//easyjson:json
type Req struct {
	Nickname string `json:"nickname,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email,omitempty"`
}

//easyjson:json
type Reqs []Req

type User struct {
	DB *sql.DB
}

func (user *User) Create(ctx *fasthttp.RequestCtx) {
	request := &Req{}
	request.Nickname = ctx.UserValue("nickname").(string)
	easyjson.Unmarshal(ctx.PostBody(), request)
	_, err := user.DB.Exec("INSERT INTO users (nickname, fullname, about, email) "+
		"VALUES ($1, $2, $3, $4)",
		request.Nickname,
		request.Fullname,
		request.About,
		request.Email,
	)
	if err != nil {
		rows, _ := user.DB.Query("SELECT nickname, fullname, about, email "+
			"FROM users "+
			"WHERE nickname=$1 or email=$2",
			request.Nickname,
			request.Email)
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Println(err)
			}
		}(rows)
		results := make(Reqs, 0)
		for rows.Next() {
			user := &Req{}
			err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
			if err != nil {
				fmt.Println(err)
			}
			results = append(results, *user)
		}
		resp, err := easyjson.Marshal(results)
		if err != nil {
			fmt.Println(err)
		}
		ctx.Response.SetBody(resp)
		ctx.SetContentType("application/json")
		ctx.Response.SetStatusCode(409)
		return
	}
	ctx.Response.SetStatusCode(201)
	resp, err := easyjson.Marshal(request)
	if err != nil {
		log.Println(err)
	}
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
}

func (user *User) Profile(ctx *fasthttp.RequestCtx) {
	request := &Req{}
	request.Nickname = ctx.UserValue("nickname").(string)
	rows, _ := user.DB.Query("SELECT fullname, about, email "+
		"FROM users "+
		"WHERE nickname=$1",
		request.Nickname)
	if rows.Next() {
		rows.Scan(&request.Fullname, &request.About, &request.Email)
		resp, err := easyjson.Marshal(request)
		if err != nil {
			fmt.Println(err)
		}
		ctx.Response.SetBody(resp)
		ctx.SetContentType("application/json")
		ctx.Response.SetStatusCode(200)
		return
	} else {
		var b bytes.Buffer
		b.Grow(100)
		fmt.Fprintf(&b, "Can't find user with nickname %s", request.Nickname)
		ctx.SetBody(b.Bytes())
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
	}
}

func (user *User) Update(ctx *fasthttp.RequestCtx) {
	request := &Req{}
	request.Nickname = ctx.UserValue("nickname").(string)
	easyjson.Unmarshal(ctx.PostBody(), request)
	result, err := user.DB.Exec("UPDATE users "+
		"SET fullname=$1, about=$2, email=$3"+
		"WHERE nickname=$4",
		request.Fullname,
		request.About,
		request.Email,
		request.Nickname,
	)
	if err != nil {
		var b bytes.Buffer
		b.Grow(100)
		fmt.Fprintf(&b, "Can't find user with nickname %s", request.Nickname)
		ctx.SetBody(b.Bytes())
		ctx.SetStatusCode(409)
		ctx.SetContentType("application/json")
		return
	}
	if res, _ := result.RowsAffected(); res == 0 {
		var b bytes.Buffer
		b.Grow(100)
		fmt.Fprintf(&b, "Can't find user with nickname %s", request.Nickname)
		ctx.SetBody(b.Bytes())
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}
	ctx.Response.SetStatusCode(200)
	resp, _ := easyjson.Marshal(request)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
}
