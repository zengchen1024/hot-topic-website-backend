package app

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

const (
	sheetNewTopics       = "本次新发现的话题"
	sheetLastTopics      = "上次的热点话题"
	sheetAppendToOld     = "上次未入选的话题且有新讨论源加入"
	sheetRemoveFromOld   = "上次未入选的话题且有讨论源移除"
	sheetUnchangedTopics = "上次未入选的话题且未发生变化"
	sheetMultiIntersects = "本次的话题与上次未入选的话题发生讨论源交叉"
)

type fileToReview struct {
	file     *excelize.File
	filePath string

	rowNewTopic        int
	rowAppendToOld     int
	rowRemoveFromOld   int
	rowUnchangedTopic  int
	rowMultiIntersects int

	totalNewTopic        int
	totalAppendToOld     int
	totalRemoveFromOld   int
	totalUnchangedTopic  int
	totalMultiIntersects int

	topicSpliterStyle             int
	removedDiscussionSourceStyle  int
	appendedDiscussionSourceStyle int
}

func newfileToReview(file string) (*fileToReview, error) {
	r := &fileToReview{
		file:     excelize.NewFile(),
		filePath: file,

		rowNewTopic:        3, // must start with 3, the first line is used for note, second line is blank
		rowAppendToOld:     3,
		rowRemoveFromOld:   3,
		rowUnchangedTopic:  3,
		rowMultiIntersects: 3,
	}
	err := r.doInit()

	return r, err
}

func (ftr *fileToReview) doInit() (err error) {
	if err = ftr.initSheets(); err != nil {
		return
	}

	if err = ftr.initStyle(); err != nil {
		return
	}

	return ftr.initLegend()
}

func (ftr *fileToReview) initSheets() error {
	all := []string{
		sheetLastTopics, // last hot topic first
		sheetNewTopics,
		sheetAppendToOld,
		sheetRemoveFromOld,
		sheetUnchangedTopics,
		sheetMultiIntersects,
	}

	for _, name := range all {
		if _, err := ftr.file.NewSheet(name); err != nil {
			return err
		}
	}

	return ftr.file.DeleteSheet("Sheet1")
}

func (ftr *fileToReview) initStyle() (err error) {
	newStyle := func(color string) (int, error) {
		return ftr.file.NewStyle(&excelize.Style{
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{color},
				Pattern: 1,
			},
		})
	}

	items := []struct {
		style *int
		value string
	}{
		{&ftr.topicSpliterStyle, "92D050"},
		{&ftr.removedDiscussionSourceStyle, "FFC7CE"},
		{&ftr.appendedDiscussionSourceStyle, "FFFF43"},
	}

	for i := range items {
		item := &items[i]

		if *item.style, err = newStyle(item.value); err != nil {
			return
		}
	}

	return
}

type legend struct {
	cell  string
	desc  string
	style int
}

func (ftr *fileToReview) initLegend() error {
	h := func(sheet string, lg legend) error {
		cell := lg.cell

		if err := ftr.file.SetCellValue(sheet, cell, lg.desc); err != nil {
			return err
		}

		return ftr.file.SetCellStyle(sheet, cell, cell, lg.style)
	}

	m := map[string]legend{
		sheetLastTopics: legend{
			cell:  "B1",
			desc:  "图例：新加入的讨论源",
			style: ftr.appendedDiscussionSourceStyle,
		},

		sheetAppendToOld: legend{
			cell:  "B1",
			desc:  "图例：新加入的讨论源",
			style: ftr.appendedDiscussionSourceStyle,
		},

		sheetRemoveFromOld: legend{
			cell:  "B1",
			desc:  "图例：移除的讨论源",
			style: ftr.removedDiscussionSourceStyle,
		},

		sheetMultiIntersects: legend{
			cell:  "B1",
			desc:  "图例：新加入的讨论源",
			style: ftr.appendedDiscussionSourceStyle,
		},
	}
	for k, lg := range m {
		if err := h(k, lg); err != nil {
			return err
		}
	}

	return nil
}

