package service

import (
	"context"
	"errors"
	"portfolyo/internal/model"
	"portfolyo/internal/repository"
	"portfolyo/internal/viewmodel"
)

type ReminderService interface {
	Create(ctx context.Context, vm *viewmodel.ReminderRequest, userID int64) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context, userID int64) ([]*viewmodel.ReminderVM, error)
}

type reminderService struct {
	rr repository.ReminderRepository
}

func NewReminderService(rr repository.ReminderRepository) ReminderService {
	return &reminderService{rr: rr}
}

func (s *reminderService) Create(ctx context.Context, vm *viewmodel.ReminderRequest, userID int64) error {
	if vm == nil {
		return errors.New("input is validate")
	}
	if userID == 0 {
		return errors.New("unauthorized")
	}
	reminde := &model.Reminder{
		UserID: userID,
		Title:  vm.Title,
	}
	date, err := parseTransactionDate(vm.Date)
	if err != nil {
		return err
	}
	reminde.Date = date

	return s.rr.Create(ctx, reminde)
}

func (s *reminderService) Delete(ctx context.Context, id int64) error {
	reminde, err := s.rr.GetByReminderID(ctx, id)
	if err != nil {
		return err
	}
	if reminde == nil {
		return errors.New("reminde not found")
	}

	return s.rr.Delete(ctx, reminde)
}

func (s *reminderService) GetAll(ctx context.Context, userID int64) ([]*viewmodel.ReminderVM, error) {
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	reminders, err := s.rr.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	return viewmodel.ToReminderVMs(reminders), nil
}
