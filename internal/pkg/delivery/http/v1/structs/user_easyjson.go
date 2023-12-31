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

func easyjson9e1087fdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(in *jlexer.Lexer, out *UserInfo) {
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
		case "username":
			out.Username = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "surname":
			out.Surname = string(in.String())
		case "about":
			out.About = string(in.String())
		case "is_subscribed":
			out.IsSubscribed = bool(in.Bool())
		case "subscribers":
			out.SubsCount = int(in.Int())
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
func easyjson9e1087fdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(out *jwriter.Writer, in UserInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"surname\":"
		out.RawString(prefix)
		out.String(string(in.Surname))
	}
	{
		const prefix string = ",\"about\":"
		out.RawString(prefix)
		out.String(string(in.About))
	}
	{
		const prefix string = ",\"is_subscribed\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsSubscribed))
	}
	{
		const prefix string = ",\"subscribers\":"
		out.RawString(prefix)
		out.Int(int(in.SubsCount))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(l, v)
}
func easyjson9e1087fdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(in *jlexer.Lexer, out *ProfileInfo) {
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
		case "username":
			out.Username = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "subscribers":
			out.SubsCount = int(in.Int())
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
func easyjson9e1087fdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(out *jwriter.Writer, in ProfileInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"subscribers\":"
		out.RawString(prefix)
		out.Int(int(in.SubsCount))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ProfileInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ProfileInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ProfileInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ProfileInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(l, v)
}
