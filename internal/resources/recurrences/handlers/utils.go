package handlers

import (
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
)

func MapRecurrenceResource(r repository.Recurrence) RecurrenceResource {
	return RecurrenceResource{
		ID:            r.ID,
		UserID:        r.UserID,
		Name:          r.Name,
		Amount:        r.Amount,
		DayOfMonth:    r.DayOfMonth,
		StartPeriod:   r.StartPeriod,
		EndPeriod:     r.EndPeriod,
		Note:          r.Note,
		CategoryID:    r.CategoryID,
		CategoryName:  r.CategoryName,
		CategoryColor: r.CategoryColor,
		CreatedAt:     r.CreatedAt,
	}
}
