package main

import (
	"time"

	"net/http"

	"bytes"

	"os"

	"fmt"

	"io/ioutil"
	"strconv"
)

func main() {
	args := os.Args

	url := args[1]
	asgId := args[2]
	nodeId := args[3]
	value := args[4]
	interval, _ := strconv.Atoi(args[5])

	client := &http.Client{}
	for {
		var jsonStr = []byte(`{"ID":"` + asgId + `","NodeID":"` + nodeId + `", "Metrics":[{"Value":` + value + `, "Time":"` + time.Now().Format(time.RFC3339Nano) + `"}]}`)
		fmt.Printf("Request: \n ---------\n %s \n ---------- \n", jsonStr)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Resp err: %s \n\n", err)
		}
		body, _ := ioutil.ReadAll(resp.Body)

		fmt.Printf("Resp: %+v \n\n", resp)
		fmt.Printf("Body: %s \n\n", body)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
