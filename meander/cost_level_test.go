package meander

import (
	"fmt"
	"testing"
)

func TestCostValues(t *testing.T) {
	for _, tc := range []struct {
		cost     Cost
		expected int
	}{
		{Cost1, 1},
		{Cost2, 2},
		{Cost3, 3},
		{Cost4, 4},
		{Cost5, 5},
	} {
		if g, w := int(tc.cost), tc.expected; g != w {
			t.Errorf("\nWrong value for Cost\ngot:  %v\nwant: %v", g, w)
		}
	}
}

func TestCostString(t *testing.T) {
	for _, tc := range []struct {
		cost     Cost
		expected string
	}{
		{Cost1, "$"},
		{Cost2, "$$"},
		{Cost3, "$$$"},
		{Cost4, "$$$$"},
		{Cost5, "$$$$$"},
		{6, "invalid"},
	} {
		if g, w := tc.cost.String(), tc.expected; g != w {
			t.Errorf("\nWrong value for Cost\ngot:  %v\nwant: %v", g, w)
		}
	}
}

func ExampleCost_String() {
	c := Cost2
	fmt.Println(c.String())

	// Output:
	// $$
}

func TestParseCost(t *testing.T) {
	for _, tc := range []struct {
		costStr  string
		expected Cost
	}{
		{"$", Cost1},
		{"$$", Cost2},
		{"$$$", Cost3},
		{"$$$$", Cost4},
		{"$$$$$", Cost5},
	} {
		if g, w := ParseCost(tc.costStr), tc.expected; g != w {
			t.Errorf("\nWrong value for Cost\ngot:  %v\nwant: %v", g, w)
		}
	}
}

func TestParseCostRange(t *testing.T) {
	for _, tc := range []struct {
		rangeStr string
		expected CostRange
	}{
		{
			"$$...$$$",
			CostRange{
				From: Cost2,
				To:   Cost3,
			},
		},
		{
			"$...$$$$$",
			CostRange{
				From: Cost1,
				To:   Cost5,
			},
		},
	} {
		actual := ParseCostRange(tc.rangeStr)
		if g, w := actual.From, tc.expected.From; g != w {
			t.Errorf("\nIncorrect 'From' cost\ngot:  %v\nwant: %v", g, w)
		}
		if g, w := actual.From, tc.expected.From; g != w {
			t.Errorf("\nIncorrect 'To' cost\ngot:  %v\nwant: %v", g, w)
		}
	}
}

func TestCostRangeString(t *testing.T) {
	for _, tc := range []struct {
		r        CostRange
		expected string
	}{
		{
			CostRange{
				From: Cost2,
				To:   Cost5,
			},
			"$$...$$$$$",
		},
		{
			CostRange{
				From: Cost3,
				To:   Cost4,
			},
			"$$$...$$$$",
		},
	} {
		if g, w := tc.r.String(), tc.expected; g != w {
			t.Errorf("\nIncorrect cost string\ngot:  %v\nwant: %v", g, w)
		}
	}
}

func ExampleCostRange_String() {
	cr := &CostRange{
		From: Cost2,
		To:   Cost5,
	}

	fmt.Println(cr.String())
	// Output:
	// $$...$$$$$
}
