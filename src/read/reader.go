package read

import (
	"fmt"
	"io/ioutil"
	"os"

	// "github.com/huydang284/fixedwidth"
	"github.com/ianlopshire/go-fixedwidth"
)

// Defines a credit record for a Person
type Credit struct {
	Name           string `fixed:"1,72"`
	SocialSecurity string `fixed:"73,81"`
	// capturing as much credit tags
	CreditTag string `fixed:"82,999999999"`
}

type CreditTag struct {
	Tags []string `fixed:"9"`
}

func ReadData() []Credit {
	f, _ := os.Open("./test.dat")
	defer f.Close()

	byteValue, _ := ioutil.ReadAll(f)

	allCredits := []Credit{}

	err := fixedwidth.Unmarshal(byteValue, &allCredits)
	// err := fixedwidth.Unmarshal(byteValue, &allCredits)
	if err != nil {
		fmt.Printf("", err)
	}
	// fmt.Println(len(allCredits))
	return allCredits
}
