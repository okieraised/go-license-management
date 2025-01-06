package config

const SuperAdminUsername = "superadmin"

const (
	ServerMode           = "server.mode"
	ServerHttpPort       = "server.http_port"
	ServerEnableTLS      = "server.enable_tls"
	ServerCertFile       = "server.cert_file"
	ServerKeyFile        = "server.key_file"
	ServerRequestTimeout = "server.request_timeout"
)

const (
	SuperAdminPassword = "superadmin.password"
)

const (
	TracerURI = "tracer.uri"
)

const (
	PostgresHost     = "postgres.host"
	PostgresPort     = "postgres.port"
	PostgresUsername = "postgres.username"
	PostgresPassword = "postgres.password"
	PostgresDatabase = "postgres.database"
)

const (
	AccessTokenTTL = "access_token.ttl"
)
