package suite

import (
	"fmt"
)

type Describe struct {
	Skip  bool
	Suite Suite
}
type It struct {
	Skip bool
	Do   func(instance map[string]interface{}) error
}
type Action struct {
	Description string
	Do          func(instance map[string]interface{}) error
}
type Spec struct {
	Skip        bool
	Description string
	It          It
}
type ActionException struct {
	Name    string
	Message string
}
type SpecResult struct {
	Name                string           `json:"name"`
	Status              string           `json:"status"`
	Message             string           `json:"message"`
	BeforeEachException *ActionException `json:"before_each_exception"`
	AfterEachException  *ActionException `json:"after_each_exception"`
}
type Result struct {
	Name               string           `json:"name"`
	BeforeAllException *ActionException `json:"before_all_exception"`
	SpecResults        []SpecResult     `json:"spec_results"`
	Children           []Result         `json:"children"`
	AfterAllException  *ActionException `json:"after_all_exception"`
	Passed             int              `json:"passed"`
	Skipped            int              `json:"skipped"`
	Failed             int              `json:"failed"`
	TotalPassed        int              `json:"total_passed"`
	TotalSkipped       int              `json:"total_skipped"`
	TotalFailed        int              `json:"total_failed"`
}
type Suite interface {
	Run() Result
	Skip() Result
	GetName() string
	BeforeEach(description string, action func(instance map[string]interface{}) error) Suite
	AfterEach(description string, action func(instance map[string]interface{}) error) Suite
	BeforeAll(description string, action func(instance map[string]interface{}) error) Suite
	AfterAll(description string, action func(instance map[string]interface{}) error) Suite
	It(description string, assertion func(instance map[string]interface{}) error) Suite
	Xit(description string, assertion func(instance map[string]interface{}) error) Suite
	Describe(children Suite) Suite
	Xdescribe(children Suite) Suite
}

func createProcessStepFn(instance map[string]interface{}) func(action *Action) error {
	return func(action *Action) error {
		if action != nil {
			fmt.Printf("RUN Action: %s\n", action.Description)
			return action.Do(instance)
		}
		return nil
	}
}
func createAssertFn(instance map[string]interface{}) func(spec *Spec) SpecResult {
	return func(spec *Spec) SpecResult {
		fmt.Printf("RUN Spec: %s\n", spec.Description)
		err := spec.It.Do(instance)
		if err != nil {
			return SpecResult{
				Name:                spec.Description,
				Status:              "FAILED",
				Message:             err.Error(),
				BeforeEachException: nil,
				AfterEachException:  nil,
			}
		} else {
			return SpecResult{
				Name:                spec.Description,
				Status:              "PASSED",
				BeforeEachException: nil,
				AfterEachException:  nil,
			}
		}
	}
}
func skipChild(child Describe) Result {
	child.Skip = true
	return child.Suite.Skip()
}
func skipSpec(spec Spec) SpecResult {
	fmt.Printf("SKIP Spec: %s\n", spec.Description)
	return SpecResult{
		Name:                spec.Description,
		Status:              "SKIPPED",
		BeforeEachException: nil,
		AfterEachException:  nil,
	}
}
func runChild(child Describe) Result {
	if child.Skip {
		fmt.Printf("SKIP Suite: %s\n", child.Suite.GetName())
		return child.Suite.Skip()
	} else {
		return child.Suite.Run()
	}
}
func runSpec(spec Spec, instance map[string]interface{}, beforeEach *Action, assert func(spec *Spec) SpecResult, afterEach *Action) SpecResult {
	var err error
	if beforeEach != nil {
		err = beforeEach.Do(instance)
		if err != nil {
			return SpecResult{
				Name:   spec.Description,
				Status: "SKIPPED",
				BeforeEachException: &ActionException{
					Name:   beforeEach.Description,
					Message: err.Error(),
				},
				AfterEachException: nil,
			}
		}
	}
	specResult := assert(&spec)
	if afterEach != nil {
		err = afterEach.Do(instance)
		if err != nil {
			specResult.AfterEachException = &ActionException{
				Name:    afterEach.Description,
				Message: err.Error(),
			}
		}
	}
	return specResult
}
func (result *Result) CalculateResults() Result {
	var passed, skipped, failed int
	for _, specResult := range result.SpecResults {
		switch specResult.Status {
		case "PASSED":
			passed += 1
		case "SKIPPED":
			skipped += 1
		case "FAILED":
			failed += 1
		}
	}
	result.Passed = passed
	result.Skipped = skipped
	result.Failed = failed
	result.TotalPassed = passed
	result.TotalSkipped = skipped
	result.TotalFailed = failed
	if len(result.Children) == 0 {
		return *result
	} else {
		for _, child := range result.Children {
			child.CalculateResults()
		}
	}
	for _, child := range result.Children {
		result.TotalPassed += child.TotalPassed
		result.TotalSkipped += child.TotalSkipped
		result.TotalFailed += child.TotalFailed
	}
	return *result
}
