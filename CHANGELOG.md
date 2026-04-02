# Changelog

## [Unreleased]

### Changed

- **BREAKING**: `Cookie()` return type changed from `*fasthttp.Cookie` to `string`
- **BREAKING**: `Cookies()` return type changed from `[]*fasthttp.Cookie` to `map[string]string`

### Fixed

- Fixed dangling pointer bug in internal `cookie()` method

## [0.9.0] - Apr 1, 2026

### Added

- AGENTS.md with comprehensive project guidelines for AI agents
- HTTP status code constants (StatusOK, StatusNotFound, StatusInternalServerError, etc.)
- HTTP method constants (MethodGet, MethodPost, MethodPut, etc.)
- MIME type constants (MIMEApplicationJSON, MIMETextHTML, etc.)
- Header key constants (HeaderContentType, HeaderAccept, HeaderAuthorization, etc.)
- Additional Context methods: IsAjax(), IsWebSocket(), AcceptedLanguages(), RemoteAddr()
- Middleware caching in route groups

### Changed

- **Major**: Refactor from net/http to fasthttp for better performance
- Upgrade Go minimum version requirement to 1.25
- Upgrade fasthttp from v1.52.0 to v1.69.0
- Upgrade validator/v10 from v10.12.0 to v10.30.2
- Upgrade golang.org/x/crypto from v0.19.0 to v0.49.0
- Upgrade golang.org/x/sys from v0.17.0 to v0.42.0
- Upgrade golang.org/x/text from v0.14.0 to v0.35.0
- Router optimization: matchChild from O(n) to O(1)
- Replace interface{} with any throughout codebase
- Simplified cookiesMap to store string values instead of http.Cookie pointers
- Update examples for fasthttp compatibility
- Update HTTP constants to match fasthttp

### Fixed

- Improved test coverage for lightning.go (Run, Shutdown, Static, Context pooling)
- Fixed route matching issues
- Fixed context pool reuse bug where data wasn't being reset properly
- Fixed X-Forwarded-For header parsing for comma-separated IP addresses

### Performance

- Test coverage improved from 90.1% to 96.0%
- Context pooling via sync.Pool for reduced GC pressure

## [0.8.0] - Mar 29, 2026

### Added

- Graceful shutdown support with `Run()`, `RunListener()`, `Shutdown()`, and `RunGraceful()` methods
- `sync.Pool` for Context reuse to reduce GC pressure and improve performance
- `MaxRequestBodySize` configuration option for limiting request body size
- Middleware caching in route groups for better performance
- Signal handling (SIGINT/SIGTERM) for graceful shutdown
- Comprehensive integration tests with dynamic port allocation

### Changed

- **BREAKING**: Debug mode default changed from `true` to `false`
- Logger middleware now uses `ctx.App.Logger.Info()` instead of `fmt.Printf`
- X-Forwarded-For header parsing now correctly handles comma-separated IP addresses
- Replaced `interface{}` with `any` throughout codebase
- Improved global logger implementation

### Fixed

- Critical shutdown issue where server wouldn't stop properly
- Config merge nil pointer panic
- Request body size limit enforcement
- Context pool reuse bug where data wasn't being reset properly

### Removed

- Deprecated and unused code
- Hardcoded test ports in favor of dynamic allocation

### Performance

- Test coverage improved from 81.2% to 90.9%

## [0.7.3] - Jun 28, 2024

### Added

- add QueryBool method to Context for parsing boolean query parameters
- add comprehensive unit tests for Context methods

## [0.7.2] - Mar 11, 2024

### Fixed

- fix route matching issue

## [0.7.1] - Oct 13, 2023

### Fixed

- Fixed empty string check in QueryInt, QueryUInt, QueryInt8, QueryUInt8, QueryInt32, QueryUInt32, QueryInt64, QueryUInt64, QueryFloat32, and QueryFloat64.

## [0.7.0] - Oct 13, 2023

### Added

- `ctx.QueryString(key string) string`: Returns the value of a given query parameter as a string.
- `ctx.QueryInt(key string) (int, error)`: Returns the value of a given query parameter as an int.
- `ctx.QueryUInt(key string) (uint, error)`: Returns the value of a given query parameter as a uint.
- `ctx.QueryInt8(key string) (int8, error)`: Returns the value of a given query parameter as an int8.
- `ctx.QueryUInt8(key string) (uint8, error)`: Returns the value of a given query parameter as a uint8.
- `ctx.QueryInt32(key string) (int32, error)`: Returns the value of a given query parameter as an int32.
- `ctx.QueryUInt32(key string) (uint32, error)`: Returns the value of a given query parameter as a uint32.
- `ctx.QueryInt64(key string) (int64, error)`: Returns the value of a given query parameter as an int64.
- `ctx.QueryUInt64(key string) (uint64, error)`: Returns the value of a given query parameter as a uint64.
- `ctx.QueryFloat32(key string) (float32, error)`: Returns the value of a given query parameter as a float32.
- `ctx.QueryFloat64(key string) (float64, error)`: Returns the value of a given query parameter as a float64.

## [0.6.0] - Sep 22, 2023

### Added

- ParamInt(key string) (int, error): Returns the value of a URL parameter as an integer.
- ParamUInt(key string) (uint, error): Returns the value of a URL parameter as a uint.
- ParamInt64(key string) (int64, error): Returns the value of a URL parameter as an int64.
- ParamUInt64(key string) (uint64, error): Returns the value of a URL parameter as a uint64.
- ParamFloat32(key string) (float32, error): Returns the value of a URL parameter as a float32.
- ParamFloat64(key string) (float64, error): Returns the value of a URL parameter as a float64.
- ParamString(key string) string: Returns the value of a URL parameter as a string.

## [0.5.0] - May 14, 2023

### Added

- `ctx.HTML(code int, name string, data interface{})` writes an HTML response with the given status code, template name, and data.
- `ctx.Static(root string, prefix string)` serves static files from the given root directory with the given prefix.
- `app.LoadHTMLGlob(pattern string)` loads HTML templates from a glob pattern and sets them in the Application struct.
- `app.SetFuncMap(funcMap template.FuncMap)` sets the funcMap in the Application struct to the funcMap passed in as an argument.

## [0.4.1] - May 13, 2023

### Changed

- refactor router to prefix tree

## [0.4.0] - Apr 20, 2023

### Added

- JSON serialization supports customization

### Changed

- `JSONBody(v interface{}, valid ...bool)` method supports parameter validation
- `recovery` middleware supports custom handler
- `NewApp` supports custom configuration

### Removed

- remove `ctx.Bind` and `ctx.BindAndValidate` method

### Other

- update examples

## [0.3.0] - Apr 8, 2023

### Added

- Add SkipFlush method and original Req/Res to context

## [0.2.0] - Apr 7, 2023

### Added

- context support Bind and BindAndValidate method

## [0.1.5] - Apr 5, 2023

### Added

- context expose Body() / SetBody([]byte body) method

## [0.1.4] - Apr 5, 2023

### Fixed

- fix route register

## [0.1.3] - Apr 3, 2023

### Other

- chore: republish

## [0.1.2] - Apr 3, 2023

### Changed

- expose `NewContext` function

## [0.1.1] - Apr 3, 2023

### Added

- add type alias `Middleware`

## [0.1.0] - Apr 2, 2023

### Added

- Implement `lightning` v1 version
