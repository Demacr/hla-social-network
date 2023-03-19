FROM node AS frontend
RUN mkdir /build
WORKDIR /build
COPY frontend/package*.json /build/
RUN npm install
COPY frontend /build/
RUN npm run build

FROM golang AS migrator
RUN mkdir /build
WORKDIR /build
COPY cmd/migrator ./cmd/migrator/
COPY migrations ./migrations
COPY go.mod go.sum ./
RUN go mod tidy
RUN CGO_ENABLED=0 go build -a -o migrator cmd/migrator/migrator.go

FROM golang AS backend
RUN mkdir /build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -a -o social-network cmd/social-network/social-network.go

FROM alpine
COPY migrations ./migrations
COPY --from=frontend /build/dist frontend/dist/
COPY --from=backend /build/social-network .
COPY --from=migrator /build/migrator .
CMD ./migrator -dir migrations/ "${MYSQL_USER}:${MYSQL_PASSWORD}@${MYSQL_HOST}/${MYSQL_DATABASE}?parseTime=true" up; ./social-network
