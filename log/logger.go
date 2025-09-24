package log

import (
	"fmt"
	"github.com/jack0829/letsgo/common/fs"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	baseDir       string         // 日志根路径
	tc            *time.Ticker   // 定期检查文件
	x             *sync.RWMutex  // RW-Lock
	w             io.WriteCloser // 日志 writer
	path          string         // 当前日志路径
	once          *sync.Once     // once
	formatsOfTime []string       // 日志切割规则 time.Format，例：[]string{"2006","01","02","15.log"} 每天一个目录，每小时一个文件
}

func NewLogger(baseDir string, path ...string) (logger *Logger, err error) {

	if err = fs.MustDir(baseDir); err != nil {
		return
	}

	logger = &Logger{
		baseDir: baseDir,
		tc:      time.NewTicker(time.Second * 10),
		x:       &sync.RWMutex{},
		once:    &sync.Once{},
	}

	logger.rotateByTime(path...)
	return
}

func (l *Logger) rotateByTime(f ...string) {
	l.formatsOfTime = f
	l.checkRotate(time.Now())
	l.once.Do(func() {
		go l.rotate()
	})
}

func (l *Logger) rotate() {

	fmt.Println("滚动日志开始")
	defer func() {
		if l.w != nil {
			l.w.Close()
		}
		fmt.Println("滚动日志停止")
	}()

	for {
		t, ok := <-l.tc.C
		l.checkRotate(t)
		if !ok {
			break
		}
	}
}

func (l *Logger) GetLogPath(t time.Time) (path, dir string) {

	L := len(l.formatsOfTime)

	// 默认切割规则
	if L == 0 {
		dir = l.baseDir + t.Format("/2006/01/02")
		path = l.baseDir + t.Format("/2006/01/02/15.log")
		return
	}

	// 自定义规则
	if L == 1 {
		dir = l.baseDir
	} else {
		dir = l.baseDir + "/" + t.Format(strings.Join(l.formatsOfTime[:L-1], "/"))
	}

	path = l.baseDir + "/" + t.Format(strings.Join(l.formatsOfTime, "/"))
	return
}

func (l *Logger) checkRotate(t time.Time) error {

	l.x.Lock()
	defer l.x.Unlock()

	path, dir := l.GetLogPath(t)
	if path == l.path {
		return nil
	}

	if err := fs.MustDir(dir); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return err
	}

	if l.w != nil {
		if err := l.w.Close(); err != nil {
			return err
		}
	}

	l.w = file
	l.path = path
	return nil
}

func (l *Logger) Write(p []byte) (int, error) {
	l.x.RLock()
	defer l.x.RUnlock()
	return l.w.Write(p)
}

func (l *Logger) Close() error {
	if l.tc != nil {
		l.tc.Stop()
	}
	return nil
}
