// Package main provides ...
package main

import (
	"fmt"
)

func testint() {

	var u8 uint8 = 255
	var u32 uint32 = 60500
	var f32 float32 = -4321

	fmt.Println(u8, uint32(u8))
	fmt.Println(u32, uint8(u32))
	fmt.Println(f32, uint8(f32))

	// f := uint8(200) + uint8(200)
	// fmt.Println(f)

	fmt.Println(int32(uint32(40)))
}
