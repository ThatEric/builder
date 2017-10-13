package builder

import (
	"testing"
)

func TestRunEnvironment(t *testing.T) {
	RunEnvironment()
}

func BenchmarkRunEvenironment(b *testing.B) {
        for n := 0; n < b.N; n++ {
         	RunEnvironment()
        }
}
