FROM --platform=linux/amd64/v2 golang:1.21.5 AS build

WORKDIR /go/bin/app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

# Because of the fiber prefork usage, we need to occupy the PID 1
# https://github.com/gofiber/fiber/issues/1036#issuecomment-1605889848
# RUN apt-get update && apt-get install -y dumb-init

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app

# Distrolless image
FROM gcr.io/distroless/static-debian11:nonroot
WORKDIR /
# Copy our static executable for PID 1
# COPY --from=build /usr/bin/dumb-init /usr/bin/dumb-init
# Copy our static executable
COPY --from=build /go/bin/app/web-service-cluster_management .
# Change the user to non-root
USER 65532:65532
EXPOSE 3000
# Run the app as PID 1 with dumb-init to occupy the PID 1
# ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/web-service-cluster_management", "-t",  "inCluster"]
