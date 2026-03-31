# ================================
# Stage 1 — Build (Alpine)
# ================================
FROM public.ecr.aws/docker/library/golang:1.25.4-alpine AS build
WORKDIR /src

# cache modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# static build, output named bootstrap
ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /src/bootstrap ./cmd/boleto_online_cancel

# ================================
# Stage 2 — Runtime (AWS Lambda provided image)
# ================================
FROM public.ecr.aws/lambda/provided:al2023.2025.11.05.13

# Copy bootstrap to locations esperados pelo Lambda/SAM
COPY --from=build /src/bootstrap /var/runtime/bootstrap
COPY --from=build /src/bootstrap /var/task/bootstrap

# give execute permissions
RUN chmod +x /var/runtime/bootstrap /var/task/bootstrap

# Não expor ENTRYPOINT/CMD customizados: o rapid (SAM) e o Lambda irão executar /var/runtime/bootstrap
# Porém podemos garantir um ENTRYPOINT direto caso alguém queira rodar a imagem sem rapid:
ENTRYPOINT ["/var/runtime/bootstrap"]
CMD []
