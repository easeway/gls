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
able to get the request using `gls.Ctx()`.

## Usage

```go
type Context struct {
    ...
}

func myWork() {
    ctx := gls.Ctx().(*Context)
    ...
}

func handler(req *Request) {
    context := contextFromRequest(req)
    gls.WithCtx(context, myWork)
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
  context := gls.Ctx()
  go func() {
      gls.WithCtx(context, workFn)
  }()
  ```

- Limited buffer for stack trace

  `gls` uses `runtime.Stack` to dump formatted stack.
  The default buffer size is defined in `gls.StackBufferSize`.
  If the stack will become very deep, please tune `gls.StackBufferSize`
  to a large value in the beginning of your program.

## License
MIT
