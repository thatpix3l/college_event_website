// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package gen_sql

import (
	json "encoding/json"
	fmt "fmt"
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

func easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql(in *jlexer.Lexer, out *CreateUniversityParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	var TitleSet bool
	var CoordinateSet bool
	var AboutSet bool
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
		case "Title":
			out.Title = string(in.String())
			TitleSet = true
		case "Coordinate":
			out.Coordinate = int32(in.Int32())
			CoordinateSet = true
		case "About":
			out.About = string(in.String())
			AboutSet = true
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
	if !TitleSet {
		in.AddError(fmt.Errorf("key 'Title' is required"))
	}
	if !CoordinateSet {
		in.AddError(fmt.Errorf("key 'Coordinate' is required"))
	}
	if !AboutSet {
		in.AddError(fmt.Errorf("key 'About' is required"))
	}
}
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql(out *jwriter.Writer, in CreateUniversityParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Title\":"
		out.RawString(prefix[1:])
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"Coordinate\":"
		out.RawString(prefix)
		out.Int32(int32(in.Coordinate))
	}
	{
		const prefix string = ",\"About\":"
		out.RawString(prefix)
		out.String(string(in.About))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateUniversityParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateUniversityParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateUniversityParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateUniversityParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql(l, v)
}
func easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql1(in *jlexer.Lexer, out *CreateUniversityMemberParams) {
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
		case "Column1":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Column1).UnmarshalJSON(data))
			}
		case "Column2":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Column2).UnmarshalJSON(data))
			}
		case "Column3":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Column3).UnmarshalJSON(data))
			}
		case "Column4":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Column4).UnmarshalJSON(data))
			}
		case "Column5":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Column5).UnmarshalJSON(data))
			}
		case "Column6":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Column6).UnmarshalJSON(data))
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
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql1(out *jwriter.Writer, in CreateUniversityMemberParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Column1\":"
		out.RawString(prefix[1:])
		out.Raw((in.Column1).MarshalJSON())
	}
	{
		const prefix string = ",\"Column2\":"
		out.RawString(prefix)
		out.Raw((in.Column2).MarshalJSON())
	}
	{
		const prefix string = ",\"Column3\":"
		out.RawString(prefix)
		out.Raw((in.Column3).MarshalJSON())
	}
	{
		const prefix string = ",\"Column4\":"
		out.RawString(prefix)
		out.Raw((in.Column4).MarshalJSON())
	}
	{
		const prefix string = ",\"Column5\":"
		out.RawString(prefix)
		out.Raw((in.Column5).MarshalJSON())
	}
	{
		const prefix string = ",\"Column6\":"
		out.RawString(prefix)
		out.Raw((in.Column6).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateUniversityMemberParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateUniversityMemberParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateUniversityMemberParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateUniversityMemberParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql1(l, v)
}
func easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql2(in *jlexer.Lexer, out *CreateCoordinateParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	var TitleSet bool
	var LatitudeSet bool
	var LongitudeSet bool
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
		case "Title":
			out.Title = string(in.String())
			TitleSet = true
		case "Latitude":
			out.Latitude = float64(in.Float64())
			LatitudeSet = true
		case "Longitude":
			out.Longitude = float64(in.Float64())
			LongitudeSet = true
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
	if !TitleSet {
		in.AddError(fmt.Errorf("key 'Title' is required"))
	}
	if !LatitudeSet {
		in.AddError(fmt.Errorf("key 'Latitude' is required"))
	}
	if !LongitudeSet {
		in.AddError(fmt.Errorf("key 'Longitude' is required"))
	}
}
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql2(out *jwriter.Writer, in CreateCoordinateParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Title\":"
		out.RawString(prefix[1:])
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"Latitude\":"
		out.RawString(prefix)
		out.Float64(float64(in.Latitude))
	}
	{
		const prefix string = ",\"Longitude\":"
		out.RawString(prefix)
		out.Float64(float64(in.Longitude))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateCoordinateParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateCoordinateParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateCoordinateParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateCoordinateParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql2(l, v)
}