package forum

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

type reqToAddComment struct {
	TopicId int    `json:"topic_id"`
	Comment string `json:"raw"`
}

type responseOfGettingPost struct {
	PostStream struct {
		Posts []struct {
			Cooked string `json:"cooked"`
		} `json:"posts"`
	} `json:"post_stream"`
}

func (resp *responseOfGettingPost) parse(sc solutionComment) []string {
	urls := []string{}

	items := resp.PostStream.Posts
	for i := range items {
		if v := sc.ParseURL(items[i].Cooked); v != "" {
			urls = append(urls, v)
		}
	}

	return urls
}

type Config struct {
	User     string `json:"user"     required:"true"`
	ApiKey   string `json:"api_key"  required:"true"`
	Endpoint string `json:"endpoint" required:"true"`
}

type solutionComment interface {
	ParseURL(comment string) string
}

func NewClient(cfg *Config, sc solutionComment) *clientImpl {
	return &clientImpl{
		cli:             NewHttpClient(3),
		solutionComment: sc,

		user:          cfg.User,
		apiKey:        cfg.ApiKey,
		getPostURL:    fmt.Sprintf("%s/t/", cfg.Endpoint),
		addCommentURL: fmt.Sprintf("%s/posts.json", cfg.Endpoint),
	}
}

type clientImpl struct {
	cli HttpClient
	solutionComment

	user          string
	apiKey        string
	getPostURL    string
	addCommentURL string
}

func (impl *clientImpl) CountCommentedSolutons(ds *domain.DiscussionSource) ([]string, error) {
	req, err := http.NewRequest(
		http.MethodGet, fmt.Sprintf("%s%s.json", impl.getPostURL, ds.SourceId), nil,
	)
	if err != nil {
		return nil, err
	}

	impl.setHeaderForReq(req)

	resp := responseOfGettingPost{}
	if _, err := impl.cli.ForwardTo(req, &resp); err != nil {
		return nil, err
	}

	return resp.parse(impl.solutionComment), nil
}

func (impl *clientImpl) AddSolution(ds *domain.DiscussionSource, comment string) error {
	topicId, err := strconv.Atoi(ds.SourceId)
	if err != nil {
		return err
	}

	body, err := JsonMarshal(reqToAddComment{
		TopicId: topicId,
		Comment: comment,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, impl.addCommentURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	impl.setHeaderForReq(req)

	_, err = impl.cli.ForwardTo(req, nil)

	return err
}

func (impl *clientImpl) setHeaderForReq(req *http.Request) {
	req.Header.Set("Api-Key", impl.apiKey)
	req.Header.Set("Api-Username", impl.user)
	req.Header.Set("content-type", "application/json")
}
