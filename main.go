package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

// 定义JSON数据结构
type Topic struct {
	Summary    string `json:"summary"`
	Discussion []struct {
		Title      string  `json:"title"`
		URL        string  `json:"url"`
		CreatedAt  string  `json:"created_at"`
		SourceType string  `json:"source_type"`
		SourceID   string  `json:"source_id"`
		Cosine     float64 `json:"cosine"`
	} `json:"discussion"`
}

func main() {
	// JSON数据
	jsonData := `
	{
		"0": {
			"summary": "BMC开发者如何优化供电模式切换特性",
			"discussion": [
				{
					"title": "【需求】支持优化供电模式切换特性",
					"url": "https://gitcode.com/openUBMC/power_mgmt/issues/9",
					"created_at": "2025-05-08 13:05:50+00",
					"source_type": "issue",
					"source_id": "gitcode-3085646",
					"cosine": 0.6814956069909915
				},
				{
					"title": "【需求】支持查询冗余供电信息能力",
					"url": "https://gitcode.com/openUBMC/power_mgmt/issues/27",
					"created_at": "2025-05-08 13:06:47+00",
					"source_type": "issue",
					"source_id": "gitcode-3085677",
					"cosine": 0.5606464980268354
				}
			]
		},
		"1": {
			"summary": "BMC开发者如何统一组件模型与接口属性约束，避免内部实现暴露并满足异构算力规范要求",
			"discussion": [
				{
					"title": "【需求】支持异构算力满足资源树协作接口关键字规范要求",
					"url": "https://gitcode.com/openUBMC/bios/issues/22",
					"created_at": "2025-05-09 06:30:18+00",
					"source_type": "issue",
					"source_id": "gitcode-3087012",
					"cosine": 0.6222606560660725
				},
				{
					"title": "【需求】支持异构算力满足资源树协作接口关键字规范要求",
					"url": "https://gitcode.com/openUBMC/pcie_device/issues/7",
					"created_at": "2025-05-09 06:30:20+00",
					"source_type": "issue",
					"source_id": "gitcode-3087013",
					"cosine": 0.6124624867792491
				}
			]
		}
	}`

	// 解析JSON数据
	var topics map[string]Topic
	if err := json.Unmarshal([]byte(jsonData), &topics); err != nil {
		log.Fatal("JSON解析错误:", err)
	}

	// 创建Excel文件
	f := excelize.NewFile()
	sheet := "Sheet1"

	// 设置初始行号
	row := 1

	// 遍历所有话题
	for _, topic := range topics {
		// 第1行：是否入选（标记为1）
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), 1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), 1)
		row++

		// 第2行：热度值（这里使用discussion中cosine的平均值乘以100）
		avgCosine := 0.0
		for _, disc := range topic.Discussion {
			avgCosine += disc.Cosine
		}
		avgCosine = avgCosine / float64(len(topic.Discussion)) * 100
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), avgCosine)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), avgCosine)
		row++

		// 第3行：话题描述
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), topic.Summary)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), topic.Summary)
		row++

		// 第4行开始：讨论源（title & url）
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "讨论源（title&url）")
		row++

		// 写入每个讨论项
		for _, disc := range topic.Discussion {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "")
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), disc.Title)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), disc.URL)
			row++
		}

		// 空一行并设置浅黄色背景
		style, _ := f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Color: []string{"FFFF00"}, Pattern: 1},
		})
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), style)
		row++
	}

	// 保存Excel文件
	if err := f.SaveAs("output.xlsx"); err != nil {
		log.Fatal("保存Excel文件错误:", err)
	}

	fmt.Println("Excel文件已成功生成: output.xlsx")
}
