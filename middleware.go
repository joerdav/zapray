package zapray

import (
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
)

// NewMiddlware creates X-Ray middleware that creates a subsegment
// for each HTTP request.
func NewMiddleware(appName string, next http.Handler) http.Handler {
	return Middleware{
		Next: next,
		name: appName,
	}
}

// Middlware applies X-Ray segements to the wrapped handler.
type Middleware struct {
	Next http.Handler
	name string
}

func (h Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, s := xray.BeginSubsegment(r.Context(), "xrayMiddleware")
	defer s.Close(nil)
	xray.HandlerWithContext(ctx, xray.NewFixedSegmentNamer(h.name+"Handler"), setParent(s, h.Next)).ServeHTTP(w, r.WithContext(ctx))
}

func setParent(s *xray.Segment, h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		seg := xray.GetSegment(r.Context())
		if seg != nil {
			seg.ParentID = s.ID
		}
		h.ServeHTTP(rw, r)
	})
}
