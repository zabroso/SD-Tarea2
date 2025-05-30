package main

import (
	"context"
	"fmt"
	"log"
	"net"

	da "tarea2d/proto/dronesasignacion"
	pb "tarea2d/proto/emergencia" // importa el paquete generado por el .proto

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

// Función cliente gRPC para llamar a DronService y notificar el resultado
func notificarAlDron() {
	conn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al dron: %v", err)
	}
	defer conn.Close()

	client := da.NewAsignacionServiceClient(conn)

	res, err := client.NotificarResultado(context.Background(), &da.ResultadoEmergencia{
		DroneId:   "dron01",
		Nombre:    "Incendio Forestal Sur",
		Resultado: "Extinguido",
	})
	if err != nil {
		log.Fatalf("Error al notificar al dron: %v", err)
	}

	fmt.Println("✅ Respuesta del dron:", res.Mensaje)
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAsignacionServiceServer(grpcServer, &servidorAsignacion{})

	fmt.Println("Servidor de asignación escuchando en el puerto 50052...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
