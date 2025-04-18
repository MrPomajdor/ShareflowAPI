package permissions

type Permission int

const (
	// Represents a guest. Lowest permission
	GUEST Permission = iota
	// Represents a Student
	STUDENT
	// Represents a Teacher
	TEACHER
	// Represents a School Administrator
	ADMINISTRATOR
	// Represents a Website Administrator. Highest permission
	FULL
)

// Compares perm1 with perm2 and returns true if perm1 if higher, or the same as perm2
func ComparePermissions(perm1, perm2 Permission) bool {
	return int(perm1) >= int(perm2)
}
