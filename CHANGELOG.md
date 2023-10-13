# Changelog

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
