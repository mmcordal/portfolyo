package viewmodel

import "portfolyo/internal/model"

type ReminderVM struct {
	CreatedAt  string  `json:"created_at"`
	ID         int64   `json:"id"`
	UserID     int64   `json:"user_id"`
	User       *UserVM `json:"user,omitempty"`
	Title      string  `json:"title"`
	ReminderAt string  `json:"date"`
}
type ReminderRequest struct {
	Title string `json:"title" validate:"required" labelName:"başlık"`
	Date  string `json:"date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00" labelName:"tarih"`
}

func ToReminderVM(k *model.Reminder) *ReminderVM {
	vm := &ReminderVM{
		CreatedAt:  k.CreatedAt.Format("2006-01-02 15:04:05"),
		ID:         k.ID,
		UserID:     k.UserID,
		Title:      k.Title,
		ReminderAt: k.ReminderAt.Format("2006-01-02 15:04:05"),
	}
	if k.User != nil {
		vm.User = ToUserVM(k.User)
	}
	return vm
}

func ToReminderVMs(reminders []*model.Reminder) []*ReminderVM {
	vms := make([]*ReminderVM, 0, len(reminders))
	for _, k := range reminders {
		vms = append(vms, ToReminderVM(k))
	}
	return vms
}
