package gopack

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/mschenk42/gopack/action"
	"github.com/mschenk42/gopack/color"
)

var (
	delayedNotify taskRunSet = taskRunSet{}
	tasksRun      []string   = []string{}
	indentLevel   int
)

type BaseTask struct {
	OnlyIf      guardFunc
	NotIf       guardFunc
	ContOnError bool

	notify actionTaskRunSet
}

type ActionRunStatus map[action.Name]bool

type guardFunc func() (bool, error)
type actionTaskRunSet map[action.Name]map[string]func()
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
	Run(actions ...action.Name) ActionRunStatus
}

type Task interface {
	Runner
	fmt.Stringer
}

func (b BaseTask) RunActions(task Task, regActions action.Funcs, runActions []action.Name) ActionRunStatus {
	var (
		f         action.Func
		found     bool
		err       error
		canRun    bool
		reason    string
		runStatus ActionRunStatus = ActionRunStatus{}
		timeStart time.Time       = time.Now()
	)

	indentLevel += 1

	if len(runActions) == 0 {
		b.logError(task, action.NewSlice(action.Nil), fmt.Errorf("unable to run, no action given"), timeStart)
		indentLevel -= 1
		return runStatus
	}

	if canRun, reason = b.canRun(); !canRun {
		b.logSkipped(task, runActions, reason, timeStart)
		indentLevel -= 1
		return runStatus
	}

	for _, a := range runActions {
		timeStart = time.Now()
		if f, found = regActions.Func(a); !found {
			b.logError(task, action.NewSlice(a), errors.New("action not registered with task"), timeStart)
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
		if indentLevel == 1 {
			Log.Println()
		}
	}
	indentLevel -= 1
	return runStatus
}

func (b *BaseTask) SetNotify(notify Task, forAction, whenAction action.Name, delayed bool) {
	if b.notify == nil {
		b.notify = actionTaskRunSet{}
	}
	if b.notify[whenAction] == nil {
		b.notify[whenAction] = map[string]func(){}
	}
	if delayed {
		b.notify[whenAction][fmt.Sprintf("%s:%s", notify, forAction)] = func() {
			delayedNotify[fmt.Sprintf("%s:%s", notify, forAction)] = func() { notify.Run(forAction) }
		}
	} else {
		b.notify[whenAction][fmt.Sprintf("%s:%s", notify, forAction)] = func() {
			notify.Run(forAction)
		}
	}
}

func (b *BaseTask) SetOnlyIf(f guardFunc) {
	b.OnlyIf = f
}

func (b *BaseTask) SetNotIf(f guardFunc) {
	b.NotIf = f
}

func (b BaseTask) notifyTasks(action action.Name) {
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
	logStartFmt     = color.Cyan("%s%s: %s (%s)")
	logErrHeaderFmt = color.Red("%s%s: %s (%s) %s")
	logRunFmt       = color.Cyan("%s%s: %s (%s) %s")
	logErrFmt       = color.Red("%s! %s")
	logInfoFmt      = "%s%s"
)

func logIndent() string {
	if indentLevel-1 <= 0 {
		return ""
	}
	return strings.Repeat(" ", (indentLevel-1)*2)
}

func NewTaskInfoWriter() io.Writer {
	return taskInfoWriter{}
}

type taskInfoWriter struct {
}

func (t taskInfoWriter) Write(b []byte) (int, error) {
	buf := bytes.NewBuffer(b)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		Log.Printf(logInfoFmt, logIndent(), scanner.Text())
	}
	return len(b), nil
}

func (b BaseTask) logStart(task Task, a action.Name) {
	Log.Printf(logStartFmt, logIndent(), task, a, "started")
}

func (b BaseTask) logRun(task Task, a action.Name, hasRun bool, reason string, t time.Time) {
	s := fmt.Sprintf(logRunFmt, logIndent(), task, a, "up to date", time.Since(t))
	if hasRun {
		status := "has run"
		if reason != "" {
			status = fmt.Sprintf("%s %s", status, reason)
		}
		s = fmt.Sprintf(logRunFmt, logIndent(), task, a, status, time.Since(t))
		// let's just track the top most tasks
		if indentLevel == 1 {
			tasksRun = append(tasksRun, fmt.Sprintf(logRunFmt, "", task, a, status, time.Since(t)))
		}
	}
	Log.Printf(s)
}

func (b BaseTask) logSkipped(task Task, a []action.Name, reason string, t time.Time) {
	Log.Printf(logRunFmt, logIndent(), task, a, fmt.Sprintf("skipped %s", reason), time.Since(t))
}

func (b BaseTask) logError(task Task, a []action.Name, err error, t time.Time) {
	Log.Printf(logErrHeaderFmt, logIndent(), task, a, "error", time.Since(t))
	b.handleTaskError(err)
}

func (b BaseTask) handleTaskError(err error) {
	if err == nil {
		return
	}
	if b.ContOnError {
		lines := strings.Split(err.Error(), "\n")
		for _, s := range lines {
			Log.Printf(logErrFmt, logIndent(), s)
		}
	} else {
		Log.Panicf(logErrFmt, logIndent(), err)
	}
}
