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

	columnA = "A"
	columnB = "B"
	columnC = "C"
	columnD = "D"
	columnE = "E"

	rowStart = 3
)

func cellId(column string, row int) string {
	return fmt.Sprintf("%s%d", column, row)
}

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

		rowNewTopic:        rowStart, // must start with 3, the first line is used for note, second line is blank
		rowAppendToOld:     rowStart,
		rowRemoveFromOld:   rowStart,
		rowUnchangedTopic:  rowStart,
		rowMultiIntersects: rowStart,
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

func (ftr *fileToReview) saveTopic(topicTitle string, row1 *int, sheet string, saveDS func(*int) (string, error)) (err error) {
	f := ftr.file
	row := *row1

	setCell := func(cell, value string) error {
		return f.SetCellValue(sheet, cellId(cell, row), value)
	}

	if err = setCell(columnA, "话题描述"); err != nil {
		return
	}
	if err = setCell(columnB, topicTitle); err != nil {
		return
	}
	row++

	if err = setCell(columnA, "相关讨论源（title&url）"); err != nil {
		return
	}
	row++

	// save discussion sources
	column, err := saveDS(&row)
	if err != nil {
		return
	}

	err = f.SetCellStyle(sheet, cellId(columnA, row), cellId(column, row), ftr.topicSpliterStyle)
	if err != nil {
		return
	}
	row++

	*row1 = row

	return
}

type dsInfo struct {
	*domain.DiscussionSource

	title  string
	Closed bool
}

func (ds *dsInfo) metaData() string {
	return fmt.Sprintf(
		"{\"id\":%d, \"source_type\":\"%s\", \"source_id\":\"%s\", \"created_at\":\"%s\"}",
		ds.Id, ds.Type, ds.SourceId, ds.CreatedAt,
	)
}

func (ftr *fileToReview) saveOneDS(ds *dsInfo, row int, sheet string, color bool, colorValue int) (err error) {
	f := ftr.file

	bcell := cellId(columnB, row)
	if err = f.SetCellValue(sheet, bcell, ds.title); err != nil {
		return
	}

	if err = f.SetCellValue(sheet, cellId(columnC, row), ds.URL); err != nil {
		return
	}

	cell := cellId(columnD, row)
	if err = f.SetCellValue(sheet, cell, ds.metaData()); err != nil {
		return
	}

	if ds.Closed {
		cell := cellId(columnE, row)
		if err = f.SetCellValue(sheet, cell, "Closed"); err != nil {
			return
		}
	}

	if color {
		err = f.SetCellStyle(sheet, bcell, cell, colorValue)
	}

	return
}

func (ftr *fileToReview) saveLastHotTopics(oldTopics []domain.HotTopic, current map[string]*OptionalTopic) error {
	if len(oldTopics) == 0 {
		return nil
	}

	f := ftr.file

	if err := f.SetCellValue(sheetLastTopics, "A1", "话题总数："+strconv.Itoa(len(oldTopics))); err != nil {
		return err
	}

	row := rowStart

	for i := range oldTopics {
		old := &oldTopics[i]

		item, ok := current[old.Title]
		if !ok {
			return errors.New("can't find topic from current ones")
		}

		item.updateAppended(old.GetDSSet())

		if err := ftr.saveHotTopic(item, old, &row, sheetLastTopics); err != nil {
			return err
		}
	}

	return nil
}

func (ftr *fileToReview) saveHotTopic(topic *OptionalTopic, oldTopic *domain.HotTopic, row1 *int, sheet string) error {
	setSubDS := func(row2 *int, infos DiscussionSourceInfos) (err error) {
		row := *row2

		ds := dsInfo{}
		items := infos.sort()
		for _, item := range items {
			ds = dsInfo{
				DiscussionSource: &item.DiscussionSource,
				title:            item.Title,
				Closed:           item.Closed,
			}

			err = ftr.saveOneDS(&ds, row, sheet, item.Appended, ftr.appendedDiscussionSourceStyle)
			if err != nil {
				return
			}

			row++
		}

		*row2 = row

		return
	}

	setDS := func(row2 *int) (column string, err error) {
		column = columnC

		for _, items := range topic.DiscussionSources {
			if err = setSubDS(row2, items); err != nil {
				return
			}

			(*row2)++
		}

		return
	}

	setCell := func(cell, value string) error {
		return ftr.file.SetCellValue(sheet, cellId(cell, *row1), value)
	}

	if err := setCell(columnA, "Id"); err != nil {
		return err
	}
	if err := setCell(columnB, oldTopic.Id); err != nil {
		return err
	}
	(*row1)++

	if err := setCell(columnA, "顺序"); err != nil {
		return err
	}
	if err := setCell(columnB, strconv.Itoa(oldTopic.Order)); err != nil {
		return err
	}
	(*row1)++

	return ftr.saveTopic(topic.Title, row1, sheet, setDS)
}

