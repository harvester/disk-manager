/*
Copyright 2024 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package fake

import (
	"context"

	v1beta2 "github.com/longhorn/longhorn-manager/k8s/pkg/apis/longhorn/v1beta2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBackingImageManagers implements BackingImageManagerInterface
type FakeBackingImageManagers struct {
	Fake *FakeLonghornV1beta2
	ns   string
}

var backingimagemanagersResource = v1beta2.SchemeGroupVersion.WithResource("backingimagemanagers")

var backingimagemanagersKind = v1beta2.SchemeGroupVersion.WithKind("BackingImageManager")

// Get takes name of the backingImageManager, and returns the corresponding backingImageManager object, and an error if there is any.
func (c *FakeBackingImageManagers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta2.BackingImageManager, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(backingimagemanagersResource, c.ns, name), &v1beta2.BackingImageManager{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.BackingImageManager), err
}

// List takes label and field selectors, and returns the list of BackingImageManagers that match those selectors.
func (c *FakeBackingImageManagers) List(ctx context.Context, opts v1.ListOptions) (result *v1beta2.BackingImageManagerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(backingimagemanagersResource, backingimagemanagersKind, c.ns, opts), &v1beta2.BackingImageManagerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta2.BackingImageManagerList{ListMeta: obj.(*v1beta2.BackingImageManagerList).ListMeta}
	for _, item := range obj.(*v1beta2.BackingImageManagerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested backingImageManagers.
func (c *FakeBackingImageManagers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(backingimagemanagersResource, c.ns, opts))

}

// Create takes the representation of a backingImageManager and creates it.  Returns the server's representation of the backingImageManager, and an error, if there is any.
func (c *FakeBackingImageManagers) Create(ctx context.Context, backingImageManager *v1beta2.BackingImageManager, opts v1.CreateOptions) (result *v1beta2.BackingImageManager, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(backingimagemanagersResource, c.ns, backingImageManager), &v1beta2.BackingImageManager{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.BackingImageManager), err
}

// Update takes the representation of a backingImageManager and updates it. Returns the server's representation of the backingImageManager, and an error, if there is any.
func (c *FakeBackingImageManagers) Update(ctx context.Context, backingImageManager *v1beta2.BackingImageManager, opts v1.UpdateOptions) (result *v1beta2.BackingImageManager, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(backingimagemanagersResource, c.ns, backingImageManager), &v1beta2.BackingImageManager{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.BackingImageManager), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBackingImageManagers) UpdateStatus(ctx context.Context, backingImageManager *v1beta2.BackingImageManager, opts v1.UpdateOptions) (*v1beta2.BackingImageManager, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(backingimagemanagersResource, "status", c.ns, backingImageManager), &v1beta2.BackingImageManager{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.BackingImageManager), err
}

// Delete takes name of the backingImageManager and deletes it. Returns an error if one occurs.
func (c *FakeBackingImageManagers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(backingimagemanagersResource, c.ns, name, opts), &v1beta2.BackingImageManager{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBackingImageManagers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(backingimagemanagersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta2.BackingImageManagerList{})
	return err
}

// Patch applies the patch and returns the patched backingImageManager.
func (c *FakeBackingImageManagers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta2.BackingImageManager, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(backingimagemanagersResource, c.ns, name, pt, data, subresources...), &v1beta2.BackingImageManager{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.BackingImageManager), err
}
