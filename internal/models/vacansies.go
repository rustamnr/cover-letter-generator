package models

// Vacancy представляет вакансию с hh.ru
type Vacancy struct {
	ID                     string                `json:"id"`
	Name                   string                `json:"name"`
	AlternateURL           string                `json:"alternate_url"`
	ApplyAlternateURL      string                `json:"apply_alternate_url"`
	Archived               bool                  `json:"archived"`
	Approved               bool                  `json:"approved"`
	CreatedAt              string                `json:"created_at"`
	ExpiresAt              string                `json:"expires_at"`
	PublishedAt            string                `json:"published_at"`
	InitialCreatedAt       string                `json:"initial_created_at"`
	Description            string                `json:"description"`
	BrandedDescription     *string               `json:"branded_description,omitempty"`
	Employment             Employment            `json:"employment"`
	Experience             Experience            `json:"experience"`
	Salary                 *ResumeSalary         `json:"salary,omitempty"`
	Employer               Employer              `json:"employer"`
	Area                   Area                  `json:"area"`
	Address                *Address              `json:"address,omitempty"`
	Contacts               *Contacts             `json:"contacts,omitempty"`
	ProfessionalRoles      []ProfessionalRole    `json:"professional_roles"`
	WorkScheduleByDays     []WorkScheduleByDays  `json:"work_schedule_by_days"`
	WorkingDays            []WorkingDay          `json:"working_days"`
	WorkingHours           []WorkingHour         `json:"working_hours"`
	WorkingTimeIntervals   []WorkingTimeInterval `json:"working_time_intervals"`
	WorkingTimeModes       []WorkingTimeMode     `json:"working_time_modes"`
	WorkFormat             []WorkFormat          `json:"work_format"`
	DriverLicenseTypes     []DriverLicense       `json:"driver_license_types"`
	Languages              []Language            `json:"languages"`
	KeySkills              []KeySkill            `json:"key_skills"`
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
	Schedule               Schedule              `json:"schedule"`
	BillingType            BillingType           `json:"billing_type"`
	CanUpgradeBillingType  bool                  `json:"can_upgrade_billing_type"`
	Code                   *string               `json:"code,omitempty"`
	Manager                *Manager              `json:"manager,omitempty"`
	Department             *Department           `json:"department,omitempty"`
	InsiderInterview       *InsiderInterview     `json:"insider_interview,omitempty"`
}

// Работодатель
type Employer struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	AlternateURL string   `json:"alternate_url"`
	URL          string   `json:"url"`
	Trusted      bool     `json:"trusted"`
	Blacklisted  bool     `json:"blacklisted"`
	Badges       []Badge  `json:"badges"`
	LogoURLs     LogoURLs `json:"logo_urls"`
}

// Адрес
type Address struct {
	City          string         `json:"city"`
	Street        string         `json:"street"`
	Building      string         `json:"building"`
	Description   string         `json:"description"`
	Lat           float64        `json:"lat"`
	Lng           float64        `json:"lng"`
	MetroStations []MetroStation `json:"metro_stations"`
}

// Метро
type MetroStation struct {
	ID       string  `json:"station_id"`
	Name     string  `json:"station_name"`
	LineID   string  `json:"line_id"`
	LineName string  `json:"line_name"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

// Контакты
type Contacts struct {
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Phones []Phone `json:"phones"`
}

// Телефон
type Phone struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Number  string `json:"number"`
}

// Зарплата
type Salary struct {
	From     *int   `json:"from,omitempty"`
	To       *int   `json:"to,omitempty"`
	Currency string `json:"currency"`
	Gross    bool   `json:"gross"`
}

// Опыт работы
type Experience struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Тип занятости
type Employment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// График работы
type Schedule struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Ключевые навыки
type KeySkill struct {
	Name string `json:"name"`
}

// Категория водительских прав
type DriverLicense struct {
	ID string `json:"id"`
}

// Языки
type Language struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Level Level  `json:"level"`
}

type Level struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Бейджи работодателя
type Badge struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Year        int    `json:"year"`
	URL         string `json:"url"`
	Position    string `json:"position"`
}

// Логотипы работодателя
type LogoURLs struct {
	Small    string `json:"90"`
	Medium   string `json:"240"`
	Original string `json:"original"`
}

// Профессиональные роли
type ProfessionalRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Расписание работы
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

// Видео-вакансия
type VideoVacancy struct {
	VideoURL     string       `json:"video_url"`
	CoverPicture CoverPicture `json:"cover_picture"`
}

type CoverPicture struct {
	ResizedPath   string `json:"resized_path"`
	ResizedWidth  int    `json:"resized_width"`
	ResizedHeight int    `json:"resized_height"`
}

// Тесты
type Test struct {
	Required bool `json:"required"`
}

// Тип вакансии
type Type struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Тип биллинга
type BillingType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Департамент
type Department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Менеджер
type Manager struct {
	ID string `json:"id"`
}

// Интервью
type InsiderInterview struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// Area содержит информацию о регионе
type Area struct {
	ID   string `json:"id"`   // ID региона
	Name string `json:"name"` // Название региона
	URL  string `json:"url"`  // Ссылка на описание региона
}
