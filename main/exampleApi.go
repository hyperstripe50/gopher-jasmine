package main

import (
	"gopher-jasmine/api"
	"gopher-jasmine/suite"
)

func main() {
	s := suite.NewSequentialSuite("parent suite").
		It("should run one top level test", func(instance map[string]interface{}) error {
			return nil
		}).
		Describe(suite.NewSequentialSuite("first child suite").
			It("should run one child test", func(instance map[string]interface{}) error {
				return nil
			}),
		)
	api.NewApi([]suite.Suite{s}).ListenAndServe(":9091")
}
