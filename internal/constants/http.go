package constants

const (
	ContentDispositionInline     = "inline"
	ContentDispositionAttachment = "attachment; filename=%s"
)
const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
	ContentTypeImage  = "image/%s"
	ContentTypeXML    = "text/xml"
	ContentTypePDF    = "application/pdf"
)

const (
	MimeTypeJPEG      = "image/jpeg"
	MimeTypePNG       = "image/png"
	MimeTypePlainText = "text/plain"
	MimeTypeZip       = "application/zip"
	MimeTypeRar       = "application/x-rar-compressed"
	MimeTypeXML       = "application/xml"
)

const (
	AllowAllOrigins = "*"
)

const (
	AcceptHeader                    = "Accept"
	AcceptLanguageHeader            = "Accept-Language"
	AccessControlAllowHeadersHeader = "Access-Control-Allow-Headers"
	AuthorizationHeader             = "Authorization"
	ContentLengthHeader             = "Content-Length"
	ContentTypeHeader               = "Content-Type"
	ContentDispositionHeader        = "Content-Disposition"
	ContentTransferEncodingHeader   = "Content-Transfer-Encoding"
	ContentDescriptionHeader        = "Content-Description"
	OriginHeader                    = "Origin"
	XRequestIDHeader                = "X-Request-ID"
	XRequestedWithHeader            = "X-Requested-With"
	XAPIKeyHeader                   = "X-API-Key"
)

const (
	AuthorizationTypeBearer = "Bearer"
)
