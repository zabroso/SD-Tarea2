package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "SD-Tarea2/vm1/emergencia"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Estructura usada para parsear el archivo JSON
type EmergenciaArchivo struct {
	Name      string `json:"name"`
	Latitude  int32  `json:"latitude"`
	Longitude int32  `json:"longitude"`
	Magnitude int32  `json:"magnitude"`
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Uso: %s archivo.json\n", os.Args[0])
	}

	// Leer archivo JSON con emergencias
	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Error al leer el archivo: %v", err)
	}

	// Parsear emergencias desde JSON
	var emergenciasInput []EmergenciaArchivo
	if err := json.Unmarshal(file, &emergenciasInput); err != nil {
		log.Fatalf("Error al parsear el JSON: %v", err)
	}

	// Convertir a formato gRPC
	var lista pb.EmergenciaList
	for _, e := range emergenciasInput {
		lista.Emergencias = append(lista.Emergencias, &pb.Emergencia{
			Name:      e.Name,
			Latitude:  e.Latitude,
			Longitude: e.Longitude,
			Magnitude: e.Magnitude,
		})
	}

	// Conectarse al servicio de asignación (gRPC)
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al servicio de asignación: %v", err)
	}
	defer conn.Close()

	client := pb.NewAsignacionClient(conn)

	// Enviar lista de emergencias
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	respuesta, err := client.EnviarEmergencias(ctx, &lista)
	if err != nil {
		log.Fatalf("Error al enviar emergencias: %v", err)
	}

	fmt.Println("Respuesta del servicio de asignación:", respuesta.Mensaje)
}
