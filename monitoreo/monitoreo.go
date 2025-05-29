package main

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "tarea2d/proto"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type monitoreoServer struct {
	pb.UnimplementedServicioMonitoreoServer
}

func (s *monitoreoServer) RecibirActualizaciones(req *pb.Actualizacion, stream pb.ServicioMonitoreo_RecibirActualizacionesServer) error {
	// Simulación de updates periódicos (en vez de RabbitMQ real por ahora)
	for i := 0; i < 3; i++ {
		update := &pb.Actualizacion{
			Nombre: "Incendio Simulado",
			Estado: fmt.Sprintf("Actualización %d", i+1),
		}
		if err := stream.Send(update); err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
	return nil
}

func main() {
	// Conexión simulada a RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("⚠️  No se pudo conectar a RabbitMQ (simulación de monitoreo):", err)
	} else {
		defer conn.Close()
		log.Println("✅ Conectado a RabbitMQ (simulado)")
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterServicioMonitoreoServer(s, &monitoreoServer{})

	fmt.Println("📡 Servicio de Monitoreo corriendo en puerto 50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
