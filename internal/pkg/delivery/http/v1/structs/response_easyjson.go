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

func easyjson6ff3ac1dDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(in *jlexer.Lexer, out *JsonResponse) {
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
		case "status":
			out.Status = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "body":
			if m, ok := out.Body.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Body.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Body = in.Interface()
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
func easyjson6ff3ac1dEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(out *jwriter.Writer, in JsonResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix[1:])
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"body\":"
		out.RawString(prefix)
		if m, ok := in.Body.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Body.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Body))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v JsonResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v JsonResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *JsonResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *JsonResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs(l, v)
}
func easyjson6ff3ac1dDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(in *jlexer.Lexer, out *JsonErrResponse) {
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
		case "status":
			out.Status = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "code":
			out.Code = string(in.String())
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
func easyjson6ff3ac1dEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(out *jwriter.Writer, in JsonErrResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix[1:])
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"code\":"
		out.RawString(prefix)
		out.String(string(in.Code))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v JsonErrResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v JsonErrResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *JsonErrResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *JsonErrResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryHttpV1Structs1(l, v)
}