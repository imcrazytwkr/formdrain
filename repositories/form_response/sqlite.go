package form_response

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	fr "github.com/imcrazytwkr/formdrain/models/form_response"
	"github.com/imcrazytwkr/formdrain/repositories"
)

type sqliteFormResponseRepository struct {
	db *sql.DB
}

func NewSqliteFormResponseRepository(db *sql.DB) repositories.FormResponseRepository {
	return &sqliteFormResponseRepository{db: db}
}

const insertFormResponse = `
INSERT INTO form_responses (
	id,
	form_id,
	created_at,
	schema_version,
	client_ip,
	payload
) VALUES (?, ?, ?, ?, ?, ?)
`

func (r *sqliteFormResponseRepository) SaveFormResponse(ctx context.Context, response *fr.FormResponse) (string, error) {
	if response == nil {
		return "", fmt.Errorf("form response is nil")
	}

	// @NOTE: do not modify source object until insertion has succeeded
	id := response.Id
	if uuid.Validate(id) != nil {
		id = uuid.New().String()
	}

	createdAt := response.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	payloadJSON, err := json.Marshal(response.Payload)
	if err != nil {
		return "", err
	}

	var clientIP any
	if response.ClientIP.IsValid() {
		clientIP = response.ClientIP.String()
	}

	_, err = r.db.ExecContext(
		ctx,
		insertFormResponse,
		id,
		response.FormId,
		createdAt.UTC().Format(time.RFC3339),
		response.SchemaVersion,
		clientIP,
		string(payloadJSON),
	)

	if err != nil {
		return "", err
	}

	response.Id = id
	response.CreatedAt = createdAt.UTC()
	return id, nil
}
