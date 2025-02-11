package hash

// member is a consistent hash member.
type member string

// String returns the string representation of the member.
func (m member) String() string {
	return string(m)
}
