package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"mxshop/user_srv/global"
	"mxshop/user_srv/model"
	"mxshop/user_srv/proto"
	"time"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Model2UserInfo(user model.User) proto.UserInfoResponse {
	userInfoRep := proto.UserInfoResponse{
		Id:       user.ID,
		Mobile:   user.Mobile,
		Password: user.Password,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRep.Birthday = uint64(user.Birthday.Unix())
	}
	return userInfoRep
}

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	global.DB.Scopes(Paginate(int(req.Pn), int(req.Psize))).Find(&users)
	userListRep := proto.UserListResponse{}
	for _, user := range users {
		userInfoRep := Model2UserInfo(user)
		userListRep.Data = append(userListRep.Data, &userInfoRep)
	}
	userListRep.Total = int32(len(users))
	return &userListRep, nil
}

func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var users []model.User
	res := global.DB.Table("user").Where("mobile = ?", req.Mobile).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	if len(users) == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	userInfoRep := Model2UserInfo(users[0])
	return &userInfoRep, nil
}

func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var users []model.User
	res := global.DB.Table("user").Where("id = ?", req.Id).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	if len(users) == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	userInfoRep := Model2UserInfo(users[0])
	return &userInfoRep, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	res := global.DB.Table("user").Where("mobile=?", req.Nickname).First(&user)
	if res.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已经存在")
	}
	user.Mobile = req.Nickname
	user.Nickname = req.Nickname
	user.Password = req.Password
	res = global.DB.Create(&user)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
	}
	userInfoRep := Model2UserInfo(user)
	return &userInfoRep, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*proto.CheckResponse, error) {
	var user model.User
	var check proto.CheckResponse
	res := global.DB.Table("user").Where("nickname = ?", req.Nickname).First(&user)
	if res.RowsAffected == 0 {
		check.Ok = false
		return &check, status.Errorf(codes.NotFound, "用户不存在")
	}
	user.Nickname = req.Nickname
	birthday := time.Unix(int64(req.Birthday), 0)
	user.Birthday = &birthday
	user.Gender = req.Gender
	res = global.DB.Save(&user)
	if res.Error != nil {
		check.Ok = false
		return &check, status.Errorf(codes.Internal, res.Error.Error())
	}
	check.Ok = true
	return &check, nil
}

func (s *UserServer) CheckPassword(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	var check proto.CheckResponse
	if req.EncryptedPassword == req.Password {
		check.Ok = true
	} else {
		check.Ok = false
	}
	return &check, nil
}
