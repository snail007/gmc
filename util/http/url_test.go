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
}
