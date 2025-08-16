package service

import (
	"github.com/massivemadness/schedule-service/internal/entity"
	"github.com/massivemadness/schedule-service/internal/repository"
)

type InstructorService struct {
	instructorRepo repository.InstructorRepository
	groupRepo      repository.InstructorGroupRepository
}

func NewInstructorService(
	instructorRepo repository.InstructorRepository,
	groupRepo repository.InstructorGroupRepository,
) *InstructorService {
	return &InstructorService{
		instructorRepo: instructorRepo,
		groupRepo:      groupRepo,
	}
}

func (s *InstructorService) Register(instructorID int64, instructorName string) error {
	return s.instructorRepo.CreateIfNotExists(instructorID, instructorName)
}

func (s *InstructorService) CheckIsRegistered(instructorID int64) error {
	instructor, err := s.instructorRepo.GetByUserID(instructorID)
	if err != nil {
		return err
	}
	if instructor == nil {
		return entity.ErrNotFound
	}

	group, err := s.groupRepo.GetByUserID(instructorID)
	if err != nil {
		return err
	}
	if group == nil {
		return entity.ErrNotLinked
	}
	return nil
}

func (s *InstructorService) LinkGroup(instructorID int64, groupID int64) error {
	// Проверяем, что пользователь — зарегистрированный инструктор
	instructor, err := s.instructorRepo.GetByUserID(instructorID)
	if err != nil {
		return err
	}
	if instructor == nil {
		return entity.ErrNotFound
	}

	// Проверяем, что группа ещё не связана с другим инструктором
	existingGroup, err := s.groupRepo.GetByGroupID(groupID)
	if err != nil {
		return err
	}
	if existingGroup != nil {
		if existingGroup.InstructorID == instructor.ID {
			return entity.ErrAlreadyLinked
		} else {
			return entity.ErrOtherUserLinked
		}
	}

	// Связываем группу с инструктором
	return s.groupRepo.LinkGroup(instructorID, groupID)
}
