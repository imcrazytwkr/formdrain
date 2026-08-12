package mappers

import (
	"github.com/google/uuid"
	fr "github.com/imcrazytwkr/formdrain/models/form_response"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
)

func FormResponse(src *fr.FormResponse) (api.FormResponse, error) {
	id, err := uuid.Parse(src.Id)
	if err != nil {
		return api.FormResponse{}, err
	}

	payload := src.Payload
	if payload == nil {
		payload = map[string]any{}
	}

	return api.FormResponse{
		Id:            id,
		FormId:        src.FormId,
		CreatedAt:     src.CreatedAt,
		SchemaVersion: src.SchemaVersion,
		ClientIp:      clientIP(src),
		Payload:       payload,
	}, nil
}

func clientIP(src *fr.FormResponse) *string {
	if !src.ClientIP.IsValid() {
		return nil
	}

	s := src.ClientIP.String()
	return &s
}