func (ftr *fileToReview) saveToFile() error {
	h := func(sheet string, num int) error {
		return ftr.file.SetCellValue(sheet, "A1", "话题总数："+strconv.Itoa(num))
	}

	m := map[string]int{
		sheetNewTopics:       ftr.totalNewTopic,
		sheetAppendToOld:     ftr.totalAppendToOld,
		sheetRemoveFromOld:   ftr.totalRemoveFromOld,
		sheetUnchangedTopics: ftr.totalUnchangedTopic,
		sheetMultiIntersects: ftr.totalMultiIntersects,
	}
	for k, v := range m {
		if err := h(k, v); err != nil {
			return err
		}
	}

	return ftr.file.SaveAs(ftr.filePath)
}

func (ftr *fileToReview) saveLastTopics(oldTopics []domain.HotTopic, current map[int]*OptionalTopic) error {
	if len(oldTopics) == 0 {
		return nil
	}

	f := ftr.file

	if err := f.SetCellValue(sheetLastTopics, "A1", "话题总数："+strconv.Itoa(len(oldTopics))); err != nil {
		return err
	}

	row := 3

	for i := range oldTopics {
		old := &oldTopics[i]

		item, ok := current[old.Order]
		if !ok {
			return errors.New("can't find topic from current ones")
		}

		item.updateAppended(old.GetDSSet())

		//fmt.Printf("save hot topic, row:%d, title: %s\n", row, item.Title)

		if err := ftr.saveAppendedTopic(item, &row, sheetLastTopics); err != nil {
			return err
		}
	}

	return nil
}

func (ftr *fileToReview) saveAppendedTopic(topic *OptionalTopic, row1 *int, sheet string) (err error) {
	f := ftr.file
	row := *row1

	// 话题描述
	if err = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "话题描述"); err != nil {
		return
	}

	if err = f.SetCellValue(sheet, fmt.Sprintf("B%d", row), topic.Title); err != nil {
		return
	}
	row++

	// 讨论源（title & url）
	if err = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "相关讨论源（title&url）"); err != nil {
		return
	}
	row++

	// 写入每个讨论项
	items := topic.sort()
	for _, item := range items {
		if err = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), ""); err != nil {
			return
		}

		bcell := fmt.Sprintf("B%d", row)
		ccell := fmt.Sprintf("C%d", row)
		if err = f.SetCellValue(sheet, bcell, item.Title); err != nil {
			return
		}
		if err = f.SetCellValue(sheet, ccell, item.URL); err != nil {
			return
		}

		if item.appended {
			fmt.Println("set appended style")
			if err = f.SetCellStyle(sheet, bcell, ccell, ftr.appendedDiscussionSourceStyle); err != nil {
				return
			}
		}

		row++
	}

	// 空一行并设置浅黄色背景
	err = f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), ftr.topicSpliterStyle)
	row++

	*row1 = row

	return
}

func (ftr *fileToReview) saveNewTopic(topic *OptionalTopic) error {
	ftr.totalNewTopic++

	return ftr.saveTopic(topic, &ftr.rowNewTopic, sheetNewTopics)
}

func (ftr *fileToReview) saveUnchangedTopic(topic *OptionalTopic) error {
	ftr.totalUnchangedTopic++

	return ftr.saveTopic(topic, &ftr.rowUnchangedTopic, sheetUnchangedTopics)
}

func (ftr *fileToReview) saveTopic(topic *OptionalTopic, row1 *int, sheet string) error {
	f := ftr.file
	row := *row1

	// 话题描述
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "话题描述")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), topic.Title)
	row++

	// 讨论源（title & url）
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "相关讨论源（title&url）")
	row++

	// 写入每个讨论项
	for i := range topic.DiscussionSources {
		item := &topic.DiscussionSources[i]

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "")

		bcell := fmt.Sprintf("B%d", row)
		ccell := fmt.Sprintf("C%d", row)
		f.SetCellValue(sheet, bcell, item.Title)
		f.SetCellValue(sheet, ccell, item.URL)

		row++
	}

	// 空一行并设置浅黄色背景
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), ftr.topicSpliterStyle)
	row++

	*row1 = row

	return nil
}

