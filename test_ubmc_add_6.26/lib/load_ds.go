package lib

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

var (
	idRegex         = regexp.MustCompile(`"id":(\d+)`)
	sourceTypeRegex = regexp.MustCompile(`"source_type":([^,}]+)`)
	sourceIdRegex   = regexp.MustCompile(`"source_id":([^,}]+)`)
	createdAtRegex  = regexp.MustCompile(`"created_at":([^,}]+)`)
)

func genDiscussionSource(str string) (ds domain.DiscussionSource, err error) {
	matches := idRegex.FindStringSubmatch(str)
	if len(matches) < 1 {
		err = errors.New("no id")

		return
	}
	if ds.Id, err = strconv.Atoi(matches[1]); err != nil {
		return
	}

	matches = sourceTypeRegex.FindStringSubmatch(str)
	if len(matches) < 1 {
		err = errors.New("no source type")

		return
	}
	ds.Type = matches[1]

	matches = sourceIdRegex.FindStringSubmatch(str)
	if len(matches) < 1 {
		err = errors.New("no source id")

		return
	}
	ds.SourceId = matches[1]

	matches = createdAtRegex.FindStringSubmatch(str)
	if len(matches) < 1 {
		err = errors.New("no created at")

		return
	}
	ds.CreatedAt = matches[1]

	return
}
