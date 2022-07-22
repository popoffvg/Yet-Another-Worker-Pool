# Yet Another Worker Pool

`YAWP` is implementation worker pool with generic.

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

## License

[See](LICENSE) 