func (ftr *fileToReview) saveAppendedTopic(topic *OptionalTopic, row1 *int, sheet string) error {
	setDS := func(row2 *int) (column string, err error) {
		row := *row2
		column = columnC

		ds := dsInfo{}
		items := topic.sort()
		for _, item := range items {
			ds = dsInfo{
				DiscussionSource: &item.DiscussionSource,
				title:            item.Title,
				Closed:           item.Closed,
			}

			err = ftr.saveOneDS(&ds, row, sheet, item.Appended, ftr.appendedDiscussionSourceStyle)
			if err != nil {
				return
			}

			row++
		}

		*row2 = row

		return
	}

	return ftr.saveTopic(topic.Title, row1, sheet, setDS)
}

func (ftr *fileToReview) saveNewTopic(topic *OptionalTopic) error {
	ftr.totalNewTopic++

	return ftr.saveTopicDirectly(topic, &ftr.rowNewTopic, sheetNewTopics)
}

func (ftr *fileToReview) saveUnchangedTopic(topic *OptionalTopic) error {
	ftr.totalUnchangedTopic++

	return ftr.saveTopicDirectly(topic, &ftr.rowUnchangedTopic, sheetUnchangedTopics)
}

func (ftr *fileToReview) saveTopicDirectly(topic *OptionalTopic, row1 *int, sheet string) error {
	setDS := func(row2 *int) (column string, err error) {
		row := *row2
		column = columnC

		ds := dsInfo{}
		for i := range topic.discussionSources {
			item := topic.discussionSources[i]

			ds = dsInfo{
				DiscussionSource: &item.DiscussionSource,
				title:            item.Title,
				Closed:           item.Closed,
			}
			if err = ftr.saveOneDS(&ds, row, sheet, false, 0); err != nil {
				return
			}

			row++
		}

		*row2 = row

		return
	}

	return ftr.saveTopic(topic.Title, row1, sheet, setDS)
}

func (ftr *fileToReview) saveTopicThatAppendToOld(topic *OptionalTopic, dsIdsOfOldTopic map[int]bool) error {
	ftr.totalAppendToOld++

	topic.updateAppended(dsIdsOfOldTopic)

	return ftr.saveAppendedTopic(topic, &ftr.rowAppendToOld, sheetAppendToOld)
}

func (ftr *fileToReview) saveTopicThatRemoveFromOld(oldTopic *domain.NotHotTopic, dsIdsOfNewTopic map[int]bool) error {
	ftr.totalRemoveFromOld++

	sheet := sheetRemoveFromOld

	setDS := func(row2 *int) (column string, err error) {
		row := *row2
		column = columnC

		oldTopic.UpdateRemoved(dsIdsOfNewTopic)

		ds := dsInfo{}
		source := domain.DiscussionSource{}
		items := oldTopic.Sort()
		for _, item := range items {
			source.Id = item.Id
			source.URL = item.URL

			ds = dsInfo{
				title:            item.Title,
				DiscussionSource: &source,
			}
			err = ftr.saveOneDS(&ds, row, sheet, item.Removed(), ftr.removedDiscussionSourceStyle)
			if err != nil {
				return
			}

			row++
		}

		*row2 = row

		return
	}

	return ftr.saveTopic(oldTopic.Title, &ftr.rowRemoveFromOld, sheet, setDS)
}

func (ftr *fileToReview) saveTopicThatIntersectWithMultiOlds(topic *OptionalTopic, oldIds []int, oldTopics []domain.NotHotTopic) error {
	ftr.totalMultiIntersects++

	setDS := func(row2 *int) (column string, err error) {
		column = columnC

		newSets := topic.getDSSet()
		for _, i := range oldIds {
			newSets = getDifferentiation(newSets, oldTopics[i].GetDSSet())
		}
		if len(newSets) > 0 {
			if err = ftr.saveIntersectedDS(topic, newSets, row2, true); err != nil {
				return
			}
		}

		newSets = topic.getDSSet()
		for _, i := range oldIds {
			item := &oldTopics[i]

			err = ftr.file.SetCellValue(sheetMultiIntersects, cellId(columnA, *row2), "以下讨论源属于旧话题："+item.Title)
			if err != nil {
				return
			}

			common := getIntersection(item.GetDSSet(), newSets)
			if err = ftr.saveIntersectedDS(topic, common, row2, false); err != nil {
				return
			}
		}

		return
	}

	return ftr.saveTopic(topic.Title, &ftr.rowMultiIntersects, sheetMultiIntersects, setDS)
}

func (ftr *fileToReview) saveIntersectedDS(topic *OptionalTopic, topicIds map[int]bool, row1 *int, color bool) (err error) {
	row := *row1
	sheet := sheetMultiIntersects

	ds := dsInfo{}
	for i := range topic.discussionSources {
		item := topic.discussionSources[i]

		if !topicIds[item.Id] {
			continue
		}

		ds = dsInfo{
			DiscussionSource: &item.DiscussionSource,
			title:            item.Title,
			Closed:           item.Closed,
		}
		err = ftr.saveOneDS(&ds, row, sheet, color, ftr.appendedDiscussionSourceStyle)
		if err != nil {
			return
		}

		row++
	}

	*row1 = row + 1

	return nil
}
