package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type Customer struct {
	Surname          string
	Name             string
	Patronymic       string
	Age              string
	DateRegistration string
}

type Shop struct {
	Name      string
	Adress    string
	IsWorking string
	Owner     string
}

type C map[string]Customer
type S map[string]Shop

type server struct {
	pb.UnimplementedGreeterServer
	Cah Cash
}

type Cash struct {
	Cust C
	Shp  S
}

var ser = server{
	Cah: Cash{
		Cust: C{
			"Ilya": Customer{
				Surname:          "Surname",
				Name:             "Name",
				Patronymic:       "Patronymic",
				Age:              "Age",
				DateRegistration: "DateRegistration",
			},
			"Ilya1": Customer{
				Surname:          "Surname1",
				Name:             "Name1",
				Patronymic:       "allWorked))",
				Age:              "Age1",
				DateRegistration: "DateRegistration1",
			},
		},
		Shp: S{
			"Shop": Shop{
				Name:      "Name",
				Adress:    "Adress",
				IsWorking: "IsWorking",
				Owner:     "Owner",
			},
			"Shop1": Shop{
				Name:      "Name1",
				Adress:    "Adress1",
				IsWorking: "IsWorking1",
				Owner:     "Owner1",
			},
		},
	},
}

func (s server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	var outStr string
	var str interface{}
	var ok bool
	input := in.GetName()
	var argument []string
	argument = strings.Split(input, ",")

	if str, ok = ser.Cah.Cust[argument[0]]; !ok {
		str = ser.Cah.Shp[argument[0]]
	}

	fmt.Println("-str-")
	fmt.Println(str)
	fmt.Println("-")

	//numfield := reflect.ValueOf(str).Elem().NumField()
	numfield := reflect.ValueOf(str).NumField()
	if argument[1] != "" {
		for i := 0; i < numfield; i++ {
			name := reflect.TypeOf(str).Field(i).Name
			fmt.Println("-name-")
			fmt.Println(name)
			fmt.Println("-")

			if argument[1] == name {
				v := reflect.ValueOf(str).FieldByName(name)
				s := v.Interface().(string)
				fmt.Println("-")
				fmt.Println(v) //fix me
				fmt.Println(s)
				fmt.Println("-")

				outStr = name + ": " + s
				fmt.Println(outStr)
				return &pb.HelloReply{Message: outStr}, nil
			}
		}
	} else {
		for i := 0; i < numfield; i++ {
			name := reflect.TypeOf(str).Field(i).Name
			v := reflect.ValueOf(str).FieldByName(name)
			outStr += fmt.Sprintf("%s: %s\n", name, v)
		}
		return &pb.HelloReply{Message: outStr}, nil
	}
	return &pb.HelloReply{Message: ""}, nil

}

/*
	func (s Shop) GetData(name string, field string) string{
		var outStr string
		str := cash.Shp[name]

		numfield := reflect.ValueOf(str).Elem().NumField()
		if field != "" {
		for i := 0; i < numfield; i++{
			name := reflect.TypeOf(str).Elem().Field(i).Name
			if field == name{
				 v := reflect.ValueOf(str).FieldByName(name)
				 s := v.Interface().(string)
				outStr = name + s
				return outStr
			}
		}
		}else {
			for i := 0; i < numfield; i++{
				name := reflect.TypeOf(str).Elem().Field(i).Name
				v := reflect.ValueOf(str).FieldByName(name)
				outStr += fmt.Sprintf("%s %s\n", name, v)
			}
			return outStr
		}
	}
*/
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
