1 - gerar cliente grcp:
  cd grpc/
  mkdir discount
  protoc --go_out=./discount --go_opt=paths=source_relative --go-grpc_out=./discount --go-grpc_opt=paths=source_relative discount.proto


imagem docker do desconto:
  docker run -d -p 50051:50051 hashorg/hash-mock-discount-service




post:
curl -X POST http://localhost:8080/checkout -H 'Content-Type: application/json' -d '{"products": [{"id": 1,"quantity": 1}]}'

