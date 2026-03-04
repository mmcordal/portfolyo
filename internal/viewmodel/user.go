package viewmodel

import "portfolyo/internal/model"

type UserVM struct {
	ID         int64          `json:"id" labelName:"id"`
	CreatedAt  string         `json:"created_at" labelName:"oluşturulma tarihi"`
	UpdatedAt  string         `json:"updated_at" labelName:"güncellenme tarihi"`
	DeletedAt  *string        `json:"deleted_at,omitempty" labelName:"silinme tarihi"`
	FullName   string         `json:"full_name" labelName:"tam isim"`
	Email      string         `json:"email" labelName:"e-posta"`
	UserAssets []*UserAssetVM `json:"user_assets,omitempty" labelName:"kullanıcı varlıkları"`
	Reminders  []*ReminderVM  `json:"reminders,omitempty" labelName:"hatırlatıcı"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50" labelName:"isim"`
	Surname  string `json:"surname" validate:"required,min=2,max=50" labelName:"soyisim"`
	Email    string `json:"email" validate:"required,email" labelName:"e-posta"`
	Password string `json:"password" validate:"required,min=8,max=72" labelName:"şifre"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" labelName:"e-posta"`
	Password string `json:"password" validate:"required" labelName:"şifre"`
}

type UpdateRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=50" labelName:"isim"`
	Surname  string `json:"surname" validate:"omitempty,min=2,max=50" labelName:"soyisim"`
	Email    string `json:"email" validate:"omitempty,email" labelName:"e-posta"`
	Password string `json:"password" validate:"omitempty,min=8,max=72" labelName:"şifre"`
}

type LoginResponse struct {
	Token string `json:"token" labelName:"token"`
}

func ToUserVM(u *model.User) *UserVM {
	var deletedAt *string
	if u.DeletedAt != nil {
		deletedAtStr := u.DeletedAt.Format("2006-01-02 15:04:05")
		deletedAt = &deletedAtStr
	}

	vm := &UserVM{
		ID:        u.ID,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt: deletedAt,
		FullName:  u.String(),
		Email:     u.Email,
	}
	return vm
}
