package slog

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	log := GetLogger()
	log.Trace("are you prety?", true)
	log.Debugf("are you prety? %t", true)
	log.Info("how old are you? ", nil)
	log.Infof("i'm %010d", 18)
	log.Warn("you aren't honest! ")
	log.Warnf("haha%02d %v", 1000, nil)
	log.Trace("set level to warn!!!!!")
	log.Trace("what?")
	log.Info("what?")
	log.Error("what?")
	log.Errorf("what?..$%s$", "XD")
	log.Fatalf("import cycle not allowed! %s", "shit...")
	log.Fatal("never reach here")
	log.Panic("panic...")

	SetContextField("uid", rand.Uint64())
	SetContextFields(Fields{"admin": true, "username": "jackson"})
	log.Infof("with uid, admin, username context variable")

	log.BindFields(Fields{"tid": "bytedance"})
	log.Infof("with tid logger binded fields")

	DelContextField("uid")
	log.Infof("no uid context variable")
}

func TestDefaultLogger(t *testing.T) {
	SetLevel(TraceLevel)
	Trace("are you prety?", true)
	Debugf("are you prety? %t", true)
	Debug("okkkkkk")
	Info("how old are you? ", nil)
	Infof("i'm %010d", 18)
	Warn("you aren't honest! ")
	Warnf("haha%02d %v", 1000, nil)
	Trace("set level to warn!!!!!")
	Tracef("what: %d", 1230)
	Info("what?")
	Error("what?")
	Errorf("what?..$%s$", "XD")
	Panic("panic")
	Fatalf("import cycle not allowed! %s", "shit...")
	Fatal("never reach here")
	time.Sleep(time.Millisecond)
}

func TestAsyncLogger(t *testing.T) {
	EnableAsyncDriver()

	log := GetLogger()
	log.Infof("no fields")

	SetContextField("uid", rand.Uint64())
	SetContextFields(Fields{"admin": true, "username": "jackson"})
	log.Infof("with uid, admin, username context variable")

	log.BindFields(Fields{"tid": "bytedance"})
	log.Infof("with tid logger binded fields")

	DelContextField("uid")
	log.Infof("no uid context variable")

	time.Sleep(time.Millisecond)
	SetDriver(new(StdDriver))
}

func TestLoggerLevelFilter(t *testing.T) {
	SetLevel(WarnLevel)
	SetLoggerLevel("debug", DebugLevel)
	SetLoggerLevelMap(map[string]Level{
		"info":  InfoLevel,
		"error": ErrorLevel,
	})

	debugLog := NewLogger("debug")
	infoLog := NewLogger("info")
	errorLog := NewLogger("error")
	tmpLog := NewLogger("xxxxx")

	var debugCount, infoCount, errorCount int32
	RegisterHook(func(log *Log) {
		switch log.Level {
		case DebugLevel:
			atomic.AddInt32(&debugCount, 1)
		case InfoLevel:
			atomic.AddInt32(&infoCount, 1)
		case ErrorLevel:
			atomic.AddInt32(&errorCount, 1)
		}
	})

	debugLog.Trace("debug.trace, invisible")
	debugLog.Debug("debug.debug, visible")

	infoLog.Debug("info.debug, invisible")
	infoLog.Info("info.info, visible")
	infoLog.Error("info.error, visible")

	errorLog.Info("error.info, invisible")
	errorLog.Warn("error.warn, invisble")
	errorLog.Error("error.error, visible")

	tmpLog.Info("tmp.info, invisible")
	tmpLog.Warn("tmp.warn, visible")
	tmpLog.Error("tmp.error, visible")

	time.Sleep(time.Millisecond)
	assert.True(t, atomic.LoadInt32(&debugCount) == 1, atomic.LoadInt32(&debugCount))
	assert.True(t, atomic.LoadInt32(&infoCount) == 1)
	assert.True(t, atomic.LoadInt32(&errorCount) == 3)
}

func TestConcurrency(t *testing.T) {
	log := NewLogger("concurrency")
	d := newAsyncDriver(1 << 12)
	d.stdout = nil
	SetDriver(d)

	const threadNum = 64
	var wait sync.WaitGroup
	wait.Add(threadNum)
	for i := 0; i < threadNum; i++ {
		threadId := i
		go func() {
			for x := 0; x < 1000; x++ {
				log.Infof("threadId=%v, seq=%d", threadId, x)
				log.Error(threadId, x, " xxxxxxxxxxxxxxxx")
				time.Sleep(time.Microsecond * 100)
			}
			wait.Done()
		}()
	}
	wait.Wait()
	time.Sleep(time.Second)

	// use default StdDriver for avoiding break other tests
	SetDriver(new(StdDriver))
}
