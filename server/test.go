package main

import (
	"strconv"
	"strings"
	"sync"
)

var mapData map[string]int
var lock sync.RWMutex
var goW sync.WaitGroup

func main() {

	mapData = map[string]int{
		"aa": 1,
		"b":  2,
		"cc": 3,
		"dd": 4,
		"ff": 5,
	}
	i := 0
	for i < 100000 {
		i += 2
		goW.Add(4)
		go read()
		go writer()
		go forData()
		go TestGet()

	}

	goW.Wait()
	println("ok", i)
}

func TestGet() {
	data := getMap()
	_ = data["a"]
}

func getMap() map[string]int {

	defer goW.Done()

	lock.RLock()
	defer lock.RUnlock()

	return map[string]int{
		"aa": mapData["aa"],
		"b":  mapData["b"],
		"cc": mapData["cc"],
		"dd": mapData["dd"],
		"ff": mapData["ff"],
	}
	//return data
}

func forData() {

	_, _ = getAll()

	goW.Done()
}
func read() {
	lock.RLock()
	defer lock.RUnlock()
	_ = mapData["a"]
	goW.Done()
}

func writer() {

	lock.Lock()
	defer lock.Unlock()
	mapData["a"] = 12
	goW.Done()

}

func getAll() (service []int, user []int) {
	lock.RLock()
	defer lock.RUnlock()

	for socketUserId, _ := range mapData {
		temp := strings.Split(socketUserId, ":")
		RoleId, err := strconv.Atoi(temp[0])
		if err != nil {
			continue
		}

		if temp[0] == "service" {
			service = append(service, RoleId)

		}

		if temp[0] == "user" {
			user = append(user, RoleId)
		}
	}
	return service, user
}
