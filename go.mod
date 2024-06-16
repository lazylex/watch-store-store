module github.com/lazylex/watch-store/store

go 1.22.4

require (
	github.com/go-chi/chi v1.5.4
	github.com/go-chi/render v1.0.3
	github.com/go-sql-driver/mysql v1.7.1
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/golang/mock v1.6.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/lazylex/watch-store/store/pkg/secure v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.17.0
	github.com/segmentio/kafka-go v0.4.43
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/prometheus/client_model v0.4.1-0.20230718164431-9a2bf3000d16 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	golang.org/x/sys v0.11.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/lazylex/watch-store/store/pkg/db-viewer => ./pkg/db-viewer

replace github.com/lazylex/watch-store/store/pkg/secure => ./pkg/secure
