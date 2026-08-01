package hcaptcha

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/services/captcha_validation/validators"
	"github.com/imcrazytwkr/formdrain/services/captcha_validation/validators/common"
	"github.com/imcrazytwkr/formdrain/utils/maputil"
	"github.com/imcrazytwkr/formdrain/utils/utf8util"
)

type hcaptchaValidator struct {
	secret     string
	httpClient *http.Client
}

func NewHcaptchaValidator(secret string, httpClient *http.Client) validators.CaptchaValidator {
	return &hcaptchaValidator{
		secret:     secret,
		httpClient: httpClient,
	}
}

func (v *hcaptchaValidator) Validate(ctx context.Context, form map[string]any, hostname string, userIP string) error {
	log := common.GetLoggerForProvider(ctx, providerHcaptcha, common.ApiFormatHttp)

	responseToken, ok := maputil.GetString(form, hcaptchaKey)
	if !ok {
		return ErrNoHcaptchaToken
	}

	payload := url.Values{}
	payload.Set("secret", v.secret)
	payload.Set("response", responseToken)
	if len(userIP) > 0 {
		payload.Set("remoteip", userIP)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, hcaptchaUrl, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set(constants.HeaderContentType, constants.ContentTypeForm)
	response, err := v.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("hcaptcha backend responded with %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return fmt.Errorf("hcaptcha backend responded with malformed body: %w", err)
	}

	body, err = utf8util.FixBytes(body)
	if err != nil {
		return fmt.Errorf("hcaptcha backend responded with malformed body: %w", err)
	}

	var data hcaptchaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("hcaptcha backend responded with malformed body: %w", err)
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
