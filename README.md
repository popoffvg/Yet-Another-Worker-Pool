# Yet Another Worker Pool

`YAWP` is implementation worker pool with generic.

Design for not think about:
- result type
- error handling
- goroutine leaks
- panic handling

## Questions

### Pub, Run return Result type. Is it promise?

Yep, it's promise. Method `Get` lock until result is not ready.

If you need get result only than it's ready use `RedirectOutput` option.

## Hot to install

```bash
go get github.com/popoffvg/Yet-Another-Worker-Pool
```

## Usage

### Use pool with function stream

```go
wp := New[int, int](3)

for i := 0; i < 1_000; i++ {
    number := i + 1
    wp.Run(func(ctx context.Context) (int, error) {
        return number, nil
    })
}
```

### Use pool with value stream

```go
wp := New(3,
    WithWorker(func(ctx context.Context, i int) (int, error) {
        return i + 1, nil
    }),
)

for i := 0; i < 1_000; i++ {
    wp.Pub(number)
}
```

### Use Pool Fun-In Fun-Out

```go
wp := New(3, RedirectOutput[int, int]())
go func(){
    for {
        wp.Run(func(ctx context.Context) (int, error) {
            return 1, nil
        })
    }
}()

for v := range wp.Stream() {
   // get result as value stream
}
```

## License

[See](LICENSE) 