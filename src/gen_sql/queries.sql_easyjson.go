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
	var UniversityTitleSet bool
	var UniversityAboutSet bool
	var CoordTitleSet bool
	var CoordLatitudeSet bool
	var CoordLongitudeSet bool
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
		case "UniversityTitle":
			out.UniversityTitle = string(in.String())
			UniversityTitleSet = true
		case "UniversityAbout":
			out.UniversityAbout = string(in.String())
			UniversityAboutSet = true
		case "CoordTitle":
			out.CoordTitle = string(in.String())
			CoordTitleSet = true
		case "CoordLatitude":
			out.CoordLatitude = float64(in.Float64())
			CoordLatitudeSet = true
		case "CoordLongitude":
			out.CoordLongitude = float64(in.Float64())
			CoordLongitudeSet = true
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
	if !UniversityTitleSet {
		in.AddError(fmt.Errorf("key 'UniversityTitle' is required"))
	}
	if !UniversityAboutSet {
		in.AddError(fmt.Errorf("key 'UniversityAbout' is required"))
	}
	if !CoordTitleSet {
		in.AddError(fmt.Errorf("key 'CoordTitle' is required"))
	}
	if !CoordLatitudeSet {
		in.AddError(fmt.Errorf("key 'CoordLatitude' is required"))
	}
	if !CoordLongitudeSet {
		in.AddError(fmt.Errorf("key 'CoordLongitude' is required"))
	}
}
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql(out *jwriter.Writer, in CreateUniversityParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"UniversityTitle\":"
		out.RawString(prefix[1:])
		out.String(string(in.UniversityTitle))
	}
	{
		const prefix string = ",\"UniversityAbout\":"
		out.RawString(prefix)
		out.String(string(in.UniversityAbout))
	}
	{
		const prefix string = ",\"CoordTitle\":"
		out.RawString(prefix)
		out.String(string(in.CoordTitle))
	}
	{
		const prefix string = ",\"CoordLatitude\":"
		out.RawString(prefix)
		out.Float64(float64(in.CoordLatitude))
	}
	{
		const prefix string = ",\"CoordLongitude\":"
		out.RawString(prefix)
		out.Float64(float64(in.CoordLongitude))
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
func easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql1(in *jlexer.Lexer, out *CreateTaggedRsoParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	var TagSet bool
	var RsoSet bool
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
		case "Tag":
			out.Tag = int32(in.Int32())
			TagSet = true
		case "Rso":
			out.Rso = int32(in.Int32())
			RsoSet = true
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
	if !TagSet {
		in.AddError(fmt.Errorf("key 'Tag' is required"))
	}
	if !RsoSet {
		in.AddError(fmt.Errorf("key 'Rso' is required"))
	}
}
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql1(out *jwriter.Writer, in CreateTaggedRsoParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Tag\":"
		out.RawString(prefix[1:])
		out.Int32(int32(in.Tag))
	}
	{
		const prefix string = ",\"Rso\":"
		out.RawString(prefix)
		out.Int32(int32(in.Rso))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateTaggedRsoParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateTaggedRsoParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateTaggedRsoParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateTaggedRsoParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql1(l, v)
}
func easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql2(in *jlexer.Lexer, out *CreateTaggedEventParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	var TagSet bool
	var BaseEventSet bool
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
		case "Tag":
			out.Tag = int32(in.Int32())
			TagSet = true
		case "BaseEvent":
			out.BaseEvent = int32(in.Int32())
			BaseEventSet = true
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
	if !TagSet {
		in.AddError(fmt.Errorf("key 'Tag' is required"))
	}
	if !BaseEventSet {
		in.AddError(fmt.Errorf("key 'BaseEvent' is required"))
	}
}
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql2(out *jwriter.Writer, in CreateTaggedEventParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Tag\":"
		out.RawString(prefix[1:])
		out.Int32(int32(in.Tag))
	}
	{
		const prefix string = ",\"BaseEvent\":"
		out.RawString(prefix)
		out.Int32(int32(in.BaseEvent))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateTaggedEventParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateTaggedEventParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateTaggedEventParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateTaggedEventParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql2(l, v)
}
func easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql3(in *jlexer.Lexer, out *CreateRsoParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	var TitleSet bool
	var UniversitySet bool
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
		case "University":
			out.University = int32(in.Int32())
			UniversitySet = true
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
	if !UniversitySet {
		in.AddError(fmt.Errorf("key 'University' is required"))
	}
}
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql3(out *jwriter.Writer, in CreateRsoParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Title\":"
		out.RawString(prefix[1:])
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"University\":"
		out.RawString(prefix)
		out.Int32(int32(in.University))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateRsoParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateRsoParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateRsoParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateRsoParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql3(l, v)
}
func easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql4(in *jlexer.Lexer, out *CreateRsoEventParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	var IDSet bool
	var RsoSet bool
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
		case "ID":
			out.ID = int32(in.Int32())
			IDSet = true
		case "Rso":
			out.Rso = int32(in.Int32())
			RsoSet = true
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
	if !IDSet {
		in.AddError(fmt.Errorf("key 'ID' is required"))
	}
	if !RsoSet {
		in.AddError(fmt.Errorf("key 'Rso' is required"))
	}
}
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql4(out *jwriter.Writer, in CreateRsoEventParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ID\":"
		out.RawString(prefix[1:])
		out.Int32(int32(in.ID))
	}
	{
		const prefix string = ",\"Rso\":"
		out.RawString(prefix)
		out.Int32(int32(in.Rso))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateRsoEventParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateRsoEventParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateRsoEventParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateRsoEventParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql4(l, v)
}
func easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql5(in *jlexer.Lexer, out *CreateMemberParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	var IDSet bool
	var UniversitySet bool
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
		case "ID":
			out.ID = int32(in.Int32())
			IDSet = true
		case "University":
			out.University = int32(in.Int32())
			UniversitySet = true
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
	if !IDSet {
		in.AddError(fmt.Errorf("key 'ID' is required"))
	}
	if !UniversitySet {
		in.AddError(fmt.Errorf("key 'University' is required"))
	}
}
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql5(out *jwriter.Writer, in CreateMemberParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ID\":"
		out.RawString(prefix[1:])
		out.Int32(int32(in.ID))
	}
	{
		const prefix string = ",\"University\":"
		out.RawString(prefix)
		out.Int32(int32(in.University))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateMemberParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateMemberParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateMemberParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateMemberParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql5(l, v)
}
func easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql6(in *jlexer.Lexer, out *CreateCoordinateParams) {
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
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql6(out *jwriter.Writer, in CreateCoordinateParams) {
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
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateCoordinateParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateCoordinateParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateCoordinateParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql6(l, v)
}
func easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql7(in *jlexer.Lexer, out *CreateCommentParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	var BodySet bool
	var PostedBySet bool
	var BaseEventSet bool
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
		case "Body":
			out.Body = string(in.String())
			BodySet = true
		case "PostedBy":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.PostedBy).UnmarshalJSON(data))
			}
			PostedBySet = true
		case "BaseEvent":
			out.BaseEvent = int32(in.Int32())
			BaseEventSet = true
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
	if !BodySet {
		in.AddError(fmt.Errorf("key 'Body' is required"))
	}
	if !PostedBySet {
		in.AddError(fmt.Errorf("key 'PostedBy' is required"))
	}
	if !BaseEventSet {
		in.AddError(fmt.Errorf("key 'BaseEvent' is required"))
	}
}
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql7(out *jwriter.Writer, in CreateCommentParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Body\":"
		out.RawString(prefix[1:])
		out.String(string(in.Body))
	}
	{
		const prefix string = ",\"PostedBy\":"
		out.RawString(prefix)
		out.Raw((in.PostedBy).MarshalJSON())
	}
	{
		const prefix string = ",\"BaseEvent\":"
		out.RawString(prefix)
		out.Int32(int32(in.BaseEvent))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateCommentParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateCommentParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateCommentParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateCommentParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql7(l, v)
}
func easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql8(in *jlexer.Lexer, out *CreateBaseUserParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	var NameFirstSet bool
	var NameMiddleSet bool
	var NameLastSet bool
	var EmailSet bool
	var PasswordHashSet bool
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
		case "NameFirst":
			out.NameFirst = string(in.String())
			NameFirstSet = true
		case "NameMiddle":
			out.NameMiddle = string(in.String())
			NameMiddleSet = true
		case "NameLast":
			out.NameLast = string(in.String())
			NameLastSet = true
		case "Email":
			out.Email = string(in.String())
			EmailSet = true
		case "PasswordHash":
			out.PasswordHash = string(in.String())
			PasswordHashSet = true
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
	if !NameFirstSet {
		in.AddError(fmt.Errorf("key 'NameFirst' is required"))
	}
	if !NameMiddleSet {
		in.AddError(fmt.Errorf("key 'NameMiddle' is required"))
	}
	if !NameLastSet {
		in.AddError(fmt.Errorf("key 'NameLast' is required"))
	}
	if !EmailSet {
		in.AddError(fmt.Errorf("key 'Email' is required"))
	}
	if !PasswordHashSet {
		in.AddError(fmt.Errorf("key 'PasswordHash' is required"))
	}
}
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql8(out *jwriter.Writer, in CreateBaseUserParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"NameFirst\":"
		out.RawString(prefix[1:])
		out.String(string(in.NameFirst))
	}
	{
		const prefix string = ",\"NameMiddle\":"
		out.RawString(prefix)
		out.String(string(in.NameMiddle))
	}
	{
		const prefix string = ",\"NameLast\":"
		out.RawString(prefix)
		out.String(string(in.NameLast))
	}
	{
		const prefix string = ",\"Email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"PasswordHash\":"
		out.RawString(prefix)
		out.String(string(in.PasswordHash))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateBaseUserParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateBaseUserParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateBaseUserParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateBaseUserParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql8(l, v)
}
func easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql9(in *jlexer.Lexer, out *CreateBaseEventParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	var TitleSet bool
	var BodySet bool
	var UniversitySet bool
	var OccurrenceTimeSet bool
	var OccurrenceLocationSet bool
	var ContactPhoneSet bool
	var ContactEmailSet bool
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
		case "Body":
			out.Body = string(in.String())
			BodySet = true
		case "University":
			out.University = int32(in.Int32())
			UniversitySet = true
		case "OccurrenceTime":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.OccurrenceTime).UnmarshalJSON(data))
			}
			OccurrenceTimeSet = true
		case "OccurrenceLocation":
			out.OccurrenceLocation = int32(in.Int32())
			OccurrenceLocationSet = true
		case "ContactPhone":
			out.ContactPhone = string(in.String())
			ContactPhoneSet = true
		case "ContactEmail":
			out.ContactEmail = string(in.String())
			ContactEmailSet = true
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
	if !BodySet {
		in.AddError(fmt.Errorf("key 'Body' is required"))
	}
	if !UniversitySet {
		in.AddError(fmt.Errorf("key 'University' is required"))
	}
	if !OccurrenceTimeSet {
		in.AddError(fmt.Errorf("key 'OccurrenceTime' is required"))
	}
	if !OccurrenceLocationSet {
		in.AddError(fmt.Errorf("key 'OccurrenceLocation' is required"))
	}
	if !ContactPhoneSet {
		in.AddError(fmt.Errorf("key 'ContactPhone' is required"))
	}
	if !ContactEmailSet {
		in.AddError(fmt.Errorf("key 'ContactEmail' is required"))
	}
}
func easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql9(out *jwriter.Writer, in CreateBaseEventParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Title\":"
		out.RawString(prefix[1:])
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"Body\":"
		out.RawString(prefix)
		out.String(string(in.Body))
	}
	{
		const prefix string = ",\"University\":"
		out.RawString(prefix)
		out.Int32(int32(in.University))
	}
	{
		const prefix string = ",\"OccurrenceTime\":"
		out.RawString(prefix)
		out.Raw((in.OccurrenceTime).MarshalJSON())
	}
	{
		const prefix string = ",\"OccurrenceLocation\":"
		out.RawString(prefix)
		out.Int32(int32(in.OccurrenceLocation))
	}
	{
		const prefix string = ",\"ContactPhone\":"
		out.RawString(prefix)
		out.String(string(in.ContactPhone))
	}
	{
		const prefix string = ",\"ContactEmail\":"
		out.RawString(prefix)
		out.String(string(in.ContactEmail))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateBaseEventParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateBaseEventParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba41af8EncodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateBaseEventParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateBaseEventParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba41af8DecodeGithubComThatpix3lCollgeEventWebsiteSrcGenSql9(l, v)
}
