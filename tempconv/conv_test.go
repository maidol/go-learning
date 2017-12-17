package tempconv

import "testing"
import "fmt"

func TestName(t *testing.T) {
	cf := CtoF(100)
	fmt.Println(cf)
}

func BenchmarkName(b *testing.B) {
	CtoF(100)
}

func ExampleName() {
	CtoF(100)
}
