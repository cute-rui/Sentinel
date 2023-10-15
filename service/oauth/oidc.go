package oauth

import (
	"Sentinel/utils/config"
	utils "Sentinel/utils/string"
	"crypto/sha256"
	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"
	"net/http"
)

type OIDC struct {
	storage *Storage
	issuer  string

	provider op.OpenIDProvider
	key      [32]byte
}

// Todo: refactor all of this package
func SetupOIDC() *OIDC {
	RegisterClients(
		WebClient(
			config.Conf.GetString("OIDC.WebID"),
			config.Conf.GetString("OIDC.WebSecret"),
			config.Conf.GetStringSlice("OIDC.RedirectAllowedList")...,
		),
	)

	storage := NewStorage(userStore{})

	issuer := config.Conf.GetString("OIDC.Issuer")
	key := sha256.Sum256(utils.StringToByte(config.Conf.GetString("OIDC.32BKey")))

	provider, err := newOP(storage, issuer, key)
	if err != nil {
		panic(err)
	}

	return &OIDC{
		storage:  storage,
		issuer:   issuer,
		provider: provider,
		key:      key,
	}

}

func (i *OIDC) DiscoveryHandler() http.HandlerFunc {
	return op.NewIssuerInterceptor(i.provider.IssuerFromRequest).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		op.Discover(w, op.CreateDiscoveryConfig(r.Context(), i.provider, i.provider.Storage()))
	})
}

func (i *OIDC) GetMainHandler() http.Handler {
	return http.Handler(i.provider)
}

func newOP(storage op.Storage, issuer string, key [32]byte) (op.OpenIDProvider, error) {
	c := &op.Config{
		CryptoKey: key,

		// will be used if the end_session endpoint is called without a post_logout_redirect_uri
		//DefaultLogoutRedirectURI: pathLoggedOut,

		// enables code_challenge_method S256 for PKCE (and therefore PKCE in general)
		CodeMethodS256: true,

		// enables additional client_id/client_secret authentication by form post (not only HTTP Basic Auth)
		AuthMethodPost: true,

		// enables additional authentication by using private_key_jwt
		AuthMethodPrivateKeyJWT: true,

		// enables refresh_token grant use
		GrantTypeRefreshToken: true,

		// enables use of the `request` Object parameter
		RequestObjectSupported: true,

		// this example has only static texts (in English), so we'll set the here accordingly
		SupportedUILocales: []language.Tag{language.English},

		/*DeviceAuthorization: op.DeviceAuthorizationConfig{
			Lifetime:     5 * time.Minute,
			PollInterval: 5 * time.Second,
			UserFormPath: "/device",
			UserCode:     op.UserCodeBase20,
		},*/
	}
	handler, err := op.NewOpenIDProvider(issuer, c, storage,
		append([]op.Option{
			//we must explicitly allow the use of the http issuer
			op.WithAllowInsecure(),
			// as an example on how to customize an endpoint this will change the authorization_endpoint from /authorize to /auth
			op.WithCustomEndpoints(
				op.NewEndpoint("oauth2/oidc/authorize"),
				op.NewEndpoint("oauth2/oidc/token"),
				op.NewEndpoint("oauth2/oidc/userinfo"),
				op.NewEndpoint("oauth2/oidc/revoke"),
				op.NewEndpoint("oauth2/oidc/end_session"),
				op.NewEndpoint("oauth2/oidc/keys"),
			),
			op.WithCustomIntrospectionEndpoint(op.NewEndpoint("oauth2/oidc/introspect")),

			op.WithCustomDeviceAuthorizationEndpoint(op.NewEndpoint("oauth2/oidc/device")),

			// Pass our logger to the OP
			//op.WithLogger(logger.WithGroup("op")),
		})...,
	)
	if err != nil {
		return nil, err
	}
	return handler, nil
}
