namespace go download

struct DownloadRequest {
    1: string filepath (api.path="filepath")
}

struct DownloadResponse {
    1: binary content
    2: string error
}

service DownloadService {
    DownloadResponse Download(1: DownloadRequest req) (api.get="/download/*filepath")
}
