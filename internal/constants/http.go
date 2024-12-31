package constants

import "go-license-management/internal/utils"

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
	ContentDigestHeader             = "Content-Digest"
	ContentTransferEncodingHeader   = "Content-Transfer-Encoding"
	ContentDescriptionHeader        = "Content-Description"
	OriginHeader                    = "Origin"
	XRequestIDHeader                = "X-Request-ID"
	XRequestedWithHeader            = "X-Requested-With"
	XAPIKeyHeader                   = "X-API-Key"
	XRateLimitWindowHeader          = "X-RateLimit-Window" // The current rate limiting window that is closest to being reached, percentage-wise.
	XRateLimitCountHeader           = "X-RateLimit-Count"  // The number of requests that have been performed within the current rate limit window.
	XRateLimitLimitHeader           = "X-RateLimit-Limit"  // 	The maximum number of requests that the IP is permitted to make for the current window.
	RetryAfterHeader                = "Retry-After"
	XRateLimitRemainingHeader       = "X-RateLimit-Remaining" //	The number of requests remaining in the current rate limit window.
	XRateLimitResetHeader           = "X-RateLimit-Reset"     //	The time at which the current rate limit window resets in UTC epoch seconds.
)

const (
	AuthorizationTypeBearer = "Bearer"
)

const (
	ContextValuePermissions = "permissions"
	ContextValueTenant      = "tenant"
	ContextValueSubject     = "subject"
	ContextValueAudience    = "audience"
)

type QueryCommonParam struct {
	Limit  *int `form:"limit" validate:"optional" example:"10"`
	Offset *int `form:"offset" validate:"optional" example:"10"`
}

func (req *QueryCommonParam) Validate() {
	if req.Limit == nil {
		req.Limit = utils.RefPointer(100)
	}

	if req.Offset == nil {
		req.Offset = utils.RefPointer(0)
	}
}
