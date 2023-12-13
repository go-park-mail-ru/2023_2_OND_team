// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package structs

import (
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

func easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(in *jlexer.Lexer, out *DeletePinFromBoard) {
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
		case "pin_id":
			out.PinID = int(in.Int())
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
func easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(out *jwriter.Writer, in DeletePinFromBoard) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"pin_id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.PinID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v DeletePinFromBoard) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v DeletePinFromBoard) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *DeletePinFromBoard) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *DeletePinFromBoard) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(l, v)
}
func easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(in *jlexer.Lexer, out *CertainBoardWithUsername) {
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
		case "board_id":
			out.ID = int(in.Int())
		case "author_id":
			out.AuthorID = int(in.Int())
		case "author_username":
			out.AuthorUsername = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "created_at":
			out.CreatedAt = string(in.String())
		case "pins_number":
			out.PinsNumber = int(in.Int())
		case "pins":
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
		case "tags":
			if in.IsNull() {
				in.Skip()
				out.Tags = nil
			} else {
				in.Delim('[')
				if out.Tags == nil {
					if !in.IsDelim(']') {
						out.Tags = make([]string, 0, 4)
					} else {
						out.Tags = []string{}
					}
				} else {
					out.Tags = (out.Tags)[:0]
				}
				for !in.IsDelim(']') {
					var v2 string
					v2 = string(in.String())
					out.Tags = append(out.Tags, v2)
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
func easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(out *jwriter.Writer, in CertainBoardWithUsername) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"board_id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"author_id\":"
		out.RawString(prefix)
		out.Int(int(in.AuthorID))
	}
	{
		const prefix string = ",\"author_username\":"
		out.RawString(prefix)
		out.String(string(in.AuthorUsername))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.String(string(in.CreatedAt))
	}
	{
		const prefix string = ",\"pins_number\":"
		out.RawString(prefix)
		out.Int(int(in.PinsNumber))
	}
	{
		const prefix string = ",\"pins\":"
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
		const prefix string = ",\"tags\":"
		out.RawString(prefix)
		if in.Tags == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Tags {
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
func (v CertainBoardWithUsername) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CertainBoardWithUsername) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CertainBoardWithUsername) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CertainBoardWithUsername) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(l, v)
}
func easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs2(in *jlexer.Lexer, out *CertainBoard) {
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
		case "board_id":
			out.ID = int(in.Int())
		case "author_id":
			out.AuthorID = int(in.Int())
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "created_at":
			out.CreatedAt = string(in.String())
		case "pins_number":
			out.PinsNumber = int(in.Int())
		case "pins":
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
					var v7 string
					v7 = string(in.String())
					out.Pins = append(out.Pins, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "tags":
			if in.IsNull() {
				in.Skip()
				out.Tags = nil
			} else {
				in.Delim('[')
				if out.Tags == nil {
					if !in.IsDelim(']') {
						out.Tags = make([]string, 0, 4)
					} else {
						out.Tags = []string{}
					}
				} else {
					out.Tags = (out.Tags)[:0]
				}
				for !in.IsDelim(']') {
					var v8 string
					v8 = string(in.String())
					out.Tags = append(out.Tags, v8)
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
func easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs2(out *jwriter.Writer, in CertainBoard) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"board_id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"author_id\":"
		out.RawString(prefix)
		out.Int(int(in.AuthorID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.String(string(in.CreatedAt))
	}
	{
		const prefix string = ",\"pins_number\":"
		out.RawString(prefix)
		out.Int(int(in.PinsNumber))
	}
	{
		const prefix string = ",\"pins\":"
		out.RawString(prefix)
		if in.Pins == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v9, v10 := range in.Pins {
				if v9 > 0 {
					out.RawByte(',')
				}
				out.String(string(v10))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"tags\":"
		out.RawString(prefix)
		if in.Tags == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v11, v12 := range in.Tags {
				if v11 > 0 {
					out.RawByte(',')
				}
				out.String(string(v12))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CertainBoard) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CertainBoard) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CertainBoard) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CertainBoard) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs2(l, v)
}
func easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs3(in *jlexer.Lexer, out *BoardData) {
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
		case "title":
			if in.IsNull() {
				in.Skip()
				out.Title = nil
			} else {
				if out.Title == nil {
					out.Title = new(string)
				}
				*out.Title = string(in.String())
			}
		case "description":
			if in.IsNull() {
				in.Skip()
				out.Description = nil
			} else {
				if out.Description == nil {
					out.Description = new(string)
				}
				*out.Description = string(in.String())
			}
		case "public":
			if in.IsNull() {
				in.Skip()
				out.Public = nil
			} else {
				if out.Public == nil {
					out.Public = new(bool)
				}
				*out.Public = bool(in.Bool())
			}
		case "tags":
			if in.IsNull() {
				in.Skip()
				out.Tags = nil
			} else {
				in.Delim('[')
				if out.Tags == nil {
					if !in.IsDelim(']') {
						out.Tags = make([]string, 0, 4)
					} else {
						out.Tags = []string{}
					}
				} else {
					out.Tags = (out.Tags)[:0]
				}
				for !in.IsDelim(']') {
					var v13 string
					v13 = string(in.String())
					out.Tags = append(out.Tags, v13)
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
func easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs3(out *jwriter.Writer, in BoardData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix[1:])
		if in.Title == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Title))
		}
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		if in.Description == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Description))
		}
	}
	{
		const prefix string = ",\"public\":"
		out.RawString(prefix)
		if in.Public == nil {
			out.RawString("null")
		} else {
			out.Bool(bool(*in.Public))
		}
	}
	{
		const prefix string = ",\"tags\":"
		out.RawString(prefix)
		if in.Tags == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v14, v15 := range in.Tags {
				if v14 > 0 {
					out.RawByte(',')
				}
				out.String(string(v15))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v BoardData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BoardData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *BoardData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BoardData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs3(l, v)
}
