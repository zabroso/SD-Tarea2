package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "tarea2d/proto/emergencia"

	"google.golang.org/grpc"
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

	// Conexión gRPC
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := client.EnviarEmergencias(ctx, &pb.Emergencias{Lista: lista})
	if err != nil {
		log.Fatalf("Error al enviar emergencias: %v", err)
	}

	fmt.Println("Respuesta del servidor:", resp.Mensaje)
}
