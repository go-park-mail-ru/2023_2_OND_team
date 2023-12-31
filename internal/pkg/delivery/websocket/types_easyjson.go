// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package websocket

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

func easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket(in *jlexer.Lexer, out *ResponseOnRequest) {
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
		case "requestID":
			out.ID = int(in.Int())
		case "type":
			out.Type = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "code":
			out.Code = string(in.String())
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
func easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket(out *jwriter.Writer, in ResponseOnRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"requestID\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	if in.Code != "" {
		const prefix string = ",\"code\":"
		out.RawString(prefix)
		out.String(string(in.Code))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	if in.Body != nil {
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
func (v ResponseOnRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResponseOnRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResponseOnRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResponseOnRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket(l, v)
}
func easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket1(in *jlexer.Lexer, out *ResponseMessage) {
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
		case "code":
			out.Code = string(in.String())
		case "messageText":
			out.MessageText = string(in.String())
		case "eventType":
			out.Type = string(in.String())
		case "message":
			(out.Message).UnmarshalEasyJSON(in)
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
func easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket1(out *jwriter.Writer, in ResponseMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix[1:])
		out.String(string(in.Status))
	}
	if in.Code != "" {
		const prefix string = ",\"code\":"
		out.RawString(prefix)
		out.String(string(in.Code))
	}
	if in.MessageText != "" {
		const prefix string = ",\"messageText\":"
		out.RawString(prefix)
		out.String(string(in.MessageText))
	}
	if in.Type != "" {
		const prefix string = ",\"eventType\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		(in.Message).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResponseMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResponseMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResponseMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResponseMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket1(l, v)
}
func easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket2(in *jlexer.Lexer, out *PublishRequest) {
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
		case "requestID":
			out.ID = int(in.Int())
		case "message":
			(out.Message).UnmarshalEasyJSON(in)
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
func easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket2(out *jwriter.Writer, in PublishRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"requestID\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		(in.Message).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PublishRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PublishRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PublishRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PublishRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket2(l, v)
}
func easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket3(in *jlexer.Lexer, out *Object) {
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
		case "eventType":
			out.Type = string(in.String())
		case "message":
			(out.Message).UnmarshalEasyJSON(in)
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
func easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket3(out *jwriter.Writer, in Object) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Type != "" {
		const prefix string = ",\"eventType\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"message\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Message).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Object) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Object) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Object) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Object) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket3(l, v)
}
func easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket4(in *jlexer.Lexer, out *MessageFromChannel) {
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
		case "type":
			out.Type = string(in.String())
		case "message":
			(out.Message).UnmarshalEasyJSON(in)
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
func easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket4(out *jwriter.Writer, in MessageFromChannel) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix[1:])
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		(in.Message).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessageFromChannel) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessageFromChannel) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6601e8cdEncodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessageFromChannel) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessageFromChannel) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6601e8cdDecodeGithubComGoParkMailRu20232ONDTeamInternalPkgDeliveryWebsocket4(l, v)
}
