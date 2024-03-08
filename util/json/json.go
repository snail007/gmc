package gjson

import (
	"encoding/json"
	"errors"
	"fmt"
	gcore "github.com/snail007/gmc/core"
	gcast "github.com/snail007/gmc/util/cast"
	gvalue "github.com/snail007/gmc/util/value"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"strconv"
	"strings"
)

var (
	// AddModifier binds a custom modifier command to the GJSON syntax.
	// This operation is not thread safe and should be executed prior to
	// using all other gjson function.
	AddModifier = gjson.AddModifier

	// ModifierExists returns true when the specified modifier exists.
	ModifierExists = gjson.ModifierExists

	// Escape returns an escaped path component.
	//
	//	json := `{
	//	  "user":{
	//	     "first.name": "Janet",
	//	     "last.name": "Prichard"
	//	   }
	//	}`
	//	user := gjson.Get(json, "user")
	//	println(user.Get(gjson.Escape("first.name"))
	//	println(user.Get(gjson.Escape("last.name"))
	//	// Output:
	//	// Janet
	//	// Prichard
	Escape = gjson.Escape

	// ForEachLine iterates through lines of JSON as specified by the JSON Lines
	// format (http://jsonlines.org/).
	// Each line is returned as a GJSON Result.
	ForEachLine = gjson.ForEachLine

	// Parse parses the json and returns a result.
	//
	// This function expects that the json is well-formed, and does not validate.
	// Invalid json will not panic, but it may return back unexpected results.
	// If you are consuming JSON from an unpredictable source then you may want to
	// use the Valid function first.
	Parse = func(json string) Result {
		r := gjson.Parse(json)
		return Result{
			Result: r,
			path:   r.Path(json),
			paths:  r.Paths(json),
		}
	}
	// ParseBytes parses the json and returns a result.
	// If working with bytes, this method preferred over Parse(string(data))
	ParseBytes = gjson.ParseBytes

	// Valid returns true if the input is valid json.
	//
	//	if !gjson.Valid(json) {
	//		return errors.New("invalid json")
	//	}
	//	value := gjson.Get(json, "name.last")
	Valid = gjson.Valid

	// ValidBytes returns true if the input is valid json.
	//
	//	if !gjson.Valid(json) {
	//		return errors.New("invalid json")
	//	}
	//	value := gjson.Get(json, "name.last")
	//
	// If working with bytes, this method preferred over ValidBytes(string(data))
	ValidBytes = gjson.ValidBytes
)

type Options = sjson.Options

type Result struct {
	gjson.Result
	path  string
	paths []string
}

func (s Result) Path() string {
	return s.path
}

func (s Result) Paths() []string {
	return s.paths
}

func (s Result) ToJSONObject() *JSONObject {
	return NewJSONObject(s.Raw)
}

func (s Result) ToJSONArray() *JSONArray {
	return NewJSONArray(s.Raw)
}

func (s Result) ForEach(f func(k, v Result) bool){
	s.Result.ForEach(func(key, value gjson.Result) bool {
		return f(Result{
			Result:key,
		},Result{
			Result:value,
		})
	}) 
}

type Builder struct {
	json string
}

func NewBuilderE(v interface{}) (*Builder, error) {
	str, err := getJSONStr(v, "")
	if err != nil {
		return nil, err
	}
	return &Builder{json: str}, nil
}

func NewBuilder(v interface{}) *Builder {
	obj, _ := NewBuilderE(v)
	return obj
}

// Delete deletes a value from json for the specified path.
// path syntax: https://github.com/tidwall/sjson?tab=readme-ov-file#path-syntax
func (s *Builder) Delete(path string) error {
	j, err := sjson.Delete(s.json, path)
	if err == nil {
		s.json = j
	}
	return err
}

// Set sets a json value for the specified path.
// A path is in dot syntax, such as "name.last" or "age".
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return back unexpected results.
// An error is returned if the path is not valid.
//
// A path is a series of keys separated by a dot.
//
//	{
//	  "name": {"first": "Tom", "last": "Anderson"},
//	  "age":37,
//	  "children": ["Sara","Alex","Jack"],
//	  "friends": [
//	    {"first": "James", "last": "Murphy"},
//	    {"first": "Roger", "last": "Craig"}
//	  ]
//	}
//	"name.last"          >> "Anderson"
//	"age"                >> 37
//	"children.1"         >> "Alex"
//
// path syntax: https://github.com/tidwall/sjson?tab=readme-ov-file#path-syntax
func (s *Builder) Set(path string, value interface{}) error {
	j, err := sjson.Set(s.json, path, value)
	if err == nil {
		s.json = j
	}
	return err
}

