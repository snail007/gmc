package gurl

import (
	gmap "github.com/snail007/gmc/util/map"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURLBuilder_Query(t *testing.T) {
	link := NewBuilder().HTTP().Host("www.example.com").Path("/abc").Query(gmap.Mss{"a": "b", "c": "d"})
	assert.Equal(t, "http://www.example.com/abc?a=b&c=d", link.String())
	link = NewBuilder().HTTPS().Host("www.example.com").Path("/abc").Query(gmap.Mss{"a": "b", "c": "d"})
	assert.Equal(t, "https://www.example.com/abc?a=b&c=d", link.String())
	link = NewBuilder().Scheme("FTP").Host("www.example.com").Path("/abc").Query(gmap.Mss{"a": "b", "c": "d"})
	assert.Equal(t, "FTP://www.example.com/abc?a=b&c=d", link.String())
	link = NewBuilder().Host("www.example.com").Path("/abc").Query(gmap.Mss{"a": "b", "c": "d"})
	assert.Equal(t, "http://www.example.com/abc?a=b&c=d", link.String())
	links := NewBuilder().Path("/abc").HostsURL([]string{"www.a.com", "www.b.com"})
	assert.Len(t, links, 2)
	assert.Equal(t, links, []string{"http://www.a.com/abc", "http://www.b.com/abc"})
	links = NewBuilder().Host("www.a.com").PathsURL([]string{"/abc", "/def"})
	assert.Len(t, links, 2)
	assert.Equal(t, links, []string{"http://www.a.com/abc", "http://www.a.com/def"})
	links = NewBuilder().Host("www.a.com").QueriesURL([]gmap.Mss{{"a": "1"}, {"a": "2"}})
	assert.Len(t, links, 2)
	assert.Equal(t, links, []string{"http://www.a.com/?a=1", "http://www.a.com/?a=2"})
	links = NewBuilder().Host("www.a.com").Path("/-key-/-id-").
		Holders("-key-", "-id-").HolderValuesURL([]string{"a", "b"}, []string{"d", "e"})
	assert.Len(t, links, 2)
	assert.Equal(t, links, []string{"http://www.a.com/a/d", "http://www.a.com/b/e"})
	link = NewBuilder().Host("a.com:81").Port("88")
	assert.Equal(t, "http://a.com:88/", link.String())
	link = NewBuilder().Host("a.com").Port("88")
	assert.Equal(t, "http://a.com:88/", link.String())
	link = NewBuilder().Host("a.com").Port("88").SetQuery("a", "1")
	assert.Equal(t, "http://a.com:88/?a=1", link.String())
	link = NewBuilder().Parse("http://a.com:88/?a=1")
	assert.Equal(t, "http://a.com:88/?a=1", link.String())
	assert.Nil(t, link.Error())
	link = NewBuilder().Parse("http://a.com:abc/?a=1")
	assert.NotNil(t, link.Error())
}

func TestGetConcatChar(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "contains ?", args: args{URL: "a.com/?a=c"}, want: "&"},
		{name: "not contains ?", args: args{URL: "a.com"}, want: "?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetConcatChar(tt.args.URL), "GetConcatChar(%v)", tt.args.URL)
		})
	}
}
