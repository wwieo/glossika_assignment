package utils

import (
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
)

var singletonSnowflakeNode *snowflake.Node

func init() {
	snowflakeNode, err := snowflake.NewNode(1)
	if err != nil {
		panic(errors.Wrap(err, "create snowflake failed"))
	}
	singletonSnowflakeNode = snowflakeNode
}

func GetSnowflakeIDInt64() int64 {
	return singletonSnowflakeNode.Generate().Int64()
}
