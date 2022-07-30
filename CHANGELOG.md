# Changelog

## v1.0.1 -> v2.0.0

### Breaking changes:

1. `Attempt` was renamed to `AttemptFunc`.
2. `ErrorHandler` was renamed to `ErrorHandlerFunc`.

### New Features:

1. `AttemptFunc` can return an `ExitError` to stop the current retry cycle.
2. `ForceExit` wraps any error in an `ExitError`.

## Refactoring:

1. Rename `master` branch to `main`.