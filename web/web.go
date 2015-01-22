package web

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	sigKill = 0
)

type IqRunner struct {
	running        bool
	controlSignals chan int
	errorSignals   chan error
	handlers       map[string]handlerF
	port           int
}

type handlerF func(http.ResponseWriter, *http.Request)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s!", r.URL.Path[1:])
}

func CreateRunner(port int, fns map[string]handlerF) *IqRunner {
	result := &IqRunner{
		handlers:       fns,
		port:           port,
		controlSignals: make(chan int),
		errorSignals:   make(chan error),
	}

	if len(fns) == 0 {
		result.handlers = map[string]handlerF{
			"/default": defaultHandler,
		}
	}
	return result
}

func (iqr *IqRunner) StartRun() {
	for k, f := range iqr.handlers {
		http.HandleFunc(k, f)
	}

	iqr.running = true
	go func() { iqr.errorSignals <- http.ListenAndServe(":"+strconv.Itoa(iqr.port), nil) }()
	go func() {
		for {
			sig, ok := <-iqr.controlSignals
			// if it's closed, it's been cleaned up already
			if !ok {
				return
			}
			switch sig {
			case sigKill:
				iqr.dieCleanly()
				return
			default:
				iqr.errorSignals <- fmt.Errorf("unknown signal %#v")
				iqr.dieHard()
				return
			}
		}
	}()
}

func (iqr *IqRunner) Errors() <-chan error {
	return iqr.errorSignals
}

func (iqr *IqRunner) Kill() {
	iqr.controlSignals <- sigKill
}

func (iqr *IqRunner) cleanup() {
	iqr.cleanupError()
	iqr.cleanupControl()
	iqr.running = false
}

func (iqr *IqRunner) cleanupError() {
	close(iqr.errorSignals)
}

func (iqr *IqRunner) cleanupControl() {
	close(iqr.controlSignals)
}

func (iqr *IqRunner) dieHard() {
	iqr.errorSignals <- fmt.Errorf("die hard")
	iqr.cleanup()
}

func (iqr *IqRunner) dieCleanly() {
	iqr.errorSignals <- nil
	iqr.cleanup()
}
