package recaptcha

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"net/url"
	"strings"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/services/captcha_validation/validators"
	"github.com/imcrazytwkr/formdrain/services/captcha_validation/validators/common"
	"github.com/imcrazytwkr/formdrain/utils/utf8util"
)

type recaptchaValidator struct {
	client *http.Client
}

func NewRecaptchaValidator(client *http.Client) validators.CaptchaValidator {
	return &recaptchaValidator{
		client: client,
	}
}

func (v *recaptchaValidator) Validate(
	ctx context.Context,
	secret string,
	responseToken string,
	hostname string,
	userIP netip.Addr,
) error {
	log := common.GetLoggerForProvider(ctx, providerRecaptcha, common.ApiFormatHttp)

	if len(responseToken) < 1 {
		return constants.ErrCaptchaNotPassed
	}

	payload := url.Values{}
	payload.Set("secret", secret)
	payload.Set("response", responseToken)
	if userIP.IsValid() {
		payload.Set("remoteip", userIP.String())
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, recaptchaUrl, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set(constants.HeaderContentType, constants.ContentTypeForm)
	response, err := v.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("recaptcha backend responded with %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("recaptcha backend responded with malformed body: %w", err)
	}

	body, err = utf8util.FixBytes(body)
	if err != nil {
		return fmt.Errorf("recaptcha backend responded with malformed body: %w", err)
	}

	var data recaptchaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("recaptcha backend responded with malformed body: %w", err)
	}

	if !data.Success {
		return constants.ErrCaptchaNotPassed
	}

	if len(data.Hostname) > 0 && data.Hostname != hostname {
		log.Warn().Msgf("token hostname mismatch: expected %q but got %q", hostname, data.Hostname)
		return constants.ErrCaptchaNotPassed
	}

	return nil
}
