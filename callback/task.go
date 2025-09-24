package callback

import (
	"github.com/google/uuid"
	"io"
)

type (
	Task struct {
		ID     string `json:"id"`
		URL    string `json:"url"`
		Body   string `json:"body"`
		Status status `json:"status"`
		body   io.Reader
	}
	TaskOption func(t *Task)
)

func NewTask(url string, ops ...TaskOption) *Task {

	t := &Task{
		URL: url,
	}

	for _, op := range ops {
		op(t)
	}

	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return t
}
