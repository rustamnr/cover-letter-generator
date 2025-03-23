package models

// UserInfo содержит основную информацию о пользователе
type UserInfo struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	BirthDate  string `json:"birth_date"`
	Location   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"area"`
	Contact struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"contact"`
}

// WorkExperience описывает опыт работы пользователя
type WorkExperience struct {
	Company     string `json:"company"`
	Position    string `json:"position"`
	StartDate   string `json:"start"`
	EndDate     string `json:"end"`
	Description string `json:"description"`
}

// Education содержит данные об образовании
type Education struct {
	Level          string `json:"level"`
	Institution    string `json:"institution"`
	Faculty        string `json:"faculty"`
	GraduationYear int    `json:"year"`
}

// Skills представляет список навыков
type Skills struct {
	SkillList []string `json:"skills"`
}

// Language содержит информацию о языках
type Language struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

// Salary описывает зарплатные ожидания
type Salary struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

// Resume объединяет все данные резюме
type Resume struct {
	UserInfo   UserInfo         `json:"user"`
	Experience []WorkExperience `json:"experience"`
	Education  Education        `json:"education"`
	Skills     Skills           `json:"skills"`
	Languages  []Language       `json:"languages"`
	Salary     Salary           `json:"salary"`
	Employment string           `json:"employment"`
}
