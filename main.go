package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func copyHeader(dest *http.Header, src http.Header) {
	for k, h := range src {
		for _, v := range h {
			dest.Add(k, v)
		}
	}
}

func passThrough(wr http.ResponseWriter, req *http.Request) {
	//passthrough request
	client := &http.Client{}
	req.RequestURI = ""

	if resp, err := client.Do(req); err == nil {

		defer resp.Body.Close()

		fmt.Println(req.RemoteAddr, " ", resp.Status)

		//copy header
		header := wr.Header()
		copyHeader(&header, resp.Header)

		//set status code
		wr.WriteHeader(resp.StatusCode)

		//TODO: cache downloaded file, only oct-stream
		io.Copy(wr, resp.Body)
	} else {
		fmt.Println("passthrough error:", err)
	}
}

func proxyHandler(wr http.ResponseWriter, req *http.Request) {

	//log any request
	fmt.Println("From:", req.RemoteAddr, " ", req.Method, " ", req.URL)

	//TODO: match specified files(regex)
	//TODO：support, Header: Accept-Ranges = bytes
	//TODO: only binary stream, Header: Content-Type = application/octet-stream

	if u, err := url.Parse(req.URL.String()); err == nil {
		filename := path.Base(u.Path)

		//check if exist
		filePath := fmt.Sprintf("./%s", filename)

		fmt.Println(filePath)

		if fileInfo, err := os.Stat(filePath); err == nil {
			fmt.Printf("Local file exist {%d}. return from local file stream", fileInfo.Size())

			if file, _ := os.Open(filePath); err == nil {
				defer file.Close()
				fmt.Println("Return from local file stream")

				header := wr.Header()
				header.Add("Content-Length", fmt.Sprint(fileInfo.Size()))
				wr.WriteHeader(200)
				io.Copy(wr, file)
				return
			}
		}

	} else {
		panic(err)
	}

	passThrough(wr, req)
}

func main() {
	fmt.Println("Start Proxy...")
	http.HandleFunc("/", proxyHandler)
	http.ListenAndServe(":8000", nil)
}
