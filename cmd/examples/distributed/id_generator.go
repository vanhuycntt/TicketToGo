package main

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/pkg/errors"
	"math/rand"
)

func main() {
	workerIds := []int{1, 2, 3, 4}

	var nodes = make(map[int]*snowflake.Node)
	for _, id := range workerIds {
		node, _ := snowflake.NewNode(int64(id))

		nodes[id] = node
	}
	setIDs := hashset.New()
	for _, v := range nodes {
		for i := 0; i < 100; i++ {
			uniqueID := v.Generate()
			fmt.Println(fmt.Sprintf("%d|%d|%d", uniqueID.Node(), uniqueID.Time(), uniqueID.Step()))

			if setIDs.Contains(uniqueID) {
				panic(errors.New("Duplicated ID"))
			}

			setIDs.Add(uniqueID)
		}
	}

	fmt.Println(setIDs)
	for i := 0; i < 10; i++ {
		fmt.Println(rand.Perm(10))
	}

}
