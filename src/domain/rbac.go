package domain

import "fmt"

type RBACUser string
type RBACRole string
type RBACObject string
type RBACAction string

// func NewSpaceWriterRole(spaceID SpaceID) RBACRole {
// 	return RBACRole(fmt.Sprintf("space_%d_writer", uint(spaceID)))
// }

// func NewSpaceObject(spaceID SpaceID) RBACObject {
// 	return RBACObject(fmt.Sprintf("space_%d", uint(spaceID)))
// }

func NewUserObject(appUserID AppUserID) RBACUser {
	return RBACUser(fmt.Sprintf("user_%d", appUserID.Int()))
}

var PrivilegeRead = RBACAction("read")
var PrivilegeUpdate = RBACAction("update")
var PrivilegeRemove = RBACAction("remove")
