// Package auth provides authentication and authorization support.
// Authentication: You are who you say you are.
// Authorization:  You have permission to do what you are requesting to do.
package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userstore/userdb"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"strings"
)

// ErrForbidden is returned when a auth issue is identified.
var ErrForbidden = errors.New("attempted action is not allowed")

// Claims represents the authorization claims transmitted via a JWT.
type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

// KeyLookup declares a method set of behavior for looking up
// private and public keys for JWT use. The return could be a
// PEM encoded string or a JWS based key.
type KeyLookup interface {
	PrivateKey(kid string) (key string, err error)
	PublicKey(kid string) (key string, err error)
}

// Config represents information required to initialize auth.
type Config struct {
	Log       *logger.Logger
	DB        *sqlx.DB
	KeyLookup KeyLookup
	Issuer    string
}

// Auth is used to authenticate clients. It can generate a token for a
// set of user claims and recreate the claims by parsing the token.
type Auth struct {
	keyLookup KeyLookup
	userBus   *userbus.Business
	method    jwt.SigningMethod
	parser    *jwt.Parser
	issuer    string
}

// New creates an Auth to support authentication/authorization.
func New(cfg Config) (*Auth, error) {
	// If a database connection is not provided, we won't perform the user enabled check.
	var userBus *userbus.Business
	if cfg.DB != nil {
		userBus = userbus.NewBusiness(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB))
	}

	a := Auth{
		keyLookup: cfg.KeyLookup,
		userBus:   userBus,
		method:    jwt.GetSigningMethod(jwt.SigningMethodRS256.Name),
		parser:    jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name})),
		issuer:    cfg.Issuer,
	}

	return &a, nil
}

// Issuer provides the configured issuer used to authenticate tokens.
func (a *Auth) Issuer() string {
	return a.issuer
}

// GenerateToken generates a signed JWT token string representing the user Claims.
func (a *Auth) GenerateToken(kid string, claims Claims) (string, error) {
	token := jwt.NewWithClaims(a.method, claims)
	token.Header["kid"] = kid

	privateKeyPEM, err := a.keyLookup.PrivateKey(kid)
	if err != nil {
		return "", fmt.Errorf("private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return "", fmt.Errorf("parsing private pem: %w", err)
	}

	str, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return str, nil
}

// Authenticate processes the token to validate the sender's token is valid.
func (a *Auth) Authenticate(ctx context.Context, bearerToken string) (Claims, error) {
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		return Claims{}, errors.New("expected authorization header format: Bearer <token>")
	}

	jwtTokenStr := bearerToken[7:]

	var claims Claims
	token, _, err := a.parser.ParseUnverified(jwtTokenStr, &claims)
	if err != nil {
		return Claims{}, fmt.Errorf("error parsing token: %w", err)
	}

	kidRaw, exists := token.Header["kid"]
	if !exists {
		return Claims{}, fmt.Errorf("kid missing from header: %w", err)
	}

	kid, ok := kidRaw.(string)
	if !ok {
		return Claims{}, fmt.Errorf("kid malformed: %w", err)
	}

	pem, err := a.keyLookup.PublicKey(kid)
	if err != nil {
		return Claims{}, fmt.Errorf("failed to fetch public key: %w", err)
	}

	_, err = a.parser.Parse(jwtTokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
	})
	if err != nil {
		return Claims{}, fmt.Errorf("authentication failed : %w", err)
	}

	// Check the database for this user to verify they are still enabled.
	if err := a.isValidUser(ctx, claims); err != nil {
		return Claims{}, fmt.Errorf("invalid user : %w", err)
	}

	return claims, nil
}

// Authorize attempts to authorize the user with the provided input roles, if
// none of the input roles are within the user's claims, we return an error
// otherwise the user is authorized.
func (a *Auth) Authorize(_ context.Context, claims Claims, userID uuid.UUID, rule Rule) error {
	var (
		roles userbus.RolesList
		err   error
	)
	roles, err = userbus.ParseRoles(claims.Roles)
	if err != nil {
		return fmt.Errorf("parsing roles: %w", err)
	}

	switch rule {
	case Rules.Any:
		if !roles.Contains(userbus.Roles.Admin) && !roles.Contains(userbus.Roles.User) {
			return fmt.Errorf("rule_any: %w", ErrForbidden)
		}
	case Rules.Admin:
		if !roles.Contains(userbus.Roles.Admin) {
			return fmt.Errorf("rule_admin_only: %w", ErrForbidden)
		}
	case Rules.User:
		if !roles.Contains(userbus.Roles.User) {
			return fmt.Errorf("rule_user_only: %w", ErrForbidden)
		}
	case Rules.Owner:
		if claims.Subject != userID.String() {
			return fmt.Errorf("rule_owner: %w", ErrForbidden)
		}
	case Rules.AdminOrOwner:
		if !roles.Contains(userbus.Roles.Admin) && claims.Subject != userID.String() {
			return fmt.Errorf("rule_admin_or_owner: %w", ErrForbidden)
		}
	default:
		return fmt.Errorf("unknown rule: %s", rule)
	}

	return nil
}

// isValidUser checks the user is still valid in the system: not disabled and email confirmed.
// If userBus is not provided, we skip this check.
func (a *Auth) isValidUser(ctx context.Context, claims Claims) error {
	if a.userBus == nil {
		return nil
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return fmt.Errorf("parse user: %w", err)
	}

	usr, err := a.userBus.QueryByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("query user: %w", err)
	}

	if !usr.Enabled {
		return fmt.Errorf("user disabled")
	}

	if usr.EmailConfirmToken != "" {
		return fmt.Errorf("user not confirm email")
	}

	return nil
}
