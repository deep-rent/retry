# deep-rent/retry

![Logo](https://raw.githubusercontent.com/deep-rent/retry/master/logo.png)

This library provides a retry mechanism for Go based on highly configurable backoff strategies. Unlike similar libraries, the retry mechanism is initiated through a reusable struct type, which allows for easy mocking and sharing of backoff configuration.

[![Test Status](https://github.com/deep-rent/retry/actions/workflows/test.yml/badge.svg)](https://github.com/deep-rent/retry/actions/workflows/test.yml) [![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/deep-rent/retry) [![Code Quality](https://goreportcard.com/badge/github.com/nanomsg/mangos)](https://goreportcard.com/report/github.com/deep-rent/retry)


## Installation

Download the libary using `go get`:

```
go get github.com/deep-rent/retry@latest
```

Add the following imports to your project:

```go
import (
    "github.com/deep-rent/retry"
    "github.com/deep-rent/retry/backoff"
)
```

## Usage

First, define a function to be retried. The signature must match `retry.Attempt`.

```go
attempt := func(n int) error {
    if n == 5 {
        // succeed after 5 attempts
        return nil
    } else {
        // force retry
        return errors.New("whoops")
    }
}
```

Next, configure a `retry.Cycler`. Once configured, the cycler can be reused across your project.

```go
// exponentially increase delays between consecutive attempts
cycler := retry.NewCycler(backoff.Exponential(5 * time.Second, 1.5))
cycler.Cap(3 * time.Minute)      // cap backoff to 3 minutes
cycler.Limit(25)                 // stop after 25 attempts
cycler.Timeout(10 * time.Minute) // time out after 10 minutes
cycler.Jitter(0.5)               // introduce 50% random jitter
```

Register a `retry.ErrorHandler` to catch intermediate errors.

```go
cycler.OnError(func(n int, delay time.Duration, err error) {
    s := delay.Seconds()
    fmt.Printf("attempt #%d failed: %v => wait %4.f s\n", n, err, s)
})
```

Finally, retry `attempt` according to the previous configuration.

```go
err := cycler.Try(attempt)
if err != nil {
    fmt.Printf("failed after retries: %v", err)
}
```

Another way to exit a retry cycle early is to flag the returned error using `ForceExit`.

```go
attempt := func(int) error {
    // exit the retry cycle immediately
    return retry.ForceExit(errors.New("unrecoverable")) 
}
```

## License

Licensed under the Apache 2.0 License. For the full copyright and licensing information, please view the `LICENSE` file that was distributed with this source code.