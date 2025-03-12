package model

type TableName string

func (c TableName) String() string {
	return string(c)
}

type ColumnName string

func (c ColumnName) String() string {
	return string(c)
}
