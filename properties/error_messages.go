package properties

const (
	RECORD_NOT_FOUND          = "record not found"
	GENERIC_DATABASE_ERROR    = "error occured in processing to database"
	LACKING_EMAIL_OR_PASSWORD = "email or password cannot be empty"
	DESERIALIZING_ISSUE       = "cannot deserialize request"
	SERIALIZING_ISSUE         = "cannot serialize request"
	HASHING_PASSWORD_ISSUE    = "cannot hash password"
	INCORRECT_INPUT           = "incorrect input"
	TOKEN_GENERIC_ERROR       = "cannot create token"
	INVALID_TOKEN             = "invalid token"
	EXPIRED_TOKEN             = "token has expired"
	NO_PERMISSIONS            = "you do not have permissions for this operation"
	CANNOT_PROCESS            = "cannot process request"
	INVALID_PAYLOAD_TYPE      = "invalid payload_type"
)
