package Resources

import (
	"github.com/vahidlotfi71/Task_Manager/Models"
)

type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Assignee    string `json:"assignee"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func Single(t Models.Task) Task {
	return Task{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Assignee:    t.Assignee,
		CreatedAt:   t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func Collection(tasks []Models.Task) []Task {
	out := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, Single(t))
	}
	return out
}
