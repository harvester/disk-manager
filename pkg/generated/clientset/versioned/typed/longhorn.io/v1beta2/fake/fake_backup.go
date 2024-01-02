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
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBackups implements BackupInterface
type FakeBackups struct {
	Fake *FakeLonghornV1beta2
	ns   string
}

var backupsResource = schema.GroupVersionResource{Group: "longhorn.io", Version: "v1beta2", Resource: "backups"}

var backupsKind = schema.GroupVersionKind{Group: "longhorn.io", Version: "v1beta2", Kind: "Backup"}

// Get takes name of the backup, and returns the corresponding backup object, and an error if there is any.
func (c *FakeBackups) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta2.Backup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(backupsResource, c.ns, name), &v1beta2.Backup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.Backup), err
}

// List takes label and field selectors, and returns the list of Backups that match those selectors.
func (c *FakeBackups) List(ctx context.Context, opts v1.ListOptions) (result *v1beta2.BackupList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(backupsResource, backupsKind, c.ns, opts), &v1beta2.BackupList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta2.BackupList{ListMeta: obj.(*v1beta2.BackupList).ListMeta}
	for _, item := range obj.(*v1beta2.BackupList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested backups.
func (c *FakeBackups) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(backupsResource, c.ns, opts))

}

// Create takes the representation of a backup and creates it.  Returns the server's representation of the backup, and an error, if there is any.
func (c *FakeBackups) Create(ctx context.Context, backup *v1beta2.Backup, opts v1.CreateOptions) (result *v1beta2.Backup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(backupsResource, c.ns, backup), &v1beta2.Backup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.Backup), err
}

// Update takes the representation of a backup and updates it. Returns the server's representation of the backup, and an error, if there is any.
func (c *FakeBackups) Update(ctx context.Context, backup *v1beta2.Backup, opts v1.UpdateOptions) (result *v1beta2.Backup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(backupsResource, c.ns, backup), &v1beta2.Backup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.Backup), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBackups) UpdateStatus(ctx context.Context, backup *v1beta2.Backup, opts v1.UpdateOptions) (*v1beta2.Backup, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(backupsResource, "status", c.ns, backup), &v1beta2.Backup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.Backup), err
}

// Delete takes name of the backup and deletes it. Returns an error if one occurs.
func (c *FakeBackups) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(backupsResource, c.ns, name, opts), &v1beta2.Backup{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBackups) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(backupsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta2.BackupList{})
	return err
}

// Patch applies the patch and returns the patched backup.
func (c *FakeBackups) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta2.Backup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(backupsResource, c.ns, name, pt, data, subresources...), &v1beta2.Backup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.Backup), err
}
