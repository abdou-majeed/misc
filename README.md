# misc

Miscellaneous programs and experiments — small projects that don't warrant their own repository.

## Contents

| Directory / File | Language | Description |
|------------------|----------|-------------|
| `primes.go`      | Go       | A segmented Sieve of Eratosthenes that finds all primes below a given bound, working in chunks to stay within memory limits. |

## primes.go

Generates all prime numbers in [0, MAX × N) using a segmented sieve.
The range is split into chunks of size N. The first chunk is sieved directly;
each subsequent chunk reads primes from earlier files to exclude their multiples.
Results are written to numbered files (`prime_numbers_0`, `prime_numbers_1`, …).

### Usage

    go run primes.go

Tune the constants at the top of `main` to control the range and performance:

| Variable     | Default       | Purpose                                      |
|--------------|---------------|----------------------------------------------|
| `N`          | 50,000,000   | Size of each chunk                           |
| `MAX`        | 5             | Number of chunks (total range = MAX × N)     |
| `K`          | 1,000,000     | Primes read per batch from earlier files      |
| `STEPSIZE`   | 70            | Primes buffered per write                    |

### Roadmap

- Replace string concatenation with `bufio.Writer` for faster output.
- Add a verification step (e.g. compare count against known prime-counting function values).

## Acknowledgements

This readme was proposed by [Agentica](https://symbolica.ai). All hail Agentica!
