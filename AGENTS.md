# AGENTS.md — structconfig

This document provides guidance for AI agents working on this codebase.

## Overview

`structconfig` is a Go library that maps **default values**, **environment variables**, and **command-line flags** to struct fields using struct tags. It uses reflection and a Visitor pattern to provide a clean, decoupled pipeline for populating and merging configuration structs.

---

## Repository Layout

```
.
├── structconfig.go                        # Public API: StructConfig[T], Merger[T], New, NewMerger
├── builder.go                             # Builder[T]: chains multiple config generators with merge
├── util.go                                # NewConfigWithMerge: convenience helper (default + env + flags)
├── structconfig_config_generated.go       # Generated: Option/ConfigBuilder (via goconfig)
├── internal/
│   ├── field.go                           # Receptor interface, StructField, typed sub-interfaces
│   ├── type.go                            # Type: struct metadata, Fields(), Accept(Receptor)
│   ├── tag.go                             # Tag: parses name/default/usage/short struct tags
│   ├── switch.go                          # Switch/Call: dispatches StructField to correct Receptor method
│   ├── pair.go                            # ParsePair, PairsReceptor, PairsSynthReceptor
│   ├── typed.go                           # TypedReceptor interface, DefaultTypedReceptor
│   ├── conv.go                            # Converter interface, DefaultConverter (strconv wrappers)
│   ├── set.go                             # SetReceptor, SetTypedReceptor: reflect-based field writer
│   ├── default.go                         # DefaultReceptor: populates fields from `default` tags
│   ├── env.go                             # EnvVar type: normalizes env var names
│   ├── envset.go                          # EnvReceptor: populates fields from environment variables
│   ├── flagset.go                         # FlagSetReceptor, FlagSetTypedReceptor: defines pflag flags
│   ├── pflag.go                           # PFlagSetReceptor, PFlagGetReceptor, PFlagGetConverter
│   ├── merge.go                           # Merger[T]: field-level merge based on default values
│   ├── error.go                           # Sentinel errors and helpers (JoinErrors, Errorf)
│   ├── util.go                            # TryParse, ParseInt/Uint/Float, IsSupportedKind
│   └── structfield_dataclass_generated.go # Generated: StructField data class (via dataclass)
```

---

## Architecture

The library is built around a **Visitor pattern** and a **pipeline of three roles**:

### 1. Reflection & Metadata Layer

- [`internal/type.go`](internal/type.go) — `Type` wraps a `reflect.Type` for a struct. `Type.Fields()` returns a `[]StructField`, each holding the field name, `reflect.Kind`, and parsed `Tag`.
- [`internal/tag.go`](internal/tag.go) — `Tag` reads the four supported struct tags: `name`, `default`, `usage`, `short`. An optional `prefix` can be prepended to all tag keys (set via `WithPrefix`).
- [`internal/field.go`](internal/field.go) — Defines the `Receptor` interface (the Visitor) and its sub-interfaces (`BoolReceptor`, `IntReceptor`, etc.).

### 2. Dispatch Layer (Visitor routing)

- [`internal/switch.go`](internal/switch.go) — `Switch(r Receptor, kind)` returns the correct method of `r` based on the field's `reflect.Kind`. Unknown kinds route to `r.Any`.
- `Type.Accept(r Receptor)` iterates all fields and calls `Call(r, f)` for each, which invokes `Switch` internally.

### 3. Value Pipeline (Get → Convert → Set)

Every concrete Receptor is assembled as a [`PairsSynthReceptor`](internal/pair.go) that composes three roles:

| Role | Interface | Default Implementation |
|---|---|---|
| **Get** | `func(StructField) (string, error)` | Tag lookup or env var or flag value |
| **Convert** | [`Converter`](internal/conv.go) | `DefaultConverter` (strconv-based) |
| **Set** | [`TypedReceptor`](internal/typed.go) | `DefaultTypedReceptor` (reflect field writer) |

The three concrete receptors built on this pipeline:

- [`internal/default.go`](internal/default.go) — `DefaultReceptor`: Get from `default` tag → Convert → Set.
- [`internal/envset.go`](internal/envset.go) — `EnvReceptor`: Get from `os.LookupEnv` (falls back to `default` tag) → Convert → Set.
- [`internal/pflag.go`](internal/pflag.go) — `PFlagGetReceptor`: Get from parsed `pflag.FlagSet` → Convert → Set.

### 4. Merge & Priority Layer

- [`internal/merge.go`](internal/merge.go) — `Merger[T].Merge(left, right T)`: For each field tagged with `name`, if `right`'s value differs from the default, `right` wins; otherwise `left` wins; otherwise the default is used.
- [`builder.go`](builder.go) — `Builder[T]` chains generator functions (each producing a `*T`) and sequentially merges them, left-to-right, into a base default config.

