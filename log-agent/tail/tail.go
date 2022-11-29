package tail

import (
	"github.com/hpcloud/tail"
)

type Tail struct {
	tail *tail.Tail
}

// ReadChan 从通道中读日志
func (t *Tail) ReadChan() <-chan *tail.Line {
	return t.tail.Lines
}

func New(filename string) *Tail {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	tailF, err := tail.TailFile(filename, config)
	if err != nil {
		panic(err)
	}
	return &Tail{tail: tailF}
}
