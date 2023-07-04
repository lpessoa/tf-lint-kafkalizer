package rules

import (
	"github.com/stretchr/testify/assert"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func TestSchemaNameRule_Check(t *testing.T) {
	rule := NewSchemaRule()
	content := `
variable "topics_list_entry" {
  description = "Topics to be created using the provided values in 'vars.tfvars' file"
  type = list(
    object({
      topic_name             = string,
      partitions_count       = number,
      config                 = map(string),
      should_retry           = bool,
      with_dlq               = bool,
      retry_count            = number,
      retry_partitions_count = number,
      schema                 = string,
      schema_type            = string,
    })
  )
}

topics_list_entry = [
  {
    topic_name             = "sandbox.user"
    partitions_count       = 1
    config                 = {}
    should_retry           = true
    with_dlq               = true
    retry_count            = 1
    retry_partitions_count = 1
    schema                 = "user.json"
    schema_type            = "JSON"
  },
{
    topic_name             = "sandbox.user.bananas"
    partitions_count       = 1
    config                 = {}
    should_retry           = true
    with_dlq               = true
    retry_count            = 1
    retry_partitions_count = 1
    schema                 = "user.json"
    schema_type            = "JSON"
  },
]

locals {
  topics_list = [for topic in var.topics_list_entry : {
    topic_name       = topic.topic_name
    partitions_count = topic.partitions_count
    config           = topic.config
    schema_type      = topic.schema_type
    schema           = topic.schema
  }]
}

module "topics" {
  topics_list            = local.all_topics_list
}
`
	runner := helper.TestRunner(t, map[string]string{"resource.tf": content})
	err := rule.Check(runner)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	assert.Len(t, runner.Issues, 0)
}
