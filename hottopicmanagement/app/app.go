package app

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain/repository"
	"github.com/opensourceways/hot-topic-website-backend/utils"
)

type AppService interface {
	ToReview(string, CmdToUploadOptionalTopics) error
	GenReport(string) ([]byte, error)
}

func NewAppService(
	config *Config,
	repoHotTopic repository.RepoHotTopic,
	repoNotHotTopic repository.RepoNotHotTopic,
	repoTopicReport repository.RepoTopicReport,
) *appService {
	return &appService{
		filePath:        config.FilePath,
		repoHotTopic:    repoHotTopic,
		repoNotHotTopic: repoNotHotTopic,
		repoTopicReport: repoTopicReport,
	}
}

type appService struct {
	filePath        string
	repoHotTopic    repository.RepoHotTopic
	repoNotHotTopic repository.RepoNotHotTopic
	repoTopicReport repository.RepoTopicReport
}

func (s *appService) reviewFile(community string) string {
	return filepath.Join(s.filePath, fmt.Sprintf("%s_%s.xlsx", community, utils.Date()))
}

func (s *appService) ToReview(community string, cmd CmdToUploadOptionalTopics) error {
	cmd.init()

	file, err := newfileToReview(s.reviewFile(community))
	if err != nil {
		return fmt.Errorf("new excel failed, err:%s", err.Error())
	}

	newOnes, err := s.handleOldTopics(community, cmd, file)
	if err != nil {
		return err
	}

	if err := s.handleNewOptionalTopics(community, newOnes, file); err != nil {
		return err
	}

	if err := file.saveToFile(); err != nil {
		return err
	}

	// TODO send email

	return nil
}

func (s *appService) handleOldTopics(community string, cmd CmdToUploadOptionalTopics, file *fileToReview) ([]*OptionalTopic, error) {
	oldTopics, err := s.repoHotTopic.FindOpenOnes(community)
	if err != nil {
		return nil, err
	}
	if len(oldTopics) == 0 {
		newOnes := make([]*OptionalTopic, len(cmd))
		for i := range cmd {
			newOnes[i] = &cmd[i]
		}

		fmt.Println("no old topics")

		return newOnes, nil
	}

	oldTopicsSets := make(map[string]bool, len(oldTopics))
	for i := range oldTopics {
		item := &oldTopics[i]
		oldTopicsSets[item.Title] = true
	}

	oldOnes := make(map[string]*OptionalTopic, len(cmd))
	newOnes := make([]*OptionalTopic, 0, len(cmd))

	for i := range cmd {
		item := &cmd[i]

		if _, ok := oldTopicsSets[item.Title]; ok {
			oldOnes[item.Title] = item
		} else {
			newOnes = append(newOnes, item)
		}
	}

	if n := len(oldTopics); len(oldOnes) != n {
		return nil, fmt.Errorf("the count of old topics is not matched, expect :%d, actual:%d", n, len(oldOnes))
	}

	if err := file.saveLastHotTopics(oldTopics, oldOnes); err != nil {
		return nil, err
	}

	return newOnes, nil
}

func (s *appService) handleNewOptionalTopics(community string, newOnes []*OptionalTopic, file *fileToReview) error {
	if len(newOnes) == 0 {
		return errors.New("no new topics")
	}

	oldTopics, err := s.repoNotHotTopic.FindAll(community)
	if err != nil {
		return err
	}
	if len(oldTopics) == 0 {
		for i := range newOnes {
			if err := file.saveNewTopic(newOnes[i]); err != nil {
				return err
			}
		}

		return nil
	}

	oldTopicsSet := make([]map[int]bool, len(oldTopics))
	for i := range oldTopics {
		oldTopicsSet[i] = oldTopics[i].GetDSSet()
	}

	newTopicsSet := make([]map[int]bool, len(newOnes))
	for i := range newOnes {
		newTopicsSet[i] = newOnes[i].getDSSet()
	}

	new2old, old2new := findRelationsBetweenCategories(newTopicsSet, oldTopicsSet)

	for i, v := range new2old {
		n := len(v)

		if n == 0 {
			if err := file.saveNewTopic(newOnes[i]); err != nil {
				return err
			}

			continue
		}

		if j := v[0]; n == 1 && len(old2new[j]) == 1 {
			ns := newTopicsSet[i]
			os := oldTopicsSet[j]

			switch parseRelationshipBetweenSets(ns, os) {
			case setsRelationSame:
				if err := file.saveUnchangedTopic(newOnes[i]); err != nil {
					return err
				}

			case setsRelationLeftIncludesRight:
				if err := file.saveTopicThatAppendToOld(newOnes[i], os); err != nil {
					return err
				}

			case setsRelationRightIncludesLeft:
				if err := file.saveTopicThatRemoveFromOld(&oldTopics[j], ns); err != nil {
					return err
				}

			default:
				if err := file.saveTopicThatIntersectWithMultiOlds(newOnes[i], v, oldTopics); err != nil {
					return err
				}
			}

			continue
		}

		if err := file.saveTopicThatIntersectWithMultiOlds(newOnes[i], v, oldTopics); err != nil {
			return err
		}
	}

	for j, v := range old2new {
		if len(v) == 0 {
			if err := file.saveTopicThatRemoveFromOld(&oldTopics[j], map[int]bool{}); err != nil {
				return err
			}
		}
	}

	return nil
}
