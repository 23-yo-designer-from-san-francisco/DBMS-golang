// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package post

import (
	json "encoding/json"
	pgxpool "github.com/jackc/pgx/v4/pgxpool"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson5a72dc82DecodeDBMSHandlersPost(in *jlexer.Lexer, out *Thread) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "author":
			out.Author = string(in.String())
		case "created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
			}
		case "forum":
			out.Forum = string(in.String())
		case "id":
			out.ID = int(in.Int())
		case "message":
			out.Message = string(in.String())
		case "slug":
			out.Slug = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "votes":
			out.Votes = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeDBMSHandlersPost(out *jwriter.Writer, in Thread) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix[1:])
		out.String(string(in.Author))
	}
	{
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.Raw((in.Created).MarshalJSON())
	}
	{
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.String(string(in.Forum))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"slug\":"
		out.RawString(prefix)
		out.String(string(in.Slug))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"votes\":"
		out.RawString(prefix)
		out.Int(int(in.Votes))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Thread) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeDBMSHandlersPost(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Thread) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeDBMSHandlersPost(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Thread) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeDBMSHandlersPost(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Thread) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeDBMSHandlersPost(l, v)
}
func easyjson5a72dc82DecodeDBMSHandlersPost1(in *jlexer.Lexer, out *ResPost) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "author":
			out.Author = string(in.String())
		case "created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
			}
		case "forum":
			out.Forum = string(in.String())
		case "id":
			out.ID = int(in.Int())
		case "message":
			out.Message = string(in.String())
		case "thread":
			out.Thread = int(in.Int())
		case "isEdited":
			out.IsEdited = bool(in.Bool())
		case "parent":
			out.Parent = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeDBMSHandlersPost1(out *jwriter.Writer, in ResPost) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix[1:])
		out.String(string(in.Author))
	}
	{
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.Raw((in.Created).MarshalJSON())
	}
	{
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.String(string(in.Forum))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		out.Int(int(in.Thread))
	}
	{
		const prefix string = ",\"isEdited\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsEdited))
	}
	{
		const prefix string = ",\"parent\":"
		out.RawString(prefix)
		out.Int64(int64(in.Parent))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResPost) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeDBMSHandlersPost1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResPost) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeDBMSHandlersPost1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResPost) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeDBMSHandlersPost1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResPost) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeDBMSHandlersPost1(l, v)
}
func easyjson5a72dc82DecodeDBMSHandlersPost2(in *jlexer.Lexer, out *Res) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "post":
			easyjson5a72dc82Decode(in, &out.Post)
		case "author":
			if in.IsNull() {
				in.Skip()
				out.Author = nil
			} else {
				if out.Author == nil {
					out.Author = new(Author)
				}
				(*out.Author).UnmarshalEasyJSON(in)
			}
		case "thread":
			if in.IsNull() {
				in.Skip()
				out.Thread = nil
			} else {
				if out.Thread == nil {
					out.Thread = new(Thread)
				}
				(*out.Thread).UnmarshalEasyJSON(in)
			}
		case "forum":
			if in.IsNull() {
				in.Skip()
				out.Forum = nil
			} else {
				if out.Forum == nil {
					out.Forum = new(Forum)
				}
				(*out.Forum).UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeDBMSHandlersPost2(out *jwriter.Writer, in Res) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"post\":"
		out.RawString(prefix[1:])
		easyjson5a72dc82Encode(out, in.Post)
	}
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		if in.Author == nil {
			out.RawString("null")
		} else {
			(*in.Author).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		if in.Thread == nil {
			out.RawString("null")
		} else {
			(*in.Thread).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		if in.Forum == nil {
			out.RawString("null")
		} else {
			(*in.Forum).MarshalEasyJSON(out)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Res) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeDBMSHandlersPost2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Res) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeDBMSHandlersPost2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Res) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeDBMSHandlersPost2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Res) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeDBMSHandlersPost2(l, v)
}
func easyjson5a72dc82Decode(in *jlexer.Lexer, out *struct {
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	ID       int       `json:"id"`
	Message  string    `json:"message"`
	Thread   int       `json:"thread"`
	IsEdited bool      `json:"isEdited"`
	Parent   int64     `json:"parent"`
}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "author":
			out.Author = string(in.String())
		case "created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
			}
		case "forum":
			out.Forum = string(in.String())
		case "id":
			out.ID = int(in.Int())
		case "message":
			out.Message = string(in.String())
		case "thread":
			out.Thread = int(in.Int())
		case "isEdited":
			out.IsEdited = bool(in.Bool())
		case "parent":
			out.Parent = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82Encode(out *jwriter.Writer, in struct {
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	ID       int       `json:"id"`
	Message  string    `json:"message"`
	Thread   int       `json:"thread"`
	IsEdited bool      `json:"isEdited"`
	Parent   int64     `json:"parent"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix[1:])
		out.String(string(in.Author))
	}
	{
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.Raw((in.Created).MarshalJSON())
	}
	{
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.String(string(in.Forum))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		out.Int(int(in.Thread))
	}
	{
		const prefix string = ",\"isEdited\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsEdited))
	}
	{
		const prefix string = ",\"parent\":"
		out.RawString(prefix)
		out.Int64(int64(in.Parent))
	}
	out.RawByte('}')
}
func easyjson5a72dc82DecodeDBMSHandlersPost3(in *jlexer.Lexer, out *Post) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "DB":
			if in.IsNull() {
				in.Skip()
				out.DB = nil
			} else {
				if out.DB == nil {
					out.DB = new(pgxpool.Pool)
				}
				easyjson5a72dc82DecodeGithubComJackcPgxV4Pgxpool(in, out.DB)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeDBMSHandlersPost3(out *jwriter.Writer, in Post) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"DB\":"
		out.RawString(prefix[1:])
		if in.DB == nil {
			out.RawString("null")
		} else {
			easyjson5a72dc82EncodeGithubComJackcPgxV4Pgxpool(out, *in.DB)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Post) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeDBMSHandlersPost3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Post) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeDBMSHandlersPost3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Post) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeDBMSHandlersPost3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Post) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeDBMSHandlersPost3(l, v)
}
func easyjson5a72dc82DecodeGithubComJackcPgxV4Pgxpool(in *jlexer.Lexer, out *pgxpool.Pool) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeGithubComJackcPgxV4Pgxpool(out *jwriter.Writer, in pgxpool.Pool) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}
func easyjson5a72dc82DecodeDBMSHandlersPost4(in *jlexer.Lexer, out *Forum) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "posts":
			out.Posts = int(in.Int())
		case "slug":
			out.Slug = string(in.String())
		case "threads":
			out.Threads = int(in.Int())
		case "title":
			out.Title = string(in.String())
		case "user":
			out.User = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeDBMSHandlersPost4(out *jwriter.Writer, in Forum) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"posts\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Posts))
	}
	{
		const prefix string = ",\"slug\":"
		out.RawString(prefix)
		out.String(string(in.Slug))
	}
	{
		const prefix string = ",\"threads\":"
		out.RawString(prefix)
		out.Int(int(in.Threads))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix)
		out.String(string(in.User))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Forum) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeDBMSHandlersPost4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Forum) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeDBMSHandlersPost4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Forum) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeDBMSHandlersPost4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Forum) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeDBMSHandlersPost4(l, v)
}
func easyjson5a72dc82DecodeDBMSHandlersPost5(in *jlexer.Lexer, out *Author) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "about":
			out.About = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "fullname":
			out.Fullname = string(in.String())
		case "nickname":
			out.Nickname = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeDBMSHandlersPost5(out *jwriter.Writer, in Author) {
	out.RawByte('{')
	first := true
	_ = first
	if in.About != "" {
		const prefix string = ",\"about\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.About))
	}
	if in.Email != "" {
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	if in.Fullname != "" {
		const prefix string = ",\"fullname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Fullname))
	}
	if in.Nickname != "" {
		const prefix string = ",\"nickname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Nickname))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Author) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeDBMSHandlersPost5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Author) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeDBMSHandlersPost5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Author) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeDBMSHandlersPost5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Author) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeDBMSHandlersPost5(l, v)
}
