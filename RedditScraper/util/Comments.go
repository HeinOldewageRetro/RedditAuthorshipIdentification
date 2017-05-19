package util

import (
	"encoding/json"
	"fmt"

	"strings"
)

func GetAllComments(author string) []string {
	comments, after, more := getCommentsAfter(author, "")
	for more {
		var Morecomments []string
		Morecomments, after, more = getCommentsAfter(author, after)
		comments = append(comments, Morecomments...)
	}
	return comments
}

func getCommentsAfter(author, after string) ([]string, string, bool) {
	var params map[string]string
	params = make(map[string]string)
	if after != "" {

		params["after"] = after

	}
	params["limit"] = "100"

	thepage := DoRequest("https://oauth.reddit.com/user/"+author+"/comments", params)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("thepage", err)
		}
	}()
	decoder := json.NewDecoder(strings.NewReader(thepage))
	var out map[string]interface{}
	decoder.Decode(&out)

	posts := out["data"].(map[string]interface{})["children"].([]interface{})

	comments := make([]string, 0, len(posts))

	for k := range posts {
		//permalink := posts[k].(map[string]interface{})["data"].(map[string]interface{})["permalink"]
		//fullName := posts[k].(map[string]interface{})["data"].(map[string]interface{})["name"].(string)
		comment := posts[k].(map[string]interface{})["data"].(map[string]interface{})["body"].(string)

		comments = append(comments, comment)

	}

	aI := out["data"].(map[string]interface{})["after"]
	if aI != nil {
		return comments, aI.(string), true
	} else {
		return comments, "", false
	}

}
