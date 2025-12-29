# 001 - Clean Architecture 重构路线图

> 创建时间: 2025-12-29
> 状态: 待实施

## 背景

当前项目在目录结构上遵循 Clean Architecture，但实现层面存在若干反模式：
- Domain 层直接依赖 GORM（违反依赖规则）
- Controller 承担业务逻辑（密码校验、token 生成）
- JWT refresh token 完全无状态（不可撤销）
- Request/Response DTO 混在 Domain 层

本文档定义按 PR 粒度拆分的渐进式重构计划。

---

## 决策记录

| 问题 | 决策 | 理由 |
|------|------|------|
| `Password` 字段名 | 保留原名 | 减少变更范围 |
| DTO 位置 | 移到 `api/dto/` | 更符合 Clean Architecture |
| Session 存储 | Redis | 高频读写场景更适合 |
| 迁移工具 | 保持 AutoMigrate | 当前规模足够 |

---

## PR 总览

| Phase | PR | 主题 | 预估工时 | 依赖 |
|-------|-----|------|---------|------|
| P0-1 | #1 | Domain 与 GORM 解耦 | 2-3h | - |
| P0-2 | #2 | JWT 解析安全修复 | 1-2h | - |
| P0-3 | #3 | DTO 分离 + Controller 瘦身 + 统一错误处理 | 2-3h | #1 |
| P0-4 | #4 | Refresh Token 可撤销机制 (Redis) | 3-4h | #2, #3 |
| P1-1 | #5 | DI 集中装配 + 迁移剥离 | 1-2h | #4 |
| P2-1 | #6 | 配置校验 + 测试补齐 | 2-3h | #5 |

### 依赖关系

```
PR #1 (Domain/GORM 解耦)
   │
   ├──► PR #2 (JWT 安全修复) ──► PR #3 (Controller 瘦身)
   │                                    │
   │                                    ▼
   │                              PR #4 (Refresh Token)
   │                                    │
   └────────────────────────────────────┼──► PR #5 (DI 集中)
                                        │         │
                                        ▼         ▼
                                   PR #6 (配置校验 + 测试)
```

---

## PR #1: Domain 与 GORM 解耦

### 目标
让 `domain/` 不依赖任何外部框架（GORM），符合 Clean Architecture 依赖规则。

### 改动文件

| 文件路径 | 操作 | 变更点 |
|---------|------|--------|
| `domain/user.go` | 修改 | 移除 `gorm.Model` 和 `gorm` import |
| `repository/model/user_model.go` | 新增 | GORM 模型 + 双向映射方法 |
| `repository/user_repository.go` | 修改 | 使用 `UserModel` 操作数据库 |
| `bootstrap/app.go` | 修改 | AutoMigrate 改用 `repository/model.UserModel` |

### 代码规范

**domain/user.go（目标状态）**
```go
package domain

import (
    "context"
    "time"
)

type User struct {
    ID        uint
    Name      string
    Email     string
    Password  string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type UserRepository interface {
    Create(ctx context.Context, user *User) error
    Fetch(ctx context.Context) ([]User, error)
    GetByEmail(ctx context.Context, email string) (User, error)
    GetByID(ctx context.Context, id string) (User, error)
    Update(ctx context.Context, user *User) error
}
```

**repository/model/user_model.go（新增）**
```go
package model

import (
    "time"
    "github.com/horaoen/go-backend-clean-architecture/domain"
)

type UserModel struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:255;not null"`
    Email     string    `gorm:"size:255;uniqueIndex;not null"`
    Password  string    `gorm:"column:password;size:255;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (UserModel) TableName() string {
    return "users"
}

func (m *UserModel) ToDomain() domain.User {
    return domain.User{
        ID:        m.ID,
        Name:      m.Name,
        Email:     m.Email,
        Password:  m.Password,
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
    }
}

func ToUserModel(u *domain.User) UserModel {
    return UserModel{
        ID:        u.ID,
        Name:      u.Name,
        Email:     u.Email,
        Password:  u.Password,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}
```

### 验收标准
- [ ] `go build ./...` 通过
- [ ] `go test ./...` 通过
- [ ] `domain/` 目录无任何 `gorm` import
- [ ] 移除 `domain/user.go` 中的 `CollectionUser` 常量（MongoDB 遗留）

