package thread

import (
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

//easyjson:json
type ResThread struct {
	Id       int       `json:"id,omitempty"`
	Parent   int       `json:"parent,omitempty"`
	Author   string    `json:"author,omitempty"`
	Message  string    `json:"message,omitempty"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	Thread   int       `json:"thread,omitempty"`
	Created  time.Time `json:"created,omitempty"`
}

//easyjson:json
type ResThreads []ResThread

func Create(ctx *fasthttp.RequestCtx) {
	//SLUG := ctx.UserValue("slug_or_id").(string)
	threads := &ResThreads{}
	log.Println(string(ctx.PostBody()))

	if err := easyjson.Unmarshal(ctx.PostBody(), threads); err != nil {
		log.Println(err)
	}
	for thread := range *threads {
		fmt.Println(thread)
	}
	ctx.SetBody([]byte("[]"))
	ctx.SetStatusCode(201)
	ctx.SetContentType("application/json")
}

func Details(ctx *fasthttp.RequestCtx) {

}

func Update(ctx *fasthttp.RequestCtx) {

}

func Messages(ctx *fasthttp.RequestCtx) {

}

func Vote(ctx *fasthttp.RequestCtx) {

}
