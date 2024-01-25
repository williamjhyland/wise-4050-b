package rest_api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

func RunRESTAPI() {
	username := "root"       // Replace with your username
	password := "00000000"   // Replace with your password
	deviceIP := "10.1.14.38" // Replace with your device IP
	url := fmt.Sprintf("http://%s/di_value/slot_0", deviceIP)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Add("Authorization", "Basic "+auth)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status code %d\n", resp.StatusCode)
		return
	}

	var data struct {
		DIVal []struct {
			Ch     int `json:"Ch"`
			Md     int `json:"Md"`
			Stat   int `json:"Stat"`
			Val    int `json:"Val"`
			Cnting int `json:"Cnting"`
			ClrCnt int `json:"ClrCnt"`
			OvLch  int `json:"OvLch"`
		} `json:"DIVal"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	for _, di := range data.DIVal {
		fmt.Printf("Channel %d: Mode=%d, Status=%d, Value=%d, Counting=%d, ClearCount=%d, OverLatch=%d\n",
			di.Ch, di.Md, di.Stat, di.Val, di.Cnting, di.ClrCnt, di.OvLch)
	}
}
