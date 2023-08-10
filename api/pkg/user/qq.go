package user

import (
	"strings"

	qq "github.com/OhYee/auth_qq"
	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetQQConn QQ connect object
func GetQQConn() (conn *qq.Connect) {
	conn = qq.New("", "", "")

	var id, key, redirect string
	v, err := variable.Get("qq_id", "qq_key", "qq_redirect")
	if err != nil {
		return
	}
	if v.SetString("qq_id", &id) != nil {
		return
	}
	if v.SetString("qq_key", &key) != nil {
		return
	}
	if v.SetString("qq_redirect", &redirect) != nil {
		return
	}
	conn = qq.New(id, key, redirect)
	return
}

// QQConnect connect qq and return user data
func QQConnect(code string) (token, openID, unionID string, res qq.UserInfo, err error) {
	qqConn := GetQQConn()

	token, err = qqConn.Auth(code)
	if err != nil {
		return
	}

	_, openID, unionID, err = qqConn.OpenID(token)
	if err != nil {
		return
	}
	res, err = qqConn.Info(token, openID)
	if err != nil {
		return
	}

	return
}

func GetUserByQQUnionID(unionID string) *TypeDB {
	users := make([]TypeDB, 0)
	cnt, err := mongo.Find("blotter", "users", bson.M{
		"qq_union_id": unionID,
	}, nil, &users)
	if err == nil && cnt != 0 {
		return &users[0]
	}
	return nil
}

func NewUserFromQQConnect(token string, openID string, unionID string, userInfo qq.UserInfo) (u *TypeDB, err error) {
	objID := primitive.NewObjectID()
	username := objID.Hex()
	uu := GetUserByUsername(userInfo.Nickname)
	if uu == nil {
		username = userInfo.Nickname
	}
	u = &TypeDB{
		TypeBase: TypeBase{
			Username:       username,
			Avatar:         strings.Replace(userInfo.FigQQ, "http://", "https://", 1),
			Token:          "",
			Email:          "",
			QQ:             "",
			NintendoSwitch: "",
			Permission:     0,
		},
		ID: objID,

		Password: "",

		QQToken:   token,
		QQOpenID:  openID,
		QQUnionID: unionID,
	}

	_, err = mongo.Add("blotter", "users", nil, u)
	return
}

func (u *TypeDB) ConnectQQ(token string, openID string, unionID string, userinfo qq.UserInfo) (err error) {
	u.QQToken = token
	u.QQOpenID = openID
	u.QQUnionID = unionID
	_, err = mongo.Update("blotter", "users", bson.M{
		"_id": u.ID,
	}, bson.M{
		"$set": bson.M{
			"qq_token":    token,
			"qq_open_id":  openID,
			"qq_union_id": unionID,
		},
	}, nil)
	return
}
