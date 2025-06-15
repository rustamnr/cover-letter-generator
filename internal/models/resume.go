package models

import (
	"fmt"
	"strings"
)

type ResumesResponse struct {
	Items []ResumeFull `json:"items"`
}

type Resume struct {
	ID         string       `json:"id"`         // Идентификатор резюме
	Title      string       `json:"title"`      // Название резюме
	FirstName  string       `json:"first_name"` // Имя
	LastName   string       `json:"last_name"`  // Фамилия
	Location   string       `json:"location"`   // Локация (город)
	Contact    []Contact    `json:"contact"`
	Skills     string       `json:"skills"`     // Ключевые навыки
	SkillsSet  []string     `json:"skill_set"`  // Ключевые навыки (набор)
	Experience []Experience `json:"experience"` // Опыт работы
}

func (rs *Resume) ToString() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Резюме: %s\n", rs.Title))
	sb.WriteString(fmt.Sprintf("Имя: %s %s\n", rs.FirstName, rs.LastName))
	sb.WriteString(fmt.Sprintf("Локация: %s\n", rs.Location))

	for _, contact := range rs.Contact {
		if contact.Type.ID == "email" {
			sb.WriteString(fmt.Sprintf("Email: %s\n", contact.Value))
		}
		if contact.Type.ID == "cell" {
			if phone, ok := contact.Value.(map[string]interface{}); ok {
				if formatted, exists := phone["formatted"].(string); exists {
					sb.WriteString(fmt.Sprintf("Телефон: %s\n", formatted))
				}
			}
		}
	}

	sb.WriteString(fmt.Sprintf("Ключевые навыки: %s\n", rs.Skills))

	sb.WriteString("Опыт работы:\n")
	for _, exp := range rs.Experience {
		endDate := "по настоящее время"
		if exp.EndDate != nil && *exp.EndDate != "" {
			endDate = *exp.EndDate
		}
		sb.WriteString(fmt.Sprintf("%s - %s (%s - %s)\n", exp.Position, exp.Company, exp.StartDate, endDate))
		if exp.Description != "" {
			sb.WriteString(fmt.Sprintf("Описание: %s\n", exp.Description))
		}
	}

	return sb.String()
}
