package main

import (
	"gopher-jasmine/api"
	"gopher-jasmine/suite"
)

func main() {
	s := suite.NewSynchronousSuite("parent suite").
		It("should run one top level test", func(instance map[string]interface{}) error {
			return nil
		}).
		Describe(func() suite.Suite {
			return suite.NewSynchronousSuite("first child suite").
				It("should run one child suite", func(instance map[string]interface{}) error {
					return nil
				})
		})
	api.NewApi([]suite.Suite{s}).ListenAndServe(":9091")
}
