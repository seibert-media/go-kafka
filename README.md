# Go utils for Kafka

## AvroEncoder with Schema-ID for Sarama

To use Kafka with Avro encoded records the schema.AvroEncoder is useful.

```
encoder := &schema.AvroEncoder{
	SchemaId: 123,
	Content:  avroEncodedByteArray,
}
```
https://github.com/Shopify/sarama 

## SchemaRegistry

schema.Registry allow easy get schema ID from kafka schema registry

```
import "http"
schemaRegistry := schema.Registry{
	SchemaRegistryUrl: "http://schema-registry.example.com",
	HttpClient:        http.DefaultClient,
}
schemaId, err := schemaRegistry.SchemaId(subject, schemaJsonString)
```

https://github.com/actgardner/gogen-avro/

## Run Docker Compose

`sudo vi /etc/hosts`

```
127.0.0.1 kafka
127.0.0.1 zookeeper
127.0.0.1 schema-registry
127.0.0.1 ksql-server
```

```
docker-compose up
```

## KSQL-Cli

```
docker run -ti \
--net=kafka \
confluentinc/cp-ksql-cli:5.1.2 \
http://ksql-server:8088
```
