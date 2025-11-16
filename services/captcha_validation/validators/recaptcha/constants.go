package recaptcha

import "errors"

const providerRecaptcha = "recaptcha"

const recaptchaKey = "g-recaptcha"
const recaptchaUrl = "https://www.google.com/recaptcha/api/siteverify"

const successKey = "success"
const hostnameKey = "hostname"

var ErrNoRecaptchaToken = errors.New("no recaptcha token in request")
