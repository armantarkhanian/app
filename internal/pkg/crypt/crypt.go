// Package crypt ...
package crypt

import "strings"

func Encode(str string) string {
	return replacer.Replace(str)
}

func Decode(str string) string {
	return replacer.Replace(str)
}

var replacer = strings.NewReplacer(

	"N", "a",
	"M", "b",
	"O", "c",
	"L", "d",
	"P", "e",
	"K", "f",
	"Q", "g",
	"J", "h",
	"R", "i",
	"I", "j",
	"S", "k",
	"H", "l",
	"T", "m",
	"G", "n",
	"U", "o",
	"F", "p",
	"V", "q",
	"E", "r",
	"W", "s",
	"D", "t",
	"X", "u",
	"C", "v",
	"Y", "w",
	"B", "x",
	"Z", "y",
	"A", "z",

	"a", "N",
	"b", "M",
	"c", "O",
	"d", "L",
	"e", "P",
	"f", "K",
	"g", "Q",
	"h", "J",
	"i", "R",
	"j", "I",
	"k", "S",
	"l", "H",
	"m", "T",
	"n", "G",
	"o", "U",
	"p", "F",
	"q", "V",
	"r", "E",
	"s", "W",
	"t", "D",
	"u", "X",
	"v", "C",
	"w", "Y",
	"x", "B",
	"y", "Z",
	"z", "A",

	"%", "0",
	"$", "1",
	"^", "2",
	"#", "3",
	"&", "4",
	"@", "5",
	"*", "6",
	"!", "7",
	"(", "8",
	")", "9",

	"0", "%",
	"1", "$",
	"2", "^",
	"3", "#",
	"4", "&",
	"5", "@",
	"6", "*",
	"7", "!",
	"8", "(",
	"9", ")",
)
