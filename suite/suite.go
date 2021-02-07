package suite

type Describe struct {
	Description string
	Skip        bool
	Build       func() Suite
}
type It struct {
	Skip bool
	Do func(instance map[string]interface{}) error
}
type Action struct {
	Description string
	Do func(instance map[string]interface{}) error
}
type Spec struct {
	Skip bool
	Description string
	It It
}
type ActionException struct {
	Name string
	Status string
	Message string
}
type SpecResult struct {
	Name string								`json:"name"`
	Status string							`json:"status"`
	Message string							`json:"message"`
	BeforeEachException *ActionException	`json:"before_each_exception"`
	AfterEachException *ActionException		`json:"after_each_exception"`
}
type Result struct {
	Name string								`json:"name"`
	BeforeAllException *ActionException		`json:"before_all_exception"`
	SpecResults []SpecResult				`json:"spec_results"`
	Children []Result						`json:"children"`
	AfterAllException *ActionException		`json:"after_all_exception"`
	Passed int 								`json:"passed"`
	Skipped int 							`json:"skipped"`
	Failed int								`json:"failed"`
	TotalPassed int							`json:"total_passed"`
	TotalSkipped int						`json:"total_skipped"`
	TotalFailed int							`json:"total_failed"`
}
type Suite interface {
	Run() Result
	Skip() Result
	BeforeEach(description string, action func(instance map[string]interface{}) error) Suite
	AfterEach(description string, action func(instance map[string]interface{}) error) Suite
	BeforeAll(description string, action func(instance map[string]interface{}) error) Suite
	AfterAll(description string, action func(instance map[string]interface{}) error) Suite
	It(description string, assertion func(instance map[string]interface{}) error) Suite
	Xit(description string, assertion func(instance map[string]interface{}) error) Suite
	Describe(description string, children func() Suite) Suite
	Xdescribe(description string, children func() Suite) Suite
}
func (result *Result) CalculateResults() Result {
	var passed, skipped, failed int
	for _, specResult := range result.SpecResults {
		switch specResult.Status {
		case "PASSED": passed += 1
		case "SKIPPED": skipped += 1
		case "FAILED": failed += 1
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