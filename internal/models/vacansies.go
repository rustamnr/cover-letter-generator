package models

import "fmt"

type Vacancy struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	BrandedDescription *string            `json:"branded_description,omitempty"`
	Area               Area               `json:"area"`
	Salary             *Salary            `json:"salary,omitempty"`
	Employment         Employment         `json:"employment"`
	VacancyExperience  VacancyExperience  `json:"experience"`
	Schedule           Schedule           `json:"schedule"`
	ProfessionalRoles  []ProfessionalRole `json:"professional_roles"`
	KeySkills          []KeySkill         `json:"key_skills"`
	Languages          []Language         `json:"languages,omitempty"`
	DriverLicenseTypes []DriverLicense    `json:"driver_license_types,omitempty"`

	// Vacancy-specific fields
	AlternateURL           string                `json:"alternate_url"`
	ApplyAlternateURL      string                `json:"apply_alternate_url"`
	Archived               bool                  `json:"archived"`
	Approved               bool                  `json:"approved"`
	CreatedAt              string                `json:"created_at"`
	ExpiresAt              string                `json:"expires_at"`
	PublishedAt            string                `json:"published_at"`
	InitialCreatedAt       string                `json:"initial_created_at"`
	Employer               Employer              `json:"employer"`
	Address                *Address              `json:"address,omitempty"`
	Contacts               *Contacts             `json:"contacts,omitempty"`
	WorkScheduleByDays     []WorkScheduleByDays  `json:"work_schedule_by_days,omitempty"`
	WorkingDays            []WorkingDay          `json:"working_days,omitempty"`
	WorkingHours           []WorkingHour         `json:"working_hours,omitempty"`
	WorkingTimeIntervals   []WorkingTimeInterval `json:"working_time_intervals,omitempty"`
	WorkingTimeModes       []WorkingTimeMode     `json:"working_time_modes,omitempty"`
	WorkFormat             []WorkFormat          `json:"work_format,omitempty"`
	AllowMessages          bool                  `json:"allow_messages"`
	ResponseLetterRequired bool                  `json:"response_letter_required"`
	ResponseNotifications  bool                  `json:"response_notifications"`
	ResponseURL            *string               `json:"response_url,omitempty"`
	HasTest                bool                  `json:"has_test"`
	Test                   *Test                 `json:"test,omitempty"`
	Internship             bool                  `json:"internship"`
	NightShifts            bool                  `json:"night_shifts"`
	Premium                bool                  `json:"premium"`
	VideoVacancy           *VideoVacancy         `json:"video_vacancy,omitempty"`
	Type                   Type                  `json:"type"`
	BillingType            BillingType           `json:"billing_type"`
	CanUpgradeBillingType  bool                  `json:"can_upgrade_billing_type"`
	Code                   *string               `json:"code,omitempty"`
	Manager                *Manager              `json:"manager,omitempty"`
	Department             *Department           `json:"department,omitempty"`
	InsiderInterview       *InsiderInterview     `json:"insider_interview,omitempty"`
}

type VacancyForLLM struct {
	ID          string            `json:"id"`          // Идентификатор вакансии
	Name        string            `json:"name"`        // Название вакансии
	Description string            `json:"description"` // Описание вакансии
	Contacts    Contacts          `json:"contacts,omitempty"`
	Location    string            `json:"location"`     // Локация (город)
	Employment  string            `json:"employment"`   // Тип занятости
	Experience  VacancyExperience `json:"experience"`   // Требуемый опыт работы
	Schedule    Schedule          `json:"schedule"`     // График работы
	KeySkills   []string          `json:"key_skills"`   // Ключевые навыки
	CompanyName string            `json:"company_name"` // Название компании
}

func (v *Vacancy) VacancyToLLMModel() VacancyForLLM {
	// Преобразуем ключевые навыки
	var keySkills []string
	for _, skill := range v.KeySkills {
		keySkills = append(keySkills, skill.Name)
	}

	// Проверяем наличие контактов
	var contacts Contacts
	if v.Contacts != nil {
		contacts = *v.Contacts
	}

	return VacancyForLLM{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description, // Если нужно, можно добавить очистку HTML-тегов
		Contacts:    contacts,
		Location:    v.Area.Name,
		Employment:  v.Employment.Name,
		Experience:  v.VacancyExperience,
		Schedule:    v.Schedule,
		KeySkills:   keySkills,
		CompanyName: v.Employer.Name,
	}
}

// SimilarVacanciesResponse представляет ответ от API /resumes/{resume_id}/similar_vacancies
type SimilarVacanciesResponse struct {
	Found   int       `json:"found"`    // Количество найденных вакансий
	Items   []Vacancy `json:"items"`    // Список вакансий
	Page    int       `json:"page"`     // Текущая страница
	Pages   int       `json:"pages"`    // Общее количество страниц
	PerPage int       `json:"per_page"` // Количество элементов на странице
}

type Employer struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	AlternateURL string    `json:"alternate_url"`
	URL          string    `json:"url"`
	Trusted      bool      `json:"trusted"`
	Blacklisted  bool      `json:"blacklisted"`
	Badges       []Badge   `json:"badges,omitempty"`
	LogoURLs     *LogoURLs `json:"logo_urls,omitempty"`
}

type Address struct {
	City          string         `json:"city"`
	Street        string         `json:"street"`
	Building      string         `json:"building"`
	Description   string         `json:"description"`
	Lat           float64        `json:"lat"`
	Lng           float64        `json:"lng"`
	MetroStations []MetroStation `json:"metro_stations"`
}

type MetroStation struct {
	ID       string  `json:"station_id"`
	Name     string  `json:"station_name"`
	LineID   string  `json:"line_id"`
	LineName string  `json:"line_name"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

type Contacts struct {
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Phones []Phone `json:"phones"`
}

type Phone struct {
	Country string  `json:"country"`
	City    string  `json:"city"`
	Number  string  `json:"number"`
	Comment *string `json:"comment"` // Комментарий (может быть nil)
}

type Badge struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Year        int    `json:"year"`
	URL         string `json:"url"`
	Position    string `json:"position"`
}

type LogoURLs struct {
	Small    string `json:"90"`
	Medium   string `json:"240"`
	Original string `json:"original"`
}

type WorkScheduleByDays struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WorkingDay struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WorkingHour struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WorkingTimeInterval struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WorkingTimeMode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WorkFormat struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VideoVacancy struct {
	VideoURL     string       `json:"video_url"`
	CoverPicture CoverPicture `json:"cover_picture"`
}

type CoverPicture struct {
	ResizedPath   string `json:"resized_path"`
	ResizedWidth  int    `json:"resized_width"`
	ResizedHeight int    `json:"resized_height"`
}

type Test struct {
	Required bool `json:"required"`
}

type Type struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BillingType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Manager struct {
	ID string `json:"id"`
}

type InsiderInterview struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type VacancyExperience struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (v *Vacancy) ShortInfo() string {
	return fmt.Sprintf("%s (%s) - %s", v.Name, v.Employer.Name, v.Area.Name)
}
