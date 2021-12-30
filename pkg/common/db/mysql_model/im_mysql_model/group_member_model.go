package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"time"
)

//type GroupMember struct {
//	GroupID            string    `gorm:"column:group_id;primaryKey;"`
//	UserID             string    `gorm:"column:user_id;primaryKey;"`
//	NickName           string    `gorm:"column:nickname"`
//	FaceUrl            string    `gorm:"user_group_face_url"`
//	RoleLevel int32     `gorm:"column:role_level"`
//	JoinTime           time.Time `gorm:"column:join_time"`
//	JoinSource int32 `gorm:"column:join_source"`
//	OperatorUserID  string `gorm:"column:operator_user_id"`
//	Ex string `gorm:"column:ex"`
//}

func InsertIntoGroupMember(toInsertInfo GroupMember) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	toInsertInfo.JoinTime = time.Now()
	if toInsertInfo.RoleLevel == 0 {
		toInsertInfo.RoleLevel = constant.GroupOrdinaryUsers
	}
	err = dbConn.Table("group_member").Create(toInsertInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func GetGroupMemberListByUserID(userID string) ([]GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupMemberList []GroupMember
	err = dbConn.Table("group_member").Where("user_id=?", userID).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func GetGroupMemberListByGroupID(groupID string) ([]GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupMemberList []GroupMember
	err = dbConn.Table("group_member").Where("group_id=?", groupID).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func GetGroupMemberListByGroupIDAndRoleLevel(groupID string, roleLevel int32) ([]GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupMemberList []GroupMember
	err = dbConn.Table("group_member").Where("group_id=? and role_level=?", groupID, roleLevel).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func GetGroupMemberInfoByGroupIDAndUserID(groupID, userID string) (*GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupMember GroupMember
	err = dbConn.Table("group_member").Where("group_id=? and user_id=? ", groupID, userID).Limit(1).Find(&groupMember).Error
	if err != nil {
		return nil, err
	}
	return &groupMember, nil
}

func DeleteGroupMemberByGroupIDAndUserID(groupID, userID string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table("group_member").Where("group_id=? and user_id=? ", groupID, userID).Delete(&GroupMember{}).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateGroupMemberInfo(groupMemberInfo GroupMember) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table("group_member").Where("group_id=? and user_id=?", groupMemberInfo.GroupID, groupMemberInfo.UserID).Update(&groupMemberInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func GetOwnerManagerByGroupID(groupID string) ([]GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupMemberList []GroupMember
	err = dbConn.Table("group_member").Where("group_id=? and role_level>?", groupID, constant.GroupOrdinaryUsers).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func GetGroupMemberNumByGroupID(groupID string) uint32 {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0
	}
	var number uint32
	err = dbConn.Table("group_member").Where("group_id=?", groupID).Count(&number).Error
	if err != nil {
		return 0
	}
	return number
}

func GetGroupOwnerInfoByGroupID(groupID string) (*GroupMember, error) {
	omList, err := GetOwnerManagerByGroupID(groupID)
	if err != nil {
		return nil, err
	}
	for _, v := range omList {
		if v.RoleLevel == constant.GroupOwner {
			return &v, nil
		}
	}
	return nil, nil
}

func IsExistGroupMember(groupID, userID string) bool {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return false
	}
	var number int32
	err = dbConn.Table("group_member").Where("group_id = ? and user_id = ?", groupID, userID).Count(&number).Error
	if err != nil {
		return false
	}
	if number != 1 {
		return false
	}
	return true
}

func RemoveGroupMember(groupID string, UserID string) error {
	return DeleteGroupMemberByGroupIDAndUserID(groupID, UserID)
}

func GetMemberInfoByID(groupID string, userID string) (*GroupMember, error) {
	return GetGroupMemberInfoByGroupIDAndUserID(groupID, userID)
}

func GetGroupMemberByGroupID(groupID string, filter int32, begin int32, maxNumber int32) ([]GroupMember, error) {
	var memberList []GroupMember
	var err error
	if filter >= 0 {
		memberList, err = GetGroupMemberListByGroupIDAndRoleLevel(groupID, filter) //sorted by join time
	} else {
		memberList, err = GetGroupMemberListByGroupID(groupID)
	}

	if err != nil {
		return nil, err
	}
	if begin >= int32(len(memberList)) {
		return nil, nil
	}

	var end int32
	if begin+int32(maxNumber) < int32(len(memberList)) {
		end = begin + maxNumber
	} else {
		end = int32(len(memberList))
	}
	return memberList[begin:end], nil
}

func GetJoinedGroupIDListByUserID(userID string) ([]string, error) {
	memberList, err := GetGroupMemberListByUserID(userID)
	if err != nil {
		return nil, err
	}
	var groupIDList []string = make([]string, len(memberList))
	for _, v := range memberList {
		groupIDList = append(groupIDList, v.GroupID)
	}
	return groupIDList, nil
}

func IsGroupOwnerAdmin(groupID, UserID string) bool {
	groupMemberList, err := GetOwnerManagerByGroupID(groupID)
	if err != nil {
		return false
	}
	for _, v := range groupMemberList {
		if v.UserID == UserID && v.RoleLevel > constant.GroupOrdinaryUsers {
			return true
		}
	}
	return false
}

//
//func SelectGroupList(groupID string) ([]string, error) {
//	var groupUserID string
//	var groupList []string
//	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
//	if err != nil {
//		return groupList, err
//	}
//
//	rows, err := dbConn.Model(&GroupMember{}).Where("group_id = ?", groupID).Select("user_id").Rows()
//	if err != nil {
//		return groupList, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		rows.Scan(&groupUserID)
//		groupList = append(groupList, groupUserID)
//	}
//	return groupList, nil
//}
