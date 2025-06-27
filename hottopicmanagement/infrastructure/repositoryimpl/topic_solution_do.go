package repositoryimpl

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain/repository"
	"github.com/opensourceways/hot-topic-website-backend/utils"
)

func totopicSolutionsDO(Community string, v []domain.TopicSolution) topicSolutionsDO {
	r := make([]topicSolutionDO, len(v))
	for i := range v {
		r[i] = totopicSolutionDO(&v[i])
	}

	return topicSolutionsDO{
		Community:      Community,
		CreatedAt:      utils.Now(),
		TopicSolutions: r,
	}
}

// topicSolutionDO
type topicSolutionsDO struct {
	Id             primitive.ObjectID `bson:"_id"           json:"-"`
	Community      string             `bson:"community"     json:"community"`
	CreatedAt      int64              `bson:"created_at"    json:"created_at"`
	TopicSolutions []topicSolutionDO  `bson:"topics"        json:"topics"`
}

func (do *topicSolutionsDO) toDoc() (bson.M, error) {
	return genDoc(do)
}

func (do *topicSolutionsDO) index() string {
	return do.Id.Hex()
}

func (do *topicSolutionsDO) toTopicSolutions() repository.TopicSolutions {
	r := make([]domain.TopicSolution, len(do.TopicSolutions))
	for i := range do.TopicSolutions {
		r[i] = do.TopicSolutions[i].toTopicSolution()
	}

	return repository.TopicSolutions{
		Id:             do.index(),
		Community:      do.Community,
		TopicSolutions: r,
	}
}

// topicSolutionDO
type topicSolutionDO struct {
	TopicId   string                       `bson:"topic_id"   json:"topic_id"`
	Solutions []discussionSourceSolutionDO `bson:"solutions"  json:"solutions"`
}

func (do *topicSolutionDO) toTopicSolution() domain.TopicSolution {
	r := make([]domain.DiscussionSourceSolution, len(do.Solutions))
	for i := range do.Solutions {
		r[i] = do.Solutions[i].toDiscussionSourceSolution()
	}

	return domain.TopicSolution{
		TopicId:   do.TopicId,
		Solutions: r,
	}
}

func totopicSolutionDO(v *domain.TopicSolution) topicSolutionDO {
	r := make([]discussionSourceSolutionDO, len(v.Solutions))
	for i := range v.Solutions {
		r[i] = todiscussionSourceSolutionDO(&v.Solutions[i])
	}

	return topicSolutionDO{
		TopicId:   v.TopicId,
		Solutions: r,
	}
}

// discussionSourceSolutionDO
type discussionSourceSolutionDO struct {
	ResolvedOne int   `bson:"resolved_one"  json:"resolved_one"`
	RelatedOnes []int `bson:"related_ones"  json:"related_ones"`
}

func (do *discussionSourceSolutionDO) toDiscussionSourceSolution() domain.DiscussionSourceSolution {
	return domain.DiscussionSourceSolution{
		ResolvedOne: do.ResolvedOne,
		RelatedOnes: do.RelatedOnes,
	}
}

func todiscussionSourceSolutionDO(v *domain.DiscussionSourceSolution) discussionSourceSolutionDO {
	return discussionSourceSolutionDO{
		ResolvedOne: v.ResolvedOne,
		RelatedOnes: v.RelatedOnes,
	}
}
