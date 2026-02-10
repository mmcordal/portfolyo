package viewmodel

import "portfolyo/internal/model"

type ReminderVM struct {
	CreatedAt string `json:"created_at"`
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Title     string `json:"title"`
	Date      string `json:"date"`
}
type ReminderRequest struct {
	Title string `json:"title"`
	Date  string `json:"date"`
}

func ToReminderVM(k *model.Reminder) *ReminderVM {
	return &ReminderVM{
		CreatedAt: k.CreatedAt.Format("2006-01-02 15:04:05"),
		ID:        k.ID,
		UserID:    k.UserID,
		Title:     k.Title,
		Date:      k.Date.Format("2006-01-02 15:04:05"),
	}
}

func ToReminderVMs(reminde []*model.Reminder) []*ReminderVM {
	vms := make([]*ReminderVM, 0, len(reminde))
	for _, k := range reminde {
		vms = append(vms, ToReminderVM(k))
	}
	return vms
}
