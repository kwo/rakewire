package logging

import (
	"regexp"
	"testing"
)

func TestPattern(t *testing.T) {

	testString := "2015/12/09 16:25:27 [INFO]  [main]  caught signal"
	p := regexp.MustCompile(`^(.*)\s+\[(\w+)\]\s+\[(\w+)\]\s+(.*)$`)
	matches := p.FindStringSubmatch(testString)
	numMatches := len(matches)

	if numMatches != 5 {
		t.Fatalf("incorrect number of matches: %d", numMatches)
	}

	t.Logf("prefix: !%s!", matches[1])
	t.Logf("level:  !%s!", matches[2])
	t.Logf("name:   !%s!", matches[3])
	t.Logf("suffix: !%s!", matches[4])

}
