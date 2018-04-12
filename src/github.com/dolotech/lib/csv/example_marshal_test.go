/**********************************************************
 * Author : Michael
 * Email : dolotech@163.com
 * Last modified : 2016-07-07 23:42
 * Filename : example_marshal_test.go
 * Description :
 * *******************************************************/
package csv

import (
	"fmt"
)

func ExampleMarshal() {
	type Person struct {
		Name    string `csv:"FullName"`
		Gender  string
		Age     int
		Wallet  float32 `csv:"Bank Account"`
		Happy   bool    `true:"Yes!" false:"Sad"`
		private int     `csv:"-"`
	}

	people := []Person{
		Person{
			Name:   "Smith, Joe",
			Gender: "M",
			Age:    23,
			Wallet: 19.07,
			Happy:  false,
		},
	}

	out, _ := Marshal(people)
	fmt.Printf("%s", out)
	// Output:
	// FullName,Gender,Age,Bank Account,Happy
	// "Smith, Joe",M,23,19.07,Sad
}
