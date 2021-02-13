package suite

import "fmt"

type SequentialSuite struct {
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

func NewSequentialSuite(name string) *SequentialSuite {
	return &SequentialSuite{
		name:     name,
		instance: make(map[string]interface{}),
		result:   Result{Name: name},
	}
}
func (suite *SequentialSuite) GetName() string {
	return suite.name
}
func (suite *SequentialSuite) Skip() Result {
	suite.result.SpecResults = skipSpecsSequentially(suite.specs)
	suite.result.Children = skipChildrenSequentially(suite.children)
	return suite.result.CalculateResults()
}
func (suite *SequentialSuite) Run() Result {
	fmt.Printf("RUN Sequential Suite: %s\n", suite.name)
	suite.processStep = createProcessStepFn(suite.instance)
	suite.assert = createAssertFn(suite.instance)
	err := suite.processStep(suite.beforeAll)
	if err == nil {
		suite.result.SpecResults = runSpecsSequentially(suite.specs, suite.instance, suite.beforeEach, suite.assert, suite.afterEach)
		suite.result.Children = runChildrenSequentially(suite.children)
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
func (suite *SequentialSuite) BeforeEach(description string, action func(instance map[string]interface{}) error) Suite {
	suite.beforeEach = &Action{Description: description, Do: action}
	return suite
}
func (suite *SequentialSuite) BeforeAll(description string, action func(instance map[string]interface{}) error) Suite {
	suite.beforeAll = &Action{Description: description, Do: action}
	return suite
}
func (suite *SequentialSuite) AfterEach(description string, action func(instance map[string]interface{}) error) Suite {
	suite.afterEach = &Action{Description: description, Do: action}
	return suite
}
func (suite *SequentialSuite) AfterAll(description string, action func(instance map[string]interface{}) error) Suite {
	suite.afterAll = &Action{Description: description, Do: action}
	return suite
}
func (suite *SequentialSuite) It(description string, assertion func(instance map[string]interface{}) error) Suite {
	suite.specs = append(suite.specs, Spec{Description: description, It: It{Do: assertion}})
	return suite
}
func (suite *SequentialSuite) Xit(description string, assertion func(instance map[string]interface{}) error) Suite {
	suite.specs = append(suite.specs, Spec{Description: description, Skip: true, It: It{Do: assertion}})
	return suite
}
func (suite *SequentialSuite) Describe(children Suite) Suite {
	suite.children = append(suite.children, Describe{Suite: children})
	return suite
}
func (suite *SequentialSuite) Xdescribe(children Suite) Suite {
	suite.children = append(suite.children, Describe{Skip: true, Suite: children})
	return suite
}

func runChildrenSequentially(children []Describe) []Result {
	results := make([]Result, 0)
	for _, child := range children {
		results = append(results, runChild(child))
	}
	return results
}
func runSpecsSequentially(specs []Spec, instance map[string]interface{}, beforeEach *Action, assert func(spec *Spec) SpecResult, afterEach *Action) []SpecResult {
	results := make([]SpecResult, 0)
	for _, spec := range specs {
		if !spec.Skip {
			results = append(results, runSpec(spec, instance, beforeEach, assert, afterEach))
		} else {
			fmt.Printf("SKIP Spec: %s\n", spec.Description)
			results = append(results, SpecResult{
				Name:                spec.Description,
				Status:              "SKIPPED",
				BeforeEachException: nil,
				AfterEachException:  nil,
			})
		}
	}
	return results
}
func skipSpecsSequentially(specs []Spec) []SpecResult {
	results := make([]SpecResult, 0)
	for _, spec := range specs {
		results = append(results, skipSpec(spec))
	}
	return results
}
func skipChildrenSequentially(children []Describe) []Result {
	results := make([]Result, 0)
	for _, child := range children {
		results = append(results, skipChild(child))
	}
	return results
}
