# Spec: 增强认证流程并完善单元测试

## 背景
本项目作为一个 Go 后端模板，核心认证流程的稳定性和可测试性至关重要。目前已实现基础的 Signup, Login 和 Refresh Token，但需要通过完善的单元测试来确保其健壮性，并补充“修改密码”功能以提供完整的用户管理。

## 目标
1. 完善 Signup, Login, Refresh Token 的单元测试，确保关键逻辑覆盖率 >80%。
2. 实现“修改密码”功能，包括 API 接口、Usecase 层和 Repository 层。
3. 确保所有新增代码符合整洁架构（Clean Architecture）规范。

## 核心功能
* **测试完善**：
    * 为 `signup_controller` 和 `signup_usecase` 编写测试。
    * 为 `login_controller` 和 `login_usecase` 编写测试。
    * 为 `refresh_token_controller` 和 `refresh_token_usecase` 编写测试。
* **修改密码**：
    * 新增 `POST /profile/change-password` 接口（受 JWT 保护）。
    * 验证旧密码，更新新密码（需加密存储）。

## 验收标准
* 所有单元测试通过且覆盖率 >80%。
* 修改密码功能正常工作，旧密码验证失败时返回错误。
* 代码风格符合 Go 项目规范。
