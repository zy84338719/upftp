namespace go system

struct GetServerStatusRequest {}

struct GetServerStatusResponse {
    1: string ip
    2: i32 http_port
    3: i32 ftp_port
    4: string root_dir
    5: bool http_enabled
    6: bool ftp_enabled
    7: bool webdav_enabled
    8: bool nfs_enabled
    9: string error
}

struct StartServerRequest {
    1: string root_dir
    2: i32 http_port
    3: i32 ftp_port
    4: bool enable_ftp
    5: bool enable_webdav
    6: bool enable_nfs
}

struct StartServerResponse {
    1: bool success
    2: string error
}

struct StopServerRequest {}

struct StopServerResponse {
    1: bool success
    2: string error
}

service SystemService {
    GetServerStatusResponse GetServerStatus(1: GetServerStatusRequest req) (api.get="/api/system/status")
    StartServerResponse StartServer(1: StartServerRequest req) (api.post="/api/system/start")
    StopServerResponse StopServer(1: StopServerRequest req) (api.post="/api/system/stop")
}
