# Sprint 6: Token Revocation Implementation

## Overview
Token revocation is now fully implemented using Redis as a JWT blocklist. This allows immediate invalidation of JWT access tokens before their natural expiration.

## Architecture

### Components
1. **TokenRevocationService** - Interface for token revocation operations
2. **redisRevocationService** - Redis-based implementation with TTL
3. **noopRevocationService** - No-op implementation when Redis is disabled

### Key Features
- ✅ Individual token revocation by token ID
- ✅ Subject-based revocation (all tokens for a user/client)
- ✅ Automatic TTL based on token expiration
- ✅ Fail-open by default (if Redis is unavailable, tokens are not blocked)
- ✅ Integrated with ValidateAccessToken (checks revocation on every validation)

## Usage

### 1. Initialize Token Revocation Service

```go
import (
    "github.com/aystek/aystek-core/internal/modules/auth/service"
    "github.com/aystek/aystek-core/internal/platform/redis"
)

// Option A: With Redis (recommended for production)
redisClient := redis.Client() // Get initialized Redis client
var revocationSvc service.TokenRevocationService
if redisClient != nil {
    revocationSvc = service.NewRedisRevocationService(redisClient)
} else {
    // Fallback to no-op if Redis is disabled
    revocationSvc = service.NewNoopRevocationService()
}

// Option B: No-op (for development/testing without Redis)
revocationSvc := service.NewNoopRevocationService()
```

### 2. Create OAuth2 Service with Revocation

```go
oauth2Service := service.NewOAuth2Service(
    usersRepo,
    oauth2Repo,
    revocationSvc,  // <-- Token revocation service
    secretManager,
    rsaKeyPair,
    issuer,
)
```

### 3. Revoke Individual Token

```go
// Via OAuth2 revocation endpoint (POST /oauth2/revoke)
req := &dto.OAuth2RevocationRequest{
    Token:         "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
    TokenTypeHint: "access_token",
}

err := oauth2Service.RevokeToken(ctx, req)
// Token is now in Redis blocklist with TTL = token expiration time
```

### 4. Revoke All Tokens for Subject (Logout All Devices)

```go
// Revoke all tokens for a specific user or client
err := oauth2Service.RevokeAllTokensForSubject(ctx, "user_id_or_client_id")
// All future token validations for this subject will fail until TTL expires
```

### 5. Token Validation (Automatic Revocation Check)

```go
// ValidateAccessToken now automatically checks revocation
token, err := oauth2Service.ValidateAccessToken(ctx, tokenString)
if err != nil {
    // Error could be:
    // - "token has been revoked" (individual revocation)
    // - "all tokens for this subject have been revoked" (subject revocation)
    // - Other validation errors (expired, invalid signature, etc.)
}
```

## Redis Key Structure

### Individual Token Revocation
- **Key**: `oauth2:revoked:{tokenID}`
- **Value**: Unix timestamp of revocation
- **TTL**: Time until token naturally expires
- **Example**: `oauth2:revoked:user123:1709836800` → `1709836700`

### Subject-Based Revocation
- **Key**: `oauth2:revoked:subject:{subject}`
- **Value**: Unix timestamp of revocation
- **TTL**: Maximum token lifetime (15 minutes for access tokens)
- **Example**: `oauth2:revoked:subject:user123` → `1709836700`

## Token ID Generation

The system uses two methods to identify tokens:

1. **JTI Claim (preferred)**: If JWT includes `jti` claim, use it directly
2. **Fallback**: Generate ID as `{sub}:{exp}` (subject + expiration timestamp)

## Security Considerations

### Fail-Open vs Fail-Closed
Current implementation: **Fail-open** (if Redis is unavailable, tokens are not blocked)

**Rationale**:
- Prioritizes availability over security in edge cases
- Redis failures won't bring down the entire auth system
- Tokens still expire naturally via JWT expiration

