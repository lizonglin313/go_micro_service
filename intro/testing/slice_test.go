package testing

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {

	arr1 := [...]int{1,2,3,4}
	arr2 := [...]int{1,2,3,4}

	sli1 := arr1[0:2]
	sli2 := arr2[2:4]

	fmt.Printf("len1 %d, cap1 %d\n",len(sli1), cap(sli1))
	fmt.Printf("len2 %d, cap2 %d\n",len(sli2), cap(sli2))

	fmt.Printf("arr1 %p\n", &arr1)
	fmt.Printf("arr2 %p\n", &arr2)

	fmt.Printf("s1 %p\n", &sli1[0])
	fmt.Printf("s2 %p\n", &sli2[0])

}
