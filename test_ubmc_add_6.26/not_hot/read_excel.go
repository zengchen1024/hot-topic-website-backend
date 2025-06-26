package not_hot

import (
	"fmt"
	"strconv"
	"strings"

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

	rowStart = 3
)

func cellId(column string, row int) string {
	return fmt.Sprintf("%s%d", column, row)
}

type fileToReview struct {
	file *excelize.File
}

func NewfileToReview(file string) (*fileToReview, error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		return nil, err
	}

	return &fileToReview{
		file: f,
	}, nil
}

func (ftr *fileToReview) Exit() {
	if err := ftr.file.Close(); err != nil {
		fmt.Println(err)
	}
}

func (ftr *fileToReview) loadDS(row1 *int, sheet string, blankLineNum int) (dss []domain.DiscussionSourceInfo, err error) {
	row := *row1

	vb, vc, vd := "", "", ""
	num := 0
	for {
		vb, err = ftr.file.GetCellValue(sheet, cellId(columnB, row))
		if err != nil {
			return
		}

		vc, err = ftr.file.GetCellValue(sheet, cellId(columnC, row))
		if err != nil {
			return
		}

		vd, err = ftr.file.GetCellValue(sheet, cellId(columnD, row))
		if err != nil {
			return
		}

		if vb == "" && vc == "" {
			num++
			if num >= blankLineNum {
				row++ // the line of next topic
				break
			}
		} else {
			dsId, err1 := parseDSId(vd)
			if err1 != nil {
				err = fmt.Errorf("parse %s to ds id failed, err:%s", vd, err1.Error())
				return
			}

			dss = append(dss, domain.DiscussionSourceInfo{
				Id:    dsId,
				URL:   vc,
				Title: vb,
			})

			num = 0
		}

		row++
	}

	*row1 = row
	return
}

func (ftr *fileToReview) loadNotHotTopic(row1 *int, sheet string, blankLineNum int) (topic domain.NotHotTopic, err error) {
	row := *row1

	topic.Title, err = ftr.file.GetCellValue(sheet, cellId(columnB, row))
	if err != nil {
		return
	}
	row += 2

	*row1 = row

	topic.DiscussionSources, err = ftr.loadDS(row1, sheet, blankLineNum)
	if err != nil {
		return
	}

	return
}

func (ftr *fileToReview) getTopicNum(sheet string) (int, error) {
	str, err := ftr.file.GetCellValue(sheet, cellId(columnA, 1))
	if err != nil {
		return 0, err
	}
	v := strings.Split(str, "：")
	if len(v) != 2 {
		return 0, fmt.Errorf("invalid topic num desc, v=%v", v)
	}
	num, err := strconv.Atoi(v[1])
	if err != nil {
		return 0, fmt.Errorf("can't convert topic num from:%s, err:%s", str, err.Error())
	}

	return num, nil
}

type loadTopicFunc func(row1 *int, sheet string, blankLineNum int) (topic domain.NotHotTopic, err error)

func (ftr *fileToReview) readSheet(sheet string, blankLineNum int, loadTopic loadTopicFunc) ([]domain.NotHotTopic, error) {
	num, err := ftr.getTopicNum(sheet)
	if err != nil {
		return nil, err
	}

	row := 3

	r := []domain.NotHotTopic{}
	for i := 0; i < num; i++ {
		topic, err := loadTopic(&row, sheet, blankLineNum)
		if err != nil {
			return nil, err
		}

		fmt.Printf("now row = %d\n", row)

		r = append(r, topic)
	}

	return r, nil
}

func (ftr *fileToReview) ReadOtherTopics() ([]domain.NotHotTopic, error) {
	r := []domain.NotHotTopic{}

	sheets := []string{sheetNewTopics, sheetUnchangedTopics, sheetAppendToOld}
	for _, sheet := range sheets {
		v, err := ftr.readSheet(sheet, 1, ftr.loadNotHotTopic)
		if err != nil {
			return nil, err
		}

		fmt.Printf("read from sheet:%s, got %d topics\n", sheet, len(v))

		r = append(r, v...)
	}

	v, err := ftr.readSheet(sheetMultiIntersects, 2, ftr.loadNotHotTopic)
	if err != nil {
		return nil, err
	}

	fmt.Printf("read from sheet:%s, got %d topics\n", sheetMultiIntersects, len(v))

	r = append(r, v...)

	return r, nil
}
