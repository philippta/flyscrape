package flyscrape

import (
	"errors"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
)

var StopWatch = errors.New("stop watch")

func Watch(path string, fn func(string) error) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating file watcher: %w", err)
	}
	defer watcher.Close()

	if err := watcher.Add(path); err != nil {
		return fmt.Errorf("watching file %q: %w", path, err)
	}

	update := func() error {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return fn(string(data))
	}

	if err := update(); err != nil {
		if errors.Is(err, StopWatch) {
			return nil
		}
		return err
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Has(fsnotify.Remove) {
				return nil
			}
			if event.Has(fsnotify.Chmod) {
				continue
			}
			watcher.Remove(path)
			watcher.Add(path)
			if err := update(); err != nil {
				if errors.Is(err, StopWatch) {
					return nil
				}
				return err
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}
