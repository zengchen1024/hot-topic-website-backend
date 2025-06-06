package app

import (
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

const (
	sheetNewTopics       = "new_topics"
	sheetLastTopics      = "last_hot_topics"
	sheetAppendToOld     = "append_to_last_discarded"
	sheetRemoveFromOld   = "remove_from_last_discarded"
	sheetUnchangedTopics = "unchanged_topics"
)

var (
	topicSpliterStyle             int
	removedDiscussionSourceStyle  int
	appendedDiscussionSourceStyle int
)

func Init() (err error) {
	f := excelize.NewFile()

	newStyle := func(color string) (int, error) {
		return f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{color},
				Pattern: 1,
			},
		})
	}

	if topicSpliterStyle, err = newStyle("92D050"); err != nil {
		return
	}

	if appendedDiscussionSourceStyle, err = newStyle("FFFF43"); err != nil {
		return
	}

	if removedDiscussionSourceStyle, err = newStyle("FFC7CE"); err != nil {
		return
	}

	return
}

type fileToReview struct {
	file     *excelize.File
	filePath string

	rowNewTopic       int
	rowAppendToOld    int
	rowRemoveFromOld  int
	rowUnchangedTopic int
}

func newfileToReview(file string) (*fileToReview, error) {
	all := []string{
		sheetLastTopics, // last hot topic first
		sheetNewTopics,
		sheetAppendToOld,
		sheetRemoveFromOld,
		sheetUnchangedTopics,
	}

	f := excelize.NewFile()
	for _, name := range all {
		if _, err := f.NewSheet(name); err != nil {
			return nil, err
		}
	}

	if err := f.DeleteSheet("Sheet1"); err != nil {
		return nil, err
	}

	return &fileToReview{
		file:     f,
		filePath: file,
	}, nil
}

func (ftr *fileToReview) saveToFile() error {
	return ftr.file.SaveAs(ftr.filePath)
}

func (ftr *fileToReview) saveLastTopics(oldTopics []domain.HotTopic, current map[int]*OptionalTopic) error {
	if len(oldTopics) == 0 {
		return nil
	}

	row := 0
	for i := range oldTopics {
		old := &oldTopics[i]

		item, ok := current[old.Order]
		if !ok {
			return errors.New("can't find topic from current ones")
		}

		item.updateAppended(old.GetDSSet())

		ftr.saveAppendedTopic(item, &row, sheetLastTopics)
	}

	return nil
}

func (ftr *fileToReview) saveAppendedTopic(topic *OptionalTopic, row *int, sheet string) error {
	f := ftr.file

	// 话题描述
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "话题描述")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), topic.Title)
	(*row)++

	// 讨论源（title & url）
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "相关讨论源（title&url）")
	(*row)++

	// 写入每个讨论项
	items := topic.sort()
	for _, item := range items {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "")

		bcell := fmt.Sprintf("B%d", row)
		ccell := fmt.Sprintf("C%d", row)
		f.SetCellValue(sheet, bcell, item.Title)
		f.SetCellValue(sheet, ccell, item.URL)

		if item.appended {
			f.SetCellStyle(sheet, bcell, ccell, appendedDiscussionSourceStyle)
		}

		(*row)++
	}

	// 空一行并设置浅黄色背景
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), topicSpliterStyle)
	(*row)++

	return nil
}

func (ftr *fileToReview) saveNewTopic(topic *OptionalTopic) error {
	return ftr.saveTopic(topic, &ftr.rowNewTopic, sheetNewTopics)
}

func (ftr *fileToReview) saveUnchangedTopic(topic *OptionalTopic) error {
	return ftr.saveTopic(topic, &ftr.rowUnchangedTopic, sheetUnchangedTopics)
}

func (ftr *fileToReview) saveTopic(topic *OptionalTopic, row *int, sheet string) error {
	f := ftr.file

	// 话题描述
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "话题描述")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), topic.Title)
	(*row)++

	// 讨论源（title & url）
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "相关讨论源（title&url）")
	(*row)++

	// 写入每个讨论项
	for i := range topic.DiscussionSources {
		item := &topic.DiscussionSources[i]

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "")

		bcell := fmt.Sprintf("B%d", row)
		ccell := fmt.Sprintf("C%d", row)
		f.SetCellValue(sheet, bcell, item.Title)
		f.SetCellValue(sheet, ccell, item.URL)

		(*row)++
	}

	// 空一行并设置浅黄色背景
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), topicSpliterStyle)
	(*row)++

	return nil
}

func (ftr *fileToReview) saveTopicThatAppendToOld(topic *OptionalTopic, dsIdsOfOldTopic map[int]bool) error {
	topic.updateAppended(dsIdsOfOldTopic)

	return ftr.saveAppendedTopic(topic, &ftr.rowAppendToOld, sheetAppendToOld)
}

func (ftr *fileToReview) saveTopicThatRemoveFromOld(oldTopic *domain.NotHotTopic, dsIdsOfNewTopic map[int]bool) error {
	f := ftr.file
	row := ftr.rowRemoveFromOld
	sheet := sheetRemoveFromOld

	// 话题描述
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "话题描述")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), oldTopic.Title)
	row++

	// 讨论源（title & url）
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "相关讨论源（title&url）")
	row++

	// 写入每个讨论项
	oldTopic.UpdateRemoved(dsIdsOfNewTopic)
	items := oldTopic.Sort()
	for _, item := range items {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "")

		bcell := fmt.Sprintf("B%d", row)
		ccell := fmt.Sprintf("C%d", row)
		f.SetCellValue(sheet, bcell, item.Title)
		f.SetCellValue(sheet, ccell, item.URL)

		if item.Removed() {
			f.SetCellStyle(sheet, bcell, ccell, removedDiscussionSourceStyle)
		}

		row++
	}

	// 空一行并设置浅黄色背景
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), topicSpliterStyle)
	row++

	ftr.rowRemoveFromOld = row

	return nil
}
