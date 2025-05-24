package services

import "errors"

var (
	ErrCourseNotFound          = errors.New("курс не найден")
	ErrCourseAlreadyPurchased  = errors.New("курс уже куплен")
	ErrInsufficientPermissions = errors.New("недостаточно прав")
	ErrInvalidTestScore        = errors.New("неверный результат теста")
)