---

## PR #2: JWT 解析安全修复

### 目标
修复 token 重复解析、类型断言 panic 风险、条件逻辑错误。

### 改动文件

| 文件路径 | 操作 | 变更点 |
|---------|------|--------|
| `internal/tokenutil/tokenutil.go` | 修改 | 使用 `ParseWithClaims`；修复验证条件 |
| `api/middleware/jwt_auth_middleware.go` | 修改 | 一次解析；严格 Bearer 校验 |
| `api/middleware/jwt_auth_middleware_test.go` | 新增 | 覆盖 5 种场景 |

### 代码规范

**internal/tokenutil/tokenutil.go（关键修改）**
```go
func ExtractIDFromToken(requestToken string, secret string) (string, error) {
    claims := &domain.JwtCustomClaims{}
    token, err := jwt.ParseWithClaims(requestToken, claims, func(token *jwt.Token) (any, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })
    if err != nil {
        return "", err
    }
    if !token.Valid {
        return "", fmt.Errorf("invalid token")
    }
    return claims.ID, nil
}
```

**api/middleware/jwt_auth_middleware.go（重构后）**
```go
func JwtAuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")

        if !strings.HasPrefix(authHeader, "Bearer ") {
            c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "missing or invalid authorization header"})
            c.Abort()
            return
        }

        authToken := strings.TrimPrefix(authHeader, "Bearer ")

        userID, err := tokenutil.ExtractIDFromToken(authToken, secret)
        if err != nil {
            c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "invalid or expired token"})
            c.Abort()
            return
        }

        c.Set("x-user-id", userID)
        c.Next()
    }
}
```

### 验收标准
- [ ] 无效 token 返回 401（不 panic）
- [ ] 过期 token 返回 401
- [ ] 缺少 `Bearer ` 前缀返回 401
- [ ] 签名错误返回 401
- [ ] 中间件测试覆盖以上场景

---

## PR #3: DTO 分离 + Controller 瘦身 + 统一错误处理

### 目标
- Request/Response DTO 移到 `api/dto/`
- 业务逻辑（密码校验、token 生成）移入 Usecase
- 定义领域错误并统一 HTTP 映射

### 改动文件

| 文件路径 | 操作 | 变更点 |
|---------|------|--------|
| `api/dto/auth.go` | 新增 | Login/Signup Request/Response |
| `api/dto/profile.go` | 新增 | Profile 相关 DTO |
| `api/dto/token.go` | 新增 | RefreshToken Request/Response |
| `domain/errors.go` | 新增 | 领域错误定义 |
| `domain/token.go` | 新增 | `TokenPair` + `TokenService` 接口 |
| `domain/login.go` | 修改 | 移除 DTO，简化接口 |
| `domain/signup.go` | 修改 | 移除 DTO，简化接口 |
| `domain/refresh_token.go` | 修改 | 移除 DTO |
| `domain/profile.go` | 修改 | 移除 DTO |
| `usecase/token_service.go` | 新增 | 实现 TokenService |
| `usecase/login_usecase.go` | 修改 | 封装完整登录逻辑 |
| `usecase/signup_usecase.go` | 修改 | 封装完整注册逻辑 |
| `api/controller/*.go` | 修改 | 瘦身，使用 dto |
| `api/controller/*_test.go` | 修改 | 适配新接口 |

### 代码规范

**domain/errors.go（新增）**
```go
package domain

import "errors"

var (
    ErrUserNotFound       = errors.New("user not found")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUserAlreadyExists  = errors.New("user already exists")
    ErrInternalServer     = errors.New("internal server error")
)
```

**domain/token.go（新增）**
```go
package domain

type TokenPair struct {
    AccessToken  string
    RefreshToken string
}

type TokenService interface {
    GenerateTokenPair(user *User) (TokenPair, error)
}
```

**domain/login.go（简化后）**
```go
package domain

import "context"

type LoginUsecase interface {
    Login(ctx context.Context, email, password string) (TokenPair, error)
}
```

**api/dto/auth.go（新增）**
```go
package dto

type LoginRequest struct {
    Email    string `form:"email" binding:"required,email"`
    Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
}

type SignupRequest struct {
    Name     string `form:"name" binding:"required"`
    Email    string `form:"email" binding:"required,email"`
    Password string `form:"password" binding:"required"`
}

type SignupResponse struct {
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
}
```

