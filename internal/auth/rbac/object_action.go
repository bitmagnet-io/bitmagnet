package rbac

import "cmp"

// type ObjectAction interface {
// 	Namespace() string
// 	Object() string
// 	Action() string
// 	Compare(other ObjectAction) int
// }

type ObjectAction struct {
	Namespace string
	Object    string
	Action    string
}

// func (oa objectAction) Namespace() string {
// 	return oa.namespace
// }

// func (oa objectAction) Object() string {
// 	return oa.object
// }

// func (oa objectAction) Action() string {
// 	return oa.action
// }

func (oa ObjectAction) Compare(other ObjectAction) int {
	r := cmp.Compare(oa.Namespace, other.Namespace)
	if r != 0 {
		return r
	}

	r = cmp.Compare(oa.Object, other.Object)
	if r != 0 {
		return r
	}

	return cmp.Compare(oa.Action, other.Action)
}

func NewObjectAction(namespace string, object string, action string) ObjectAction {
	return ObjectAction{
		Namespace: namespace,
		Object:    object,
		Action:    action,
	}
}
