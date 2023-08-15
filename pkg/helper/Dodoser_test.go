package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplaceUrl(t *testing.T) {
	v := ReplaceUrl("https://docs.google.com/forms/d/e/1FAIpQLSdHTZCJvhNmw0bRuKLfOUE4JMb7X3jbhals3W4M9oKxBIdWbg/viewform")
	assert.Equal(t, "https://docs.google.com/forms/d/e/1FAIpQLSdHTZCJvhNmw0bRuKLfOUE4JMb7X3jbhals3W4M9oKxBIdWbg/formResponse", v)
}
