package form_config

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/imcrazytwkr/formdrain/models/common"
	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/repositories"
)

type sqliteFormConfigRepository struct {
	db *sql.DB
}

func NewSqliteFormConfigRepository(db *sql.DB) repositories.FormConfigRepository {
	return &sqliteFormConfigRepository{db: db}
}

const selectFormConfigById = `
SELECT
	site_id,
	captcha_type,
	captcha_field,
	redirect_to,
	field_schema,
	schema_version,
	notifiers,
	success_template,
	error_template,
	redirect_template
FROM forms
WHERE id = ?
`

func (r *sqliteFormConfigRepository) GetFormConfigById(ctx context.Context, id int64) (*fc.FormConfig, error) {
	if id < 1 {
		return nil, nil
	}

	var config fc.FormConfig
	var captchaField sql.NullString
	var redirectTo sql.NullString
	var rawCaptchaType string
	var fieldSchema string
	var notifiersJSON string
	var successTemplate sql.NullString
	var errorTemplate sql.NullString
	var redirectTemplate sql.NullString

	err := r.db.QueryRowContext(ctx, selectFormConfigById, id).Scan(
		&config.SiteId,
		&rawCaptchaType,
		&captchaField,
		&redirectTo,
		&fieldSchema,
		&config.SchemaVersion,
		&notifiersJSON,
		&successTemplate,
		&errorTemplate,
		&redirectTemplate,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	config.CaptchaType, err = fc.ParseCaptchaType(rawCaptchaType)
	if err != nil {
		return nil, err
	}

	if config.CaptchaType == fc.CaptchaTypeUndefined {
		return nil, ErrInvalidCaptchaType
	}

	if captchaField.Valid {
		config.CaptchaField = captchaField.String
	}

	if redirectTo.Valid {
		config.RedirectTo = redirectTo.String
	}

	err = json.Unmarshal([]byte(fieldSchema), &config.FieldSchema)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(notifiersJSON), &config.Notifiers)
	if err != nil {
		return nil, err
	}

	config.SuccessTemplate, err = parseOptionalTemplate(successTemplate)
	if err != nil {
		return nil, err
	}

	config.ErrorTemplate, err = parseOptionalTemplate(errorTemplate)
	if err != nil {
		return nil, err
	}

	config.RedirectTemplate, err = parseOptionalTemplate(redirectTemplate)
	if err != nil {
		return nil, err
	}

	config.FormId = id
	return &config, nil
}

func parseOptionalTemplate(raw sql.NullString) (*common.Template, error) {
	if !raw.Valid || len(raw.String) < 1 {
		return nil, nil
	}

	return common.NewTemplate(raw.String)
}
