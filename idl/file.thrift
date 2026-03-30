namespace go file

include "common.thrift"

struct ListFilesRequest {
    1: string path
}

struct ListFilesResponse {
    1: list<common.FileInfo> files
    2: string error
}

struct GetFileRequest {
    1: string path
}

struct GetFileResponse {
    1: binary content
    2: string error
}

struct UploadFileRequest {
    1: string path
    2: binary content
}

struct UploadFileResponse {
    1: bool success
    2: string error
}

struct DeleteFileRequest {
    1: string path
}

struct DeleteFileResponse {
    1: bool success
    2: string error
}

service FileService {
    ListFilesResponse ListFiles(1: ListFilesRequest req) (api.get="/api/file/list")
    GetFileResponse GetFile(1: GetFileRequest req) (api.get="/api/file/get")
    UploadFileResponse UploadFile(1: UploadFileRequest req) (api.post="/api/file/upload")
    DeleteFileResponse DeleteFile(1: DeleteFileRequest req) (api.delete="/api/file/delete")
}
