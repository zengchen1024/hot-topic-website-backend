package watch

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain/repository"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/watch/forum"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/watch/gitcodeissue"
)

const commentKey = "的信息可能对你解决本问题有所帮助，请参考，谢谢！"

type platformClient interface {
	CountCommentedSolutons(ds *domain.DiscussionSource, key string) (int, error)
	AddSolution(ds *domain.DiscussionSource, comment string) error
}

func genSolutionComment(solution *domain.DiscussionSource) string {
	return fmt.Sprintf("你好！这个[连接](%s)%s", solution.URL, commentKey)
}

func newClients(cfg *Config) clients {
	cli := clients{}

	for i := range cfg.Forums {
		item := cfg.Forums[i]

		cli[cli.key(item.Community, item.typeDesc())] = forum.NewClient(&item.Detail)
	}

	for i := range cfg.GitCodes {
		item := cfg.GitCodes[i]

		cli[cli.key(item.Community, item.typeDesc())] = gitcodeissue.NewClient(&item.Detail)
	}

	return cli
}

// key is {community}_{discussion source type}
type clients map[string]platformClient

func (cli clients) key(community, t string) string {
	return strings.ToLower(fmt.Sprintf("%s_%s", community, t))
}

func (cli clients) get(community string, ds *domain.DiscussionSource) (platformClient, error) {
	v := cli[cli.key(community, ds.Type)]
	if v == nil {
		return nil, fmt.Errorf("no client for %s and %s", community, ds.Type)
	}

	return v, nil
}

// doneCounter
type doneCounter struct {
	num       int
	expiredAt int64
}

func (c *doneCounter) add() {
	c.num++
	c.setExpired()
}

func (c *doneCounter) canDo() bool {
	return c.num < 3
}

func (c *doneCounter) isExpired(now int64) bool {
	return c.expiredAt > 0 && c.expiredAt < now
}

func (c *doneCounter) setExpired() {
	if !c.canDo() {
		c.expiredAt = expiry(24 * 5)
	}
}

func newdoneCounter(num int) doneCounter {
	v := doneCounter{num: num}
	v.setExpired()

	return v
}

func expiry(expiry int64) int64 {
	return time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
}

// key is {community}_{discussion id}; value is the num that solution was commented on it
type doneCache map[string]*doneCounter

func (c doneCache) key(community string, ds *domain.DiscussionSource) string {
	return fmt.Sprint("%s_%d", community, ds.Id)
}

func (c doneCache) refresh() {
	now := time.Now().Unix()

	for k, v := range c {
		if v.isExpired(now) {
			delete(c, k)
		}
	}
}

func (c doneCache) get(cli platformClient, community string, ds *domain.DiscussionSource) (*doneCounter, error) {
	k := c.key(community, ds)

	if counter, ok := c[k]; ok {
		return counter, nil
	}

	n, err := cli.CountCommentedSolutons(ds, commentKey)
	if err != nil {
		return nil, err
	}

	v := newdoneCounter(n)
	c[k] = &v

	return &v, nil
}

func newtopicSolutionHandler(repo repository.RepoHotTopic, cfg *Config) *topicSolutionHandler {
	return &topicSolutionHandler{
		repo:    repo,
		cache:   doneCache{},
		clients: newClients(cfg),
	}
}

// topicSolutionHandler
type topicSolutionHandler struct {
	repo    repository.RepoHotTopic
	cache   doneCache
	clients clients
}

func (h *topicSolutionHandler) handle(solution *repository.TopicSolutions, needStop func() bool) {
	for i := range solution.TopicSolutions {
		if needStop() {
			return
		}

		ts := &solution.TopicSolutions[i]

		topic, err := h.repo.Find(solution.Community, ts.TopicId)
		if err != nil {
			logrus.Warn("find the topic(%s) failed, err:%s", ts.TopicId, err.Error)

			continue
		}

		h.handleTopicSolution(solution.Community, &topic, ts.Solutions, needStop)
	}

	h.cache.refresh()
}

func (h *topicSolutionHandler) handleTopicSolution(
	community string, topic *domain.HotTopic,
	solutions []domain.DiscussionSourceSolution,
	needStop func() bool,
) {
	for i := range solutions {
		if needStop() {
			return
		}

		item := &solutions[i]

		resolvedOne := h.getDiscussionSource(topic, item.ResolvedOne)
		if resolvedOne == nil {
			continue
		}

		for _, dsId := range item.RelatedOnes {
			ds := h.getDiscussionSource(topic, dsId)
			if ds == nil {
				continue
			}

			if err := h.handleDiscussionSourceSolution(community, resolvedOne, ds); err != nil {
				logrus.Errorf(
					"handle solution(%s) for discussion source() failed, err:%s",
					resolvedOne.URL, ds.URL, err.Error(),
				)
			}

			if needStop() {
				return
			}
		}
	}
}

func (h *topicSolutionHandler) getDiscussionSource(topic *domain.HotTopic, dsId int) *domain.DiscussionSource {
	ds := topic.GetDiscussionSource(dsId)
	if ds == nil {
		logrus.Errorf(
			"can't find the DiscussionSource(%d) from topic:%s:%s",
			dsId, topic.Id, topic.Title,
		)
	}

	return ds
}

func (h *topicSolutionHandler) handleDiscussionSourceSolution(
	community string, resolvedOne, ds *domain.DiscussionSource,
) error {
	cli, err := h.clients.get(community, ds)
	if err != nil {
		return err
	}

	counter, err := h.cache.get(cli, community, ds)
	if err != nil {
		return err
	}

	if !counter.canDo() {
		return nil
	}

	if err = cli.AddSolution(ds, genSolutionComment(resolvedOne)); err != nil {
		return err
	}

	counter.add()

	return nil
}
