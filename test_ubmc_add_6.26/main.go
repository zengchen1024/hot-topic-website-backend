package main

import (
	"fmt"
	"time"

	"github.com/opensourceways/hot-topic-website-backend/common/infrastructure/mongodb"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/infrastructure/repositoryimpl"
	"github.com/opensourceways/hot-topic-website-backend/test_ubmc_add_6.26/lib"
	"github.com/opensourceways/hot-topic-website-backend/test_ubmc_add_6.26/not_hot"
)

func main() {
	err := mongodb.Init(&mongodb.Config{
		Conn:    "",
		DBName:  "community-hot-topic",
		CAFile:  "",
		Timeout: 10,
	})
	if err != nil {
		fmt.Printf("connect dst failed, err:%s", err)

		return
	}

	defer mongodb.Close()

	//handleHotTopic()
	//handleNewHotTopic()
	handleNotHotTopic()
}

func toSecond(dateStr string) (int64, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, fmt.Errorf("解析日期出错:", err)
	}

	// 转换为 Unix 时间戳（秒）
	return t.Unix(), nil
}

func handleHotTopic() {
	fmt.Println("-------- handle hot topic --------")

	date := "2025-06-20"
	closedAt, err := toSecond(date)
	if err != nil {
		fmt.Println(err)

		return
	}

	ftr, err := lib.NewfileToReview("../test_review/openubmc_2025-06-20.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer ftr.Exit()

	topics, err := ftr.ReadLastHotTopics()
	if err != nil {
		fmt.Println(err)
		return
	}
	topicMap := make(map[string]*domain.HotTopic, len(topics))
	for i := range topics {
		item := &topics[i]
		topicMap[item.Id] = item
	}
	fmt.Printf("load total %d topics\n", len(topics))
	for i := range topics {
		fmt.Println(topics[i].Id)
	}

	community := "openubmc"
	repo := repositoryimpl.NewHotTopic(
		map[string]repositoryimpl.Dao{community: mongodb.DAO("openubmc_hot_topic")},
	)

	oldOnes, err := repo.FindOpenOnes(community)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("load %d open topics\n", len(oldOnes))
	for i := range oldOnes {
		fmt.Println(oldOnes[i].Id)
	}

	resolvedIds := map[string]int{
		"684ba0f12cd21bf7988a420b": 6,
		"684ba0f12cd21bf7988a420e": 7,
		"684ba0f12cd21bf7988a4216": 8,
		"684ba0f12cd21bf7988a4220": 9,
		"6852e1ea2cd21bf798f47ef5": 10,
	}
	appendedIds := map[string]int{
		"684ba0f12cd21bf7988a4213": 5,
	}

	for i := range oldOnes {
		topic := &oldOnes[i]

		if order, ok := resolvedIds[topic.Id]; ok {
			fmt.Printf("find resolved:%s\n", topic.Id)
			current := topicMap[topic.Id]

			topic.Order = order
			topic.DiscussionSources = current.DiscussionSources
			topic.StatusTransferLog = append(topic.StatusTransferLog, domain.StatusLog{
				Time:   date,
				Status: "Resolved",
			})

			if err := repo.Resolved(community, topic, closedAt); err != nil {
				fmt.Printf("add topic failed, err:%s\n", err.Error())
				return
			}

			continue
		}

		if order, ok := appendedIds[topic.Id]; ok {
			current := topicMap[topic.Id]

			topic.Order = order
			topic.DiscussionSources = current.DiscussionSources
			topic.StatusTransferLog = append(topic.StatusTransferLog, domain.StatusLog{
				Time:   "2025-06-19",
				Status: "Appended",
			})

			if err := repo.Appended(community, topic); err != nil {
				fmt.Printf("add topic failed, err:%s\n", err.Error())
				return
			}

			continue
		}
	}

	fmt.Println("done")
}

func handleNewHotTopic() {
	fmt.Println("-------- handle new hot topic --------")

	ftr, err := lib.NewfileToReview("../test_review/openubmc_2025-06-20.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer ftr.Exit()

	topics, err := ftr.ReadOtherTopics()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("load total %d topics\n", len(topics))
	/*
		for i := range topics {
			fmt.Println(topics[i].Id)
		}
	*/

	community := "openubmc"
	repo := repositoryimpl.NewHotTopic(
		map[string]repositoryimpl.Dao{community: mongodb.DAO("openubmc_hot_topic")},
	)

	newOnes := map[string]int{
		"BMC开发者Ubuntu WSL2构建QEMU版本权限不足，如何调整文件系统权限及shared目录配置？": 1,
		"openUBMC如何管理开源漏洞及依赖修复周期":                              2,
		"非天池架构BMC中EXU升级主板CPLD时需解决I2C拓扑配置及无SMC场景下的访问路径适配问题":     3,
		"BMC开发者如何实现单个风扇独立调速配置":                                 4,
	}
	newOnesDate := map[string]string{
		"BMC开发者Ubuntu WSL2构建QEMU版本权限不足，如何调整文件系统权限及shared目录配置？": "2025-06-09",
		"openUBMC如何管理开源漏洞及依赖修复周期":                              "2025-06-09",
		"非天池架构BMC中EXU升级主板CPLD时需解决I2C拓扑配置及无SMC场景下的访问路径适配问题":     "2025-06-07",
		"BMC开发者如何实现单个风扇独立调速配置":                                 "2025-05-09",
	}

	for i := range topics {
		topic := &topics[i]

		if order, ok := newOnes[topic.Title]; ok {
			fmt.Printf("find resolved:%s\n", topic.Title)

			topic.Order = order
			topic.StatusTransferLog = []domain.StatusLog{{
				Time:   newOnesDate[topic.Title],
				Status: "New",
			}}

			if err := repo.Add(community, topic); err != nil {
				fmt.Printf("add topic failed, err:%s\n", err.Error())
				return
			}

			continue
		}
	}

	fmt.Println("done")
}

func handleNotHotTopic() {
	fmt.Println("-------- handle new hot topic --------")

	ftr, err := not_hot.NewfileToReview("../test_review/openubmc_2025-06-20.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer ftr.Exit()

	topics, err := ftr.ReadOtherTopics()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("load total %d topics\n", len(topics))
	/*
		for i := range topics {
			fmt.Println(topics[i].Id)
		}
	*/

	community := "openubmc"
	repo := repositoryimpl.NewNotHotTopic(
		map[string]repositoryimpl.Dao{community: mongodb.DAO("openubmc_not_hot_topic")},
	)

	newOnes := map[string]int{
		"BMC开发者Ubuntu WSL2构建QEMU版本权限不足，如何调整文件系统权限及shared目录配置？": 1,
		"openUBMC如何管理开源漏洞及依赖修复周期":                              2,
		"非天池架构BMC中EXU升级主板CPLD时需解决I2C拓扑配置及无SMC场景下的访问路径适配问题":     3,
		"BMC开发者如何实现单个风扇独立调速配置":                                 4,
	}

	for i := range topics {
		topic := &topics[i]

		if _, ok := newOnes[topic.Title]; !ok {
			fmt.Printf("find not hot:%s\n", topic.Title)

			if err := repo.Add(community, topic); err != nil {
				fmt.Printf("add topic failed, err:%s\n", err.Error())
				return
			}
		}
	}

	fmt.Println("done")
}
