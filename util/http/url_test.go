package ghttp

import (
	gmap "github.com/snail007/gmc/util/map"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURLBuilder_Query(t *testing.T) {
	link := NewURLBuilder().HTTP().Host("www.example.com").Path("/abc").Query(gmap.Mss{"a": "b", "c": "d"})
	assert.Equal(t, "http://www.example.com/abc?a=b&c=d", link.String())
	link = NewURLBuilder().HTTPS().Host("www.example.com").Path("/abc").Query(gmap.Mss{"a": "b", "c": "d"})
	assert.Equal(t, "https://www.example.com/abc?a=b&c=d", link.String())
	link = NewURLBuilder().Scheme("FTP").Host("www.example.com").Path("/abc").Query(gmap.Mss{"a": "b", "c": "d"})
	assert.Equal(t, "FTP://www.example.com/abc?a=b&c=d", link.String())
	link = NewURLBuilder().Host("www.example.com").Path("/abc").Query(gmap.Mss{"a": "b", "c": "d"})
	assert.Equal(t, "http://www.example.com/abc?a=b&c=d", link.String())
	links := NewURLBuilder().Path("/abc").HostsURL([]string{"www.a.com", "www.b.com"})
	assert.Len(t, links, 2)
	assert.Equal(t, links, []string{"http://www.a.com/abc", "http://www.b.com/abc"})
	links = NewURLBuilder().Host("www.a.com").PathsURL([]string{"/abc", "/def"})
	assert.Len(t, links, 2)
	assert.Equal(t, links, []string{"http://www.a.com/abc", "http://www.a.com/def"})
	links = NewURLBuilder().Host("www.a.com").QueriesURL([]gmap.Mss{{"a": "1"}, {"a": "2"}})
	assert.Len(t, links, 2)
	assert.Equal(t, links, []string{"http://www.a.com/?a=1", "http://www.a.com/?a=2"})
	links = NewURLBuilder().Host("www.a.com").Path("/-key-/-id-").
		Holders("-key-", "-id-").HolderValuesURL([]string{"a", "b"}, []string{"d", "e"})
	assert.Len(t, links, 2)
	assert.Equal(t, links, []string{"http://www.a.com/a/d", "http://www.a.com/b/e"})
}