**usecase/login_usecase.go（丰富后）**
```go
type loginUsecase struct {
    userRepo     domain.UserRepository
    tokenService domain.TokenService
    timeout      time.Duration
}

func (lu *loginUsecase) Login(ctx context.Context, email, password string) (domain.TokenPair, error) {
    ctx, cancel := context.WithTimeout(ctx, lu.timeout)
    defer cancel()

    user, err := lu.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return domain.TokenPair{}, domain.ErrInvalidCredentials
    }

    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
        return domain.TokenPair{}, domain.ErrInvalidCredentials
    }

    return lu.tokenService.GenerateTokenPair(&user)
}
```

**api/controller/login_controller.go（瘦身后）**
```go
type LoginController struct {
    LoginUsecase domain.LoginUsecase
}

func (lc *LoginController) Login(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
        return
    }

    tokens, err := lc.LoginUsecase.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        switch {
        case errors.Is(err, domain.ErrInvalidCredentials):
            c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "invalid email or password"})
        default:
            c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "internal server error"})
        }
        return
    }

    c.JSON(http.StatusOK, dto.LoginResponse{
        AccessToken:  tokens.AccessToken,
        RefreshToken: tokens.RefreshToken,
    })
}
```

### 验收标准
- [ ] Controller 无 `bcrypt` 或 `tokenutil` import
- [ ] Controller 无 `bootstrap.Env` 依赖
- [ ] 错误消息不泄露内部细节
- [ ] 现有测试适配后通过

---

## PR #4: Refresh Token 可撤销机制 (Redis)

### 目标
引入 refresh token 白名单（Redis），支持 token rotation 和 logout 撤销。

### 新增依赖
```bash
go get github.com/redis/go-redis/v9
go get github.com/google/uuid
```

### 改动文件

| 文件路径 | 操作 | 变更点 |
|---------|------|--------|
| `bootstrap/env.go` | 修改 | 增加 Redis 配置 |
| `bootstrap/redis.go` | 新增 | Redis 连接初始化 |
| `bootstrap/app.go` | 修改 | 初始化 Redis 客户端 |
| `domain/session.go` | 新增 | Session 实体 + SessionRepository 接口 |
| `domain/jwt_custom.go` | 修改 | RefreshClaims 增加 JTI 字段 |
| `repository/session_repository.go` | 新增 | Redis 实现 |
| `internal/tokenutil/tokenutil.go` | 修改 | CreateRefreshToken 接收 jti |
| `usecase/refresh_token_usecase.go` | 修改 | 校验 jti + rotation |
| `usecase/logout_usecase.go` | 新增 | 登出逻辑 |
| `api/controller/logout_controller.go` | 新增 | 登出端点 |
| `api/route/route.go` | 修改 | 添加 /logout |
| `.env.example` | 修改 | 增加 Redis 配置 |
| `docker-compose.yaml` | 修改 | 增加 Redis 服务 |

### 代码规范

**domain/session.go（新增）**
```go
package domain

import (
    "context"
    "time"
)

type Session struct {
    JTI       string
    UserID    uint
    ExpiresAt time.Time
    CreatedAt time.Time
}

type SessionRepository interface {
    Create(ctx context.Context, session *Session) error
    Exists(ctx context.Context, jti string) (bool, error)
    Delete(ctx context.Context, jti string) error
    DeleteAllByUserID(ctx context.Context, userID uint) error
}
```

**domain/jwt_custom.go（修改）**
```go
type JwtCustomRefreshClaims struct {
    ID  string `json:"id"`
    JTI string `json:"jti"` // 新增
    jwt.RegisteredClaims
}
```

