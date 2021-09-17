/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	clickhousealtinitycomv1 "github.com/altinity/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClickHouseInstallations implements ClickHouseInstallationInterface
type FakeClickHouseInstallations struct {
	Fake *FakeClickhouseV1
	ns   string
}

var clickhouseinstallationsResource = schema.GroupVersionResource{Group: "clickhouse.altinity.com", Version: "v1", Resource: "clickhouseinstallations"}

var clickhouseinstallationsKind = schema.GroupVersionKind{Group: "clickhouse.altinity.com", Version: "v1", Kind: "ClickHouseInstallation"}

// Get takes name of the clickHouseInstallation, and returns the corresponding clickHouseInstallation object, and an error if there is any.
func (c *FakeClickHouseInstallations) Get(ctx context.Context, name string, options v1.GetOptions) (result *clickhousealtinitycomv1.ClickHouseInstallation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(clickhouseinstallationsResource, c.ns, name), &clickhousealtinitycomv1.ClickHouseInstallation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*clickhousealtinitycomv1.ClickHouseInstallation), err
}

// List takes label and field selectors, and returns the list of ClickHouseInstallations that match those selectors.
func (c *FakeClickHouseInstallations) List(ctx context.Context, opts v1.ListOptions) (result *clickhousealtinitycomv1.ClickHouseInstallationList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(clickhouseinstallationsResource, clickhouseinstallationsKind, c.ns, opts), &clickhousealtinitycomv1.ClickHouseInstallationList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &clickhousealtinitycomv1.ClickHouseInstallationList{ListMeta: obj.(*clickhousealtinitycomv1.ClickHouseInstallationList).ListMeta}
	for _, item := range obj.(*clickhousealtinitycomv1.ClickHouseInstallationList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clickHouseInstallations.
func (c *FakeClickHouseInstallations) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(clickhouseinstallationsResource, c.ns, opts))

}

// Create takes the representation of a clickHouseInstallation and creates it.  Returns the server's representation of the clickHouseInstallation, and an error, if there is any.
func (c *FakeClickHouseInstallations) Create(ctx context.Context, clickHouseInstallation *clickhousealtinitycomv1.ClickHouseInstallation, opts v1.CreateOptions) (result *clickhousealtinitycomv1.ClickHouseInstallation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(clickhouseinstallationsResource, c.ns, clickHouseInstallation), &clickhousealtinitycomv1.ClickHouseInstallation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*clickhousealtinitycomv1.ClickHouseInstallation), err
}

// Update takes the representation of a clickHouseInstallation and updates it. Returns the server's representation of the clickHouseInstallation, and an error, if there is any.
func (c *FakeClickHouseInstallations) Update(ctx context.Context, clickHouseInstallation *clickhousealtinitycomv1.ClickHouseInstallation, opts v1.UpdateOptions) (result *clickhousealtinitycomv1.ClickHouseInstallation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(clickhouseinstallationsResource, c.ns, clickHouseInstallation), &clickhousealtinitycomv1.ClickHouseInstallation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*clickhousealtinitycomv1.ClickHouseInstallation), err
}

// Delete takes name of the clickHouseInstallation and deletes it. Returns an error if one occurs.
func (c *FakeClickHouseInstallations) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(clickhouseinstallationsResource, c.ns, name), &clickhousealtinitycomv1.ClickHouseInstallation{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClickHouseInstallations) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(clickhouseinstallationsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &clickhousealtinitycomv1.ClickHouseInstallationList{})
	return err
}

// Patch applies the patch and returns the patched clickHouseInstallation.
func (c *FakeClickHouseInstallations) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *clickhousealtinitycomv1.ClickHouseInstallation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(clickhouseinstallationsResource, c.ns, name, pt, data, subresources...), &clickhousealtinitycomv1.ClickHouseInstallation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*clickhousealtinitycomv1.ClickHouseInstallation), err
}
