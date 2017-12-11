// Package tempconv performs Celsius and Fahrenheit conversions.
package tempconv

import "fmt"

// Celsius :
type Celsius float64 // 摄氏温度
// Fahrenheit :
type Fahrenheit float64 // 华氏温度

const (
	// AbsoluteZeroC : a
	AbsoluteZeroC Celsius = -273.15 //绝对零度
	// FreezingC : b
	FreezingC Celsius = 0
	//BoilingC : c
	BoilingC Celsius = 100
)

func (c Celsius) String() string { return fmt.Sprintf("%g°C", c) }
