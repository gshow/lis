package point

/**

role[role-n][id-hash-n][id]
role map
roleContainer

pointHashContainer
point

*/

type RoleContainer struct {
	ptcon PointsHashContainer
}

var RoleMap = make(map[uint8]*RoleContainer)
