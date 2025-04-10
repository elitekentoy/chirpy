package commons

import "time"

const (
	ROOT                 = "."
	LISTENING_PORT       = "8080"
	CHIRP_LIMIT          = 140
	ACCESS_TOKEN_EXPIRY  = time.Duration(1) * time.Hour
	REFRESH_TOKEN_EXPIRY = time.Duration(60*24) * time.Hour
	PLATFORM_KEY         = "PLATFORM"
	DEV                  = "dev"
	DEFAULT_HITS         = 0
	USER_UPGRADED        = "user.upgraded"
	CHIRP_LIMIt          = 140
	MASK                 = "****"
)
