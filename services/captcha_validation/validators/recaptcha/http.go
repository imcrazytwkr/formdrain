package recaptcha

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/services/captcha_validation/validators"
	"github.com/imcrazytwkr/formdrain/services/captcha_validation/validators/common"
	"github.com/imcrazytwkr/formdrain/types"
	"github.com/imcrazytwkr/formdrain/utils/maputil"
	"github.com/imcrazytwkr/formdrain/utils/utf8util"
	"github.com/valyala/fastjson"
)

type recaptchaValidator struct {
	secret     string
	client     types.HttpClient
	parserPool *fastjson.ParserPool
}

func NewRecaptchaValidator(secret string, client types.HttpClient) validators.CaptchaValidator {
	return &recaptchaValidator{
		secret:     secret,
		client:     client,
		parserPool: &fastjson.ParserPool{},
	}
}

func (v *recaptchaValidator) Validate(ctx context.Context, form map[string]any, hostname string, userIP string) error {
	log := common.GetLoggerForProvider(ctx, providerRecaptcha, common.ApiFormatHttp)

	responseToken, ok := maputil.GetString(form, recaptchaKey)
	if !ok {
		return ErrNoRecaptchaToken
	}

	payload := url.Values{}
	payload.Set("secret", v.secret)
	payload.Set("response", responseToken)
	if len(userIP) > 0 {
		payload.Set("remoteip", userIP)
	}

	request, err := http.NewRequest(http.MethodPost, recaptchaUrl, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set(constants.HeaderContentType, constants.ContentTypeForm)
	response, err := v.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("recaptcha backend responded with %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return fmt.Errorf("recaptcha backend responded with malformed body: %w", err)
	}

	body, err = utf8util.FixBytes(body)
	if err != nil {
		return fmt.Errorf("recaptcha backend responded with malformed body: %w", err)
	}

	parser := v.parserPool.Get()
	defer v.parserPool.Put(parser)

	value, err := parser.ParseBytes(body)
	if err != nil {
		return fmt.Errorf("recaptcha backend responded with malformed body: %w", err)
	}

	if !value.GetBool(successKey) {
		return constants.ErrCaptchaNotPassed
	}

	responseHostname := string(value.GetStringBytes(hostnameKey))
	if len(responseHostname) > 0 && responseHostname != hostname {
		log.Warn().Msgf("token hostname mismatch: expected %q but got %q", hostname, responseHostname)
		return constants.ErrCaptchaNotPassed
	}

	return nil
}
