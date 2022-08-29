# Result

This Result type for Go is inspired by the [Rust `Result` type](https://doc.rust-lang.org/std/result/enum.Result.html). For the most part, this package implements all of the same functionality as the Rust type. A few differences are that a few methods are implemented as functions rather than methods due to Go not allowing adding a second generic type in a method signature of a generic type, and some methods that are related to Rust specific features such as memory ownership being absent. One design decision made that poses a significant change to the Rust module is that this package's `Result` type does not have a generic type for it's error case, but rather works with native Go `error`s. While a generic would have a similar implementation, it seems appropriate that the idiomatic Go `error`s should persist in a type relating to failure.

## Purpose

The `Result` gives a large extension to the possibilities for error handling. At it's core, a `Result` type represents two states: one of success (`Ok`), and one of failure (`Err`). While in Go this is usually represented by returning two values with one being an `error` that is null checked, the `Result` represents this idea in a single object which can then have behavior injected into it until we are ready to deal with the error if it exists. For example, consider the following functions:

```go
func Foo(age int) (string, error) {
  if age < 21 {
    return "", errors.New("too young to drink")
  }
  return "drink up!", nil
}

func Bar(s string) (string, error) {
  if len(s) > 63 {
    return "", errors.New("message too long")
  }
  return fmt.Sprintf("I agree! %s", s), nil
}
```

To use these functions together, we might do something like the following:

```go
func main() {
  s, err := Foo(12)
  if err != nil {
    panic(err)
  }
  concat, err := Bar(s)
  if err != nil {
    panic(err)
  }

  fmt.Println(concat)
} 
```

In this example we are using errors when we have been passed data that is not acceptable. However, this means that each step we have to do an error check and terminate the program before continuing. Such error handling appears everywhere in Go code and is idiomatic, but could be more expressive. The following is how the code could change using the `Result` type:

```go
func Foo(age int) result.Result[string] {
  if age < 21 {
    return result.Err[string](errors.New("too young to drink"))
  }
  return result.Ok("drink up!")
}

func Bar(s string) (string, error) {
  if len(s) > 63 {
    return result.Err[string](errors.New("message too long"))
  }
  return result.Ok(fmt.Sprintf("I agree! %s", s))
}

func main() {
  result.AndThen(Foo(), Bar).Inspect(func(s string) {
    fmt.Println(s)
  }).InspectErr(func(err error) {
    panic(err)
  })
}
```

Here we can see the more declarative `main` function. As the errors do not have any particular impact other than making the program fail, we can get all of the logic written before dealing with errors in one place. Another added benefit is having a more clear return type when we get to an error case. Instead of needing to instantiate a zero value for a type each time we have an error, or a nil error each time we do not, we can simply return the `Err` when we have and error, or an `Ok` when we don't. This effort to reduce `nil`s in the code gives us just one more bit of type safety. For another, checkout [`Option`s in Go](https://github.com/JustinKnueppel/go-option).

## Usage

`Result` types should be instantiated via the `Ok` and `Err` constructors only (default is `Ok` with the generic's zero value). Most features of this package are implemented as methods on the `Result` type, but a few that require a second generic type are implemented as functions instead. A few examples of how to use the package follow, but more examples on the same functionality can be found in the [Rust std::result docs](https://doc.rust-lang.org/std/result/enum.Result.html), albeit written in Rust.

### Updating values or errors

This example shows how functions can safely attempt to modify both the `Ok` case and the `Err` case for an option.

```go
func AugmentResult(r result.Result[string]) result.Result[string] {
  // This will do nothing if r was Err
  r = result.Map(r, func(s string) string {
    return fmt.Sprintf("Context: %s", s)
  })
  // This will do nothing if r was Ok
  return result.MapErr(r, func(err error) error {
    return fmt.Errorf("failed to add context: %w", err)
  })
}
```

### Handling all errors at once

This example shows how we can continually pass around a `Result` until we are ready to handle any errors that may come back. This keeps all of our success logic in one place, and our error handling together.

```go
func getResult() result.Result[int] {}
func calculate(x int) int {}
func remoteCalculation(x int) result.Result[int] {}
func updateDatabase(x int) result.Result[string] {}
func logNewValue(val string) {}

func main() {
  result.AndThen(
    result.AndThen(
      result.Map(getResult(), calculate), remoteCalculation), updateDatabase) // here we can see all success logic
    .Inspect(logNewValue)
    .InspectErr(func (err error) {
      panic(err) // if an error occured at any point, we will see it here
    })
}
```

## Functions vs Methods

Most features of this package are implemented as methods on an `Option`. However, due to the lack of generic methods on generic types in Go, some of the methods from the Rust library had to be implemented as package functions. The affected functions are:

- `Map`
- `MapOr`
- `MapOrElse`
- `And`
- `AndThen`
- `Contains`
- `Flatten`
- `Equal`

## Missing methods from Rust's `std::result`

There are quite a few methods from the Rust `std::result` type that are not implemented in this package. These methods should be methods relating to Rust specific language features such as getting a mutable reference, pinned value, or result type conversion. If there are any missng methods that make sense for a Go `Result` type, feel free to leave a Github issue detailing them.
