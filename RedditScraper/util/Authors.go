package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GetAuthorNames(before string) (authors []string, after string) {
	var params map[string]string
	params = make(map[string]string)
	if before != "" {
		params["after"] = before
	}
	params["limit"] = "100"
	thepage := DoRequest("https://oauth.reddit.com/r/all/hot", params)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("err", err)
			fmt.Println("thepage", thepage)
		}
	}()
	decoder := json.NewDecoder(strings.NewReader(thepage))
	var out map[string]interface{}
	decoder.Decode(&out)

	posts := out["data"].(map[string]interface{})["children"].([]interface{})
	after = out["data"].(map[string]interface{})["after"].(string)

	Authors := make([]string, 0, len(posts))

	for k := range posts {
		//permalink := posts[k].(map[string]interface{})["data"].(map[string]interface{})["permalink"]
		//fullName := posts[k].(map[string]interface{})["data"].(map[string]interface{})["name"].(string)
		author := posts[k].(map[string]interface{})["data"].(map[string]interface{})["author"].(string)

		Authors = append(Authors, author)

	}
	return Authors, after
}
