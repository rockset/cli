package sort

// TODO should this file be generated instead?

type WorkspaceGetter interface {
	GetWorkspace() string
}

func ByWorkspace[T WorkspaceGetter](n1, n2 T) bool {
	return n1.GetWorkspace() < n2.GetWorkspace()
}

type NameGetter interface {
	GetName() string
}

func ByName[T NameGetter](n1, n2 T) bool {
	return n1.GetName() < n2.GetName()
}

type TagNameGetter interface {
	GetTagName() string
}

func ByTagName[T TagNameGetter](n1, n2 T) bool {
	return n1.GetTagName() < n2.GetTagName()
}

type RoleNameGetter interface {
	GetRoleName() string
}

func ByRoleName[T RoleNameGetter](n1, n2 T) bool {
	return n1.GetRoleName() < n2.GetRoleName()
}

type EmailGetter interface {
	GetEmail() string
}

func ByEmail[T EmailGetter](n1, n2 T) bool {
	return n1.GetEmail() < n2.GetEmail()
}

type CollectionPathGetter interface {
	GetCollectionPath() string
}

func ByCollectionPath[T CollectionPathGetter](n1, n2 T) bool {
	return n1.GetCollectionPath() < n2.GetCollectionPath()
}
