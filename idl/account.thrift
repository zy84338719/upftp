namespace go account

struct LoginRequest {
    1: string username
    2: string password
    3: bool remember
}

struct LoginResponse {
    1: bool success
    2: string token
    3: string error
}

struct LogoutRequest {}

struct LogoutResponse {
    1: bool success
    2: string error
}

struct VoidResponse {
    1: string message
}

service AccountService {
    VoidResponse HandleLoginPage() (api.get="/login")
    LoginResponse HandleLogin(1: LoginRequest req) (api.post="/api/login")
    LogoutResponse HandleLogout(1: LogoutRequest req) (api.post="/api/logout")
}
