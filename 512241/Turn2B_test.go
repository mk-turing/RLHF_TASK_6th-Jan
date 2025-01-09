package _12241

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type ABTest interface {
	Setup()
	Run(i int)
	Teardown()
}

type MyABTest struct {
	count int
}

func (t *MyABTest) Setup() {
}

func (t *MyABTest) Run(i int) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func (t *MyABTest) Teardown() {
}

type TestVariant struct {
	test   ABTest
	weight int
}

type ABTestFramework struct {
	variants    []TestVariant
	totalWeight int
	mu          sync.RWMutex
}

func NewABTestFramework(variants []TestVariant) *ABTestFramework {
	f := &ABTestFramework{
		variants: variants,
	}
	for _, v := range variants {
		f.totalWeight += v.weight
	}
	return f
}

func (f *ABTestFramework) Run(iterations int) (result map[string]float64) {
	result = make(map[string]float64)
	for _, v := range f.variants {
		result[fmt.Sprintf("Variant%d", v.weight)] = 0.0
	}

	var wg sync.WaitGroup
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()
			variant := f.selectVariant()
			variant.test.Run(i)
			f.incrementVariantCounter(variant)
		}()
	}

	wg.Wait()

	for _, v := range f.variants {
		result[fmt.Sprintf("Variant%d", v.weight)] = float64(v.weight) / float64(f.totalWeight) * 100
	}

	return result
}

func (f *ABTestFramework) selectVariant() TestVariant {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var totalWeight int
	for _, v := range f.variants {
		totalWeight += v.weight
	}

	randomNumber := rand.Intn(totalWeight)

	currentWeight := 0
	for _, v := range f.variants {
		currentWeight += v.weight
		if randomNumber < currentWeight {
			return v
		}
	}

	// This should never happen, but returning the last variant just in case
	return f.variants[len(f.variants)-1]
}

func (f *ABTestFramework) incrementVariantCounter(variant TestVariant) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, v := range f.variants {
		if v.test == variant.test {
			f.variants[i].weight++
			f.totalWeight++
			return
		}
	}
}

func BenchmarkABTestDynamicReallocation(b *testing.B) {
	iterations := 10000
	test := &MyABTest{}
	variants := []TestVariant{
		{test: test, weight: 50},
		{test: test, weight: 50},
	}

	f := NewABTestFramework(variants)

	b.Run(fmt.Sprintf("iterations=%d", iterations), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f.Run(iterations)
		}
	})
}
