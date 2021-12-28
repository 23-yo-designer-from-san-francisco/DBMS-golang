// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package forum

import (
	sql "database/sql"
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonC8d74561DecodeDBMSHandlersForum(in *jlexer.Lexer, out *ThreadReq) {
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
		case "id":
			out.ID = int64(in.Int64())
		case "author":
			out.Author = string(in.String())
		case "created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
			}
		case "forum":
			out.Forum = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "slug":
			out.Slug = string(in.String())
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
func easyjsonC8d74561EncodeDBMSHandlersForum(out *jwriter.Writer, in ThreadReq) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	if in.Author != "" {
		const prefix string = ",\"author\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Author))
	}
	if true {
		const prefix string = ",\"created\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.Created).MarshalJSON())
	}
	if in.Forum != "" {
		const prefix string = ",\"forum\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Forum))
	}
	if in.Message != "" {
		const prefix string = ",\"message\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Message))
	}
	if in.Title != "" {
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	if in.Slug != "" {
		const prefix string = ",\"slug\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Slug))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ThreadReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8d74561EncodeDBMSHandlersForum(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ThreadReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8d74561EncodeDBMSHandlersForum(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ThreadReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8d74561DecodeDBMSHandlersForum(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ThreadReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8d74561DecodeDBMSHandlersForum(l, v)
}
func easyjsonC8d74561DecodeDBMSHandlersForum1(in *jlexer.Lexer, out *Reqs) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Reqs, 0, 1)
			} else {
				*out = Reqs{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Req
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC8d74561EncodeDBMSHandlersForum1(out *jwriter.Writer, in Reqs) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Reqs) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8d74561EncodeDBMSHandlersForum1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Reqs) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8d74561EncodeDBMSHandlersForum1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Reqs) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8d74561DecodeDBMSHandlersForum1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Reqs) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8d74561DecodeDBMSHandlersForum1(l, v)
}
func easyjsonC8d74561DecodeDBMSHandlersForum2(in *jlexer.Lexer, out *Req) {
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
		case "slug":
			out.Slug = string(in.String())
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
func easyjsonC8d74561EncodeDBMSHandlersForum2(out *jwriter.Writer, in Req) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Slug != "" {
		const prefix string = ",\"slug\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Slug))
	}
	if in.Title != "" {
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	if in.User != "" {
		const prefix string = ",\"user\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.User))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Req) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8d74561EncodeDBMSHandlersForum2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Req) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8d74561EncodeDBMSHandlersForum2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Req) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8d74561DecodeDBMSHandlersForum2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Req) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8d74561DecodeDBMSHandlersForum2(l, v)
}
func easyjsonC8d74561DecodeDBMSHandlersForum3(in *jlexer.Lexer, out *Forum) {
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
					out.DB = new(sql.DB)
				}
				easyjsonC8d74561DecodeDatabaseSql(in, out.DB)
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
func easyjsonC8d74561EncodeDBMSHandlersForum3(out *jwriter.Writer, in Forum) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"DB\":"
		out.RawString(prefix[1:])
		if in.DB == nil {
			out.RawString("null")
		} else {
			easyjsonC8d74561EncodeDatabaseSql(out, *in.DB)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Forum) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8d74561EncodeDBMSHandlersForum3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Forum) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8d74561EncodeDBMSHandlersForum3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Forum) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8d74561DecodeDBMSHandlersForum3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Forum) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8d74561DecodeDBMSHandlersForum3(l, v)
}
func easyjsonC8d74561DecodeDatabaseSql(in *jlexer.Lexer, out *sql.DB) {
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
func easyjsonC8d74561EncodeDatabaseSql(out *jwriter.Writer, in sql.DB) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}
