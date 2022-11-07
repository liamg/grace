package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var landlockFlags = map[int]string{
	unix.LANDLOCK_CREATE_RULESET_VERSION: "LANDLOCK_CREATE_RULESET_VERSION",
}

func AnnotateLandlockFlags(arg Arg, _ int) {
	var joins []string
	for flag, str := range landlockFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

var landlockRuleTypes = map[int]string{
	unix.LANDLOCK_RULE_PATH_BENEATH: "LANDLOCK_RULE_PATH_BENEATH",
}

func AnnotateLandlockRuleType(arg Arg, _ int) {
	var joins []string
	for flag, str := range landlockRuleTypes {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}