// SetRaw sets a raw json value for the specified path.
// This function works the same as Set except that the value is set as a
// raw block of json. This allows for setting premarshalled json objects.
func (s *Builder) SetRaw(path, value string) error {
	j, err := sjson.SetRaw(s.json, path, value)
	if err == nil {
		if !Valid(j) {
			return errors.New("invalid json value: " + value)
		}
		s.json = j
	}
	return err
}

// Get searches json for the specified path.
// A path is in dot syntax, such as "name.last" or "age".
// When the value is found it's returned immediately.
//
// A path is a series of keys separated by a dot.
// A key may contain special wildcard characters '*' and '?'.
// To access an array value use the index as the key.
// To get the number of elements in an array or to access a child path, use
// the '#' character.
// The dot and wildcard character can be escaped with '\'.
//
//	{
//	  "name": {"first": "Tom", "last": "Anderson"},
//	  "age":37,
//	  "children": ["Sara","Alex","Jack"],
//	  "friends": [
//	    {"first": "James", "last": "Murphy"},
//	    {"first": "Roger", "last": "Craig"}
//	  ]
//	}
//	"name.last"          >> "Anderson"
//	"age"                >> 37
//	"children"           >> ["Sara","Alex","Jack"]
//	"children.#"         >> 3
//	"children.1"         >> "Alex"
//	"child*.2"           >> "Jack"
//	"c?ildren.0"         >> "Sara"
//	"friends.#.first"    >> ["James","Roger"]
//
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return back unexpected results.
// If you are consuming JSON from an unpredictable source then you may want to
// use the Valid function first.
// path syntax: https://github.com/tidwall/gjson/blob/master/SYNTAX.md
func (s *Builder) Get(path string) Result {
	r := gjson.Get(s.json, path)
	return Result{
		paths:  r.Paths(s.json),
		path:   r.Path(s.json),
		Result: r,
	}
}

// String convert the *Builder to JSON string,
func (s *Builder) String() string {
	return s.json
}

// Interface convert the *Builder to Go DATA,
func (s *Builder) Interface() (v interface{}) {
	json.Unmarshal([]byte(s.json), &v)
	return
}

// JSONObject convert the *Builder to *JSONObject,
// if the *Builder is not a json object, nil returned.
func (s *Builder) JSONObject() *JSONObject {
	if s.json != "" && !strings.HasPrefix(s.json, "{") {
		return nil
	}
	return NewJSONObject(s.json)
}

// JSONArray convert the *Builder to *JSONArray,
// if the *Builder is not a json array, nil returned.
func (s *Builder) JSONArray() *JSONArray {
	if s.json != "" && !strings.HasPrefix(s.json, "[") {
		return nil
	}
	return NewJSONArray(s.json)
}

// GetMany searches json for the multiple paths.
// The return value is a Result array where the number of items
// will be equal to the number of input paths.
func (s *Builder) GetMany(path ...string) []Result {
	rs1 := gjson.GetMany(s.json, path...)
	var rs []Result
	for _, r := range rs1 {
		rs = append(rs,
			Result{
				Result: r,
				path:   r.Path(s.json),
				paths:  r.Paths(s.json),
			})
	}
	return rs
}

type JSONObject struct {
	*Builder
}

// NewJSONObjectE create a *JSONObject form v, returned error(if have)
// v can be json object content of []byte and string, or any data which json.Marshal can be processed.
func NewJSONObjectE(v interface{}) (*JSONObject, error) {
	str, err := getJSONStr(v, "{}")
	if err != nil {
		return nil, err
	}
	if str != "" && !strings.HasPrefix(str, "{") {
		return nil, errors.New("fail to convert v to json array")
	}
	return &JSONObject{
		Builder: NewBuilder(str),
	}, nil
}

// NewJSONObject create a *JSONObject form v, if error occurred nil returned.
// v can be json object content of []byte and string, or any data which json.Marshal can be processed.
func NewJSONObject(v interface{}) *JSONObject {
	obj, _ := NewJSONObjectE(v)
	return obj
}

type JSONArray struct {
	*Builder
}

// NewJSONArrayE create a *JSONArray form v, returned error(if have)
// v can be json array content of []byte and string, or any data which json.Marshal can be processed.
func NewJSONArrayE(v interface{}) (*JSONArray, error) {
	str, err := getJSONStr(v, "[]")
	if err != nil {
		return nil, err
	}
	if str != "" && !strings.HasPrefix(str, "[") {
		return nil, errors.New("fail to convert v to json array")
	}
	return &JSONArray{
		Builder: NewBuilder(str),
	}, nil
}

// NewJSONArray create a *JSONArray form v, if error occurred nil returned.
// v can be json array content of []byte and string, or any data which json.Marshal can be processed.
func NewJSONArray(v interface{}) *JSONArray {
	obj, _ := NewJSONArrayE(v)
	return obj
}

