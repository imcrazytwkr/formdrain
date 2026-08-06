package account

import "errors"

var ErrNilAccount = errors.New("account is nil")
var ErrNoRequiredParams = errors.New("account email and password_hash are required")
