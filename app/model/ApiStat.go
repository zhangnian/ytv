package model

import (
	"fmt"
	"sync"
)

var (
	locker  sync.Mutex
	statMap map[string]*ApiStat
)

type ApiStat struct {
	Total int
	Succ  int
	Cost  int
}

func init() {
	statMap = make(map[string]*ApiStat)
}

func Update(key string, is_succ bool, cost int) {
	fmt.Println("Update cost: ", cost)
	locker.Lock()
	defer locker.Unlock()

	data, ok := statMap[key]
	if ok {
		data.Cost = (data.Cost*data.Total + cost) / (data.Total + 1)
		data.Total += 1
	} else {
		data = &ApiStat{Total: 1, Cost: cost}
	}

	fmt.Println("data.Cost: ", data.Cost)
	if is_succ {
		data.Succ += 1
	}

	statMap[key] = data

	fmt.Printf("%s -> %v\n", key, data)
}
