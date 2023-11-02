
# Stage 1: Build TypeScript, TailwindCSS, and Go
FROM node:18 AS nodebuild

# Set working directory and copy pnpm lock and package.json
WORKDIR /app
COPY pnpm-lock.yaml package.json ./

# Install pnpm globally
RUN npm install -g pnpm@latest && pnpm install

# Copy the source files needed for TypeScript and TailwindCSS build
COPY . .

# Build TypeScript and TailwindCSS
RUN pnpm tailwind:build && pnpm build

# Go build stage
FROM golang:1.21.3 AS gobuild

WORKDIR /go/src/app

# Copy go.mod, go.sum, and necessary files from the previous stage
COPY go.mod go.sum ./
COPY --from=nodebuild /app/dist ./dist
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

# Final stage: Set up the runtime environment
FROM alpine:3.14

WORKDIR /app

# Copy the built binary from the Go build stage
COPY --from=gobuild /go/src/app/app .
COPY --from=gobuild /go/src/app/public ./public
COPY --from=gobuild /go/src/app/pkg/database/scripts ./pkg/database/scripts
COPY --from=gobuild /go/src/app/pkg/database/data.txt ./pkg/database/data.txt

# Also copy any static assets (like CSS) that the app might need at runtime
COPY --from=nodebuild /app/dist ./dist

EXPOSE 42069

CMD ["./app"]
