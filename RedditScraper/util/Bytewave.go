package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GetBytewaveAuthorNames(before string) (authors []string, after string, more bool) {
	var params map[string]string
	params = make(map[string]string)
	if before != "" {
		params["after"] = before
	}
	params["limit"] = "100"
	thepage := DoRequest("https://oauth.reddit.com/user/Bytewave/submitted", params)
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

	authors = make([]string, 0, len(posts))

	for k := range posts {

		postID := posts[k].(map[string]interface{})["data"].(map[string]interface{})["id"].(string)

		fmt.Println("GetBytewaveAuthorNames post ID ", postID)
		SomeAuthors := GetPostCommentAuthors(postID)
		authors = append(authors, SomeAuthors...)

	}

	aI := out["data"].(map[string]interface{})["after"]
	if aI != nil {
		return authors, aI.(string), true
	} else {
		return authors, "", false
	}

}

func GetPostCommentAuthors(postID string) []string {
	authors, after, more := getPostCommentAuthorsAfter(postID, "")
	for more {
		var Moreauthors []string
		Moreauthors, after, more = getPostCommentAuthorsAfter(postID, after)
		authors = append(authors, Moreauthors...)
	}
	return authors
}

func getPostCommentAuthorsAfter(postID, after string) ([]string, string, bool) {
	var params map[string]string
	params = make(map[string]string)
	if after != "" {

		params["after"] = after

	}
	//params["limit"] = "100"

	thepage := DoRequest("https://oauth.reddit.com/comments/"+postID, params)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("thepage", err)
		}
	}()
	decoder := json.NewDecoder(strings.NewReader(thepage))
	var arrout []map[string]interface{}
	err := decoder.Decode(&arrout)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(thepage)

	authors := make([]string, 0)

	for _, out := range arrout {

		data := out["data"].(map[string]interface{})

		authors = append(authors, extractAuthors(data)...)

	}

	/*aI := out["data"].(map[string]interface{})["after"]
	if aI != nil {
		return authors, aI.(string), true
	} else {
		return authors, "", false
	}*/
	return authors, "", false
}

func extractAuthors(data map[string]interface{}) []string {
	authors := make([]string, 0)

	comments := data["children"].([]interface{})
	for _, comment := range comments {
		kind := comment.(map[string]interface{})["kind"].(string)
		if kind != "t1" {
			continue
		}
		data := comment.(map[string]interface{})["data"].(map[string]interface{})

		authors = append(authors, data["author"].(string))
		if replies, ok := data["replies"]; ok {

			if _, ok := replies.(map[string]interface{}); !ok {
				continue
			}

			authors = append(authors, extractAuthors(replies.(map[string]interface{})["data"].(map[string]interface{}))...)
		}

	}
	return authors
}
