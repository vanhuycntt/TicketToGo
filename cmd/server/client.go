package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type (
	Geo struct {
		Altitude  float64 `json:"altitude"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
)

var (
	locs = []Geo{
		{-97, 37.819929, -122.478255},
		{1899, 39.096849, -120.032351},
		{2619, 37.865101, -119.538329},
		{42, 33.812092, -117.918974},
		{15, 37.77493, -122.419416},
	}
)

func main() {
	tr := http.DefaultTransport

	client := &http.Client{
		Transport: tr,
		Timeout:   0,
	}
	//GetInfo(client)
	//SubmitChunks(client)
	GetChunks(client)

	time.Sleep(10 * time.Second)
}
func GetInfo(client *http.Client) {
	fmt.Println("Get Info")
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
			Path:   "/getInfo",
		},
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
	}

	defer func() {
		if res.Body != nil {
			_ = res.Body.Close()
		}

	}()

	headerBytes, _ := httputil.DumpResponse(res, false)
	fmt.Println(string(headerBytes))

	bodyBytes, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(bodyBytes))

}
func GetChunks(client *http.Client) {
	fmt.Println("Get Chunked")
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
			Path:   "/getChunked",
		},
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	fmt.Printf("Doing request\n")
	reqHeaderBytes, _ := httputil.DumpRequest(req, false)
	fmt.Printf(string(reqHeaderBytes))
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, res.Body)
		_ = res.Body.Close()
	}()

	headerBytes, _ := httputil.DumpResponse(res, false)
	fmt.Println(string(headerBytes))

	pr01, pw01 := io.Pipe()
	pr02, pw02 := io.Pipe()
	pr03, pw03 := io.Pipe()
	pr04, pw04 := io.Pipe()
	//done :=make(chan struct {})
	defer func() {
		_ = pw01.Close()
		_ = pw02.Close()
		_ = pw03.Close()
		_ = pw04.Close()
	}()
	go func() {
		fmt.Println("channel 01")
		decoder := json.NewDecoder(pr01)

		for {
			/*select {
			case _, ok := <-done:
				if !ok {
					break forLoop
				}
			default:
				var geo Geo

				err := decoder.Decode(&geo)
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(fmt.Sprintf("%v",geo))
			}*/
			var geo Geo

			err := decoder.Decode(&geo)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(fmt.Sprintf("%v", geo))
		}
		fmt.Println("channel 01 completed")
	}()
	go func() {
		fmt.Println("channel 02")
		decoder := json.NewDecoder(pr02)
		for {
			var geo Geo

			err := decoder.Decode(&geo)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(fmt.Sprintf("%v", geo))
		}
		fmt.Println("channel 02 completed")
	}()

	go func() {
		fmt.Println("channel 03")
		decoder := json.NewDecoder(pr03)
		for {
			var geo Geo
			err := decoder.Decode(&geo)
			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(fmt.Sprintf("%v", geo))
		}
		fmt.Println("channel 03 completed")
	}()

	go func() {
		fmt.Println("channel 04")
		decoder := json.NewDecoder(pr04)
		for {
			var geo Geo

			err := decoder.Decode(&geo)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(fmt.Sprintf("%v", geo))
		}
		fmt.Println("channel 04 completed")
	}()
	mw := io.MultiWriter(pw01, pw02, pw03, pw04)
	scanner := bufio.NewScanner(res.Body)
	scanner.Split(bufio.ScanLines)
	/*for decoder.More() {
		var geo Geo
		err := decoder.Decode(&geo)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(fmt.Sprintf("%v",geo))
		json.Marshal(geo)
	}*/
	for scanner.Scan() {
		_, _ = mw.Write([]byte(scanner.Text()))
		time.Sleep(5 * time.Second)
	}
	//close(done)

	//time.Sleep(5*time.Second)
	fmt.Printf("Done request. Err: %v\n", err)
}

func SubmitChunks(client *http.Client) {
	fmt.Println("Submit chunks")
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
			Path:   "/streamText",
		},
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	pr, pw := io.Pipe()
	go func() {
		enc := json.NewEncoder(pw)
		defer pw.Close()
		for _, l := range locs {
			if err := enc.Encode(l); err != nil {
				return
			}
			time.Sleep(10 * time.Second)
		}

	}()
	req.Body = pr

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {

		if res.Body != nil {
			_, _ = io.Copy(ioutil.Discard, res.Body)
			_ = res.Body.Close()
		}
	}()

}
