# Changelog

## v1.1.0 -> v1.1.1

### Potentially breaking changes:

1. Upgrade to Go 1.19.

### Refactoring:

1. Make use of links in doc comments.

## v1.0.1 -> v1.1.0

### Potentially breaking changes:

1. `Attempt` was renamed to `AttemptFunc`.
2. `ErrorHandler` was renamed to `ErrorHandlerFunc`.

### New Features:

1. `AttemptFunc` can return an `ExitError` to stop the current retry cycle.
2. `ForceExit` wraps any error in an `ExitError`.

### Refactoring:

1. Rename `master` branch to `main`.
