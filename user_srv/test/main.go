package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mxshop/user_srv/proto"
	"strconv"
)

func TestCreateUser(c proto.UserClient) {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		user, err := c.CreateUser(context.Background(), &proto.CreateUserInfo{
			Nickname: "Kimihua" + strconv.Itoa(i),
			Password: "admin" + strconv.Itoa(i),
		})
		if err != nil {
			panic(err)
			return
		}
		fmt.Println(user.Nickname, user.Password)
	}
}

func TestGetUserList(c proto.UserClient) {
	users, err := c.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		Psize: 5,
	})
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(users.Total)
	for _, user := range users.Data {
		fmt.Println(user.Nickname, user.Password)
	}
}

func TestGetUserByMobile(c proto.UserClient) {
	user, err := c.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: "Kimihua2",
	})
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(user.Nickname, user.Password)
}

func TestGetUserById(c proto.UserClient) {
	user, err := c.GetUserById(context.Background(), &proto.IdRequest{
		Id: 3,
	})
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(user.Id, user.Nickname, user.Password)
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := proto.NewUserClient(conn)
	//TestCreateUser(c)
	//TestGetUserList(c)
	TestGetUserById(c)
	//TestGetUserByMobile(c)
}
