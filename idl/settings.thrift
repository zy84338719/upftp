namespace go settings

struct GetSettingsRequest {}

struct GetSettingsResponse {
    1: map<string, string> settings
    2: string error
}

struct SetLanguageRequest {
    1: string language
}

struct SetLanguageResponse {
    1: bool success
    2: string error
}

struct SetHTTPAuthRequest {
    1: bool enabled
    2: string username
    3: string password
}

struct SetHTTPAuthResponse {
    1: bool success
    2: string error
}

struct SetServicesRequest {
    1: bool enableFTP
    2: bool enableMCP
    3: bool enableWebDAV
    4: bool enableNFS
    5: string ftpPort
    6: string mcpPort
    7: string webDAVPort
    8: string nfsPort
    9: string mcpKey
}

struct SetServicesResponse {
    1: bool success
    2: string error
    3: string mcpKey
}

struct SetFTPRequest {
    1: string username
    2: string password
}

struct SetFTPResponse {
    1: bool success
    2: string error
}

service SettingsWebService {
    GetSettingsResponse HandleGetSettings(1: GetSettingsRequest req) (api.get="/api/settings")
    SetLanguageResponse HandleSetLanguage(1: SetLanguageRequest req) (api.post="/api/settings/language")
    SetHTTPAuthResponse HandleSetHTTPAuth(1: SetHTTPAuthRequest req) (api.post="/api/settings/http-auth")
    SetServicesResponse HandleSetServices(1: SetServicesRequest req) (api.post="/api/settings/services")
    SetFTPResponse HandleSetFTP(1: SetFTPRequest req) (api.post="/api/settings/ftp")
}
