package app

import (
	"errors"
	"fmt"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain/repository"
)

type TopicSolutionAppService interface {
	Add(community string, cmd CmdToAddTopicSolution) error
}

func NewTopicSolutionAppService(
	repoSolution repository.RepoTopicSolution,
	repoHotTopic repository.RepoHotTopic,
) *topicSolutionAppService {
	return &topicSolutionAppService{
		repoSolution: repoSolution,
		repoHotTopic: repoHotTopic,
	}
}

type topicSolutionAppService struct {
	repoSolution repository.RepoTopicSolution
	repoHotTopic repository.RepoHotTopic
}

func (s *topicSolutionAppService) Add(community string, cmd CmdToAddTopicSolution) error {
	cmd.init()

	ids, err := s.handleOldTopics(community, cmd)
	if err != nil {
		return err
	}

	topicSolutions := make([]domain.TopicSolution, 0, len(cmd))

	for i := range cmd {
		items := cmd[i].DiscussionSources

		solutions := make([]domain.DiscussionSourceSolution, 0, len(items))
		for j := range items {
			resolved, unresolved := items[j].filterout()
			if len(resolved) == 0 || len(unresolved) == 0 {
				continue
			}

			relatedOnes := make([]int, len(unresolved))
			for k := range unresolved {
				relatedOnes[k] = unresolved[k].Id
			}

			solutions = append(solutions, domain.DiscussionSourceSolution{
				ResolvedOne: resolved[0].Id,
				RelatedOnes: relatedOnes,
			})
		}

		if len(solutions) == 0 {
			continue
		}

		topicSolutions = append(topicSolutions, domain.TopicSolution{
			TopicId:   ids[cmd[i].Title],
			Solutions: solutions,
		})
	}

	return s.repoSolution.Add(community, topicSolutions)
}

func (s *topicSolutionAppService) handleOldTopics(community string, cmd CmdToAddTopicSolution) (
	map[string]string, error,
) {
	oldTopics, err := s.repoHotTopic.FindOpenOnes(community)
	if err != nil {
		return nil, err
	}
	if len(oldTopics) == 0 {
		return nil, errors.New("no old topics")
	}

	oldTopicsSets := make(map[string]*domain.HotTopic, len(oldTopics))
	for i := range oldTopics {
		item := &oldTopics[i]
		oldTopicsSets[item.Title] = item
	}

	oldOnes := make(map[string]string, len(cmd))
	for i := range cmd {
		item := &cmd[i]

		if old, ok := oldTopicsSets[item.Title]; ok {
			v := parseRelationshipBetweenSets(old.GetDSSet(), item.getDSSet())
			if v != setsRelationSame && v != setsRelationLeftIncludesRight {
				return nil, errors.New("topic includes other discussions")
			}

			oldOnes[item.Title] = old.Id
		}
	}

	if n := len(oldTopics); len(oldOnes) != n {
		return nil, fmt.Errorf("the count of old topics is not matched, expect :%d, actual:%d", n, len(oldOnes))
	}

	return oldOnes, nil
}
