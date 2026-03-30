namespace go common

struct FileInfo {
    1: string name
    2: string path
    3: i64 size
    4: bool is_dir
    5: string mod_time
    6: string perm
}
