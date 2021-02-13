package main

import (
	"gopher-jasmine/api"
	"gopher-jasmine/suite"
)

func main() {
	s1 := suite.NewConcurrentSuite("parent suite 1").
		It("should run one top level test", func(instance map[string]interface{}) error {
			return nil
		}).
		Describe(suite.NewConcurrentSuite("first child suite").
			It("should run first child test", func(instance map[string]interface{}) error {
				return nil
			}).
			It("should run second child test", func(instance map[string]interface{}) error {
				return nil
			}),
		)
	s2 := suite.NewConcurrentSuite("parent suite 2").
		It("should run one top level test", func(instance map[string]interface{}) error {
			return nil
		}).
		Describe(suite.NewConcurrentSuite("first child suite").
			It("should run first child test", func(instance map[string]interface{}) error {
				return nil
			}).
			It("should run second child test", func(instance map[string]interface{}) error {
				return nil
			}),
		)
	api.NewApi([]suite.Suite{s1, s2}).ListenAndServe(":9091")
}
