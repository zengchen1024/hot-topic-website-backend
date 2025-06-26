package lib

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

// DiscussionSourceInfo
type DiscussionSourceInfo struct {
	Title string `json:"title" required:"true"`

	domain.DiscussionSource
}

func (info *DiscussionSourceInfo) ToDiscussionSourceInfo() domain.DiscussionSourceInfo {
	return domain.DiscussionSourceInfo{
		Id:    info.Id,
		URL:   info.URL,
		Title: info.Title,
	}
}

func GetDSs(file string) ([]DiscussionSourceInfo, error) {
	jsonData, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("parse json failed, err:%s\n", err.Error())

		return nil, err
	}

	// 解析JSON数据
	var dss []DiscussionSourceInfo
	if err := json.Unmarshal([]byte(jsonData), &dss); err != nil {
		fmt.Printf("JSON解析错误:%s\n", err.Error())
		return nil, err
	}

	return dss, nil
}

func main1() {
	v, err := GetDSs("")
	if err != nil {
		return
	}
	fmt.Println(v[0])

}
