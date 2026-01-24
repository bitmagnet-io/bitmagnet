package auth

import (
	"github.com/bitmagnet-io/bitmagnet/internal/auth/api_key"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/http_auth"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/identity"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/jwt"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/user"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/postgres/migrator"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type deps struct {
	fx.In
	HTTPMiddleware http_auth.Middleware
	UserService    user.Service
	Logger         *zap.Logger
}

type service rbac.Service

var (
	Ref               = ref.Root.MustSub("auth")
	RefJWT            = Ref.MustSub("jwt")
	RefPasswordPolicy = Ref.MustSub("password_policy")
	RefRBAC           = Ref.MustSub("rbac")
	RefUser           = Ref.MustSub("user")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides user authentication and authorization services"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithConfig[deps](Ref.MustSub("jwt_secret"), jwt.ParamSecret),
		builder.WithConfig[deps](Ref.MustSub("jwt_duration"), jwt.ParamDuration),
		builder.WithConfig[deps](Ref.MustSub("anonymous_access"), rbac.ParamAnonymousAccess),
		builder.WithConfig[deps](Ref.MustSub("rbac_cache_ttl"), rbac.ParamCacheTTL),
		builder.WithConfig[deps](Ref.MustSub("invitation_required"), user.ParamInvitationRequired),
		builder.WithConfig[deps](Ref.MustSub("email_required"), user.ParamEmailRequired),
		builder.WithConfig[deps](Ref.MustSub("email_verification"), user.ParamEmailVerification),
		builder.WithConfig[deps](Ref.MustSub("password_min_entropy"), user.ParamPasswordMinEntropy),
		builder.WithConfig[deps](Ref.MustSub("password_hashing_cost"), user.ParamPasswordHashingCost),
		builder.WithConfig[deps](Ref.MustSub("login_requests_per_minute"), user.ParamLoginRequestsPerMinute),
		builder.WithConfig[deps](Ref.MustSub("login_request_burst"), user.ParamLoginRequestBurst),
		builder.WithError[deps](RefUser.MustSub("user_already_exists"), user.ErrAlreadyExists),
		builder.WithError[deps](
			RefUser.MustSub("password_insufficient_entropy"),
			user.ErrPasswordInsufficientEntropy,
		),
		builder.WithError[deps](RefUser.MustSub("username_invalid"), user.ErrUsernameInvalid),
		builder.WithError[deps](RefUser.MustSub("email_invalid"), user.ErrEmailInvalid),
		builder.WithError[deps](RefUser.MustSub("email_missing"), user.ErrEmailMissing),
		builder.WithError[deps](RefUser.MustSub("invitation_code_missing"), user.ErrInvitationCodeMissing),
		builder.WithError[deps](RefUser.MustSub("invitation_not_found"), user.ErrInvitationNotFound),
		builder.WithError[deps](RefUser.MustSub("invitation_expired"), user.ErrInvitationExpired),
		builder.WithError[deps](RefUser.MustSub("invitation_claimed"), user.ErrInvitationClaimed),
		builder.WithError[deps](RefUser.MustSub("credentials_invalid"), user.ErrCredentialsInvalid),
		builder.WithError[deps](RefUser.MustSub("account_disabled"), user.ErrDisabled),
		builder.WithFxOption[deps](
			fx.Provide(
				user.NewService,
				jwt.NewService,
				api_key.NewRepository,
				api_key.NewService,
				identity.NewAuthenticator,
				http_auth.NewMiddleware,
				fx.Annotate(
					func(providers []rbac.ObjectActionProvider) rbac.ObjectActionProvider {
						return rbac.ObjectActionProviders(providers...)
					},
					fx.ParamTags(`group:"auth_object_actions"`),
				),
				fx.Annotate(
					func() rbac.PermissionProvider {
						return rbac.CorePermissions
					},
					fx.ResultTags(`group:"auth_permissions"`),
				),
				fx.Annotate(
					rbac.VerbatimPermissions,
					fx.ResultTags(`group:"auth_permissions"`),
				),
				fx.Annotate(
					func(providers []rbac.PermissionProvider) rbac.PermissionProvider {
						return rbac.PermissionProviders(providers...)
					},
					fx.ParamTags(`group:"auth_permissions"`),
				),
				func(
					dao database.DaoTransactionProvider,
					objectActions rbac.ObjectActionProvider,
					permissions rbac.PermissionProvider,
					ttl rbac.CacheTTL,
				) service {
					return rbac.NewService(
						rbac.NewRepository(dao),
						objectActions,
						permissions,
						ttl,
					)
				},
				fx.Annotate(
					rbac.NewServiceLazy,
					fx.As(new(rbac.Enforcer)),
					fx.As(new(rbac.Repository)),
					fx.As(new(rbac.Service)),
					fx.As(new(rbac.ServiceLazy)),
				),
			),
			fx.Invoke(func(service service, lazy rbac.ServiceLazy) error {
				return lazy.SetService(service)
			}),
		),
		builder.WithGinOption(Ref, httpserver.PhasePre, func(deps deps) gin.OptionFunc {
			return func(e *gin.Engine) {
				e.Use(deps.HTTPMiddleware.AttachAuth())
			}
		}),
		// todo: Move to plugin
		builder.WithWorker(
			func(deps deps) (runner.Provider, worker.Option) {
				return &initialInvitationWorker{
						userService: deps.UserService,
						logger:      deps.Logger,
					}, worker.Options(
						worker.WithDependencies(migrator.Ref, postgres.Ref),
						worker.ShortLived(),
						worker.WithAutostart(true),
					)
			},
		),
	)
)
