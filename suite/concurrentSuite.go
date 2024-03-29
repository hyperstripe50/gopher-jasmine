package suite

import (
	"fmt"
	"sync"
)

type ConcurrentSuite struct {
	name        string
	specs       []Spec
	children    []Describe
	beforeEach  *Action
	beforeAll   *Action
	afterEach   *Action
	afterAll    *Action
	instance    map[string]interface{}
	result      Result
	processStep func(action *Action) error
	assert      func(spec *Spec) SpecResult
}

func NewConcurrentSuite(name string) *ConcurrentSuite {
	return &ConcurrentSuite{name: name, instance: make(map[string]interface{}), result: Result{Name: name}}
}
func (suite *ConcurrentSuite) GetName() string {
	return suite.name
}
func (suite *ConcurrentSuite) Skip() Result {
	suite.result.SpecResults = skipSpecsConcurrently(suite.specs)
	suite.result.Children = skipChildrenConcurrently(suite.children)
	return suite.result.CalculateResults()
}
func (suite *ConcurrentSuite) Run() Result {
	fmt.Printf("RUN Concurrent Suite: %s\n", suite.name)
	suite.processStep = createProcessStepFn(suite.instance)
	suite.assert = createAssertFn(suite.instance)
	err := suite.processStep(suite.beforeAll)
	if err == nil {
		suite.result.SpecResults = runSpecsConcurrently(suite.specs, suite.instance, suite.beforeEach, suite.assert, suite.afterEach)
		suite.result.Children = runChildrenConcurrently(suite.children)
		err = suite.processStep(suite.afterAll)
		if err != nil {
			suite.result.AfterAllException = &ActionException{
				Name:    suite.afterAll.Description,
				Message: err.Error(),
			}
		}
	} else {
		suite.result.BeforeAllException = &ActionException{
			Name:    suite.beforeAll.Description,
			Message: err.Error(),
		}
		return suite.Skip()
	}
	result := suite.result.CalculateResults()
	fmt.Printf("RESULTS for Suite '%s': passed: %d, skipped: %d, failed: %d\n", suite.GetName(), result.TotalPassed, result.TotalSkipped, result.TotalFailed)
	return result
}
func (suite *ConcurrentSuite) BeforeEach(description string, action func(instance map[string]interface{}) error) Suite {
	suite.beforeEach = &Action{Description: description, Do: action}
	return suite
}
func (suite *ConcurrentSuite) BeforeAll(description string, action func(instance map[string]interface{}) error) Suite {
	suite.beforeAll = &Action{Description: description, Do: action}
	return suite
}
func (suite *ConcurrentSuite) AfterEach(description string, action func(instance map[string]interface{}) error) Suite {
	suite.afterEach = &Action{Description: description, Do: action}
	return suite
}
func (suite *ConcurrentSuite) AfterAll(description string, action func(instance map[string]interface{}) error) Suite {
	suite.afterAll = &Action{Description: description, Do: action}
	return suite
}
func (suite *ConcurrentSuite) It(description string, assertion func(instance map[string]interface{}) error) Suite {
	suite.specs = append(suite.specs, Spec{Description: description, It: It{Do: assertion}})
	return suite
}
func (suite *ConcurrentSuite) XIt(description string, assertion func(instance map[string]interface{}) error) Suite {
	suite.specs = append(suite.specs, Spec{Description: description, Skip: true, It: It{Do: assertion}})
	return suite
}
func (suite *ConcurrentSuite) Describe(children Suite) Suite {
	suite.children = append(suite.children, Describe{Suite: children})
	return suite
}
func (suite *ConcurrentSuite) XDescribe(children Suite) Suite {
	suite.children = append(suite.children, Describe{Skip: true, Suite: children})
	return suite
}

func runSpecsConcurrently(specs []Spec, instance map[string]interface{}, beforeEach *Action, assert func(spec *Spec) SpecResult, afterEach *Action) []SpecResult {
	results := make([]SpecResult, len(specs))
	var wg sync.WaitGroup
	for index, spec := range specs {
		wg.Add(1)
		go func(s Spec, i int) {
			defer wg.Done()
			if !s.Skip {
				results[i] = runSpec(s, instance, beforeEach, assert, afterEach)
			} else {
				fmt.Printf("SKIP Spec: %s\n", s.Description)
				results = append(results, SpecResult{
					Name:                s.Description,
					Status:              "SKIPPED",
					BeforeEachException: nil,
					AfterEachException:  nil,
				})
			}
		}(spec, index)
	}
	wg.Wait()
	return results
}
func runChildrenConcurrently(children []Describe) []Result {
	results := make([]Result, len(children))
	var wg sync.WaitGroup
	for index, child := range children {
		wg.Add(1)
		go func(c Describe, i int) {
			defer wg.Done()
			results[i] = runChild(c)
		}(child, index)
	}
	wg.Wait()
	return results
}
func skipSpecsConcurrently(specs []Spec) []SpecResult {
	results := make([]SpecResult, len(specs))
	var wg sync.WaitGroup
	for index, spec := range specs {
		wg.Add(1)
		go func(s Spec, i int) {
			defer wg.Done()
			results[i] = skipSpec(s)
		}(spec, index)
	}
	wg.Wait()
	return results
}
func skipChildrenConcurrently(children []Describe) []Result {
	results := make([]Result, len(children))
	var wg sync.WaitGroup
	for index, child := range children {
		wg.Add(1)
		go func(c Describe, i int) {
			defer wg.Done()
			results[i] = skipChild(c)
		}(child, index)
	}
	wg.Wait()
	return results
}
