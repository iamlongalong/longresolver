package k8sresolver

var finishChan chan struct{}
var callback = func() {}

func SetFinishCallback(f func()) {
	callback = f
}

func startRun() chan struct{} {
	go func() {
		select {
		case <-finishChan:
			callback()
		}
	}()

	return finishChan
}
