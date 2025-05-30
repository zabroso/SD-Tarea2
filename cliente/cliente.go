package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "tarea2d/proto/emergencia"
	mon "tarea2d/proto/monitoreo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EmergenciaArchivo struct {
	Name      string `json:"name"`
	Latitude  int32  `json:"latitude"`
	Longitude int32  `json:"longitude"`
	Magnitude int32  `json:"magnitude"`
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Uso: %s archivo.json", os.Args[0])
	}
	archivo := os.Args[1]

	// Leer archivo JSON
	file, err := os.ReadFile(archivo)
	if err != nil {
		log.Fatalf("Error al leer archivo: %v", err)
	}

	var emergenciasInput []EmergenciaArchivo
	if err := json.Unmarshal(file, &emergenciasInput); err != nil {
		log.Fatalf("Error al parsear JSON: %v", err)
	}

	// Conexión gRPC a asignación
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al servicio de asignación: %v", err)
	}
	defer conn.Close()

	client := pb.NewAsignacionServiceClient(conn)

	// Convertir a estructura del proto
	var lista []*pb.Emergencia
	for _, e := range emergenciasInput {
		lista = append(lista, &pb.Emergencia{
			Name:      e.Name,
			Latitude:  e.Latitude,
			Longitude: e.Longitude,
			Magnitude: e.Magnitude,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.EnviarEmergencias(ctx, &pb.Emergencias{Lista: lista})
	if err != nil {
		log.Fatalf("Error al enviar emergencias: %v", err)
	}

	fmt.Println("📨 Respuesta del servidor de asignación:", resp.Mensaje)

	// ----------------------------
	// Conexión al servicio de monitoreo
	// ----------------------------

	monConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al servicio de monitoreo: %v", err)
	}
	defer monConn.Close()

	monClient := mon.NewServicioMonitoreoClient(monConn)

	// Enviar solicitud de actualizaciones
	stream, err := monClient.RecibirActualizaciones(context.Background(), &mon.Actualizacion{
		Nombre: "Emergencia Cliente",
	})
	if err != nil {
		log.Fatalf("Error al recibir actualizaciones: %v", err)
	}

	fmt.Println("🛰️ Escuchando actualizaciones en tiempo real:")
	for {
		update, err := stream.Recv()
		if err != nil {
			fmt.Println("🚫 Fin del stream o error:", err)
			break
		}
		fmt.Printf("📡 [%s] %s\n", update.Nombre, update.Estado)
	}
}
