// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package board

import (
	json "encoding/json"
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

func easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgEntityBoard(in *jlexer.Lexer, out *BoardWithContent) {
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
		case "BoardInfo":
			(out.BoardInfo).UnmarshalEasyJSON(in)
		case "PinsNumber":
			out.PinsNumber = int(in.Int())
		case "Pins":
			if in.IsNull() {
				in.Skip()
				out.Pins = nil
			} else {
				in.Delim('[')
				if out.Pins == nil {
					if !in.IsDelim(']') {
						out.Pins = make([]string, 0, 4)
					} else {
						out.Pins = []string{}
					}
				} else {
					out.Pins = (out.Pins)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Pins = append(out.Pins, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "TagTitles":
			if in.IsNull() {
				in.Skip()
				out.TagTitles = nil
			} else {
				in.Delim('[')
				if out.TagTitles == nil {
					if !in.IsDelim(']') {
						out.TagTitles = make([]string, 0, 4)
					} else {
						out.TagTitles = []string{}
					}
				} else {
					out.TagTitles = (out.TagTitles)[:0]
				}
				for !in.IsDelim(']') {
					var v2 string
					v2 = string(in.String())
					out.TagTitles = append(out.TagTitles, v2)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgEntityBoard(out *jwriter.Writer, in BoardWithContent) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"BoardInfo\":"
		out.RawString(prefix[1:])
		(in.BoardInfo).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"PinsNumber\":"
		out.RawString(prefix)
		out.Int(int(in.PinsNumber))
	}
	{
		const prefix string = ",\"Pins\":"
		out.RawString(prefix)
		if in.Pins == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.Pins {
				if v3 > 0 {
					out.RawByte(',')
				}
				out.String(string(v4))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"TagTitles\":"
		out.RawString(prefix)
		if in.TagTitles == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.TagTitles {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v BoardWithContent) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgEntityBoard(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BoardWithContent) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgEntityBoard(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *BoardWithContent) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgEntityBoard(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BoardWithContent) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgEntityBoard(l, v)
}
func easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgEntityBoard1(in *jlexer.Lexer, out *Board) {
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
			out.ID = int(in.Int())
		case "author_id":
			out.AuthorID = int(in.Int())
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "public":
			out.Public = bool(in.Bool())
		case "created_at":
			if in.IsNull() {
				in.Skip()
				out.CreatedAt = nil
			} else {
				if out.CreatedAt == nil {
					out.CreatedAt = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.CreatedAt).UnmarshalJSON(data))
				}
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
func easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgEntityBoard1(out *jwriter.Writer, in Board) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	if in.AuthorID != 0 {
		const prefix string = ",\"author_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.AuthorID))
	}
	{
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"public\":"
		out.RawString(prefix)
		out.Bool(bool(in.Public))
	}
	if in.CreatedAt != nil {
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((*in.CreatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Board) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgEntityBoard1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Board) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgEntityBoard1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Board) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgEntityBoard1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Board) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgEntityBoard1(l, v)
}
