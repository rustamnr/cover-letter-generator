package models

// APIApplicationsResponse представляет список откликов пользователя
type APIApplicationsResponse struct {
	Found   int                `json:"found"`    // Количество найденных откликов
	Items   []ApplicationItem  `json:"items"`    // Список откликов
	Page    int                `json:"page"`     // Текущая страница
	Pages   int                `json:"pages"`    // Общее количество страниц
	PerPage int                `json:"per_page"` // Количество элементов на странице
}

// ApplicationItem представляет отклик на вакансию
type ApplicationItem struct {
	ID              string      `json:"id"`               // ID отклика
	Vacancy         ApplicationVacancy     `json:"vacancy"`          // Вакансия
	Employer        ApplicationEmployer    `json:"employer"`         // Работодатель
	Status          Status      `json:"status"`           // Статус отклика
	CreatedAt       string      `json:"created_at"`       // Дата создания
	UpdatedAt       *string     `json:"updated_at"`       // Дата обновления (может быть пустой)
	URL             *string     `json:"url"`             // Ссылка на отклик (может быть пустой)
	DeclineAllowed  *bool       `json:"decline_allowed"`  // Можно ли отклонить приглашение
	HasUpdates      *bool       `json:"has_updates"`      // Есть ли новые обновления
	Hidden          *bool       `json:"hidden"`           // Скрыт ли отклик
	MessagingStatus *string     `json:"messaging_status"` // Статус сообщений
	State           *State      `json:"state"`            // Состояние отклика (приглашение, отклик и т. д.)
	PhoneCalls      *PhoneCalls `json:"phone_calls"`      // Информация о звонках
}

// ApplicationVacancy содержит информацию о вакансии
type ApplicationVacancy struct {
	ID           string  `json:"id"`             // ID вакансии
	Title        *string `json:"title"`          // Название вакансии (может быть пустым)
	AlternateURL *string `json:"alternate_url"`  // Ссылка на вакансию
}

// ApplicationEmployer содержит информацию о работодателе
type ApplicationEmployer struct {
	ID   *string `json:"id"`   // ID работодателя (может быть пустым)
	Name *string `json:"name"` // Название компании (может быть пустым)
}

// Status описывает статус отклика
type Status struct {
	ID   *string `json:"id"`   // ID статуса (может быть пустым)
	Name *string `json:"name"` // Название статуса (может быть пустым)
}

// State описывает состояние отклика
type State struct {
	ID   *string `json:"id"`   // ID состояния (может быть пустым)
	Name *string `json:"name"` // Название состояния (может быть пустым)
}

// PhoneCalls содержит информацию о звонках
type PhoneCalls struct {
	Items                     []PhoneCall `json:"items"`                       // Список звонков
	PickedUpPhoneByOpponent   *bool       `json:"picked_up_phone_by_opponent"` // Поднял ли трубку собеседник
}

// PhoneCall содержит информацию о конкретном звонке
type PhoneCall struct {
	ID               int     `json:"id"`                // ID звонка
	CreationTime     string  `json:"creation_time"`     // Время создания звонка
	DurationSeconds  *int    `json:"duration_seconds"`  // Длительность звонка (может быть `null`)
	LastChangeTime   *string `json:"last_change_time"`  // Время последнего изменения (может быть `null`)
	Status           *string `json:"status"`            // Статус звонка (может быть пустым)
}
