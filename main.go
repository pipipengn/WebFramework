package main

import (
	"WebFramework/orm"
	"context"
	"fmt"
)

func main() {

	type User struct {
		Name string
		Age  int
	}

	selector := orm.NewSelector[User]()
	predicate := orm.Col("name").Eq("ppp").And(orm.Col("age").Eq(18))
	res, err := selector.From("user").Where(predicate).Get(context.Background())
	fmt.Println(res, err.Error())
}
