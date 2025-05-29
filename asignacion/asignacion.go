package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "tarea2d/proto" // importa el paquete generado por el .proto

	"google.golang.org/grpc"
)

type servidorAsignacion struct {
	pb.UnimplementedAsignacionServiceServer
}

func (s *servidorAsignacion) EnviarEmergencias(ctx context.Context, req *pb.Emergencias) (*pb.Estado, error) {
	fmt.Println("== Emergencias recibidas ==")
	for _, e := range req.Lista {
		fmt.Printf("Nombre: %s | Lat: %d | Lon: %d | Magnitud: %d\n", e.Name, e.Latitude, e.Longitude, e.Magnitude)
	}
	return &pb.Estado{Mensaje: "Emergencias recibidas correctamente"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAsignacionServiceServer(grpcServer, &servidorAsignacion{})

	fmt.Println("Servidor de asignación escuchando en el puerto 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
