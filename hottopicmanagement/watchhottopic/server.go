package watchhottopic

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/app"
)

var impl *watchingImpl

func Start(
	cfg *Config,
	app app.AppService,
	communities []string,
) {
	impl = &watchingImpl{
		app:         app,
		communities: communities,

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
	app         app.AppService
	communities []string

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
		impl.handle(needStop)

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

func (impl *watchingImpl) handle(needStop func() bool) {
	for _, community := range impl.communities {
		if needStop() {
			return
		}

		if err := impl.app.ApplyToHotTopic(community); err != nil {
			logrus.Errorf("apply hot topics failed, err:%s", err.Error())
		}
	}
}
