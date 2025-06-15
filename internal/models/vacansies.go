package models

import (
	"strings"
)

type Vacancy struct {
	ID                     string            `json:"id"`          // Идентификатор вакансии
	Name                   string            `json:"name"`        // Название вакансии
	Description            string            `json:"description"` // Описание вакансии
	BrandedDescription     *string           `json:"branded_description,omitempty"`
	Contacts               Contacts          `json:"contacts,omitempty"`
	Location               string            `json:"location"`     // Локация (город)
	Employment             Employment        `json:"employment"`   // Тип занятости
	Experience             VacancyExperience `json:"experience"`   // Требуемый опыт работы
	Schedule               Schedule          `json:"schedule"`     // График работы
	KeySkills              []KeySkill        `json:"key_skills"`   // Ключевые навыки
	CompanyName            string            `json:"company_name"` // Название компании
	ResponseLetterRequired bool              `json:"response_letter_required"`
	Test                   *Test             `json:"test,omitempty"`
}

func (v Vacancy) ToString() string {
	var builder strings.Builder

	// Добавляем основную информацию о вакансии
	builder.WriteString("Вакансия:\n")
	builder.WriteString("Название: " + v.Name + "\n")
	builder.WriteString("Компания: " + v.CompanyName + "\n")
	builder.WriteString("Локация: " + v.Location + "\n")
	// builder.WriteString("Тип занятости: " + v.Employment + "\n")
	builder.WriteString("Опыт работы: " + v.Experience.Name + "\n")
	builder.WriteString("График работы: " + v.Schedule.Name + "\n")

	// Добавляем описание вакансии
	builder.WriteString("\nОписание:\n")
	builder.WriteString(v.Description + "\n")

	// Добавляем ключевые навыки
	if len(v.KeySkills) > 0 {
		builder.WriteString("\nКлючевые навыки:\n")
		// builder.WriteString(strings.Join(v.KeySkills, ", ") + "\n")
		for _, skill := range v.KeySkills {
			builder.WriteString("- " + skill.Name + "\n")
		}
	}

	// Добавляем контакты, если они есть
	if v.Contacts.Name != "" || v.Contacts.Email != "" || len(v.Contacts.Phones) > 0 {
		builder.WriteString("\nКонтакты:\n")
		if v.Contacts.Name != "" {
			builder.WriteString("Имя: " + v.Contacts.Name + "\n")
		}
		if v.Contacts.Email != "" {
			builder.WriteString("Email: " + v.Contacts.Email + "\n")
		}
		if len(v.Contacts.Phones) > 0 {
			builder.WriteString("Телефоны:\n")
			for _, phone := range v.Contacts.Phones {
				builder.WriteString("- +" + phone.Country + " (" + phone.City + ") " + phone.Number + "\n")
			}
		}
	}

	return builder.String()
}
