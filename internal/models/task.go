package models

type Task struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (t *Task) Clone() *Task {
	if t == nil {
		return nil
	}

	copy := *t
	return &copy
}