func (ftr *fileToReview) saveTopicThatAppendToOld(topic *OptionalTopic, dsIdsOfOldTopic map[int]bool) error {
	ftr.totalAppendToOld++

	topic.updateAppended(dsIdsOfOldTopic)

	return ftr.saveAppendedTopic(topic, &ftr.rowAppendToOld, sheetAppendToOld)
}

func (ftr *fileToReview) saveTopicThatRemoveFromOld(oldTopic *domain.NotHotTopic, dsIdsOfNewTopic map[int]bool) error {
	ftr.totalRemoveFromOld++

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
			f.SetCellStyle(sheet, bcell, ccell, ftr.removedDiscussionSourceStyle)
		}

		row++
	}

	// 空一行并设置浅黄色背景
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), ftr.topicSpliterStyle)
	row++

	ftr.rowRemoveFromOld = row

	return nil
}

func (ftr *fileToReview) saveTopicThatIntersectWithMultiOlds(topic *OptionalTopic, oldIds []int, oldTopics []domain.NotHotTopic) error {
	newSets := topic.getDSSet()

	ftr.totalMultiIntersects++

	f := ftr.file
	row := ftr.rowMultiIntersects
	sheet := sheetMultiIntersects

	// 话题描述
	if err := f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "话题描述"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, fmt.Sprintf("B%d", row), topic.Title); err != nil {
		return err
	}
	row++

	// 讨论源（title & url）
	if err := f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "相关讨论源（title&url）"); err != nil {
		return err
	}
	row++

	//
	for _, i := range oldIds {
		if err := ftr.helper(topic, newSets, &oldTopics[i], &row); err != nil {
			return err
		}
	}

	if err := ftr.handleLast(topic, newSets, &row); err != nil {
		return err
	}

	// 空一行并设置浅黄色背景
	if err := f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), ftr.topicSpliterStyle); err != nil {
		return err
	}
	row++

	ftr.rowMultiIntersects = row

	return nil
}

func (ftr *fileToReview) helper(topic *OptionalTopic, topicIds map[int]bool, oldTopic *domain.NotHotTopic, row1 *int) error {
	row := *row1
	sheet := sheetMultiIntersects
	f := ftr.file

	common := getIntersection(topicIds, oldTopic.GetDSSet())

	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), oldTopic.Title)

	for i := range topic.DiscussionSources {
		item := topic.DiscussionSources[i]

		if !common[item.Id] {
			continue
		}

		delete(topicIds, item.Id)

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "")

		bcell := fmt.Sprintf("B%d", row)
		ccell := fmt.Sprintf("C%d", row)
		f.SetCellValue(sheet, bcell, item.Title)
		f.SetCellValue(sheet, ccell, item.URL)
		row++
	}

	*row1 = row + 1

	return nil
}

func (ftr *fileToReview) handleLast(topic *OptionalTopic, topicIds map[int]bool, row1 *int) error {
	row := *row1
	sheet := sheetMultiIntersects
	f := ftr.file

	for i := range topic.DiscussionSources {
		item := topic.DiscussionSources[i]

		if !topicIds[item.Id] {
			continue
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "")

		bcell := fmt.Sprintf("B%d", row)
		ccell := fmt.Sprintf("C%d", row)
		f.SetCellValue(sheet, bcell, item.Title)
		f.SetCellValue(sheet, ccell, item.URL)

		f.SetCellStyle(sheet, bcell, ccell, ftr.appendedDiscussionSourceStyle)

		row++
	}

	*row1 = row

	return nil
}