// Merge *JSONArray, JSONArray or any valid slice to s
func (s *JSONArray) Merge(arr interface{}) (err error) {
	var merge = func(a *JSONArray) {
		a.Get("@this").ForEach(func(key, value Result) bool {
			err = s.Append(value.Value())
			return err == nil
		})
	}
	switch a := arr.(type) {
	case *JSONArray:
		merge(a)
	case JSONArray:
		merge(&a)
	default:
		err = s.Append(gvalue.NewAny(arr).Slice()...)
	}
	return
}

func (s *JSONArray) Append(values ...interface{}) (err error) {
	for _, value := range values {
		switch v := value.(type) {
		case *JSONObject:
			err = s.SetRaw("-1", v.json)
		case JSONObject:
			err = s.SetRaw("-1", v.json)
		case *JSONArray:
			err = s.SetRaw("-1", v.json)
		case JSONArray:
			err = s.SetRaw("-1", v.json)
		default:
			err = s.Set("-1", value)
		}
		if err != nil {
			return
		}
	}
	return nil
}

func (s *JSONArray) Len() int64 {
	return s.Get("#").Int()
}

func (s *JSONArray) Last() Result {
	idx := s.Get("#").Int() - 1
	if idx < 0 {
		idx = 0
	}
	return s.Get(strconv.FormatInt(idx, 10))
}

func (s *JSONArray) First() Result {
	return s.Get("0")
}

type JSONResult struct {
	data map[string]interface{}
	ctx  gcore.Ctx
}

// NewResultCtx Optional args: code int, message string, data interface{}
func NewResultCtx(ctx gcore.Ctx, d ...interface{}) *JSONResult {
	r := NewResult(d...)
	r.ctx = ctx
	return r
}

// NewResult Optional args: code int, message string, data interface{}
func NewResult(d ...interface{}) *JSONResult {
	if len(d) == 1 {
		var b []byte
		switch d[0].(type) {
		case string:
			b = []byte(d[0].(string))
		case []byte:
			b = d[0].([]byte)
		}
		if len(b) > 0 {
			s := &JSONResult{
				data: map[string]interface{}{},
			}
			if e := json.Unmarshal(b, &s.data); e != nil {
				return nil
			}
			return s
		}
	}
	code := 0
	message := ""
	var data interface{}
	if len(d) >= 1 {
		code = gcast.ToInt(d[0])
	}
	if len(d) >= 2 {
		message = gcast.ToString(d[1])
	}
	if len(d) >= 3 {
		data = d[2]
	}
	return &JSONResult{
		data: map[string]interface{}{
			"code":    code,
			"message": message,
			"data":    data,
		},
	}
}

func (s *JSONResult) Set(key string, value interface{}) *JSONResult {
	s.data[key] = value
	return s
}

func (s *JSONResult) ToJSON() []byte {
	j, _ := json.Marshal(s.data)
	return j
}

func (s *JSONResult) SetCode(code int) *JSONResult {
	return s.Set("code", code)
}

func (s *JSONResult) SetMessage(format string, msg ...interface{}) *JSONResult {
	return s.Set("message", fmt.Sprintf(format, msg...))
}

func (s *JSONResult) SetData(d interface{}) *JSONResult {
	return s.Set("data", d)
}

func (s *JSONResult) Code() int {
	return gcast.ToInt(s.data["code"])
}

func (s *JSONResult) Message() string {
	return gcast.ToString(s.data["message"])
}

func (s *JSONResult) Data() interface{} {
	return s.data["data"]
}

func (s *JSONResult) DataMap() interface{} {
	return s.data
}

func (s *JSONResult) WriteTo(dst io.Writer) (err error) {
	_, err = dst.Write(s.ToJSON())
	return
}

func (s *JSONResult) WriteToCtx(ctx gcore.Ctx) (err error) {
	_, err = ctx.Response().Write(s.ToJSON())
	return
}

// Success only worked with NewResultCtx()
func (s *JSONResult) Success(d ...interface{}) (err error) {
	var data interface{}
	if len(d) == 1 {
		data = d[0]
	}
	return s.SetData(data).WriteToCtx(s.ctx)
}

// Fail only worked with NewResultCtx()
func (s *JSONResult) Fail(format string, v ...interface{}) (err error) {
	return s.SetCode(1).SetMessage(format, v...).WriteToCtx(s.ctx)
}

func getJSONStr(v interface{}, nilValue string) (string, error) {
	if gvalue.IsNil(v) {
		return nilValue, nil
	}
	var str string
	switch val := v.(type) {
	case string:
		str = val
	case []byte:
		str = string(val)
	default:
		b, _ := json.Marshal(v)
		str = string(b)
	}
	if !Valid(str) {
		return "", errors.New("fail to convert to invalid json")
	}
	return strings.Trim(str, " \r\n\t"), nil
}
