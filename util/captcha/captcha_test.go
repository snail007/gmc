package gcaptcha

import (
	"github.com/golang/freetype/truetype"
	gfile "github.com/snail007/gmc/util/file"
	gtest "github.com/snail007/gmc/util/testing"
	"github.com/stretchr/testify/assert"
	"image"
	"os"
	"testing"
)

func TestCaptcha_AddFont(t *testing.T) {
	f := "a.ttf"
	assert.NoError(t, gfile.Write(f, fonts["monoton"], false))
	defer os.Remove(f)
	cap := New()
	assert.NoError(t, cap.AddFont(f))
	assert.NoError(t, cap.SetFont(f, f))
	assert.Error(t, cap.AddFont("none"))
	assert.Error(t, cap.SetFont("none"))
}

func TestCaptcha_Create(t *testing.T) {
	cap := NewDefault()
	for i := 0; i < 100; i++ {
		cap.SetDisturbance(NORMAL)
		img, str := cap.Create(4, NUM)
		assert.Regexp(t, `^[0-9]{4}$`, str)
		assert.Implements(t, (*image.Image)(nil), img)
		img, str = cap.Create(6, ALL)
		assert.Regexp(t, `^[a-zA-z0-9]{6}$`, str)
		assert.Implements(t, (*image.Image)(nil), img)
		cap.SetDisturbance(MEDIUM)
		img, str = cap.Create(6, LOWER)
		assert.Regexp(t, `^[a-z]{6}$`, str)
		assert.Implements(t, (*image.Image)(nil), img)
		img, str = cap.Create(6, UPPER)
		assert.Regexp(t, `^[A-Z]{6}$`, str)
		assert.Implements(t, (*image.Image)(nil), img)
		cap.SetDisturbance(HIGH)
		img, str = cap.Create(6, CLEAR)
		assert.Regexp(t, `^[a-zA-z0-9]{6}$`, str)
		assert.Implements(t, (*image.Image)(nil), img)
	}
}

func TestCaptcha_CreateCustom(t *testing.T) {
	cap := NewDefault()
	img := cap.CreateCustom("abcd")
	assert.Implements(t, (*image.Image)(nil), img)
}

func TestCaptcha_randFont(t *testing.T) {
	if gtest.RunProcess(t, func() {
		c := NewDefault()
		for i := 0; i < 100; i++ {
			b := c.randFont()
			assert.IsType(t, &truetype.Font{}, b)
		}
	}) {
		return
	}
	_, _, err := gtest.NewProcess(t).Verbose(false).Wait()
	assert.NoError(t, err)
}

func TestCaptcha_randStr(t *testing.T) {
	b := (&Captcha{}).randStr(4, 4)
	assert.Len(t, b, 4)
}

func TestNew(t *testing.T) {
	assert.IsType(t, &Captcha{}, New())
}

func TestNewDefault(t *testing.T) {
	assert.IsType(t, &Captcha{}, NewDefault())
}
