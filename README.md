1 - gerar cliente grcp:
  cd grpc/
  mkdir discount
  protoc --go_out=./discount --go_opt=paths=source_relative --go-grpc_out=./discount --go-grpc_opt=paths=source_relative discount.proto


imagem docker do desconto:
  docker run -d -p 50051:50051 hashorg/hash-mock-discount-service




post:
curl -X POST http://localhost:8080/checkout -H 'Content-Type: application/json' -d '{"products": [{"id": 1,"quantity": 1}]}'



para gerar a imagem docker, devemos entrar no diretorio do projeto e executar
docker build -t moisesba/ms-hash-shopping-cart .

após gerar, para executar rodamos via docker-compose dentro da mesma pasta do projeto
docker-compose up -d

