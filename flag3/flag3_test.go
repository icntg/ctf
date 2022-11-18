package main

import (
	"fmt"
	"testing"
)

func TestUtility_Serialize(t *testing.T) {
	w := Web{}
	w.init()
	fmt.Println(w)
	s := gUtility.Serialize(w)
	fmt.Println(s)
	w1 := Web{}
	x := gUtility.Deserialize(s, &w1)
	//w1 := gUtility.Deserialize(s).(Web)
	fmt.Println(x)
	fmt.Println(w1)

	key := []byte("test123")
	ss := gUtility.SerializeEncrypt(key, w)
	fmt.Println(ss)
	w2 := Web{}
	x = gUtility.DecryptDeserialize(key, ss, &w2)
	fmt.Println(x)
	fmt.Println(w2)
}
