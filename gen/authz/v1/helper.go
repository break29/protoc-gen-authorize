package authzv1

func NewRole(roleStr string) Role {
	i, ok := Role_value[roleStr]
	if !ok {
		return Role_ROLE_UNSPECIFIED
	}
	return Role(i)
}
