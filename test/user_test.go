package test

import (
	"context"
	"testing"

	"github.com/grpc-user-service/models"
	"github.com/grpc-user-service/proto"
	"github.com/grpc-user-service/service"
)

func TestUserService_CreateUser(t *testing.T) {
	userService := service.NewUserService()

	req := &proto.CreateUserReq{
		Fname:   "Vinay",
		City:    "Bangalore",
		Phone:   1234567890,
		Height:  "6.2",
		Married: false,
	}

	resp, err := userService.CreateUser(context.Background(), req)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	if resp.Id <= 0 {
		t.Error("Expected a positive user ID")
	}
}

func TestUserService_GetUser(t *testing.T) {
	userService := service.NewUserService()

	createReq := &proto.CreateUserReq{
		Fname:   "Vijay",
		City:    "Mumbai",
		Phone:   484849844,
		Height:  "5'7",
		Married: true,
	}
	createResp, err := userService.CreateUser(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create user for GetUser test: %v", err)
	}

	getReq := &proto.UserIdReq{
		Id: createResp.Id,
	}

	resp, err := userService.GetUser(context.Background(), getReq)
	if err != nil {
		t.Errorf("Failed to get user: %v", err)
	}

	expectedUser := models.User{
		ID:      createResp.Id,
		FName:   createReq.Fname,
		City:    createReq.City,
		Phone:   createReq.Phone,
		Height:  createReq.Height,
		Married: createReq.Married,
	}
	if resp.Id != expectedUser.ID || resp.Fname != expectedUser.FName || resp.City != expectedUser.City ||
		resp.Phone != expectedUser.Phone || resp.Height != expectedUser.Height || resp.Married != expectedUser.Married {
		t.Error("Retrieved user does not match expected user")
	}
}

func TestUserService_GetUsers(t *testing.T) {
	userService := service.NewUserService()

	users := []models.User{
		{FName: "Vinay", City: "Bangalore", Phone: 1234567890, Height: "6.2", Married: false},
		{FName: "VK", City: "Bangalore", Phone: 2345678901, Height: "5'6", Married: true},
	}

	var createdIDs []int32
	for _, user := range users {
		req := &proto.CreateUserReq{
			Fname:   user.FName,
			City:    user.City,
			Phone:   user.Phone,
			Height:  user.Height,
			Married: user.Married,
		}
		resp, err := userService.CreateUser(context.Background(), req)
		if err != nil {
			t.Fatalf("Failed to create user for GetUsers test: %v", err)
		}
		createdIDs = append(createdIDs, resp.Id)
	}

	getReq := &proto.UserIdsReq{Ids: createdIDs}
	resp, err := userService.GetUsers(context.Background(), getReq)
	if err != nil {
		t.Errorf("Failed to get users: %v", err)
	}

	if len(resp.Users) != len(users) {
		t.Errorf("Expected %d users, got %d", len(users), len(resp.Users))
	}
}

func TestUserService_SearchUser(t *testing.T) {
	userService := service.NewUserService()

	users := []models.User{
		{FName: "Vinay", City: "Bangalore", Phone: 1234567890, Height: "6.2", Married: false},
		{FName: "Jane", City: "San Francisco", Phone: 2345678901, Height: "5'6", Married: true},
	}

	for _, user := range users {
		req := &proto.CreateUserReq{
			Fname:   user.FName,
			City:    user.City,
			Phone:   user.Phone,
			Height:  user.Height,
			Married: user.Married,
		}
		_, err := userService.CreateUser(context.Background(), req)
		if err != nil {
			t.Fatalf("Failed to create user for SearchUser test: %v", err)
		}
	}

	tests := []struct {
		field  string
		value  interface{}
		expLen int
	}{
		{"ID", int32(1), 1},
		{"FName", "Vinay", 1},
		{"City", "Bangalore", 1},
		{"Phone", int64(2345678901), 1},
		{"Height", "5'6", 1},
		{"Married", false, 1},
	}

	for _, test := range tests {
		var req *proto.SearchReq
		switch test.field {
		case "ID":
			req = &proto.SearchReq{SearchField: &proto.SearchReq_Id{Id: test.value.(int32)}}
		case "FName":
			req = &proto.SearchReq{SearchField: &proto.SearchReq_Fname{Fname: test.value.(string)}}
		case "City":
			req = &proto.SearchReq{SearchField: &proto.SearchReq_City{City: test.value.(string)}}
		case "Phone":
			req = &proto.SearchReq{SearchField: &proto.SearchReq_Phone{Phone: test.value.(int64)}}
		case "Height":
			req = &proto.SearchReq{SearchField: &proto.SearchReq_Height{Height: test.value.(string)}}
		case "Married":
			req = &proto.SearchReq{SearchField: &proto.SearchReq_Married{Married: test.value.(bool)}}
		}

		resp, err := userService.SearchUser(context.Background(), req)
		if err != nil {
			t.Errorf("Failed to search user by %s: %v", test.field, err)
		}

		if len(resp.Users) != test.expLen {
			t.Errorf("Expected %d users for %s=%v, got %d", test.expLen, test.field, test.value, len(resp.Users))
		}
	}
}
