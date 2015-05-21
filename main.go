package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Shard struct {
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Region string `json:"region_tag"`
}

type Status struct {
	Status  string `json:"status"`
	AppName string `json:"name"`
}

type Statuses struct {
	Services []Status `json:"services"`
}

var base_url = "http://status.leagueoflegends.com/"

func main() {
	shards := getShards()
	for index, shard := range shards {
		fmt.Printf("%d : %s\n", index, shard.Name)
	}
	fmt.Println("Enter shard:")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	shard_index, conversion_error := strconv.Atoi(cleanString(text))
	if conversion_error != nil {
		panic(conversion_error)
	}
	fmt.Println("Selected ", shards[shard_index].Name)
	statuses := getStatuses(shards[shard_index].Region)
	for _, app := range statuses {
		fmt.Println(app.AppName, " is ", app.Status)
	}
}

func getShards() []Shard {
	shardUrl := base_url + "/shards"
	resp, err := http.Get(shardUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	raw_json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var shardList []Shard
	err = json.Unmarshal(raw_json, &shardList)
	fmt.Println(raw_json)
	if err != nil {
		panic(err)
	}
	return shardList
}

func getStatuses(region string) []Status {
	status_url := base_url + "/shards/" + region
	resp, _ := http.Get(status_url)
	defer resp.Body.Close()
	raw, _ := ioutil.ReadAll(resp.Body)
	var services Statuses
	json.Unmarshal(raw, &services)
	return services.Services
}

func cleanString(text string) string {
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)
	return text
}
