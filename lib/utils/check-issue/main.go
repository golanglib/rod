package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/go-rod/rod/lib/utils"
	"github.com/ysmood/kit"
)

func main() {
	data, err := kit.ReadFile(os.Getenv("GITHUB_EVENT_PATH"))
	utils.E(err)

	issue := kit.JSON(data).Get("issue")

	labels := issue.Get("labels").Array()

	for _, l := range labels {
		name := l.Get("name").Str
		if name != "question" && name != "bug" {
			kit.Log("skip", name)
			return
		}
	}

	num := issue.Get("number").Int()
	body := issue.Get("body").Str

	kit.Log("check issue", num)

	m := regexp.MustCompile(`\*\*Rod Version:\*\* v[0-9.]+`).FindString(body)
	if m == "" || m == "**Rod Version:** v0.0.0" {
		kit.Log("invalid issue format", body)

		currentVer := req("/repos/go-rod/rod/releases").
			Query("per_page", "1").
			MustJSON().Get("0.tag_name").Str

		kit.Log("current rod version", currentVer)

		q := req(fmt.Sprintf("/repos/go-rod/rod/issues/%d/comments", num)).
			Post().
			JSONBody(map[string]string{
				"body": fmt.Sprintf(
					"Please add a valid `**Rod Version:** v0.0.0` in your issue. Current version is %s",
					currentVer,
				),
			})

		if q.MustResponse().StatusCode >= 400 {
			panic(q.MustString())
		}
	}
}

func req(u string) *kit.ReqContext {
	return kit.Req("https://api.github.com"+u).Header("Authorization", "token "+os.Getenv("GH_ROBOT_TOKEN"))
}
