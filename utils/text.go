package utils

import "math/rand"

func GetRandText() string {
	return text[rand.Intn(len(text))]
}

var text = []string{
	"qoo",
	"woo",
	"roo",
	"too",
	"yoo",
	"poo",
	"soo",
	"doo",
	"foo",
	"goo",
	"hoo",
	"joo",
	"koo",
	"loo",
	"zoo",
	"xoo",
	"coo",
	"voo",
	"boo",
	"noo",
	"moo",
}
