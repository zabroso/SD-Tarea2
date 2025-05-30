package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "tarea2d/proto/dronesasignacion" // Importa el paquete generado por el .proto

	"google.golang.org/grpc"
)

type droneServer struct {
	pb.UnimplementedAsignacionServiceServer
}

func (s *droneServer) EnviarEmergencias(ctx context.Context, req *pb.AsignacionEmergencia) (*pb.Respuesta, error) {
	fmt.Println("🚁 Recibida emergencia en drones.go (esto luego se manejará con lógica)")
	return &pb.Respuesta{Mensaje: "Dron asignado y en camino"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAsignacionServiceServer(s, &droneServer{}) // Reutilizando la misma interfaz

	fmt.Println("🚁 Servicio de Drones activo en puerto 50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
