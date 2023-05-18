package bridge

import (
	"errors"
	"regexp"
	"strings"

	"bunnyshell.com/sdk"
)

var resourceSpecRegex = regexp.MustCompile(`^(?P<namespace>[a-z0-9-]+)/(?P<kind>[a-z0-9-]+)/(?P<name>[a-z0-9-]+)$`)

var ErrInvalidResourceSpec = errors.New("invalid resource spec")

type ResourceSpec struct {
	Namespace string
	Kind      string
	Name      string
}

func NewResourceSpec(spec string) *ResourceSpec {
	match := resourceSpecRegex.FindStringSubmatch(spec)

	if len(match) == 0 {
		return nil
	}

	return &ResourceSpec{
		Namespace: match[resourceSpecRegex.SubexpIndex("namespace")],
		Kind:      match[resourceSpecRegex.SubexpIndex("kind")],
		Name:      match[resourceSpecRegex.SubexpIndex("name")],
	}
}

func (r *ResourceSpec) Match(resource sdk.ComponentResourceItem) bool {
	return r.Namespace == resource.GetNamespace() &&
		r.Kind == strings.ToLower(resource.GetKind()) &&
		r.Name == resource.GetName()
}

func (r *ResourceSpec) MatchString(spec string) bool {
	resourceSpec := NewResourceSpec(spec)

	return r.Namespace == resourceSpec.Namespace &&
		r.Kind == resourceSpec.Kind &&
		r.Name == resourceSpec.Name
}

func (r *ResourceSpec) String() string {
	return r.Namespace + "/" + r.Kind + "/" + r.Name
}
