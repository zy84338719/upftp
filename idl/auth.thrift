namespace go auth

struct AuthRequest {
    1: string username
    2: string password
}

struct AuthResponse {
    1: bool success
    2: string token
    3: string error
}

service AuthService {
    AuthResponse Authenticate(1: AuthRequest req) (api.post="/api/auth")
}
