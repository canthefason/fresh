package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/howeyc/fsnotify"
)

func watchFolder(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if isWatchedFile(ev.Name) {
					watcherLog("sending event %s", ev)
					startChannel <- ev.String()
					changeChannel <- struct{}{}
				}
			case err := <-watcher.Error:
				watcherLog("error: %s", err)
			}
		}
	}()

	val := os.Getenv("RUNNER_ENABLE_LOGS")
	log, _ := strconv.ParseBool(val)
	if log {
		watcherLog("Watching %s", path)
	}

	err = watcher.Watch(path)

	if err != nil {
		fatal(err)
	}
}

func watch() {
	root := watchDir()
	env := os.Getenv("GOPATH")
	root = fmt.Sprintf("%s/src/%s", env, root)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !isTmpDir(path) {
			if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}

			watchFolder(path)
		}

		return err
	})
}
