package constant

// TokenPayloadLen
// 4 int migration
// 8 int8 migration
// 24 string
const (
	TokenUserContext     = "usr"
	TokenContentLen      = 6
	FormatTimeLayout     = "15:04"
	FormatDateLayout     = "02 January 2006"
	FormatDatetimeLayout = "02 January 2006 15:04"
	BearerSchema         = "Bearer"
	RoleAdmin            = "ADMIN"
	RoleUseradmin        = "USERADMIN"
	RoleUser             = "USER"
	AuthHeaderKey        = "Authorization"
)
