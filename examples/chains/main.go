package main

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/lists/arraylist"
)

func main() {
	var values = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	fmt.Println(sumInts(values...))
	oldDistrict := arraylist.New()
	newDistricts := arraylist.New()
	newDiffOld := diff(newDistricts, oldDistrict)
	oldDiffNew := diff(oldDistrict, newDistricts)

	fmt.Println(newDiffOld.String())
	fmt.Println(oldDiffNew.String())
	values = []int{1, 2, 3, 4}

	chains := NewSequenceChains(
		func(id int) {
			fmt.Println(id)
		}, values...)

	nilChain := (*ChainHandler)(nil)
	if chains != nilChain {
		chains.Do(context.Background())
	}

}

func diff(alist *arraylist.List, blist *arraylist.List) *arraylist.List {
	result := arraylist.New()
	for _, dt := range alist.Values() {

		has := blist.Any(func(idx int, val interface{}) bool {
			if dt != val {
				return false
			}
			return true
		})
		if !has {
			result.Add(dt)
		}
	}
	return result
}

func decorate(values ...int) {
	if len(values) == 0 {
		return
	}
	decorate(values[1:]...)
	fmt.Println(values[0])
}

func sumInts(values ...int) int {
	sum := 0
	iterInts(
		func(val int) {
			sum += val
		}, values...,
	)
	return sum
}
func iterInts(f func(intVal int), values ...int) {
	seedVal := values[0]

	var nextValues []interface{}
	for _, val := range values[1:] {
		nextValues = append(nextValues, val)
	}
	iter(
		func(val interface{}) {
			f(val.(int))
		},
		seedVal,
		nextValues,
	)
}

func iter(f func(val interface{}), seedVal interface{}, values []interface{}) {
	if len(values) == 0 {
		f(seedVal)
		return
	}
	iter(
		f,
		values[0],
		values[1:],
	)
	f(seedVal)
}

func wrapFunc(f func()) func() {
	return func() {
		f()
	}
}

type (
	Handler interface {
		Do(context.Context)
	}

	HandleFunc func(context.Context)

	DecorateHandler func(Handler) Handler

	ChainHandler struct {
		idx  int
		next Handler
		do   func(idx int)
	}
)

func (hf HandleFunc) Do(ctx context.Context) {
	hf(ctx)
}

func NewChain(id int, hf func(idx int)) Handler {
	return &ChainHandler{
		idx: id,
		do:  hf,
	}
}
func NewSequenceChains(handleF func(id int), values ...int) Handler {

	var chains, curHandler *ChainHandler

	for _, v := range values {
		if chains == nil {
			chains = &ChainHandler{
				idx: v,
				do:  handleF,
			}
			curHandler = chains
			continue
		}
		nextHandler := &ChainHandler{
			idx: v,
			do:  handleF,
		}
		curHandler.next = nextHandler
		curHandler = nextHandler
	}
	return chains
}

func (ch ChainHandler) Do(ctx context.Context) {

	ch.do(ch.idx)

	if ch.next != nil {
		ch.next.Do(ctx)
	}
}

func WithParam(params map[string]interface{}) DecorateHandler {
	return func(handler Handler) Handler {
		return HandleFunc(func(ctx context.Context) {
			handler.Do(ctx)

		})
	}
}
