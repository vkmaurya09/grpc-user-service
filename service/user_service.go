package service

import (
	"context"
	"fmt"

	"github.com/grpc-user-service/models"
	"github.com/grpc-user-service/proto"
	"github.com/grpc-user-service/repository"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
	repo repository.UserRepository

	Users  map[int32]models.User
	lastID int32
}

func NewUserService() *UserService {
	return &UserService{
		repo:   repository.NewUserRepository(),
		Users:  make(map[int32]models.User),
		lastID: 0,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *proto.CreateUserReq) (*proto.UserIdResp, error) {
	user := models.User{
		FName:   req.Fname,
		City:    req.City,
		Phone:   req.Phone,
		Height:  req.Height,
		Married: req.Married,
	}
	id, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &proto.UserIdResp{Id: id}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *proto.UserIdReq) (*proto.UserResp, error) {
	user, err := s.repo.GetUser(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.UserResp{
		Id:      user.ID,
		Fname:   user.FName,
		City:    user.City,
		Phone:   user.Phone,
		Height:  user.Height,
		Married: user.Married,
	}, nil
}

func (s *UserService) GetUsers(ctx context.Context, req *proto.UserIdsReq) (*proto.UsersResp, error) {
	users, err := s.repo.GetUsers(req.Ids)
	if err != nil {
		return nil, err
	}

	var protoUsers []*proto.UserResp
	for _, user := range users {
		protoUsers = append(protoUsers, &proto.UserResp{
			Id:      user.ID,
			Fname:   user.FName,
			City:    user.City,
			Phone:   user.Phone,
			Height:  user.Height,
			Married: user.Married,
		})
	}

	return &proto.UsersResp{
		Users: protoUsers,
	}, nil
}

func (s *UserService) searchUsers(field string, value interface{}) (*proto.UsersResp, error) {
	fmt.Println(1)
	users, err := s.repo.SearchUser(field, value)
	if err != nil {
		return nil, err
	}

	var protoUsers []*proto.UserResp
	for _, user := range users {
		protoUsers = append(protoUsers, &proto.UserResp{
			Id:      user.ID,
			Fname:   user.FName,
			City:    user.City,
			Phone:   user.Phone,
			Height:  user.Height,
			Married: user.Married,
		})
	}
	return &proto.UsersResp{Users: protoUsers}, nil
}

func (s *UserService) SearchUser(ctx context.Context, req *proto.SearchReq) (*proto.UsersResp, error) {

	switch req.SearchField.(type) {
	case *proto.SearchReq_Id:
		return s.searchUsers("ID", req.GetId())
	case *proto.SearchReq_Fname:
		return s.searchUsers("FName", req.GetFname())
	case *proto.SearchReq_City:
		return s.searchUsers("City", req.GetCity())
	case *proto.SearchReq_Phone:
		return s.searchUsers("Phone", req.GetPhone())
	case *proto.SearchReq_Height:
		return s.searchUsers("Height", req.GetHeight())
	case *proto.SearchReq_Married:
		return s.searchUsers("Married", req.GetMarried())
	default:
		return nil, fmt.Errorf("invalid search field")
	}
}