### 5. Public API

- [`structconfig.go`](structconfig.go) — Thin generic wrappers over `internal`: `StructConfig[T]` (FromDefault, FromEnv, FromFlags, SetFlags) and `Merger[T]`.
- [`util.go`](util.go) — `NewConfigWithMerge` is the highest-level helper: registers env + flags generators in a `Builder` and returns the final merged config.

### Data Flow

```
User's struct T (with tags)
        │
        ▼
   Type.Accept(Receptor)          ← dispatches each field by reflect.Kind
        │
        ▼ (for each field)
   PairsSynthReceptor
    ├── Get(StructField)  →  raw string value (tag / env / flag)
    ├── Convert(string)   →  typed value
    └── Set(field, value) →  reflect.Value.Set*(...)
                                   + optional slog.Logger callback

  Multiple populated *T values
        │
        ▼
   Merger[T].Merge(left, right)   ← field-by-field, default-aware
        │
        ▼
   Final merged *T
```

---

## Struct Tag Reference

| Tag | Purpose |
|---|---|
| `name:"<key>"` | Field name for env var lookup and flag registration. `-` to ignore the field. |
| `default:"<val>"` | Default value (string form). Used by `DefaultReceptor` and `Merger`. |
| `usage:"<text>"` | Flag usage string for pflag. |
| `short:"<char>"` | Single-character shorthand for pflag. |
| `count:"true"` | Flag that counts the number of times it is provided (e.g. `-v`, `-vv`). |

A `prefix` option (`WithPrefix`) prepends a string to all tag key names, enabling namespaced tags.

---

## Supported Field Types

Natively supported `reflect.Kind` and special types (handled by the typed pipeline):

- **Primitives**: `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `string`
- **Slices**: `[]bool`, `[]int`, `[]int32`, `[]int64`, `[]uint`, `[]float32`, `[]float64`, `[]string`
- **Time & Duration**: `time.Duration`, `time.Time` (RFC3339 format)
- **Count flags**: `int` / `uint` with `count:"true"` tag

All other kinds route to the `Any` method. Provide a custom `AnyCallback` (`WithAnyCallback`) to parse and set these fields. Supply `AnyEqual` for correct merge behavior on non-primitive fields.

---

## Code Generation

Two files are auto-generated and **must not be edited by hand**:

| File | Generator | Command |
|---|---|---|
| [`structconfig_config_generated.go`](structconfig_config_generated.go) | `github.com/berquerant/goconfig` | `go generate ./...` |
| [`internal/structfield_dataclass_generated.go`](internal/structfield_dataclass_generated.go) | `github.com/berquerant/dataclass` | `go generate ./...` |

The `//go:generate` directives are in [`structconfig.go`](structconfig.go) and [`internal/field.go`](internal/field.go).

---

## Development Workflow

This project uses `make` ([`Makefile`](Makefile)).

```shell
# Run lint + tests (default)
make

# Run tests only (with race detector and coverage)
make test

# Run go vet
make vet

# Regenerate generated files
make generate

# Tidy go.mod
make tidy
```

Underlying test command: `go test -cover -race ./...`

---

## Key Constraints & Guidelines

### Do not modify generated files
Files ending in `_generated.go` are managed by `go generate`. Modify the source structs and re-run generation instead.

### Do not modify test files unless necessary
Test files (`*_test.go`) should remain unchanged except when a code change makes it strictly necessary to update them.

### internal package is private
All types in `internal/` are implementation details. The public surface is only what is exported from the root package. New functionality should be exposed through the root package.

### Error sentinel pattern
All errors are wrapped under `ErrStructConfig` (see [`internal/error.go`](internal/error.go)) using `JoinErrors` or `Errorf`. Use `errors.Is` for checking.

### ErrSkipParse vs ErrParseAsDefault
When implementing a custom `Get` function:
- Return `ErrSkipParse` to skip the field entirely (no set, no error).
- Return `ErrParseAsDefault` to use the zero value of the target type as the parsed value (triggers `callback` with zero value).

### Logger is always optional
`*slog.Logger` is passed around and always guarded with `if logger != nil`. Never assume it is non-nil.

---

## Dependencies

| Dependency | Role |
|---|---|
| `github.com/spf13/pflag` | CLI flag parsing |
| `golang.org/x/exp/constraints` | Generic numeric type constraints |
| `github.com/stretchr/testify` | Test assertions |
| `github.com/berquerant/goconfig` | Code generation for option builder |
| `github.com/berquerant/dataclass` | Code generation for value-object structs |
