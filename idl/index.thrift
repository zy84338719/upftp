namespace go index

struct GetServerInfoRequest {}

struct GetServerInfoResponse {
    1: map<string, string> info
    2: string error
}

struct ListFilesLegacyRequest {
    1: string path
}

struct ListFilesLegacyResponse {
    1: list<map<string, string>> files
    2: string error
}

struct GetDirectoryTreeRequest {
    1: string path
    2: i32 depth
}

struct GetDirectoryTreeResponse {
    1: map<string, string> tree
    2: string error
}

struct GetQRCodeRequest {
    1: string url
}

struct GetQRCodeResponse {
    1: binary content
    2: string error
}

struct IndexPageResponse {
    1: string message
}

service IndexService {
    IndexPageResponse HandleIndexPage() (api.get="/")
    GetServerInfoResponse HandleServerInfo(1: GetServerInfoRequest req) (api.get="/api/info")
    ListFilesLegacyResponse HandleFileListAPI(1: ListFilesLegacyRequest req) (api.get="/api/files")
    GetDirectoryTreeResponse HandleDirectoryTree(1: GetDirectoryTreeRequest req) (api.get="/api/tree")
    GetQRCodeResponse HandleQRCode(1: GetQRCodeRequest req) (api.get="/api/qrcode")
}
