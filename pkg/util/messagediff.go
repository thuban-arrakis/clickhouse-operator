package util

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	corev1 "k8s.io/api/core/v1"
)

// DiffResult holds the added, removed, and modified differences
type DiffResult struct {
	Added    map[string]interface{}
	Removed  map[string]interface{}
	Modified map[string]interface{}
}

// allowUnexportedFields returns a cmp.Option that allows comparison of unexported fields
func allowUnexportedFields(types ...interface{}) cmp.Option {
	typs := make([]reflect.Type, len(types))
	for i, t := range types {
		typs[i] = reflect.TypeOf(t)
	}
	return cmp.Exporter(func(t reflect.Type) bool {
		for _, typ := range typs {
			if t == typ {
				return true
			}
		}
		return false
	})
}

// DeepDiff computes the differences between two objects
func DeepDiff(x, y interface{}) (DiffResult, bool) {
	opts := []cmp.Option{
		cmpopts.IgnoreFields(corev1.Container{}, "TerminationMessagePath", "TerminationMessagePolicy"),
		cmp.Exporter(func(t reflect.Type) bool {
			// Allow unexported fields for all types
			return true
		}),
	}

	diff := DiffResult{
		Added:    make(map[string]interface{}),
		Removed:  make(map[string]interface{}),
		Modified: make(map[string]interface{}),
	}

	diffString := cmp.Diff(x, y, opts...)
	equal := diffString == ""

	if !equal {
		lines := strings.Split(diffString, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "+") {
				diff.Added[line] = line
			} else if strings.HasPrefix(line, "-") {
				diff.Removed[line] = line
			} else {
				diff.Modified[line] = line
			}
		}
	}

	return diff, equal
}

// MessageDiffString stringifies message diff
func MessageDiffString(diff DiffResult, equal bool) string {
	if equal {
		return ""
	}

	str := ""

	if len(diff.Added) > 0 {
		str += MessageDiffItemString("added items", "none", "", diff.Added)
	}

	if len(diff.Removed) > 0 {
		str += MessageDiffItemString("removed items", "none", "", diff.Removed)
	}

	if len(diff.Modified) > 0 {
		str += MessageDiffItemString("modified items", "none", "", diff.Modified)
	}

	return str
}

// MessageDiffItemString stringifies one map[string]interface{} item
func MessageDiffItemString(bannerForDiff, bannerForNoDiff, defaultPath string, items map[string]interface{}) (str string) {
	if len(items) == 0 {
		return bannerForNoDiff
	}

	str += fmt.Sprintf("Diff start -------------------------\n")
	str += fmt.Sprintf("%s num: %d\n", bannerForDiff, len(items))

	i := 0
	for path, value := range items {
		valueShort := fmt.Sprintf("%+v", value)
		valueFull := fmt.Sprintf("%s", Dump(value)) // Assuming Dump is a function to serialize the value
		if len(valueFull) < 300 {
			str += fmt.Sprintf("diff item [%d]:'%s' = '%s'\n", i, path, valueFull)
		} else {
			str += fmt.Sprintf("diff item [%d]:'%s' = '%s'\n", i, path, valueShort)
		}
		i++
	}
	str += fmt.Sprintf("Diff end -------------------------\n")

	return str
}
