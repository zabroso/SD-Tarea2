import pika
import json
from pymongo import MongoClient

# Conexión a MongoDB
mongo = MongoClient("mongodb://localhost:27017/")
db = mongo["emergencias"]
coleccion = db["registros"]

# Conexión a RabbitMQ
conexion = pika.BlockingConnection(pika.ConnectionParameters("localhost"))
canal = conexion.channel()
canal.queue_declare(queue="registro")

def callback(ch, method, properties, body):
    data = json.loads(body)
    print("[registro] Mensaje recibido:", data)

    # Si viene con estado, es actualización
    if "status" in data:
        coleccion.update_one(
            {"emergency_id": data["emergency_id"]},
            {"$set": {"status": data["status"]}}
        )
    else:
        # Insertar emergencia nueva
        coleccion.insert_one({
            "emergency_id": data["emergency_id"],
            "name": data["name"],
            "latitude": data["latitude"],
            "longitude": data["longitude"],
            "magnitude": data["magnitude"],
            "status": "En curso"
        })

canal.basic_consume(queue="registro", on_message_callback=callback, auto_ack=True)

print("[registro] Esperando mensajes...")
canal.start_consuming()
