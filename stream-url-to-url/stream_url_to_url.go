package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)


type ImageUpload struct {
	Body io.ReadCloser
	ContentLength string
	ContentType string
}
func (iu *ImageUpload) Close() {_ = iu.Body.Close()}
func (iu *ImageUpload) ContentLengthAsInt64() int64 {
	c, err := strconv.ParseInt(iu.ContentLength, 10, 64)
	if err != nil {
		return -1
	}
	return c
}

func GetImageWithMetaData(client *http.Client, url string) (*ImageUpload, error) {
	headRsp, err := client.Head(url)
	if err != nil {
		return nil, err
	}
	imageRsp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return &ImageUpload{
		Body: imageRsp.Body,
		ContentType: headRsp.Header.Get("content-type"),
		ContentLength: headRsp.Header.Get("content-length"),
	}, nil
}

func PostImageUpload(client *http.Client, uploadUrl string, payload *ImageUpload) (*http.Response, error) {
	// Parse the url
	u, err := url.Parse(uploadUrl)
	if err != nil {
		panic(err)
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL: u,
		Body: payload.Body,
		ContentLength: payload.ContentLengthAsInt64(),
		Header: http.Header{"content-type": []string{payload.ContentType}},
	}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func CopyImageUrl(client *http.Client, fromUrl, toUrl string) error {
	imageData, err := GetImageWithMetaData(client, fromUrl)
	if err != nil {
		log.Fatalf("Error during GetImageWithMetaData: %v", err.Error())
		return err
	}
	defer imageData.Close()

	postResponse, err := PostImageUpload(client, toUrl, imageData)
	if err != nil {
		log.Fatalf("Error during PostImageUpload: %v", err.Error())
		return err
	}
	_ = postResponse.Body.Close()
	return nil
}