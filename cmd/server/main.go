package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/util/net"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

type (
	Geolocation struct {
		Altitude  float64 `json:"altitude"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
)

var (
	locations = []Geolocation{
		{-97, 37.819929, -122.478255},
		{1899, 39.096849, -120.032351},
		{2619, 37.865101, -119.538329},
		{42, 33.812092, -117.918974},
		{15, 37.77493, -122.419416},
	}
)

func main() {
	http.HandleFunc("/getChunked", func(rw http.ResponseWriter, req *http.Request) {
		reqDumper, _ := httputil.DumpRequest(req, false)

		fmt.Println(string(reqDumper))
		if req.Method != http.MethodGet {

			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(404)
			_, err := rw.Write([]byte(string(`{data: null, msg:"error method"}`)))
			if err != nil {
				log.Panic("internal error")
			}
			return
		}

		rw.WriteHeader(http.StatusOK)
		locs := make(chan Geolocation)
		enc := json.NewEncoder(rw)
		syncWait := sync.WaitGroup{}
		syncWait.Add(1)

		go func() {
			defer syncWait.Done()
			for {
				select {
				case l, ok := <-locs:
					if !ok {
						return
					}
					if err := enc.Encode(l); err != nil {
						fmt.Println("error to encode ", err)
						if net.IsProbableEOF(err) {
							return
						}
						continue
					}
					rw.(http.Flusher).Flush()
				case <-req.Context().Done():
					fmt.Println("client closed")
					return
				}
			}

		}()

		go func() {
			for _, l := range locations {
				locs <- l
			}
			close(locs)
		}()

		syncWait.Wait()

	})
	http.HandleFunc("/getInfo", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(503)
			_, _ = rw.Write([]byte(string(`{data: null, msg:"error method"}`)))
			return
		}
		_, _ = fmt.Fprint(rw, "hello info \n")
	})

	http.HandleFunc("/streamText", func(rw http.ResponseWriter, req *http.Request) {
		reqDumper, _ := httputil.DumpRequest(req, false)

		fmt.Println(string(reqDumper))
		decoder := json.NewDecoder(req.Body)

		for {
			select {
			case <-req.Context().Done():
				fmt.Println("connection closed")
				return
			default:
				var geo Geolocation
				err := decoder.Decode(&geo)
				if err != nil {
					fmt.Println(err)
					break
				}
				fmt.Println(geo)
			}

		}
		fmt.Println("TRANSMISSION COMPLETE")
	})
	authMiddleware := func(reqHandler http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {

			reqHandler(rw, req)
		}
	}
	reqHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

	})

	http.Handle("/", authMiddleware(reqHandler))

	server := http.Server{
		Addr:           "127.0.0.1:8080",
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   2 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic("Unable to start server")
	}

}
