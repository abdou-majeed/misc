package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"time"
	"errors"
)

func Handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func IsPrime(n int) bool {
	if n < 0 {
		log.Fatal("Must Be Positive")
	}
	if n == 0 || n == 1 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		div := n / i // Integer division
		if div*i == n {
			return false
		}
	}
	return true
}

// ErrOverflow is returned when an arithmetic operation exceeds integer limits.
var ErrOverflow = errors.New("integer overflow")

// Add safely adds two integers and returns an error if an overflow occurs.
// This is critical for preventing infinite loops or memory panics when 
// checking multiples of very large prime numbers.
func Add(a, b int) (int, error) {
	// Check for positive overflow
	if b > 0 && a > math.MaxInt-b {
		return 0, ErrOverflow
	}
	
	// Check for negative overflow (though unlikely in your prime sieve, it's good practice)
	if b < 0 && a < math.MinInt-b {
		return 0, ErrOverflow
	}
	
	return a + b, nil
}
func PrimesLessThan(n int) []bool {
	if n < 0 {
		log.Fatal("Must Be Positive")
	}
	isPrime := make([]bool, n)
	for i := 2; i < n; i++ {
		isPrime[i] = true
	}

	nbOfPrimes := 0
	for p := 2; p < n; p++ {
		if !isPrime[p] {
			continue
		}
		nbOfPrimes++
		// Exclude the multiples of the prime number p. Careful with overflow.
		k, err := utils.Add(p, p)
		for err == nil && k < n {
			isPrime[k] = false
			k, err = utils.Add(k, p)
		}
	}
	return isPrime
}

func WritePrimes(f *os.File, isPrime []bool, offset int) (nbOfPrimesWritten int) {
	str, strlen := "", 0
	for p, prime := range isPrime {
		if !prime {
			continue
		}
		// Should I make sure offset + p does not overflow ?
		str += strconv.Itoa(offset+p) + "\n"
		strlen++
		if strlen == STEPSIZE {
			f.WriteString(str)
			nbOfPrimesWritten += STEPSIZE
			str, strlen = "", 0
		}
	}
	f.WriteString(str)
	nbOfPrimesWritten += strlen
	return
}

func ReadKPrimes(rd *bufio.Reader, K int) (primes []int) {
	if K < 0 {
		log.Fatal("Must Be Positive")
	}
	for i := 0; i < K; i++ {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		Handle(err)
		line = line[:len(line)-1]
		p, err := strconv.Atoi(string(line))
		// WHat if line is empty?
		Handle(err)
		primes = append(primes, p)
	}
	return primes
}

func ExcludeMultiples(old []int, start, stop int, new []bool) {
	for _, p := range old {
		// Find smallest multiple of p that is in [start, \infty)
		n := start / p
		n = n * p
		err := error(nil)
		if n < start {
			n, err = utils.Add(n, p)
		}
		for err == nil && n < stop {
			(new)[n-start] = false
			n, err = utils.Add(n, p)
		}
	}
	return
}

// DESCRIPTION
// We find the primes within [0, N) and write them to a file prime_numbers_0.txt,
// then we look for the primes in [N, 2N), [2N, 3N), and so forth and write them to prime_numbers_1.txt, prime_numbers_2.txt, etc.
// Everytime, we read the previous files in order to exclude the multiples of those prime numbers within the current interval
// We find all primes in [0, MAX*N)

// Global variables. Empirically, for time efficiency:
// - Maximize N (therefore minimize MAX)
// - MAX up to 4 is okay.
// - K doesn't affect much
// - STEPSIZE of 70 is around ideal
var (
	filePrefix = "prime_numbers_"
	N          = 500_000_000 // StepsRange
	MAX        = 5          // See Description
	K          = 1_000_000  // Maximum number of primes to read at a time
	STEPSIZE   = 70         // How many primes to write at a time
)

func main() {
	fmt.Println("Started")
	begin := time.Now()
	nbOfPrimes := 0 // Total number of primes written

	// Initial range: [0, N)
	f, err := os.Create(filePrefix + "0")
	Handle(err)
	isPrimes := PrimesLessThan(N)
	nbOfPrimes += WritePrimes(f, isPrimes, 0)
	f.Close()
	fmt.Printf("Range %d: [%d, %d) \tcompleted in %s\n", 1, 0, N, time.Now().Sub(begin))

	// Remaining Range: [N, MAX*N)
	for file := 1; file < MAX; file++ {
		localBegin := time.Now()
		// Seach for primes within file*N and (file+1)*N
		start, stop := file*N, (file+1)*N
		isPrimes := make([]bool, stop-start)
		for i := range isPrimes {
			isPrimes[i] = true
		}

		// read primes from previous files and exclude their multiples
		for prev := 0; prev < file; prev++ {
			filename := filePrefix + strconv.Itoa(prev)
			fr, err := os.Open(filename)
			Handle(err)

			rd := bufio.NewReader(fr)
			oldPrimes := ReadKPrimes(rd, K)
			for len(oldPrimes) > 0 {
				ExcludeMultiples(oldPrimes, start, stop, isPrimes)
				oldPrimes = ReadKPrimes(rd, K)
			}
			fr.Close()
		}

		// Write the newly obtained primes
		filename := filePrefix + strconv.Itoa(file)
		fw, err := os.Create(filename)
		Handle(err)
		nbOfPrimes += WritePrimes(fw, isPrimes, start)			
		fw.Close()
		fmt.Printf("Range %d: [%d, %d] \tcompleted in %s\n", file+1, file*N, (file+1)*N, time.Now().Sub(localBegin))
	}

	end := time.Now()
	fmt.Println(nbOfPrimes, "prime numbers written in total.")
	fmt.Println("Total duration:", end.Sub(begin))
}