**repository/session_repository.go（Redis 实现）**
```go
package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/horaoen/go-backend-clean-architecture/domain"
)

type sessionRepository struct {
    client *redis.Client
}

func NewSessionRepository(client *redis.Client) domain.SessionRepository {
    return &sessionRepository{client: client}
}

func (r *sessionRepository) Create(ctx context.Context, session *domain.Session) error {
    key := fmt.Sprintf("session:%s", session.JTI)
    data, _ := json.Marshal(session)
    ttl := time.Until(session.ExpiresAt)
    return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *sessionRepository) Exists(ctx context.Context, jti string) (bool, error) {
    key := fmt.Sprintf("session:%s", jti)
    result, err := r.client.Exists(ctx, key).Result()
    return result > 0, err
}

func (r *sessionRepository) Delete(ctx context.Context, jti string) error {
    key := fmt.Sprintf("session:%s", jti)
    return r.client.Del(ctx, key).Err()
}

func (r *sessionRepository) DeleteAllByUserID(ctx context.Context, userID uint) error {
    // 需维护 user:sessions 索引或使用 SCAN
    // 简化版可暂不实现
    return nil
}
```

### Refresh Token 流程

```
1. 解析 refresh token -> 获取 jti, userID
2. SessionRepository.Exists(jti) -> 不存在则拒绝
3. SessionRepository.Delete(jti) -> 撤销旧 session
4. 生成新 jti (UUID)
5. SessionRepository.Create(newSession)
6. 生成新 token pair（含新 jti）
7. 返回新 tokens
```

### 验收标准
- [ ] refresh token 包含 `jti` claim
- [ ] 刷新成功后旧 token 立即失效
- [ ] `/logout` 可撤销当前会话
- [ ] Redis 服务在 docker-compose 中可用

---

## PR #5: DI 集中装配 + 迁移剥离

### 目标
- route 层不再传递 `*gorm.DB`
- 迁移逻辑独立为 `cmd/migrate`

### 改动文件

| 文件路径 | 操作 | 变更点 |
|---------|------|--------|
| `bootstrap/container.go` | 新增 | 集中装配所有依赖 |
| `api/route/route.go` | 修改 | 接收 Container |
| `api/route/*_route.go` | 修改 | 从 container 获取 controller |
| `cmd/migrate/main.go` | 新增 | 独立迁移命令 |
| `bootstrap/app.go` | 修改 | 移除 AutoMigrate |
| `cmd/main.go` | 修改 | 使用 container |

### 代码规范

**bootstrap/container.go（新增）**
```go
package bootstrap

import (
    "time"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
    "github.com/horaoen/go-backend-clean-architecture/api/controller"
    "github.com/horaoen/go-backend-clean-architecture/repository"
    "github.com/horaoen/go-backend-clean-architecture/usecase"
)

type Container struct {
    LoginController        *controller.LoginController
    SignupController       *controller.SignupController
    ProfileController      *controller.ProfileController
    RefreshTokenController *controller.RefreshTokenController
    LogoutController       *controller.LogoutController
}

func NewContainer(db *gorm.DB, redis *redis.Client, env *Env, timeout time.Duration) *Container {
    // Repositories
    userRepo := repository.NewUserRepository(db)
    sessionRepo := repository.NewSessionRepository(redis)

    // Services
    tokenService := usecase.NewTokenService(
        env.AccessTokenSecret, env.RefreshTokenSecret,
        env.AccessTokenExpiryHour, env.RefreshTokenExpiryHour,
    )

    // Usecases
    loginUC := usecase.NewLoginUsecase(userRepo, tokenService, timeout)
    signupUC := usecase.NewSignupUsecase(userRepo, sessionRepo, tokenService, timeout)
    // ...

    return &Container{
        LoginController:  &controller.LoginController{LoginUsecase: loginUC},
        SignupController: &controller.SignupController{SignupUsecase: signupUC},
        // ...
    }
}
```

**cmd/migrate/main.go（新增）**
```go
package main

import (
    "log"
    "github.com/horaoen/go-backend-clean-architecture/bootstrap"
    "github.com/horaoen/go-backend-clean-architecture/repository/model"
)

func main() {
    env := bootstrap.NewEnv()
    db := bootstrap.NewPostgres(env)

    log.Println("Running migrations...")
    err := db.AutoMigrate(&model.UserModel{})
    if err != nil {
        log.Fatal("Migration failed:", err)
    }
    log.Println("Migrations completed successfully")
}
```

### 验收标准
- [ ] `route` 目录无 `gorm` import
- [ ] `go run cmd/migrate/main.go` 可独立执行迁移
- [ ] 主应用启动不再执行 AutoMigrate

---

## PR #6: 配置校验 + 测试补齐

### 目标
- 启动时校验必填配置
- 配置连接池
- 补齐中间件和集成测试

