package voom

import (
	"context"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
)

func fields(ctx context.Context, c *govmomi.Client) (map[int32]string, error) {
	fm, err := object.GetCustomFieldsManager(c.Client)
	if err != nil {
		return nil, err
	}
	f, err := fm.Field(ctx)
	if err != nil {
		return nil, err
	}

	m := make(map[int32]string)
	for _, x := range f {
		m[x.Key] = x.Name
	}
	return m, nil
}
