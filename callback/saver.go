package callback

type Saver interface {
	Load() []*Task
	Save(t []*Task)
}