**To change to fail-closed** (reject tokens if Redis is unavailable):
```go
// In ValidateAccessToken, change:
if err != nil {
    // Current: log and continue
    // Fail-closed: return nil, fmt.Errorf("revocation check failed: %w", err)
}
```

### TTL Optimization
- Individual token revocations only persist until natural expiration
- Subject revocations persist for max token lifetime (15 min for access tokens)
- No need for manual cleanup - Redis TTL handles expiration

### Performance
- Redis GET operations: O(1) lookup
- Minimal latency impact on token validation
- Consider Redis clustering for high-volume scenarios

## Testing

### Test Token Revocation
```go
// 1. Generate token
tokenResp, _ := oauth2Service.HandleTokenRequest(ctx, &dto.OAuth2TokenRequest{
    GrantType:    "client_credentials",
    ClientID:     "test_client",
    ClientSecret: "secret",
})

// 2. Validate token (should succeed)
token, err := oauth2Service.ValidateAccessToken(ctx, tokenResp.AccessToken)
assert.NoError(t, err)
assert.True(t, token.Valid)

// 3. Revoke token
err = oauth2Service.RevokeToken(ctx, &dto.OAuth2RevocationRequest{
    Token: tokenResp.AccessToken,
})
assert.NoError(t, err)

// 4. Validate again (should fail)
token, err = oauth2Service.ValidateAccessToken(ctx, tokenResp.AccessToken)
assert.Error(t, err)
assert.Contains(t, err.Error(), "revoked")
```

### Test Subject Revocation
```go
// 1. Generate tokens for user
token1, _ := oauth2Service.ExchangeAuthorizationCode(ctx, req1)
token2, _ := oauth2Service.ExchangeAuthorizationCode(ctx, req2)

// 2. Revoke all tokens for user
err := oauth2Service.RevokeAllTokensForSubject(ctx, "user123")
assert.NoError(t, err)

// 3. Both tokens should now fail validation
_, err1 := oauth2Service.ValidateAccessToken(ctx, token1.AccessToken)
_, err2 := oauth2Service.ValidateAccessToken(ctx, token2.AccessToken)
assert.Error(t, err1)
assert.Error(t, err2)
assert.Contains(t, err1.Error(), "subject have been revoked")
```

## Monitoring

### Key Metrics to Track
1. **Revocation Rate**: Number of tokens revoked per hour
2. **Revocation Check Latency**: Time taken for Redis lookup
3. **Redis Errors**: Failed revocation checks
4. **Subject Revocations**: Frequency of logout-all-devices operations

### Recommended Alerts
- Redis connection failures
- High revocation check latency (> 50ms)
- Unusual spike in revocations (potential security incident)

## Migration Notes

### Upgrading from No Revocation
1. Deploy Redis if not already available
2. Update OAuth2Service initialization to include revocationSvc
3. Existing tokens will continue to work (not retroactively revoked)
4. New revocations will be stored in Redis going forward

### Rolling Back
If you need to disable revocation:
```go
// Use no-op service instead of Redis service
revocationSvc := service.NewNoopRevocationService()
```
All revocation calls will succeed but have no effect.

## Future Enhancements

### Potential Improvements
1. **MongoDB Fallback**: Use MongoDB if Redis is unavailable
2. **Distributed Cache**: Use Redis Cluster for geo-distributed deployments
3. **Revocation Events**: Emit events when tokens are revoked for audit trails
4. **Batch Revocation**: Revoke multiple tokens in a single Redis pipeline
5. **JTI Enforcement**: Require JTI claim in all JWTs for better tracking

### Refresh Token Revocation
Note: OAuth2 refresh tokens (authorization code flow) are already revoked via:
- `RotatedAt` field in `OAuth2RefreshToken` model
- MongoDB tracking (not Redis)
- Checked during `HandleRefreshTokenRequest()`

This Redis revocation is specifically for JWT access tokens.
