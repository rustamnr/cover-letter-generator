package services

import (
	"github.com/rustamnr/cover-letter-generator/internal/clients"
	"github.com/rustamnr/cover-letter-generator/internal/models"
)

type HHProvider struct {
	client *clients.HHClient
}

func NewHHProvider(client *clients.HHClient) *HHProvider {
	return &HHProvider{client: client}
}

func (p *HHProvider) GetResumes(accessToken string) (*models.APIResumeResponse, error) {
	return p.client.GetResumes(accessToken)
}

func (p *HHProvider) GetResume(accessToken, resumeID string) (*models.Resume, error) {
	return p.client.GetResume(accessToken, resumeID)
}
