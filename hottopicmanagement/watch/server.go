package watch

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain/repository"
)

var impl *watchingImpl

func Start(
	cfg *Config,
	repo repository.RepoTopicSolution,
) {
	impl = &watchingImpl{
		repo:     repo,
		stop:     make(chan struct{}),
		stopped:  make(chan struct{}),
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
	repo     repository.RepoTopicSolution
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
			logrus.Error("failed to get oldest solution, err: %s", err.Error())
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

}
