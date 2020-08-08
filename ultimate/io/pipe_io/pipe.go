package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

//A demo on how to use sync.Cond and io.Pipe to push the same data to multiple readers.
func main() {
	inspector := UrlInspection{
		Callbacks: make([]ResponseCallback, 0),
	}
	locker := sync.NewCond(&sync.Mutex{})
	inspector.RegisterFunc(locker, func(reader io.Reader) {
		if rc, ok := reader.(io.ReadCloser); ok {
			defer rc.Close()
		}
		file, _ := os.Create("channel01.html")
		defer file.Close()
		_, err := io.Copy(file, reader)
		if err != nil {
			fmt.Println("channel 01 => " + err.Error())
		}
		fmt.Println("<= ending 01")

	})
	inspector.RegisterFunc(locker, func(reader io.Reader) {
		if rc, ok := reader.(io.ReadCloser); ok {
			defer rc.Close()
		}

		var buff = new(bytes.Buffer)
		_, err := io.Copy(buff, reader)
		if err != nil {
			fmt.Println("channel 02 => " + err.Error())
		}
		fmt.Println("channel 02 =>" + buff.String())
	})
	inspector.RegisterFunc(locker, func(reader io.Reader) {
		if rc, ok := reader.(io.ReadCloser); ok {
			defer rc.Close()
		}
		fmt.Println("channel 03 => ")
		_, err := io.Copy(os.Stdout, reader)
		if err != nil {
			fmt.Println("channel 03 => " + err.Error())
		}
		fmt.Println("<= ending 03")

	})

	err := inspector.BroadcastFrom(locker, "https://coinmarketcap.com")
	if err != nil {
		fmt.Println(err)
	}
}

type ResponseCallback func(reader io.Reader)

type Register interface {
	RegisterFunc(callback ResponseCallback)
}

type UrlInspection struct {
	Callbacks   []ResponseCallback
	lastSuccess []byte
	version     *uint
}

func (urlCrawler *UrlInspection) RegisterFunc(locker *sync.Cond, callback ResponseCallback) {

	urlCrawler.Callbacks = append(urlCrawler.Callbacks, callback)
}

func (urlCrawler UrlInspection) BroadcastFrom(locker *sync.Cond, url string) error {
	var (
		success []byte
		err     error
		readers []io.Reader
		writers []io.Writer
	)
	for i := 0; i < len(urlCrawler.Callbacks); i++ {
		r, w := io.Pipe()
		readers = append(readers, r)
		writers = append(writers, w)
	}

	mutiWriter := io.MultiWriter(writers...)
	locker.L.Lock()
	go func() {
		var crawlErr error
		locker.L.Lock()
		defer locker.L.Unlock()
		success, crawlErr = urlCrawler.inspectUrl(url)
		if crawlErr != nil {
			err = errors.Wrap(err, "Inspect to URL:"+url)
		}
		locker.Broadcast()

	}()

	for success == nil && err == nil {
		locker.Wait()
	}
	locker.L.Unlock()

	if err != nil {
		return err
	}
	var syncWait sync.WaitGroup
	for idx, callback := range urlCrawler.Callbacks {
		syncWait.Add(1)
		go func(id int, caller ResponseCallback) {
			defer syncWait.Done()
			caller(readers[id])
		}(idx, callback)
	}
	_, err = mutiWriter.Write(success)
	for _, w := range writers {
		_ = w.(*io.PipeWriter).CloseWithError(err)
	}
	syncWait.Wait()
	return err
}

func (urlCrawler UrlInspection) inspectUrl(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Create URL request: "+url)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "Request to : "+url)
	}
	defer func() {
		if res.Body != nil {
			_, _ = io.Copy(ioutil.Discard, res.Body)
			_ = res.Body.Close()
		}
	}()

	respInBytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, errors.Wrap(err, "Parse response from URL"+url)
	}
	return respInBytes, nil
}
