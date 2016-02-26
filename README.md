# Goroutine Local Storage

Similar to Thread Local Storage, but for Go routines.

## Problem solved

### Example: tracing HTTP Request

In the request handler, it will perform complicated work invoking
many functions.
Sometimes, it's difficult and even impossible (through third-party layer)
to pass `*http.Request` all the way down to every function called.
But we want most of the functions be able to access the HTTP request,
or at least be aware of some information in the HTTP request (e.g. RequestId).
With `gls`, we don't need to pass request to every function which is still
able to get the request using `gls.Get()`.

## Usage

```go
type Context struct {
    ...
}

func myWork() {
    context := gls.Get().(*Context)
    ...
}

func handler(req *Request) {
    context := contextFromRequest(req)
    gls.With(context, myWork)
}
```

`gls.Get()` requires `gls.With` called, otherwise it will `panic`.
So it's guaranteed the returned value is present.
To prevent `panic`, use `gls.GetSafe()` which returns `nil` if `gls.With` is never called.

The context doesn't go across go routines.
If the request must be handled in separate go routines, the context should be forwarded.
A helper `gls.Go` is provided to spawn a go routine with current context forwarded:

```go

func workerFn() {
    // this is in a separate go routine
    context := gls.Get().(*Context)
}

func myHandler() {
    gls.Go(workerFn)
}

```

This above code is equivalent to:

```go
func workerFn() {
    context := gls.Get().(*Context)
}

func myHandler() {
    context := gls.Get().(*Context)
    go gls.With(context, workerFn)
}
```

## How it works

Because Go routine runs on arbitrary OS thread, the functions in the same
Go routine may run on different OS threads from time to time, the Thread Local Storage doesn't solve the problem.
The implementation uses stack trace of current Go routine to locate a special
function call with magic function name (containing a UUID).
By parsing the arguments (pointer values) passed to the function, we can find
the context associated with current Go routine.

## Limitations

- Not cross Go routines

  The context is for one Go routine, if some functions are invoked using `go fn()`, the context should be explicitly passed:

  ```
  context := gls.Get()
  go gls.With(context, workFn)
  ```

  or with helper

  ```
  gls.Go(workFn)
  ```

- Limited buffer for stack trace

  `gls` uses `runtime.Stack` to dump formatted stack.
  The default buffer size is defined in `gls.StackBufferSize`.
  If the stack will become very deep, please tune `gls.StackBufferSize`
  to a large value in the beginning of your program.

- Supported platforms and architectures

  The implementation is a little different on 32-bit/64-bit architectures,
  because the size of uint passed in the function is different.

## License
MIT
