# Kratos 项目模板

## 安装 Kratos
```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

## 创建一个服务
```
# 创建一个模板项目
kratos new server

cd server
# 添加一个 proto 模板
kratos proto add api/server/server.proto
# 生成 proto 代码
kratos proto client api/server/server.proto
# 根据 proto 文件生成服务的源代码
kratos proto server api/server/server.proto -t internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```

## 通过 Makefile 生成其他辅助文件
```
# 下载并更新依赖项
make init
# 通过 proto 文件生成 API 文件（包括：pb.go、HTTP、gRPC、验证、Swagger）
make api
# 生成所有文件
make all
```

## 自动化初始化（wire）
```
# 安装 wire
go get github.com/google/wire/cmd/wire

# 生成 wire
cd cmd/server
wire
```

## 运行
```
# 运行
# 本环境运行配置依赖apollo,如无apollo将local_config.yaml的apollo相关配置删除,在配置文件中配置好数据库和redis即可运行,
# 如无数据库和缓存通过修改/helloworld/internal/data/data.go亦可运行,当然相关演示功能也将受限
# 如无需接入apm删除local_config.yaml的tracing相关配置
# 如无需接入eureka删除local_config.yaml的eureka相关配置

// ProviderSet is service providers.
var MapperSet = wire.NewSet(
    model.NewGreeterMapper,
    NewDB,  //替换注入为NewFakeCache
)

var CacheSet = wire.NewSet(
    NewGreeterCache,
    NewCache, //替换注入为NewFakeCache
)

make run
```

## Docker
```bash
# 构建 Docker 镜像
docker build -t <your-docker-image-name> .

# 运行 Docker 容器
docker run --rm -p 8000:8083 -p 9000:9003 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```


本服务基于 Kratos 框架，集成了 Apollo、Swagger、Eureka、Redis 和 GORM，同时封装了 Mapper 和 Cache对Repo的CRUD提供理想模型下的通用实现，支持缓存注入；
