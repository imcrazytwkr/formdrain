package hcaptcha

import "errors"

const providerHcaptcha = "hcaptcha"

const hcaptchaKey = "h-captcha"
const hcaptchaUrl = "https://hcaptcha.com/siteverify"

const successKey = "success"
const hostnameKey = "hostname"

var ErrNoHcaptchaToken = errors.New("no hcaptcha token in request")
