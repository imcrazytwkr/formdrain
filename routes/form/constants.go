package form

import "errors"

// @NOTE: 4kb seems like a reasonable global max form size
const maxBodySize int64 = 4096

const actionSend = "send"

const keyFormConfig = "form_config"
const keySiteConfig = "site_config"

var errInvalidOrigin = errors.New("CORS origin is not valid")
