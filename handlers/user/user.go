package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
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

//easyjson:json
type ErrMsg struct {
	Message string `json:"message,omitempty"`
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
		rows, err := user.DB.Query("SELECT nickname, fullname, about, email "+
			"FROM users "+
			"WHERE nickname=$1 or email=$2",
			request.Nickname,
			request.Email)
		if err != nil {
			log.Println(err)
		}
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
	nickname := ctx.UserValue("nickname")
	rows, _ := user.DB.Query("SELECT nickname, fullname, about, email "+
		"FROM users "+
		"WHERE nickname=$1",
		nickname)
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&request.Nickname, &request.Fullname, &request.About, &request.Email)
		resp, err := json.Marshal(request)
		if err != nil {
			fmt.Println(err)
		}
		ctx.Response.SetBody(resp)
		ctx.SetContentType("application/json")
		ctx.Response.SetStatusCode(200)
		return
	} else {
		errMsg := &ErrMsg{Message: fmt.Sprintf("Can't find user with nickname %s", nickname)}
		response, _ := json.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
	}
}

func (user *User) Update(ctx *fasthttp.RequestCtx) {
	request := &Req{}
	nickname := ctx.UserValue("nickname").(string)
	err := easyjson.Unmarshal(ctx.PostBody(), request)
	if err != nil {
		log.Println(err)
	}
	if len(request.About) == 0 && len(request.Email) == 0 && len(request.Nickname) == 0 && len(request.Fullname) == 0 {
		rows, _ := user.DB.Query("SELECT nickname, fullname, about, email FROM users WHERE nickname=$1", nickname)
		rows.Next()
		rows.Scan(&request.Nickname, &request.Fullname, &request.About, &request.Email)
		defer rows.Close()
		response, _ := easyjson.Marshal(request)
		ctx.SetBody(response)
		ctx.SetStatusCode(200)
		ctx.SetContentType("application/json")
		return
	}
	request.Nickname = nickname
	var row *sql.Row
	if len(request.Email) != 0 {
		row = user.DB.QueryRow("UPDATE users "+
			"SET fullname=CASE WHEN $1 <> '' THEN $1 ELSE fullname END,"+
			"about=CASE WHEN $2 <> '' THEN $2 ELSE about END,"+
			"email = $3"+
			"WHERE nickname=$4 RETURNING nickname, fullname, about, email",
			request.Fullname,
			request.About,
			request.Email,
			nickname,
		)
		err = row.Scan(&request.Nickname, &request.Fullname, &request.About, &request.Email)
	} else {
		row = user.DB.QueryRow("UPDATE users "+
			"SET fullname=CASE WHEN $1 <> '' THEN $1 ELSE fullname END,"+
			"about=CASE WHEN $2 <> '' THEN $2 ELSE about END "+
			"WHERE nickname=$3 RETURNING nickname, fullname, about, email",
			request.Fullname,
			request.About,
			nickname,
		)
		err = row.Scan(&request.Nickname, &request.Fullname, &request.About, &request.Email)
	}
	if err, ok := err.(*pq.Error); ok {
		switch err.Code {
		case "23505":
			errMsg := &ErrMsg{Message: fmt.Sprintf("This email is already registered by user: %s", nickname)}
			response, _ := easyjson.Marshal(errMsg)
			ctx.SetBody(response)
			ctx.SetStatusCode(409)
			ctx.SetContentType("application/json")
			return
		}
	}
	if err == sql.ErrNoRows { // No such user
		errMsg := &ErrMsg{Message: fmt.Sprintf("Can't find user with nickname %s", nickname)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}
	ctx.Response.SetStatusCode(200)
	resp, _ := easyjson.Marshal(request)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
}
