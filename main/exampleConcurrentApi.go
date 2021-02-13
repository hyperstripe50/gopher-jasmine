package main

import (
	"fmt"
	"gopher-jasmine/api"
	"gopher-jasmine/suite"
)

func main() {
	s1 := suite.NewConcurrentSuite("parent suite 1").
		It("should run one top level test", func(instance map[string]interface{}) error {
			return nil
		}).
		Describe(suite.NewConcurrentSuite("first child suite").
			BeforeAll("error out", func(instance map[string]interface{}) error {
				return fmt.Errorf("error!")
			}).
			It("should run first child test", func(instance map[string]interface{}) error {
				return nil
			}).
			It("should run second child test", func(instance map[string]interface{}) error {
				return nil
			}),
		).
		Describe(suite.NewConcurrentSuite("second child suite").
			BeforeEach("do something before", func(instance map[string]interface{}) error {
				return nil
			}).
			It("should run first child test", func(instance map[string]interface{}) error {
				return nil
			}).
			It("should run second child test", func(instance map[string]interface{}) error {
				return nil
			}).
			AfterAll("error out", func(instance map[string]interface{}) error {
				return fmt.Errorf("error!")
			}),
		).
		Describe(suite.NewConcurrentSuite("third child suite").
			BeforeEach("error out", func(instance map[string]interface{}) error {
				return fmt.Errorf("error!")
			}).
			It("should run first child test", func(instance map[string]interface{}) error {
				return nil
			}).
			It("should run second child test", func(instance map[string]interface{}) error {
				return nil
			}).
			AfterEach("error out", func(instance map[string]interface{}) error {
				return fmt.Errorf("error!")
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
