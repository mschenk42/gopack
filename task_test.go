package gopack

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/mschenk42/gopack/action"
	"github.com/stretchr/testify/assert"
)

const (
	failKeywords = "\\(up to date\\)"
	passKeywords = "\\(run\\) "
)

func TestRunTask(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Log
	buf := &bytes.Buffer{}
	Log = log.New(buf, "", 0)
	defer func() { Log = saveLogger }()

	t1 := Task1{
		Name: "task1",
	}

	assert.NotPanics(func() { t1.Run(action.Create) })
	assert.Regexp(fmt.Sprintf(`task1.*create.*%s`, passKeywords), buf.String())
	fmt.Print(buf.String())
}

func TestGuards(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Log
	buf := &bytes.Buffer{}
	Log = log.New(buf, "", 0)
	defer func() { Log = saveLogger }()

	t1 := Task1{
		Name:     "task1",
		BaseTask: BaseTask{NotIf: func() (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(action.Create) })
	assert.Regexp(fmt.Sprintf(`task1.*create.*%s`, "skipped due to not_if"), buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{NotIf: func() (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(action.Create) })
	assert.Regexp(fmt.Sprintf(`task1.*create.*%s`, "run due to not_if"), buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{OnlyIf: func() (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(action.Create) })
	assert.Regexp(fmt.Sprintf(`task1.*create.*%s`, "run due to only_if"), buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{OnlyIf: func() (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(action.Create) })
	assert.Regexp(fmt.Sprintf(`task1.*create.*%s`, "skipped due to only_if"), buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name: "task1",
		BaseTask: BaseTask{
			OnlyIf: func() (bool, error) { return true, nil },
			NotIf:  func() (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(action.Create) })
	assert.Regexp(fmt.Sprintf(`task1.*create.*%s`, "skipped due to not_if"), buf.String())
	fmt.Print(buf.String())
}

func TestContOnError(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Log
	buf := &bytes.Buffer{}
	Log = log.New(buf, "", 0)
	defer func() { Log = saveLogger }()

	t1 := Task1{
		Name: "task1",
		BaseTask: BaseTask{
			ContOnError: true},
	}

	assert.NotPanics(func() { t1.Run() })
	assert.Regexp(`! unable to run, no action given`, buf.String())
	fmt.Print(buf.String())
}

func TestWhenRun(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Log
	buf := &bytes.Buffer{}
	Log = log.New(buf, "", 0)
	defer func() { Log = saveLogger }()

	t1 := Task1{
		Name: "task1",
	}

	t2 := Task2{
		Name: "task2 notified",
	}

	t1.SetNotify(t2, action.Create, action.Create, false)

	assert.NotPanics(func() { t1.Run(action.Create) })
	assert.Regexp(fmt.Sprintf(`task1.*create.*%s`, passKeywords), buf.String())
	assert.Regexp(fmt.Sprintf(`task2 notified.*create.*%s`, passKeywords), buf.String())
	fmt.Print(buf.String())
}

func TestDelayedRun(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Log
	buf := &bytes.Buffer{}
	Log = log.New(buf, "", 0)
	defer func() { Log = saveLogger }()

	t1 := Task1{
		Name: "task1",
	}

	t2 := Task2{
		Name: "task2 notified",
	}

	t3 := Task1{
		Name: "task3",
	}

	t1.SetNotify(t2, action.Create, action.Create, true)
	t3.SetNotify(t2, action.Create, action.Create, true)

	assert.NotPanics(func() { t1.Run(action.Create) })
	assert.NotPanics(func() { t3.Run(action.Create) })
	assert.NotPanics(func() { DelayedNotify.Run() })
	assert.Regexp(fmt.Sprintf(`task1.*create.*%s`, passKeywords), buf.String())
	assert.Regexp(fmt.Sprintf(`task3.*create.*%s`, passKeywords), buf.String())
	re := regexp.MustCompile(fmt.Sprintf(`task2 notified.*create.*%s`, passKeywords))
	matches := re.FindAllString(buf.String(), -1)
	assert.Exactly(1, len(matches), "task 2 notified more than once")

	fmt.Print(buf.String())
}

func TestDelayedChainedRun(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Log
	buf := &bytes.Buffer{}
	Log = log.New(buf, "", 0)
	defer func() { Log = saveLogger }()

	t1 := Task1{
		Name: "task1 notified",
	}

	t2 := Task2{
		Name: "task2 notified",
	}

	t3 := Task1{
		Name: "task3",
	}

	t2.SetNotify(t1, action.Create, action.Create, true)
	t3.SetNotify(t2, action.Create, action.Create, true)

	assert.NotPanics(func() { t3.Run(action.Create) })
	assert.NotPanics(func() { DelayedNotify.Run() })
	assert.Regexp(fmt.Sprintf(`task1.*create.*%s`, passKeywords), buf.String())
	assert.Regexp(fmt.Sprintf(`task3.*create.*%s`, passKeywords), buf.String())
	re := regexp.MustCompile(fmt.Sprintf(`task(1|2) notified.*create.*%s`, passKeywords))
	matches := re.FindAllString(buf.String(), -1)
	assert.Exactly(2, len(matches), "notifications should be 2")

	fmt.Print(buf.String())
}

type Task1 struct {
	Name string

	BaseTask
}

func (t Task1) Run(runActions ...action.Enum) ActionRunStatus {
	regActions := action.Methods{
		action.Create: t.create,
	}
	return t.BaseTask.RunActions(&t, regActions, runActions)
}

func (t Task1) String() string {
	return t.Name
}

func (t Task1) create() (bool, error) {
	return true, nil
}

type Task2 struct {
	Name string

	BaseTask
}

func (t Task2) Run(runActions ...action.Enum) ActionRunStatus {
	regActions := action.Methods{
		action.Nothing: t.nothing,
		action.Create:  t.create,
	}
	return t.BaseTask.RunActions(&t, regActions, runActions)
}

func (t Task2) String() string {
	return t.Name
}

func (t Task2) create() (bool, error) {
	return true, nil
}

func (t Task2) nothing() (bool, error) {
	return false, nil
}
