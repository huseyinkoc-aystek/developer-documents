# TASK 055 - Auth Screens Loading UI

## Scope

`aystek-core-next` auth ekranlarında login, şifre hatırlatma ve şifre sıfırlama akışlarının loading/feedback görsel dili toparlandı.

## Changed Areas

- `/workspace/node/aystek-core-next/app/(auth)/login/page.tsx`
- `/workspace/node/aystek-core-next/app/(auth)/forgot-password/page.tsx`
- `/workspace/node/aystek-core-next/app/(auth)/reset-password/page.tsx`
- `/workspace/node/aystek-core-next/components/auth/AuthLogo.tsx`
- `/workspace/node/aystek-core-next/components/auth/AuthNotice.tsx`
- `/workspace/node/aystek-core-next/components/auth/AuthSubmitButton.tsx`
- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`

## Implementation Notes

- Login, forgot password and reset password cards now share the same card surface, logo treatment and input focus/loading behavior.
- Submit actions use `AuthSubmitButton` with a spinner state and a directional icon for idle state.
- Success, warning, info and error states use `AuthNotice` with lucide icons instead of plain text blocks.
- Logo rendering was centralized in `AuthLogo`, including fallback behavior for missing tenant/public logo assets.
- Reset password no longer creates a second full-screen layout inside the auth route wrapper.
- New user-visible text key `login.loading.session` was added to `translation_seed.go` with 28 language entries.

## Validation

- `npx eslint app/(auth)/login/page.tsx app/(auth)/forgot-password/page.tsx app/(auth)/reset-password/page.tsx components/auth/AuthLogo.tsx components/auth/AuthNotice.tsx components/auth/AuthSubmitButton.tsx`
- `go test ./internal/seeds`
- Dev server route smoke checks:
  - `GET /login` -> `200`
  - `GET /forgot-password` -> `200`
  - `GET /reset-password?token=test` -> `200`

## Known Repository-Wide Issues

- Full `npm run lint` still fails because of existing unrelated lint errors outside the auth files.
- Full `npx tsc --noEmit` still fails because of existing unrelated TypeScript errors in dashboard and API modules.
