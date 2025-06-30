package gitcodeissue

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/opensourceways/go-gitcode/openapi"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

type Config struct {
	Token string `json:"token"`
}

func NewClient(cfg *Config) *clientImpl {
	return &clientImpl{
		cli: openapi.NewAPIClientWithAuthorization([]byte(cfg.Token)),
	}
}

type clientImpl struct {
	cli *openapi.APIClient
}

type issueInfo struct {
	owner string
	repo  string
	num   string
}

func parseIssue(ds *domain.DiscussionSource) (issueInfo, error) {
	v := strings.Split(strings.TrimSpace(ds.URL), "/")
	n := len(v) - 1
	if n < 3 {
		return issueInfo{}, errors.New("invalid ds url")
	}

	return issueInfo{
		owner: v[n-3],
		repo:  v[n-2],
		num:   v[n],
	}, nil
}

func (impl *clientImpl) CountCommentedSolutons(
	ds *domain.DiscussionSource, parseSolutionComment func(string) string,
) ([]string, error) {
	issue, err := parseIssue(ds)
	if err != nil {
		return nil, err
	}

	urls := []string{}
	for page := 1; ; page++ {
		items, ok, err := impl.cli.Issues.ListIssueComments(
			context.Background(), issue.owner, issue.repo, issue.num, strconv.Itoa(page), "", "",
		)
		if err != nil {
			return nil, err
		}

		if !ok || len(items) == 0 {
			break
		}

		for i := range items {
			if v := parseSolutionComment(*items[i].Body); v != "" {
				urls = append(urls, v)
			}
		}
	}

	return urls, nil
}

func (impl *clientImpl) AddSolution(ds *domain.DiscussionSource, comment string) error {
	issue, err := parseIssue(ds)
	if err != nil {
		return err
	}

	v := openapi.IssueComment{
		Body: &comment,
	}

	_, _, err = impl.cli.Issues.CreateIssueComment(context.Background(), issue.owner, issue.repo, issue.num, &v)

	return err
}
