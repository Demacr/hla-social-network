# FROM node AS frontend
# RUN mkdir /build
# WORKDIR /build
# COPY frontend/package*.json /build/
# RUN npm install
# COPY frontend /build/
# RUN npm run build

FROM golang AS backend
RUN mkdir /build
WORKDIR /build
COPY . .
RUN go mod tidy
RUN go env
RUN CGO_ENABLED=0 go build -a -o social-network cmd/social-network.go

FROM alpine
# COPY --from=frontend /build/dist frontend/dist/
COPY --from=backend /build/social-network .
CMD ["./social-network"]
