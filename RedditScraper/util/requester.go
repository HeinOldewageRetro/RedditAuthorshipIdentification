package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var token = "g5Fc09JHchSlQOzO6vFR4EtzeFF" //"g5Fc09JHchSlQOzO6vFR4EtzeFE"

func getToken() string {
	for {

		body := &bytes.Buffer{}
		fmt.Fprint(body, "grant_type=password&username=NeoinKarasuYuu&password=Neoin123456&duration=permanent")
		req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", body)
		if err != nil {
			panic(err)
		}
		req.SetBasicAuth("Gm4SgDPOnuR31Q", "FpRas9I5PR-K5OH8PhQDGUQOo0g")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != 200 {
			continue
		}
		b, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		val := make(map[string]interface{})
		json.Unmarshal(b, &val)
		return val["access_token"].(string)
	}
}

func DoRequest(url string, params map[string]string) string {

	for k := 0; k < 10; k++ {
		client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))

		r, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
		r.Header["Authorization"] = []string{"bearer  " + token}
		r.Header["User-Agent"] = []string{"PC:Gm4SgDPOnuR31Q:1.0 (by /u/NeoinKarasuYuu)"}
		if params != nil {
			q := r.URL.Query()
			for k, v := range params {
				q.Set(k, v)
			}
			r.URL.RawQuery = q.Encode()
		}

		if params != nil {
			for k, v := range params {
				r.Header[k] = []string{v}
			}
		}
		var resp *http.Response
		indented := &bytes.Buffer{}

		resp, err = client.Do(r)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("resp.StatusCode =", resp.StatusCode)
		if resp.StatusCode != 200 {
			fmt.Println(string(body))
			fmt.Println("resp.StatusCode != 200", resp.StatusCode)
			token = getToken()
			fmt.Println("Got new token", token)
		} else {

			err = json.Indent(indented, body, "", "\t")
			if err != nil {
				fmt.Println(string(body))
				panic(err)
			}

			remaining, err := strconv.ParseFloat(resp.Header.Get("X-Ratelimit-Remaining"), 64)
			if err != nil {
				fmt.Println("Error parsing X-Ratelimit-Remaining", err)
				fmt.Println(resp.Header)

			}
			if remaining < 1 || err != nil {
				waitTime, _ := strconv.Atoi(resp.Header.Get("X-Ratelimit-Reset"))
				fmt.Println("Waiting", waitTime)
				time.Sleep((time.Duration(waitTime) + 5) * time.Second)
			}
			return string(indented.Bytes())
		}
	}
	return ""
}
