package lib

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

func (ftr *fileToReview) loadDS(row1 *int, sheet string, blankLineNum int) (dss []domain.DiscussionSource, err error) {
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
			/*
				ds := domain.DiscussionSource{}
				if err = json.Unmarshal([]byte(vd), &ds); err != nil {
					err = fmt.Errorf("unmarshal %s to DiscussionSource failed, err:%s", vd, err.Error())
					return
				}
			*/

			ds, err1 := genDiscussionSource(vd)
			if err1 != nil {
				err = fmt.Errorf("unmarshal %s to DiscussionSource failed, err:%s", vd, err1.Error())
				return
			}

			ds.URL = vc

			dss = append(dss, ds)

			num = 0
		}

		row++
	}

	*row1 = row
	return
}

func (ftr *fileToReview) loadHotTopic(row1 *int, sheet string, blankLineNum int) (topic domain.HotTopic, err error) {
	row := *row1

	topic.Id, err = ftr.file.GetCellValue(sheet, cellId(columnB, row))
	if err != nil {
		return
	}
	row++

	order, err := ftr.file.GetCellValue(sheet, cellId(columnB, row))
	if err != nil {
		return
	}
	topic.Order, err = strconv.Atoi(order)
	if err != nil {
		err = fmt.Errorf("convert %s to int failed, row:%d, err:%s", order, row, err.Error())
		return
	}
	row++

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

func (ftr *fileToReview) loadNotHotTopic(row1 *int, sheet string, blankLineNum int) (topic domain.HotTopic, err error) {
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

type loadTopicFunc func(row1 *int, sheet string, blankLineNum int) (topic domain.HotTopic, err error)

func (ftr *fileToReview) readSheet(sheet string, blankLineNum int, loadTopic loadTopicFunc) ([]domain.HotTopic, error) {
	num, err := ftr.getTopicNum(sheet)
	if err != nil {
		return nil, err
	}

	row := 3

	r := []domain.HotTopic{}
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

func (ftr *fileToReview) ReadLastHotTopics() ([]domain.HotTopic, error) {
	return ftr.readSheet(sheetLastTopics, 2, ftr.loadHotTopic)
}

func (ftr *fileToReview) ReadOtherTopics() ([]domain.HotTopic, error) {
	r := []domain.HotTopic{}

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

func (ftr *fileToReview) ReadIntersectTopics() ([]domain.HotTopic, error) {
	return ftr.readSheet(sheetMultiIntersects, 2, ftr.loadNotHotTopic)
}
