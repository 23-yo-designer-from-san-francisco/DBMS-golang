package post

import (
	"database/sql"
	"github.com/valyala/fasthttp"
)

type Post struct {
	db *sql.DB
}

func Details(ctx *fasthttp.RequestCtx) {

}

func UpdateMessage(ctx *fasthttp.RequestCtx) {

}
