
package oset_test

import (
	"testing"
	"math/rand"
	"time"
	"github.com/dlmc/ids/oset"
	"github.com/dlmc/ids/global"
	
)

//go test cpu=4 -benchmem -benchtime=5s -bench "$*.*JSON"
//TBD...

func BenchmarkOrderSet(b *testing.B) {

	rand.Seed(time.Now().UTC().UnixNano())
	a := rand.Perm(1000)
	
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ost := oset.New("oset", global.KTStr)
			for i := range a {
				ost.Create(string(i), string(i))
			}
		}
	})

}