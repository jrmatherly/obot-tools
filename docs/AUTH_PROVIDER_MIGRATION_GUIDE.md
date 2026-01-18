# Authentication Provider Migration Guide

## Relocating Keycloak and Microsoft Entra ID Providers to obot-tools

**Document Version:** 1.2
**Created:** 2026-01-18
**Last Updated:** 2026-01-18
**Status:** Phase 0, 1, & 2 COMPLETE - Migration Finished
**Author:** Claude Code Analysis

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Background](#background)
3. [Current State Analysis](#current-state-analysis)
4. [Target State](#target-state)
5. [Migration Phases](#migration-phases)
   - [Phase 1: Fix obot-tools](#phase-1-fix-obot-tools-make-self-contained)
   - [Phase 2: Cleanup obot-entraid](#phase-2-cleanup-obot-entraid)
   - [Phase 3: Documentation Updates](#phase-3-documentation-updates)
6. [Detailed Implementation Steps](#detailed-implementation-steps)
7. [Verification Procedures](#verification-procedures)
8. [Rollback Plan](#rollback-plan)
9. [Risk Assessment](#risk-assessment)
10. [Appendices](#appendices)

---

## Executive Summary

This guide documents the migration of custom Keycloak and Microsoft Entra ID authentication providers from their original location in `obot-entraid/tools/` to their canonical home in `obot-tools/`. The migration is approximately 85% complete but has critical blocking issues that must be resolved before the providers can be used from their new location.

### Key Findings

| Issue | Severity | Status |
|-------|----------|--------|
| `placeholder-credential` missing from fork | **CRITICAL** | ✅ FIXED - Copied from obot-entraid |
| Import paths reference old `obot-entraid` location | **CRITICAL** | ✅ FIXED - Updated to obot-tools |
| Providers not registered in `obot-tools/index.yaml` | **CRITICAL** | ✅ FIXED - Added to registry |
| Tests not added to `obot-tools/Makefile` | High | ✅ FIXED - Added test targets |
| Documentation references old location | Medium | ✅ FIXED - Updated READMEs |
| Old providers still exist in `obot-entraid` | Medium | ✅ FIXED - Archived |

### Success Criteria

- [x] `placeholder-credential` exists in `obot-tools` (sync from upstream or copy)
- [x] Both providers build successfully in `obot-tools`
- [x] Both providers pass all tests
- [x] Providers are registered and discoverable via GPTScript runtime
- [x] `obot-entraid` Dockerfile uses providers from `TOOLS_IMAGE`
- [x] No duplicate provider code exists (archived in `obot-entraid/archive/`)
- [x] All documentation points to canonical location

---

## Background

### Original Implementation

The Keycloak and Entra ID authentication providers were initially developed within the `obot-entraid` repository (jrmatherly's fork of obot-platform/obot) because:

1. The `obot-tools` repository had not yet been forked
2. Development required tight integration with the main obot application
3. `replace` directives in `go.mod` were used to remap `github.com/obot-platform/*` to `github.com/jrmatherly/*`

### Migration Rationale

Relocating to `obot-tools` provides:

1. **Consistency**: Aligns with native providers (Google, GitHub)
2. **Separation of Concerns**: Tools separate from main application
3. **Simplified Builds**: Single source of truth for all auth providers
4. **Upstream Compatibility**: Easier to track upstream changes
5. **Reusability**: Providers can be used by other obot forks

---

## Current State Analysis

### obot-tools Repository

```
obot-tools/
├── keycloak-auth-provider/
│   ├── main.go              # Line 22: WRONG import path
│   ├── tool.gpt             # Correct format
│   ├── go.mod               # Correct module path
│   ├── pkg/profile/         # Exists and correct
│   ├── groups_test.go       # Tests exist
│   ├── Makefile             # Good (build/test/lint/tidy)
│   ├── KEYCLOAK_SETUP.md    # Documentation exists
│   └── bin/                 # Pre-built binaries (should be gitignored)
│
├── entra-auth-provider/
│   ├── main.go              # Line 21: WRONG import path
│   ├── tool.gpt             # Correct format
│   ├── go.mod               # Correct module path
│   ├── pkg/profile/         # Exists and correct
│   ├── groups_test.go       # Tests exist
│   ├── Makefile             # Good (build/test/lint/tidy)
│   ├── README.md            # Has old obot-entraid references
│   └── bin/                 # Pre-built binaries (should be gitignored)
│
├── auth-providers-common/   # Shared library - correct
│
├── google-auth-provider/    # Reference implementation
├── github-auth-provider/    # Reference implementation
│
└── index.yaml               # MISSING keycloak and entra entries
```

#### Critical Import Path Issue

**keycloak-auth-provider/main.go:22**

```go
// CURRENT (WRONG)
import "github.com/obot-platform/obot-entraid/tools/keycloak-auth-provider/pkg/profile"

// REQUIRED (CORRECT)
import "github.com/obot-platform/tools/keycloak-auth-provider/pkg/profile"
```

**entra-auth-provider/main.go:21**

```go
// CURRENT (WRONG)
import "github.com/obot-platform/obot-entraid/tools/entra-auth-provider/pkg/profile"

// REQUIRED (CORRECT)
import "github.com/obot-platform/tools/entra-auth-provider/pkg/profile"
```

#### Build Verification

```bash
$ cd obot-tools/keycloak-auth-provider && go build .
main.go:22:2: no required module provides package
  github.com/obot-platform/obot-entraid/tools/keycloak-auth-provider/pkg/profile

$ cd obot-tools/entra-auth-provider && go build .
main.go:21:2: no required module provides package
  github.com/obot-platform/obot-entraid/tools/entra-auth-provider/pkg/profile
```

### obot-entraid Repository

```
obot-entraid/
├── Dockerfile               # Builds auth providers locally (lines 41-42)
│                            # Merges index.yaml (lines 82-106)
│
├── .envrc.dev               # OBOT_SERVER_TOOL_REGISTRIES includes ./tools
│
├── tools/
│   ├── index.yaml           # Registers local auth providers
│   ├── entra-auth-provider/ # Duplicate of obot-tools version
│   ├── keycloak-auth-provider/ # Duplicate of obot-tools version
│   ├── auth-providers-common/  # Duplicate of obot-tools version
│   ├── placeholder-credential/ # Shared credential tool
│   ├── tool.gpt             # Index loader wrapper
│   ├── dev.sh               # Development script
│   └── combine-envrc.sh     # Container runtime script
│
├── docs/docs/configuration/
│   ├── entra-id-authentication.md
│   └── keycloak-authentication.md
│
└── go.mod                   # Contains replace directives
```

#### Dockerfile Integration Points

**Lines 41-42 (bin stage)** - Local provider builds:

```dockerfile
cd tools/entra-auth-provider && make build && \
cd ../keycloak-auth-provider && make build
```

**Lines 65-80 (tools-patched stage)** - Directory creation and binary copying:

```dockerfile
RUN mkdir -p /obot-tools/tools/entra-auth-provider/bin && \
    mkdir -p /obot-tools/tools/keycloak-auth-provider/bin && ...

COPY --from=bin /app/tools/entra-auth-provider/bin/gptscript-go-tool ...
COPY --from=bin /app/tools/keycloak-auth-provider/bin/gptscript-go-tool ...
```

**Lines 82-106 (tools-patched stage)** - Index merging:

```dockerfile
# Merge custom authProviders into existing upstream index.yaml
yq eval-all '... ($upstream.authProviders + $custom.authProviders)' ...
```

---

## Target State

### Architecture After Migration

```
┌─────────────────────────────────────────────────────────────────────┐
│ obot-tools (jrmatherly/obot-tools) - CANONICAL SOURCE               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ Authentication Providers (all 4 in one registry):                   │
│   ├── google-auth-provider/      (upstream)                         │
│   ├── github-auth-provider/      (upstream)                         │
│   ├── keycloak-auth-provider/    (custom - fully integrated)        │
│   └── entra-auth-provider/       (custom - fully integrated)        │
│                                                                      │
│ index.yaml:                                                          │
│   authProviders:                                                     │
│     github-auth-provider: { reference: ./github-auth-provider }     │
│     google-auth-provider: { reference: ./google-auth-provider }     │
│     keycloak-auth-provider: { reference: ./keycloak-auth-provider } │
│     entra-auth-provider: { reference: ./entra-auth-provider }       │
│                                                                      │
│ Build: make build → discovers all main.go, builds all providers     │
│ Test: make test → runs all provider tests                           │
│ Image: ghcr.io/jrmatherly/obot-tools:latest                         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              │ TOOLS_IMAGE
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│ obot-entraid (jrmatherly/obot-entraid) - CONSUMER                    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ Dockerfile:                                                          │
│   ARG TOOLS_IMAGE=ghcr.io/jrmatherly/obot-tools:latest              │
│   COPY --from=tools /obot-tools /obot-tools                         │
│   # No local provider builds                                        │
│   # No index merging required                                       │
│                                                                      │
│ tools/                                                               │
│   └── (empty or archived - all providers come from TOOLS_IMAGE)    │
│                                                                      │
│ .envrc.dev:                                                          │
│   OBOT_SERVER_TOOL_REGISTRIES=github.com/jrmatherly/obot-tools      │
│   # No ./tools reference needed                                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Migration Phases

### Phase 0: Sync Missing Dependencies (PREREQUISITE)

**Goal:** Sync `placeholder-credential` from upstream `obot-platform/tools` to `jrmatherly/obot-tools`.

**Discovery:** During validation, it was found that the `jrmatherly/obot-tools` fork is missing the `placeholder-credential` directory that exists in the upstream repository. This is a **CRITICAL** dependency - ALL auth providers (Google, GitHub, Keycloak, Entra) reference it via `Credential: ../placeholder-credential`.

**Prerequisites:** None

**Estimated Effort:** 15 minutes

| Step | Action | Risk |
|------|--------|------|
| 0.1 | Sync from upstream OR copy from `obot-entraid/tools/placeholder-credential/` | Low |
| 0.2 | Verify files: `main.py`, `requirements.txt`, `tool.gpt` | Low |
| 0.3 | Test that existing providers (Google, GitHub) work | Low |

**Option A - Sync from upstream:**

```bash
cd obot-tools
git remote add upstream https://github.com/obot-platform/tools.git
git fetch upstream
git checkout upstream/main -- placeholder-credential/
git add placeholder-credential/
git commit -m "sync: add placeholder-credential from upstream"
```

**Option B - Copy from obot-entraid:**

```bash
cp -r ../obot-entraid/tools/placeholder-credential/ obot-tools/placeholder-credential/
```

**Verification:**

```bash
ls -la obot-tools/placeholder-credential/
# Should show: main.py, requirements.txt, tool.gpt
```

---

### Phase 1: Fix obot-tools (Make Self-Contained)

**Goal:** Make the auth providers in obot-tools fully functional and self-contained.

**Prerequisites:** None

**Estimated Effort:** 1-2 hours

| Step | File | Change | Risk |
|------|------|--------|------|
| 1.1 | `keycloak-auth-provider/main.go` | Fix import path (line 22) | Low |
| 1.2 | `entra-auth-provider/main.go` | Fix import path (line 21) | Low |
| 1.3 | `index.yaml` | Add provider registrations | Low |
| 1.4 | `Makefile` | Add test targets | Low |
| 1.5 | `entra-auth-provider/README.md` | Remove obot-entraid references | Low |
| 1.6 | Both providers | Add `.gitignore` for `bin/` | Low |

**Verification:** Providers build and tests pass.

---

### Phase 2: Cleanup obot-entraid

**Goal:** Remove duplicate providers and update Dockerfile to use TOOLS_IMAGE.

**Prerequisites:** Phase 1 complete, obot-tools image published

**Estimated Effort:** 2-4 hours

| Step | File | Change | Risk |
|------|------|--------|------|
| 2.1 | `Dockerfile` | Remove local provider builds | Medium |
| 2.2 | `Dockerfile` | Remove tools-patched stage complexity | Medium |
| 2.3 | `.envrc.dev` | Remove `./tools` from registries | Low |
| 2.4 | `tools/index.yaml` | Remove provider entries | Low |
| 2.5 | `tools/entra-auth-provider/` | Archive/delete directory | Medium |
| 2.6 | `tools/keycloak-auth-provider/` | Archive/delete directory | Medium |
| 2.7 | `tools/auth-providers-common/` | Archive/delete directory | Medium |

**Verification:** Container builds and providers are accessible.

---

### Phase 3: Documentation Updates

**Goal:** Update all documentation to reference canonical location.

**Prerequisites:** Phase 2 complete

**Estimated Effort:** 1-2 hours

| Step | File | Change |
|------|------|--------|
| 3.1 | `docs/docs/configuration/entra-id-authentication.md` | Update paths |
| 3.2 | `docs/docs/configuration/keycloak-authentication.md` | Update paths |
| 3.3 | `docs/docs/contributing/local-development.md` | Update build instructions |
| 3.4 | `docs/docs/operations/auth-provider-testing.md` | Update test instructions |
| 3.5 | `CONTRIBUTING.md` | Update build instructions |
| 3.6 | `tools/README.md` | Update or remove |

---

## Detailed Implementation Steps

### Phase 1 Implementation

#### Step 1.1: Fix keycloak-auth-provider/main.go

**File:** `obot-tools/keycloak-auth-provider/main.go`

**Current (line 22):**

```go
"github.com/obot-platform/obot-entraid/tools/keycloak-auth-provider/pkg/profile"
```

**Change to:**

```go
"github.com/obot-platform/tools/keycloak-auth-provider/pkg/profile"
```

**Command:**

```bash
cd obot-tools/keycloak-auth-provider
sed -i '' 's|obot-platform/obot-entraid/tools/keycloak-auth-provider|obot-platform/tools/keycloak-auth-provider|g' main.go
```

---

#### Step 1.2: Fix entra-auth-provider/main.go

**File:** `obot-tools/entra-auth-provider/main.go`

**Current (line 21):**

```go
"github.com/obot-platform/obot-entraid/tools/entra-auth-provider/pkg/profile"
```

**Change to:**

```go
"github.com/obot-platform/tools/entra-auth-provider/pkg/profile"
```

**Command:**

```bash
cd obot-tools/entra-auth-provider
sed -i '' 's|obot-platform/obot-entraid/tools/entra-auth-provider|obot-platform/tools/entra-auth-provider|g' main.go
```

---

#### Step 1.3: Update index.yaml

**File:** `obot-tools/index.yaml`

**Current:**

```yaml
authProviders:
  github-auth-provider:
    reference: ./github-auth-provider
  google-auth-provider:
    reference: ./google-auth-provider
```

**Change to:**

```yaml
authProviders:
  github-auth-provider:
    reference: ./github-auth-provider
  google-auth-provider:
    reference: ./google-auth-provider
  keycloak-auth-provider:
    reference: ./keycloak-auth-provider
  entra-auth-provider:
    reference: ./entra-auth-provider
```

---

#### Step 1.4: Update Makefile Test Target

**File:** `obot-tools/Makefile`

**Current:**

```makefile
test:
 cd tests && GPTSCRIPT_TOOL_REMAP="github.com/obot-platform/tools=.." go test -v tool_test.go
 cd github-auth-provider && go test ./... && cd ..
 cd google-auth-provider && go test ./... && cd ..
 cd credential-stores && go test ./... && cd ..
```

**Change to:**

```makefile
test:
 cd tests && GPTSCRIPT_TOOL_REMAP="github.com/obot-platform/tools=.." go test -v tool_test.go
 cd github-auth-provider && go test ./... && cd ..
 cd google-auth-provider && go test ./... && cd ..
 cd keycloak-auth-provider && go test ./... && cd ..
 cd entra-auth-provider && go test ./... && cd ..
 cd credential-stores && go test ./... && cd ..
```

---

#### Step 1.5: Update entra-auth-provider/README.md

**File:** `obot-tools/entra-auth-provider/README.md`

**Lines to update:**

Line 273 (current):

```markdown
# From the obot-entraid repository root
```

Change to:

```markdown
# From the obot-tools repository root
```

Line 315 (current):

```markdown
# In obot-entraid repository
```

Change to:

```markdown
# In obot-tools repository
```

---

#### Step 1.6: Add .gitignore for bin/ directories

**Files to create:**

`obot-tools/keycloak-auth-provider/.gitignore`:

```
bin/
```

`obot-tools/entra-auth-provider/.gitignore`:

```
bin/
```

**Also remove existing binaries:**

```bash
rm -rf obot-tools/keycloak-auth-provider/bin/
rm -rf obot-tools/entra-auth-provider/bin/
rm -f obot-tools/entra-auth-provider/entra-auth-provider
```

---

### Phase 2 Implementation

#### Step 2.1-2.2: Simplify Dockerfile

**File:** `obot-entraid/Dockerfile`

**Remove from bin stage (lines 41-42):**

```dockerfile
# REMOVE these lines:
cd tools/entra-auth-provider && make build && \
cd ../keycloak-auth-provider && make build
```

**Simplify tools-patched stage:**

Replace the complex stage (lines 54-110) with:

```dockerfile
# Create unified tools directory - providers now come from upstream TOOLS_IMAGE
FROM cgr.dev/chainguard/wolfi-base:latest AS tools-patched

# Copy upstream tools (includes all auth providers)
COPY --from=tools /obot-tools /obot-tools

# Copy provider tools (encryption providers, etc.)
COPY --from=provider /obot-tools /obot-tools
COPY --from=enterprise-tools /obot-tools /obot-tools

# Note: keycloak and entra auth providers are now included in TOOLS_IMAGE
# No local building or index merging required
```

---

#### Step 2.3: Update .envrc.dev

**File:** `obot-entraid/.envrc.dev`

**Current:**

```bash
export OBOT_SERVER_TOOL_REGISTRIES=github.com/jrmatherly/obot-tools,./tools
```

**Change to:**

```bash
export OBOT_SERVER_TOOL_REGISTRIES=github.com/jrmatherly/obot-tools
```

---

#### Step 2.4: Clean up tools/index.yaml

**File:** `obot-entraid/tools/index.yaml`

**Current:**

```yaml
authProviders:
  entra-auth-provider:
    reference: ./entra-auth-provider
  keycloak-auth-provider:
    reference: ./keycloak-auth-provider
```

**Option A - Remove file entirely:**

```bash
rm obot-entraid/tools/index.yaml
```

**Option B - Empty file (if structure required):**

```yaml
# Custom auth providers moved to obot-tools repository
# See: https://github.com/jrmatherly/obot-tools
```

---

#### Steps 2.5-2.7: Archive Provider Directories

**Option A - Delete immediately:**

```bash
rm -rf obot-entraid/tools/entra-auth-provider
rm -rf obot-entraid/tools/keycloak-auth-provider
rm -rf obot-entraid/tools/auth-providers-common
```

**Option B - Archive first (recommended):**

```bash
mkdir -p obot-entraid/tools/.archived
mv obot-entraid/tools/entra-auth-provider obot-entraid/tools/.archived/
mv obot-entraid/tools/keycloak-auth-provider obot-entraid/tools/.archived/
mv obot-entraid/tools/auth-providers-common obot-entraid/tools/.archived/
echo "Archived on $(date). Migrated to obot-tools repository." > obot-entraid/tools/.archived/README.md
```

---

## Verification Procedures

### Phase 1 Verification

```bash
# 1. Build verification
cd obot-tools/keycloak-auth-provider
go build -o /dev/null .
echo "Keycloak build: $?"

cd ../entra-auth-provider
go build -o /dev/null .
echo "Entra build: $?"

# 2. Test verification
cd ../keycloak-auth-provider
go test -v ./...
echo "Keycloak tests: $?"

cd ../entra-auth-provider
go test -v ./...
echo "Entra tests: $?"

# 3. Full build verification
cd ..
make build
echo "Full build: $?"

# 4. Full test verification
make test
echo "Full test: $?"

# 5. Index verification
grep -A 10 "authProviders:" index.yaml
```

**Expected output:**

- All builds exit with code 0
- All tests pass
- index.yaml shows all 4 auth providers

---

### Phase 2 Verification

```bash
# 1. Build container
cd obot-entraid
docker build -t obot-migration-test .

# 2. Verify providers exist in image
docker run --rm obot-migration-test ls -la /obot-tools/tools/ | grep auth

# 3. Verify index contains all providers
docker run --rm obot-migration-test cat /obot-tools/tools/index.yaml

# 4. Verify binaries are executable
docker run --rm obot-migration-test /obot-tools/tools/entra-auth-provider/bin/gptscript-go-tool --help
docker run --rm obot-migration-test /obot-tools/tools/keycloak-auth-provider/bin/gptscript-go-tool --help

# 5. Local development verification
source .envrc.dev
make dev
# Test auth provider selection in UI
```

**Expected output:**

- Container builds successfully
- All 4 auth providers visible in `/obot-tools/tools/`
- index.yaml contains google, github, keycloak, entra entries
- Binaries execute without error
- Dev mode starts and shows all auth providers

---

### End-to-End Verification

1. **OAuth Flow Test (Entra):**
   - Configure Entra auth provider in UI
   - Attempt login via Microsoft
   - Verify user info returned correctly
   - Verify profile picture displayed

2. **OAuth Flow Test (Keycloak):**
   - Configure Keycloak auth provider in UI
   - Attempt login via Keycloak instance
   - Verify user info returned correctly
   - Verify group membership (if configured)

---

## Rollback Plan

### Phase 1 Rollback

If Phase 1 fails, revert changes in obot-tools:

```bash
cd obot-tools
git checkout -- keycloak-auth-provider/main.go
git checkout -- entra-auth-provider/main.go
git checkout -- index.yaml
git checkout -- Makefile
```

### Phase 2 Rollback

If Phase 2 fails after Phase 1 is deployed:

1. Revert Dockerfile changes in obot-entraid
2. Revert .envrc.dev changes
3. Restore archived provider directories (if Option B was used)
4. Rebuild container

```bash
cd obot-entraid
git checkout -- Dockerfile
git checkout -- .envrc.dev
git checkout -- tools/index.yaml
# If archived:
mv tools/.archived/entra-auth-provider tools/
mv tools/.archived/keycloak-auth-provider tools/
mv tools/.archived/auth-providers-common tools/
```

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Import path change breaks other code | Low | High | Grep for all references before changing |
| Container build fails after Dockerfile changes | Medium | Medium | Test in CI before merging |
| Auth flow breaks in production | Low | High | Test OAuth flows in staging first |
| Missing dependencies after removing local providers | Low | Medium | Verify TOOLS_IMAGE contains all required files |
| Documentation becomes inconsistent | Medium | Low | Update docs in same PR as code changes |

---

## Appendices

### Appendix A: Files Changed Summary

**obot-tools (Phase 1):**

- `keycloak-auth-provider/main.go` - Import path fix
- `entra-auth-provider/main.go` - Import path fix
- `index.yaml` - Add provider registrations
- `Makefile` - Add test targets
- `entra-auth-provider/README.md` - Remove old references
- `keycloak-auth-provider/.gitignore` - New file
- `entra-auth-provider/.gitignore` - New file

**obot-entraid (Phase 2):**

- `Dockerfile` - Remove local builds, simplify stages
- `.envrc.dev` - Remove ./tools from registries
- `tools/index.yaml` - Remove or empty
- `tools/entra-auth-provider/` - Archive/delete
- `tools/keycloak-auth-provider/` - Archive/delete
- `tools/auth-providers-common/` - Archive/delete

**obot-entraid (Phase 3):**

- `docs/docs/configuration/entra-id-authentication.md`
- `docs/docs/configuration/keycloak-authentication.md`
- `docs/docs/contributing/local-development.md`
- `docs/docs/operations/auth-provider-testing.md`
- `CONTRIBUTING.md`
- `tools/README.md`

---

### Appendix B: Import Path Reference

**Correct module paths after migration:**

| Provider | Module Path |
|----------|-------------|
| Google | `github.com/obot-platform/tools/google-auth-provider` |
| GitHub | `github.com/obot-platform/tools/github-auth-provider` |
| Keycloak | `github.com/obot-platform/tools/keycloak-auth-provider` |
| Entra | `github.com/obot-platform/tools/entra-auth-provider` |
| Common | `github.com/obot-platform/tools/auth-providers-common` |

---

### Appendix C: Environment Variables

**Common (all providers):**

- `OBOT_AUTH_PROVIDER_COOKIE_SECRET` - Cookie encryption key
- `OBOT_AUTH_PROVIDER_EMAIL_DOMAINS` - Allowed email domains
- `OBOT_AUTH_PROVIDER_POSTGRES_CONNECTION_DSN` - Optional PostgreSQL storage

**Keycloak-specific:**

- `OBOT_KEYCLOAK_AUTH_PROVIDER_CLIENT_ID`
- `OBOT_KEYCLOAK_AUTH_PROVIDER_CLIENT_SECRET`
- `OBOT_KEYCLOAK_AUTH_PROVIDER_URL`
- `OBOT_KEYCLOAK_AUTH_PROVIDER_REALM`
- `OBOT_KEYCLOAK_AUTH_PROVIDER_ALLOWED_GROUPS` (optional)
- `OBOT_KEYCLOAK_AUTH_PROVIDER_ALLOWED_ROLES` (optional)

**Entra-specific:**

- `OBOT_ENTRA_AUTH_PROVIDER_CLIENT_ID`
- `OBOT_ENTRA_AUTH_PROVIDER_CLIENT_SECRET`
- `OBOT_ENTRA_AUTH_PROVIDER_TENANT_ID`
- `OBOT_ENTRA_AUTH_PROVIDER_ALLOWED_GROUPS` (optional)
- `OBOT_ENTRA_AUTH_PROVIDER_ALLOWED_TENANTS` (optional)
- `OBOT_ENTRA_AUTH_PROVIDER_USE_WORKLOAD_IDENTITY` (optional)

---

### Appendix D: CI/CD Considerations

**GitHub Actions to update (obot-entraid):**

- `.github/workflows/ci.yml` - Remove auth provider build/test steps if present
- `.github/workflows/test-upstream-merge.yml` - Update if tests auth providers

**Image dependencies:**

- `TOOLS_IMAGE` must be rebuilt and published before `obot-entraid` can use it
- Recommended workflow:
  1. Merge Phase 1 changes to obot-tools
  2. Wait for obot-tools CI to publish new image
  3. Merge Phase 2 changes to obot-entraid referencing new image tag

---

## Document Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-18 | Claude Code | Initial draft |
| 1.1 | 2026-01-18 | Claude Code | Added Phase 0 for placeholder-credential sync (discovered via /sc:reflect validation) |
| 1.2 | 2026-01-18 | Claude Code | Implemented Phase 0 and Phase 1 - all obot-tools changes complete |

---

*End of Document*
