module github.com/kadisoka/kadisoka-framework

go 1.15

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0 // indirect
	github.com/OneOfOne/xxhash v1.2.8
	github.com/alloyzeus/go-azfl v0.0.0-20210306140744-dcc1bd4a25a3
	github.com/anthonynsimon/bild v0.13.0
	github.com/aws/aws-sdk-go v1.37.32
	github.com/doug-martin/goqu/v9 v9.11.0
	github.com/emicklei/go-restful v2.15.0+incompatible
	github.com/emicklei/go-restful-openapi v1.4.1
	github.com/gabriel-vasile/mimetype v1.2.0
	github.com/go-openapi/spec v0.20.3
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/uuid v1.2.0
	github.com/gopherjs/gopherjs v0.0.0-20210202160940-bed99a852dfe // indirect
	github.com/gorilla/schema v1.2.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/jmoiron/sqlx v1.3.1
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/lib/pq v1.10.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v6 v6.0.57
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.67
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
	github.com/rez-go/crock32 v0.0.0-20210224111353-8fde4331511d
	github.com/rez-go/crux-apis v0.0.0-20200519131450-aab8ff73963b
	github.com/rez-go/stev v0.0.0-20200515184012-e0723a6f3c19
	github.com/richardlehane/crock32 v1.0.1
	github.com/rs/zerolog v1.20.0
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/square/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/thoas/stats v0.0.0-20190407194641-965cb2de1678
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	golang.org/x/crypto v0.0.0-20210314154223-e6e6c4f2bb5b
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb // indirect
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4 // indirect
	golang.org/x/sys v0.0.0-20210316164454-77fc1eacc6aa // indirect
	golang.org/x/text v0.3.5
	google.golang.org/genproto v0.0.0-20210315173758-2651cd453018 // indirect
	google.golang.org/grpc v1.36.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	// Remove when this is solved: https://github.com/etcd-io/etcd/issues/12650
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.2
	github.com/pressly/chi => github.com/go-chi/chi v0.0.0
)
