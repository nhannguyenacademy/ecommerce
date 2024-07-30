# Build the Go Binary.
FROM golang:1.22 AS builder
ENV CGO_ENABLED=0
ARG BUILD_REF

# Create the service directory and the copy the module files first and then
# download the dependencies. If this doesn't change, we won't need to do this
# again in future builds.

RUN mkdir /service
COPY go.* /service/
WORKDIR /service
RUN go mod download

# Copy the source code into the container.
COPY . /service

# Build the admin binary.
WORKDIR /service/tools/admin
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# Build the service binary.
WORKDIR /service/cmd/ecommerce
RUN go build -ldflags "-X main.build=${BUILD_REF}"


# Run the Go Binary in Alpine.
FROM alpine:3.20
ARG BUILD_DATE
ARG BUILD_REF
RUN addgroup -g 1000 -S service_grp && \
    adduser -u 1000 -h /service -G service_grp -S service_usr
COPY --from=builder --chown=service_grp:service_usr /service/tools/admin/admin /service/admin
COPY --from=builder --chown=service_grp:service_usr /service/cmd/ecommerce/ecommerce /service/ecommerce
COPY --from=builder --chown=service_grp:service_usr /service/configs/keys /service/configs/keys
COPY --from=builder --chown=service_grp:service_usr /service/internal/sdkbus/migrate/migrations /service/internal/sdkbus/migrate/migrations
WORKDIR /service
USER service_usr
CMD ["./ecommerce"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="ecommerce-api" \
      org.opencontainers.image.authors="Nhan Nguyen Academy <https://nhannguyen.academy>" \
      org.opencontainers.image.source="https://github.com/nhannguyenacademy/ecommerce" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Nhan Nguyen Academy"
