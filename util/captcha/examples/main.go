package main

import (
	gcaptcha "github.com/snail007/gmc/util/captcha"
	"image/color"
	"image/png"
	"net/http"
)

var cap *gcaptcha.Captcha

func main() {

	cap = gcaptcha.NewDefault()
	cap.SetSize(128, 64)
	cap.SetDisturbance(gcaptcha.MEDIUM)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})

	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
		img, str := cap.Create(6, gcaptcha.ALL)
		png.Encode(w, img)
		println(str)
	})

	http.HandleFunc("/c", func(w http.ResponseWriter, r *http.Request) {
		str := r.URL.RawQuery
		img := cap.CreateCustom(str)
		png.Encode(w, img)
	})

	http.ListenAndServe(":8085", nil)

}