### 改动文件

| 文件路径 | 操作 | 变更点 |
|---------|------|--------|
| `bootstrap/env.go` | 修改 | 增加 Validate() |
| `bootstrap/database.go` | 修改 | 配置连接池 |
| `api/middleware/jwt_auth_middleware_test.go` | 修改 | 确保覆盖完整 |
| `tests/integration/auth_flow_test.go` | 新增 | 端到端测试 |

### 代码规范

**bootstrap/env.go（增加校验）**
```go
func (e *Env) Validate() error {
    var missing []string

    if e.PostgresHost == "" {
        missing = append(missing, "POSTGRES_HOST")
    }
    if e.AccessTokenSecret == "" {
        missing = append(missing, "ACCESS_TOKEN_SECRET")
    }
    if e.RefreshTokenSecret == "" {
        missing = append(missing, "REFRESH_TOKEN_SECRET")
    }
    if e.RedisHost == "" {
        missing = append(missing, "REDIS_HOST")
    }

    if len(missing) > 0 {
        return fmt.Errorf("missing required config: %v", missing)
    }
    return nil
}
```

**bootstrap/database.go（连接池配置）**
```go
func NewPostgres(env *Env) *gorm.DB {
    // ... 现有连接代码 ...

    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(5)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)

    return db
}
```

### 验收标准
- [ ] 缺少必填配置时启动失败并给出明确提示
- [ ] 中间件测试覆盖完整
- [ ] 集成测试覆盖：注册 -> 登录 -> 访问受保护资源 -> 刷新 -> 登出

---

## 最终目录结构

```
gin-template/
├── api/
│   ├── controller/
│   │   ├── login_controller.go
│   │   ├── logout_controller.go          # PR #4 新增
│   │   ├── profile_controller.go
│   │   ├── refresh_token_controller.go
│   │   └── signup_controller.go
│   ├── dto/                               # PR #3 新增
│   │   ├── auth.go
│   │   ├── profile.go
│   │   └── token.go
│   ├── middleware/
│   │   ├── jwt_auth_middleware.go
│   │   └── jwt_auth_middleware_test.go    # PR #2 新增
│   └── route/
├── bootstrap/
│   ├── app.go
│   ├── container.go                       # PR #5 新增
│   ├── database.go
│   ├── env.go
│   ├── logger.go
│   └── redis.go                           # PR #4 新增
├── cmd/
│   ├── main.go
│   └── migrate/                           # PR #5 新增
│       └── main.go
├── domain/
│   ├── errors.go                          # PR #3 新增
│   ├── jwt_custom.go
│   ├── login.go
│   ├── profile.go
│   ├── refresh_token.go
│   ├── session.go                         # PR #4 新增
│   ├── signup.go
│   ├── token.go                           # PR #3 新增
│   └── user.go
├── internal/
│   └── tokenutil/
│       └── tokenutil.go
├── repository/
│   ├── model/                             # PR #1 新增
│   │   └── user_model.go
│   ├── session_repository.go              # PR #4 新增
│   └── user_repository.go
├── specs/
│   └── 001-architecture-refactor.md       # 本文档
├── tests/                                 # PR #6 新增
│   └── integration/
│       └── auth_flow_test.go
└── usecase/
    ├── login_usecase.go
    ├── logout_usecase.go                  # PR #4 新增
    ├── profile_usecase.go
    ├── refresh_token_usecase.go
    ├── signup_usecase.go
    └── token_service.go                   # PR #3 新增
```

---

## 执行检查清单

### PR #1 完成后
- [ ] `domain/` 无 `gorm` import
- [ ] 所有测试通过
- [ ] 应用可正常启动

### PR #2 完成后
- [ ] JWT 中间件测试覆盖 5 种场景
- [ ] 无 panic 风险

### PR #3 完成后
- [ ] Controller 无业务逻辑
- [ ] DTO 在 `api/dto/`
- [ ] 错误处理统一

### PR #4 完成后
- [ ] Redis 集成
- [ ] Token rotation 工作正常
- [ ] Logout 端点可用

### PR #5 完成后
- [ ] Route 层无 gorm 依赖
- [ ] 迁移独立可执行

### PR #6 完成后
- [ ] 配置校验生效
- [ ] 连接池配置
- [ ] 集成测试通过
