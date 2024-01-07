package entity

type TaskStatus int

const (
	TaskIncompleted TaskStatus = iota
	TaskCompleted
)

func (ts TaskStatus) String() string {
	return [...]string{
		"incompleted",
		"completed",
	}[ts]
}

type Task struct {
	ID     int
	Name   string
	Status TaskStatus
}
