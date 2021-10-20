FROM golang:1.14-alpine AS build

WORKDIR /src/ms-hash-shopping-cart/

COPY ./ /src/ms-hash-shopping-cart/

EXPOSE 8080
 
RUN CGO_ENABLED=0 go build -o /bin/ms-hash-shopping-cart

FROM scratch
COPY --from=build /bin/ms-hash-shopping-cart /bin/ms-hash-shopping-cart
COPY ./products.json /products.json
ENTRYPOINT ["/bin/ms-hash-shopping-cart"]