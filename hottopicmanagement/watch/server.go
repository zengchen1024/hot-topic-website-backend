package watch

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain/repository"
)

var impl *watchingImpl

func Start(
	cfg *Config,
	ht repository.RepoHotTopic,
	repo repository.RepoTopicSolution,
) {
	impl = &watchingImpl{
		repo:     repo,
		stop:     make(chan struct{}),
		stopped:  make(chan struct{}),
		handler:  newtopicSolutionHandler(ht, cfg),
		interval: cfg.intervalDuration(),
	}

	impl.start()

	logrus.Info("start to watch resolved issues")
}

func Stop() {
	if impl != nil {
		impl.exit()

		logrus.Info("stop watching resolved issues")
	}
}

// watchingImpl
type watchingImpl struct {
	handler *topicSolutionHandler
	repo    repository.RepoTopicSolution

	stop     chan struct{}
	stopped  chan struct{}
	interval time.Duration
}

func (impl *watchingImpl) start() {
	go impl.watch()
}

func (impl *watchingImpl) exit() {
	close(impl.stop)

	<-impl.stopped
}

func (impl *watchingImpl) watch() {
	needStop := func() bool {
		select {
		case <-impl.stop:
			return true
		default:
			return false
		}
	}

	var timer *time.Timer

	defer func() {
		if timer != nil {
			timer.Stop()
		}

		close(impl.stopped)
	}()

	for {
		triggered, err := impl.repo.FindOldest()
		if err != nil {
			logrus.Errorf("failed to get oldest solution, err: %s", err.Error())
		}

		impl.handle(&triggered, needStop)

		// time starts.
		if timer == nil {
			timer = time.NewTimer(impl.interval)
		} else {
			timer.Reset(impl.interval)
		}

		select {
		case <-impl.stop:
			return

		case <-timer.C:
		}
	}
}

func (impl *watchingImpl) handle(solution *repository.TopicSolutions, needStop func() bool) {
	retry := impl.handler.handle(solution, needStop)
	solution.RetryNum++

	if len(retry) == 0 || solution.RetryNum >= 3 {
		var err error
		for i := 0; i < 3; i++ {
			if err = impl.repo.Remove(solution.Id); err == nil {
				return
			}
		}

		if err != nil {
			logrus.Errorf("delete solution:%s failed, err:%s", solution.Id, err.Error())
		}
	}

	solution.TopicSolutions = retry

	var err error
	for i := 0; i < 3; i++ {
		if err = impl.repo.Save(solution); err == nil {
			return
		}
	}

	if err != nil {
		logrus.Errorf("save solution:%s failed, err:%s", solution.Id, err.Error())
	}

}
