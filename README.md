# pool

Golang implementation of worker pool pattern.

## About The Project

The library provides an API to limit the number of goroutines by scheduling tasks in queue.

Features:

- laconic api
- safe for concurrent use

## Usage

```
executor := pool.New(10)
executor.Exec(func() {
	fmt.Println("done")
})
```

## License

Pool is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENCE.md)
for the full license text.

## Contact

- Email: `cherkashin.evgeny.viktorovich@gmail.com`
- Telegram: `@evgeny_cherkashin`