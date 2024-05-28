package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/grpc-user-service/models"
	"github.com/grpc-user-service/proto"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
	Users  map[int32]models.User
	lastID int32
	mu     sync.Mutex
}

func NewUserService() *UserService {
	return &UserService{
		Users:  make(map[int32]models.User),
		lastID: 0,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *proto.CreateUserReq) (*proto.UserIdResp, error) {
	// Lock the mutex to ensure safe access to lastID.
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastID++
	// check from here - vinay
	s.Users[s.lastID] = models.User{
		ID:      s.lastID,
		FName:   req.Fname,
		City:    req.City,
		Phone:   req.Phone,
		Height:  req.Height,
		Married: req.Married,
	}
	return &proto.UserIdResp{Id: s.lastID}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *proto.UserIdReq) (*proto.UserResp, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.Users[req.Id]
	if !ok {
		return nil, fmt.Errorf("user not found")
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
	s.mu.Lock()
	defer s.mu.Unlock()
	var users []*proto.UserResp
	for _, id := range req.Ids {
		user, ok := s.Users[id]
		if ok {
			users = append(users, &proto.UserResp{
				Id:      user.ID,
				Fname:   user.FName,
				City:    user.City,
				Phone:   user.Phone,
				Height:  user.Height,
				Married: user.Married,
			})
		}
	}
	return &proto.UsersResp{
		Users: users,
	}, nil
}

func (s *UserService) SearchUser(ctx context.Context, req *proto.SearchReq) (*proto.UsersResp, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var users []*proto.UserResp
	for _, user := range s.Users {
		match := true
		if req.Phone != 0 && user.Phone != req.Phone {
			match = false
		}
		if req.City != "" && user.City != req.City {
			match = false
		}
		if req.Married && !user.Married {
			match = false
		}
		if match {
			users = append(users, &proto.UserResp{
				Id:      user.ID,
				Fname:   user.FName,
				City:    user.City,
				Phone:   user.Phone,
				Height:  user.Height,
				Married: user.Married,
			})
		}

	}
	return &proto.UsersResp{
		Users: users,
	}, nil
}
