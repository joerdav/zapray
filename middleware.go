package zapray

import (
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
)

// NewMiddlware creates X-Ray middleware that creates a subsegment
// for each HTTP request.
func NewMiddleware(next http.Handler) http.Handler {
	return Middleware{
		Next: next,
	}
}

// Middlware applies X-Ray segements to the wrapped handler.
type Middleware struct {
	Next http.Handler
}

func (h Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, seg := xray.BeginSubsegment(r.Context(), "zaprayTrace")
	h.Next.ServeHTTP(w, r.WithContext(ctx))
	seg.Close(nil)
}
