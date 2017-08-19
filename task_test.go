package gopack

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunTask(t *testing.T) {
	assert := assert.New(t)
	logger, buf := newLoggerBuffer()

	t1 := Task1{
		Name: "task1",
	}

	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*[NOT RUN]", buf.String())
	fmt.Print(buf.String())
}

func TestGuards(t *testing.T) {
	assert := assert.New(t)
	logger, buf := newLoggerBuffer()

	t1 := Task1{
		Name:     "task1",
		BaseTask: BaseTask{NotIf: func(p Properties, x *log.Logger) (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*[RUN]", buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{NotIf: func(p Properties, x *log.Logger) (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*[NOT RUN]", buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{OnlyIf: func(p Properties, x *log.Logger) (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*[NOT RUN]", buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{OnlyIf: func(p Properties, x *log.Logger) (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*[RUN]", buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name: "task1",
		BaseTask: BaseTask{
			OnlyIf: func(p Properties, x *log.Logger) (bool, error) { return true, nil },
			NotIf:  func(p Properties, x *log.Logger) (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*[NOT RUN]", buf.String())
	fmt.Print(buf.String())
}

func TestContOnError(t *testing.T) {
	assert := assert.New(t)
	logger, buf := newLoggerBuffer()

	t1 := Task1{
		Name: "task1",
		BaseTask: BaseTask{
			ContOnError: true},
	}

	assert.NotPanics(func() { t1.Run(nil, logger) })
	assert.Regexp("unable to run task1", buf.String())
	fmt.Print(buf.String())
}
