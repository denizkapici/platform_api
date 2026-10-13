# --- Aşama 1: Derleme (Build Stage) ---
FROM golang:1.27-alpine AS builder

# Güvenlik sertifikalarını ve saat dilimi paketlerini yükle
RUN apk --no-cache add ca-certificates tzdata

# Uygulamanın çalışacağı root olmayan kullanıcıyı oluştur
RUN adduser -D -g '' -u 10001 appuser

WORKDIR /app

# Önce bağımlılıkları kopyala (Docker cache optimizasyonu)
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyala
COPY . .

# CGO'yu kapatarak tamamen statik, harici kütüphanesiz binary üret
# -ldflags: binary boyutunu küçültmek için debug sembollerini (DWARF) siler
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/api-binary ./cmd/api

# --- Aşama 2: Dağıtım (Production Stage) ---
# Tamamen boş imaj (scratch)
FROM scratch

# Root olmayan kullanıcı bilgilerini ve sertifikaları kopyala
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Derlenen binary dosyasını kopyala
COPY --from=builder /app/api-binary /api-binary

# Root yetkisi olmadan çalıştır
USER 10001

EXPOSE 8080

ENTRYPOINT ["/api-binary"]