package gopack

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/mschenk42/gopack/action"
)

var (
	DelayedNotify          taskRunSet = taskRunSet{}
	TasksRun               []string   = []string{}
	ErrActionNotRegistered            = errors.New("action not registered with task")
)

type BaseTask struct {
	OnlyIf      GuardFunc
	NotIf       GuardFunc
	ContOnError bool

	notify actionTaskRunSet
}

type GuardFunc func() (bool, error)
type ActionRunStatus map[action.Enum]bool
type actionTaskRunSet map[action.Enum]map[string]func()
type taskRunSet map[string]func()

func (d *taskRunSet) Run() {
	// running task funcs can result in new map entries
	for len(*d) > 0 {
		for k, f := range *d {
			f()
			// per Go spec it's safe to delete while ranging over
			delete(*d, k)
		}
	}
}

type Runner interface {
	Run(actions ...action.Enum) ActionRunStatus
}

type Task interface {
	Runner
	fmt.Stringer
}

func (b BaseTask) RunActions(task Task, regActions action.Methods, runActions []action.Enum) ActionRunStatus {
	var (
		f         action.MethodFunc
		found     bool
		err       error
		canRun    bool
		reason    string
		runStatus ActionRunStatus = ActionRunStatus{}
		timeStart time.Time       = time.Now()
	)

	if len(runActions) == 0 {
		b.logError(task, action.NewSlice(action.Nil), fmt.Errorf("unable to run, no action given"), timeStart)
		return runStatus
	}

	if canRun, reason = b.canRun(); !canRun {
		b.logSkipped(task, runActions, reason, timeStart)
		return runStatus
	}

	for _, a := range runActions {
		if f, found = regActions.Method(a); !found {
			b.logError(task, action.NewSlice(a), ErrActionNotRegistered, timeStart)
			continue
		}
		b.logStart(task, a)
		if runStatus[a], err = f(); err == nil {
			b.logRun(task, a, runStatus[a], reason, timeStart)
			if runStatus[a] {
				b.notifyTasks(a)
			}
		} else {
			b.handleTaskError(err)
		}
	}
	return runStatus
}

func (b *BaseTask) SetNotify(notify Task, forAction, whenAction action.Enum, delayed bool) {
	if b.notify == nil {
		b.notify = actionTaskRunSet{}
	}
	if b.notify[whenAction] == nil {
		b.notify[whenAction] = map[string]func(){}
	}
	if delayed {
		b.notify[whenAction][fmt.Sprintf("%s:%s", notify, forAction)] = func() {
			DelayedNotify[fmt.Sprintf("%s:%s", notify, forAction)] = func() { notify.Run(forAction) }
		}
	} else {
		b.notify[whenAction][fmt.Sprintf("%s:%s", notify, forAction)] = func() {
			notify.Run(forAction)
		}
	}
}

func (b *BaseTask) SetOnlyIf(f GuardFunc) {
	b.OnlyIf = f
}

func (b *BaseTask) SetNotIf(f GuardFunc) {
	b.NotIf = f
}

func (b BaseTask) notifyTasks(action action.Enum) {
	funcs, found := b.notify[action]
	if found {
		for _, f := range funcs {
			f()
		}
	}
}

func (b BaseTask) canRun() (bool, string) {
	var (
		err    error
		run    bool = true
		reason string
	)
	if b.OnlyIf != nil {
		reason = "due to only_if"
		run, err = b.OnlyIf()
		b.handleTaskError(err)
	}
	if b.NotIf != nil {
		reason = "due to not_if"
		run, err = b.NotIf()
		run = !run
		b.handleTaskError(err)
	}
	return run, reason
}

var (
	logStartFmt     string = colorize.Cyan("* %s: %s (%s)")
	logErrHeaderFmt string = colorize.Red("* %s: %s (%s) %s")
	logRunFmt       string = colorize.Cyan("  ~ %s: %s (%s) %s")
	logErrFmt       string = colorize.Red("   ! %s")
	logInfoFmt      string = "%s"
)

var colorize = ColorFormat{}

func NewTaskInfoWriter() io.Writer {
	return taskInfoWriter{}
}

type taskInfoWriter struct {
}

func (t taskInfoWriter) Write(b []byte) (int, error) {
	Log.Printf(logInfoFmt, b)
	return len(b), nil
}

func (b BaseTask) logStart(task Task, a action.Enum) {
	Log.Printf(logStartFmt, task, a, "started")
}

func (b BaseTask) logRun(task Task, a action.Enum, hasRun bool, reason string, t time.Time) {
	s := fmt.Sprintf(logRunFmt, task, a, "up to date", time.Since(t))
	if hasRun {
		status := "has run"
		if reason != "" {
			status = fmt.Sprintf("%s %s", status, reason)
		}
		s = fmt.Sprintf(logRunFmt, task, a, status, time.Since(t))
		TasksRun = append(TasksRun, s)
	}
	Log.Printf(s)
}

func (b BaseTask) logSkipped(task Task, a []action.Enum, reason string, t time.Time) {
	Log.Printf(logRunFmt, task, a, fmt.Sprintf("skipped %s", reason), time.Since(t))
}

func (b BaseTask) logError(task Task, a []action.Enum, err error, t time.Time) {
	Log.Printf(logErrHeaderFmt, task, a, "error", time.Since(t))
	b.handleTaskError(err)
}

func (b BaseTask) handleTaskError(err error) {
	if err == nil {
		return
	}
	if b.ContOnError {
		Log.Printf(logErrFmt, err)
	} else {
		Log.Fatalf(logErrFmt, err)
	}
}
