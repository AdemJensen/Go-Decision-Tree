package config

import (
	"sync"

	"github.com/gosuri/uiprogress"
)

var (
	uiProgressOnce sync.Once
	uiProgress     *uiprogress.Progress
)

func GetUiProgress() *uiprogress.Progress {
	uiProgressOnce.Do(func() {
		uiProgress = uiprogress.New()
		uiProgress.Start()
	})
	return uiProgress
}
