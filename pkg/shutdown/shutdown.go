package shutdown

import (
	"go-task/pkg/logging"
	"io"
	"os"
	"os/signal"
)

func Graceful(signals []os.Signal, closerItems ...io.Closer) {
	logger := logging.GetLogger()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)

	sig := <-sigc

	logger.Infof("Caught signals: %s. Shutting down...", sig)

	for _, closer := range closerItems {
		if err := closer.Close(); err != nil {
			logger.Errorf("failed to close %v: %s", closer, err)
		}
	}

}
