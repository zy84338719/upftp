namespace go fileops

struct WebUploadFileRequest {
    1: string path
    2: binary content
}

struct WebUploadFileResponse {
    1: bool success
    2: list<string> uploaded
    3: i32 count
    4: string error
}

struct WebCreateFolderRequest {
    1: string path
    2: string name
}

struct WebCreateFolderResponse {
    1: bool success
    2: string error
}

struct WebDeleteFileRequest {
    1: string path
}

struct WebDeleteFileResponse {
    1: bool success
    2: string error
}

struct WebRenameFileRequest {
    1: string path
    2: string newName
}

struct WebRenameFileResponse {
    1: bool success
    2: string error
}

service FileOpsService {
    WebUploadFileResponse HandleUpload(1: WebUploadFileRequest req) (api.post="/api/upload")
    WebCreateFolderResponse HandleCreateFolder(1: WebCreateFolderRequest req) (api.post="/api/create-folder")
    WebDeleteFileResponse HandleDelete(1: WebDeleteFileRequest req) (api.post="/api/delete")
    WebRenameFileResponse HandleRename(1: WebRenameFileRequest req) (api.post="/api/rename")
}
