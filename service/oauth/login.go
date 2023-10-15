package oauth

import (
	"Sentinel/utils/auth"
	"Sentinel/utils/config"
	"github.com/zitadel/oidc/v3/pkg/op"
	"net/http"
	"strconv"
)

// Todo:
var defaultFrontEndLoginURL = func(id string) string {
	return "/login/username?authRequestID=" + id
}

// Todo: refactor dirty method
func (i *OIDC) LoginFunc() http.HandlerFunc {
	return op.NewIssuerInterceptor(i.provider.IssuerFromRequest).HandlerFunc(i.checkLoginHandler)
}

func (i *OIDC) checkLoginHandler(w http.ResponseWriter, r *http.Request) {
	jwt := r.Header.Get(`Authorization`)
	if jwt == `` {
		http.Error(w, `Bad Request`, http.StatusBadRequest)
		return
	}

	reqID := r.Header.Get(`ReqID`)
	if reqID == `` {
		http.Error(w, `Bad Request`, http.StatusBadRequest)
		return
	}

	claim, err := auth.StringToJWTClaim(auth.TrimBearerScheme(jwt), config.Conf.GetString(`JWT.AccessSecret`))
	if err != nil {
		http.Error(w, `Token Invalid`, http.StatusBadRequest)
		return
	}
	err = i.storage.AddUIDToReqPayload(reqID, strconv.Itoa(claim.UserID))
	if err != nil {
		http.Error(w, `Token Invalid`, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, op.AuthCallbackURL(i.provider)(r.Context(), reqID), http.StatusFound)
}
