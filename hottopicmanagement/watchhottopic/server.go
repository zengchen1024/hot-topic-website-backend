package watchhottopic

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/app"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain/repository"
)

var impl *watchingImpl

func Start(
	cfg *Config,
	app app.AppService,
	repo repository.RepoNotHotTopic,
	communities []string,
) {
	impl = &watchingImpl{
		handler:  newHandler(app, repo, communities),
		stop:     make(chan struct{}),
		stopped:  make(chan struct{}),
		interval: cfg.intervalDuration(),
	}

	impl.start()

	logrus.Info("start to watch applying hot topics")
}

func Stop() {
	if impl != nil {
		impl.exit()

		logrus.Info("stop watching applying hot topics")
	}
}

// watchingImpl
type watchingImpl struct {
	handler *handler

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
		impl.handler.handle(needStop)

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
