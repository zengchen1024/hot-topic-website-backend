package app

import (
	"errors"
	"fmt"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain/repository"
)

type AppService interface {
	ToReview(CmdToUploadOptionalTopics) error
}

func NewAppService(
	config *Config,
	repoHotTopic repository.RepoHotTopic,
	repoNotHotTopic repository.RepoNotHotTopic,
) *appService {
	return &appService{
		filePath:        config.FilePath,
		repoHotTopic:    repoHotTopic,
		repoNotHotTopic: repoNotHotTopic,
	}
}

type appService struct {
	filePath        string
	repoHotTopic    repository.RepoHotTopic
	repoNotHotTopic repository.RepoNotHotTopic
}

func (s *appService) ToReview(cmd CmdToUploadOptionalTopics) error {
	file, err := newfileToReview(s.filePath)
	if err != nil {
		return fmt.Errorf("new excel failed, err:%s", err.Error())
	}

	newOnes, err := s.handleOldTopics(cmd, file)
	if err != nil {
		return err
	}

	if err := s.handleNewOptionalTopics(newOnes, file); err != nil {
		return err
	}

	if err := file.saveToFile(); err != nil {
		return err
	}

	// TODO send email

	return nil
}

func (s *appService) handleOldTopics(cmd CmdToUploadOptionalTopics, file *fileToReview) ([]*OptionalTopic, error) {
	oldTopics, err := s.repoHotTopic.FindOpenOnes() // TODO check if ordered by HotTopic's Order
	if err != nil {
		return nil, err
	}
	if len(oldTopics) == 0 {
		return nil, fmt.Errorf("no old topics")
	}

	oldTopicsMap := make(map[string]*domain.HotTopic, len(oldTopics))
	for i := range oldTopics {
		item := &oldTopics[i]
		oldTopicsMap[item.Title] = item
	}

	oldOnes := make(map[int]*OptionalTopic, len(cmd))
	newOnes := make([]*OptionalTopic, 0, len(cmd))

	for i := range cmd {
		item := &cmd[i]

		old, ok := oldTopicsMap[item.Title]
		if ok {
			//fmt.Println(item.Title)
			oldOnes[old.Order] = item
		} else {
			newOnes = append(newOnes, item)
		}
	}

	if n := len(oldTopics); len(oldOnes) != n {
		return nil, fmt.Errorf("the count of old topics is not matched, expect :%d, actual:%d", n, len(oldOnes))
	}

	fmt.Printf("oldTopics = %d, oldOnes=%d\n", len(oldTopics), len(oldOnes))
	if err := file.saveLastTopics(oldTopics, oldOnes); err != nil {
		return nil, err
	}

	return newOnes, nil
}

func (s *appService) handleNewOptionalTopics(newOnes []*OptionalTopic, file *fileToReview) error {
	if len(newOnes) == 0 {
		return errors.New("no new topics")
	}

	fmt.Printf("newOnes : %d\n", len(newOnes))

	oldTopics, err := s.repoNotHotTopic.FindAll()
	if err != nil {
		return err
	}
	if len(oldTopics) == 0 {
		// TODO write new

		return nil
	}

	fmt.Printf("oldTopics: %d\n", len(oldTopics))

	oldTopicsSet := make([]map[int]bool, len(oldTopics))
	for i := range oldTopics {
		oldTopicsSet[i] = oldTopics[i].GetDSSet()
	}

	newTopicsSet := make([]map[int]bool, len(newOnes))
	for i := range newOnes {
		newTopicsSet[i] = newOnes[i].getDSSet()
	}

	fmt.Printf("oldTopicsSet:%d\nnewTopicsSet:%d\n", len(oldTopicsSet), len(newTopicsSet))

	new2old, old2new := findRelationsBetweenCategories(newTopicsSet, oldTopicsSet)

	fmt.Printf("new2old: %d\n", len(new2old))
	fmt.Printf("old2new: %d\n", len(old2new))

	for i, v := range new2old {
		n := len(v)

		if n == 0 {
			fmt.Println("find new")
			file.saveNewTopic(newOnes[i])

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
