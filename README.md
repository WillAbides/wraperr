# wraperrs 

Package wraperrs provides some of the useful helpers from github.com/pkg/errors that are missing from xerrors. Most of
the code in this repository is copied from github.com/pkg/errors.

`go get github.com/WillAbides/wraperr`

## Adding context to an error

The wraperrs.Wrap function returns a new error that adds context to the original error. For example
```go
_, err := ioutil.ReadAll(r)
if err != nil {
        return wraperrs.Wrap(err, "read failed")
}
```
## Retrieving the cause of an error

Using `wraperrs.Wrap` constructs a stack of errors, adding context to the preceding error. Depending on the nature of
the error it may be necessary to reverse the operation of wraperrs.Wrap to retrieve the original error for inspection.
Any error value which implements this interface can be inspected by `xerrors.Wrapper`.

`wraperrs.Cause` will recursively retrieve the topmost error which does not implement `xerrors.Wrapper`, which is assumed to be the original cause. For example:
```go
switch err := wraperrs.Cause(err).(type) {
case *MyError:
        // handle specifically
default:
        // unknown error
}
```

[Read the package documentation for more information](https://godoc.org/github.com/WillAbides/wraperrs).

## License

BSD-2-Clause