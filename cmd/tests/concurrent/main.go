package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type ItemURL string

type URLEngine struct {
	batchSize int
	urls      []ItemURL
}

func (en URLEngine) Collect() []*ItemURLResponse {
	var (
		urlResults []*ItemURLResponse
	)
	nb := en.numBatches()
	itemsChan := make(chan []*ItemURLResponse, nb)
	for i := 0; i < nb; i++ {
		bc := BatchCrawler{
			id: i,
		}
		bc.itemURLs = en.partitionAt(i)
		bc.startIndex = i * en.batchSize
		bc.endIndex = bc.startIndex + len(bc.itemURLs)

		go func() {
			resArr := bc.Crawl()
			itemsChan <- resArr
		}()
	}
	batchCompleted := 0
breakLoop:
	for {
		select {
		case resArr := <-itemsChan:
			urlResults = append(urlResults, resArr...)
			batchCompleted++
			if batchCompleted == nb {
				break breakLoop
			}

		default:

		}
	}
	return urlResults
}
func (en URLEngine) partitionAt(batchSeq int) []ItemURL {
	var batchesURL []ItemURL
	for i := batchSeq * en.batchSize; i < (batchSeq+1)*en.batchSize; i++ {
		if i < len(en.urls) {
			batchesURL = append(batchesURL, en.urls[i])
		}
	}
	return batchesURL
}

func (en URLEngine) numBatches() int {
	l := len(en.urls)
	nb := l / en.batchSize
	if l%en.batchSize > 0 {
		nb++
	}
	return nb
}

type BatchCrawler struct {
	id         int
	itemURLs   []ItemURL
	startIndex int
	endIndex   int
}

func (bc BatchCrawler) Crawl() []*ItemURLResponse {
	var (
		syncWait    sync.WaitGroup
		batchesResp = make([]*ItemURLResponse, bc.endIndex-bc.startIndex)
	)
	for idx, it := range bc.itemURLs {

		siteURL, err := url.Parse(string(it))
		if err != nil {
			fmt.Println(err)
			continue
		}
		seqIdx := idx + bc.startIndex
		syncWait.Add(1)
		go func(atIdx int) {
			defer syncWait.Done()
			rs := bc.inspectURL(seqIdx, *siteURL)

			batchesResp[atIdx] = rs

		}(idx)
	}

	syncWait.Wait()
	return batchesResp
}
func (bc BatchCrawler) GlobalIndex(indexAt int) int {
	return bc.startIndex + indexAt
}
func (bc BatchCrawler) inspectURL(idx int, url url.URL) *ItemURLResponse {
	var (
		err error
	)
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	resp, err := http.DefaultClient.Do(req)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &ItemURLResponse{
		idx:      idx,
		resType:  "string",
		response: string(buff),
	}
}

type ItemURLResponse struct {
	idx      int
	resType  string
	response string
}

func main() {
	/*bc := BatchCrawler{
		id: 0,
		startIndex: 0,
		endIndex: 3,
		itemURLs:[] ItemURL {

		},
	}
	bc01 := BatchCrawler{
		id: 1,
		startIndex: 3,
		endIndex: 5,
		itemURLs:[] ItemURL {
			"http://google.com",
			"http://microsoft.com",
		},
	}
	itResp := bc.Crawl()
	bc01.Crawl()*/
	urlExtractor := URLEngine{
		batchSize: 2,
		urls: []ItemURL{
			"http://vnexpress.net",
			"http://baomoi.com",
			"http://tinhte.vn",
			"http://google.com",
			"http://microsoft.com",
		},
	}
	itResp := urlExtractor.Collect()
	/*sort.Slice(itResp, func(i, j int) bool {
		if itResp[i].idx < itResp[j].idx {
			return true
		}
		return false
	})*/
	for _, rs := range itResp {
		fmt.Println(fmt.Sprintf("Index %d", rs.idx))
	}

}

func newSyncMap() syncMapLock {
	mu := sync.Mutex{}
	cd := sync.NewCond(&mu)
	return syncMapLock{
		cond: cd,
		mu:   &mu,
		keys: make(map[string]interface{}),
	}
}

type syncMapLock struct {
	cond *sync.Cond
	mu   *sync.Mutex
	keys map[string]interface{}
}

func (sml *syncMapLock) Lock(k string) {
	checkCond := func() bool {
		_, ok := sml.keys[k]
		return ok
	}

	sml.mu.Lock()
	defer sml.mu.Unlock()

	for checkCond() {
		fmt.Println(fmt.Sprintf("%d | key:%s existed", time.Now().Unix(), k))
		sml.cond.Wait()
	}
	sml.keys[k] = struct{}{}
}

func (sml *syncMapLock) UnLock(k string) {
	checkCond := func() bool {
		_, ok := sml.keys[k]
		return ok
	}
	sml.mu.Lock()
	defer sml.mu.Unlock()
	if checkCond() {
		delete(sml.keys, k)
		sml.cond.Broadcast()
	}
}
