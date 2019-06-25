package main

import (
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCopyImageUrl(t *testing.T) {

	client := http.DefaultClient

	// Get Image returns a 100 x 100 image with 2 colored pixels
	// At 1,1 and 2,2
	getImageFunc := func (w http.ResponseWriter, r *http.Request) {
		i := image.NewRGBA(image.Rect(0, 0, 100, 100))
		i.Set(1, 1, color.White)
		i.Set(2, 2, color.Black)
		_ = png.Encode(w, i)
	}

	// Image Source Server
	getImageTS := httptest.NewServer(http.HandlerFunc(getImageFunc))
	defer getImageTS.Close()

	// Image Destination Server function (handle a POST to update a closure
	var postImageContents image.Image = nil
	postImageFunc := func (w http.ResponseWriter, r *http.Request) {
		postImageContents, _ = png.Decode(r.Body)
	}

	// Image Destination Server
	postImageTS := httptest.NewServer(http.HandlerFunc(postImageFunc))
	defer postImageTS.Close()

	// Look, nothing up our sleeves.
	if postImageContents != nil {
		t.Fatal("ImageContents should be nil")
	}

	// Make the Call.
	err := CopyImageUrl(client, getImageTS.URL, postImageTS.URL)

	if err != nil {
		t.Fatal(err.Error())
	}

	if !colorsAreSame(postImageContents.At(1, 1), color.White) {
		t.Error("Post Image 1, 1 should be white!")
	}

	if !colorsAreSame(postImageContents.At(2, 2), color.Black) {
		t.Error("Post Image 2, 2 should be white!")
	}
}

// Helper Functions


func colorsAreSame(col1, col2 color.Color) bool {
	r1, g1, b1, a1 := col1.RGBA()
	r2, g2, b2, a2 := col2.RGBA()

	return r1 == r2 &&
		b1 == b2 &&
		g1 == g2 &&
		a1 == a2
}