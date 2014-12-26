package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func Upload(url string, file string, data []byte) error {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	fw, err := w.CreateFormFile("image", file)
	if err != nil {
		return err
	}
	f := bytes.NewReader(data)
	if _, err = io.Copy(fw, f); err != nil {
		return err
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return err
	}
	return nil
}

func Download(url string) ([]byte, error) {

	client := &http.Client{}

	req, _ := http.NewRequest(
		"GET",
		url,
		nil,
	)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {
	body, err := Download("http://i.imgur.com/aDORpy9.jpg")

	if err != nil {
		log.Printf("failed To Download %v", err)
		return
	}
	err = Upload(os.Args[1], "aDORpy9.jpg", body)
	if err != nil {
		log.Printf("failed To Upload %v", err)
		return
	}

}
