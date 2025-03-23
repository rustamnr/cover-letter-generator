package models

import "encoding/json"

// ParseResume парсит JSON в структуру Resume
func ParseResume(jsonData []byte) (*Resume, error) {
	var resume Resume
	err := json.Unmarshal(jsonData, &resume)
	if err != nil {
		return nil, err
	}
	return &resume, nil
}
