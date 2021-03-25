package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func DefaultHandler(wr http.ResponseWriter, req *http.Request) {

	fmt.Println(req.RemoteAddr, " ", req.Method, " ", req.URL)

	if u, err := url.Parse(req.URL.String()); err == nil {
		filename := path.Base(u.Path)

		fmt.Println("文件名:", filename)

		//查找本地文件名是否存在
		filePath := fmt.Sprintf("./%s", filename)

		fmt.Println(filePath)

		if file, err := os.Open(filePath); err == nil {
			defer file.Close()
			fmt.Println("本地文件存在")

			client := &http.Client{}
			req.RequestURI = ""



			//获取远程文件信息
			if resp, err := client.Do(req); err == nil {
				defer resp.Body.Close()
				header := wr.Header()

				for k, h := range resp.Header {
					for _, v := range h {
						fmt.Println("Header:", k, "=", v)
						header.Add(k, v)

					}
				}

				wr.WriteHeader(resp.StatusCode)
				io.Copy(wr, file)
			}

			return
		}

	} else {
		panic(err)
	}

	//透传请求
	client := &http.Client{}
	req.RequestURI = ""

	if resp, err := client.Do(req); err == nil {

		defer resp.Body.Close()

		fmt.Println(req.RemoteAddr, " ", resp.Status)

		//copy header
		header := wr.Header()

		for k, h := range resp.Header {
			for _, v := range h {
				fmt.Println("Header:", k, "=", v)
				header.Add(k, v)

			}
		}

		fmt.Println("new response headers:")
		fmt.Println(header)

		wr.WriteHeader(resp.StatusCode)
		io.Copy(wr, resp.Body)

	} else {

		fmt.Println(err)
	}

}

func main() {
	fmt.Println("start...")

	http.HandleFunc("/", DefaultHandler)
	http.ListenAndServe(":8000", nil)
}
