# helloworld-pbt

A tiny Go "Hello, World!" that demonstrates **property-based testing (PBT)**.

## Idea (from LESSON-4/RADME.md)

Instead of a few example-based tests, PBT enforces a **general rule that must
always hold** and runs hundreds of generated cases against it. Here we use the
Go standard library `testing/quick` — no external dependencies.

## Layout

```
helloworld-pbt/
├── go.mod
├── cmd/helloworld/main.go   # prints the greeting (CLI)
└── greeter/
    ├── greeter.go           # Greet(name) + NormalizeName(name)
    └── greeter_test.go      # property-based + example tests
```

## The properties enforced

For **any** input string, these rules always hold:

1. `Greet` output starts with `"Hello, "` and ends with `"!"`.
2. The normalized name appears exactly between the prefix and the suffix.
3. `NormalizeName` is idempotent: `Normalize(Normalize(x)) == Normalize(x)`.
4. Any empty/whitespace-only name greets `"World"`.
5. A normalized name never has leading/trailing whitespace.

Each property is checked with `quick.Check` over hundreds of random inputs.
There are also a few explicit example tests to document intent.

## Run

```bash
# from this directory
go build ./...
go test ./... -v      # runs property + example tests
go run ./cmd/helloworld            # -> Hello, World!
go run ./cmd/helloworld Ada Lovelace  # -> Hello, Ada Lovelace!
```

## Note on Kiro PBT

The RADME.md notes Kiro's built-in PBT generation (extracting properties from
spec requirements) is **IDE-only**. This project shows the same *concept*
implemented directly in Go with `testing/quick`, so it runs anywhere without
the IDE.
