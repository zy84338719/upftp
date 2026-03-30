var LANG = {
    en: { all_files:"All Files", name:"Name", size:"Size", modified:"Modified", actions:"Actions",
          search:"Search files...", upload:"Upload", copy_link:"Copy Link", download:"Download",
          preview:"Preview", empty:"This folder is empty", root:"Root", loading:"Loading...",
          preview_err:"Preview not available", copied:"Link copied!", explorer:"EXPLORER",
          qr_download:"QR Download", qr_title:"Scan to Download", zip_download:"ZIP Download",
          settings:"Settings", save_success:"Saved!", save_error:"Save failed" },
    zh: { all_files:"所有文件", name:"文件名", size:"大小", modified:"修改时间", actions:"操作",
          search:"搜索文件...", upload:"上传", copy_link:"复制链接", download:"下载",
          preview:"预览", empty:"此文件夹为空", root:"根目录", loading:"加载中...",
          preview_err:"无法预览此文件", copied:"链接已复制！", explorer:"文件树",
          qr_download:"二维码下载", qr_title:"扫码下载", zip_download:"打包下载",
          settings:"设置", save_success:"已保存！", save_error:"保存失败" }
};

var serverLang = window.__UPFTP_CONFIG__.language;
var httpAuthOn = window.__UPFTP_CONFIG__.httpAuthOn;
var httpAuthUser = window.__UPFTP_CONFIG__.httpAuthUser;
var httpAuthPass = window.__UPFTP_CONFIG__.httpAuthPass;
var curLang = serverLang || localStorage.getItem('upftp-lang') || 'en';
var allFiles = [];
var curPath = '/';

function T(k) { return (LANG[curLang] || LANG.en)[k] || k; }

function navigateTo(p) {
    curPath = p;
    window.history.pushState({path: p}, '', p === '/' ? '/' : p);
    fetch('/api/files?path=' + encodeURIComponent(p)).then(function(r) { return r.json(); }).then(function(data) {
        if (data.files) {
            allFiles = data.files;
            curPath = data.path || p;
            updateBreadcrumb();
            render(allFiles);
            updateTreeActive();
            document.getElementById('searchInput').value = '';
        }
    });
}

function updateTreeActive() {
    var nodes = document.querySelectorAll('.tree-node');
    for (var i = 0; i < nodes.length; i++) {
        var p = nodes[i].getAttribute('data-path');
        if (p === curPath) {
            nodes[i].classList.add('active');
        } else {
            nodes[i].classList.remove('active');
        }
    }
}

function loadTree() {
    fetch('/api/tree').then(function(r) { return r.json(); }).then(function(data) {
        treeData = data;
        var container = document.getElementById('fileTree');
        container.innerHTML = '';
        renderTreeNode(container, data, 0);
    }).catch(function() {
        document.getElementById('fileTree').innerHTML = '<div class="tree-loading">~</div>';
    });
}

function renderTreeNode(parent, node, depth) {
    var row = document.createElement('div');
    row.className = 'tree-node' + (node.path === curPath ? ' active' : '');
    row.style.paddingLeft = (8 + depth * 14) + 'px';
    row.setAttribute('data-path', node.path);

    var arrow = document.createElement('span');
    arrow.className = 'tree-arrow';
    var hasChildren = node.children && node.children.length > 0;
    if (hasChildren) {
        arrow.textContent = '\u25B6';
        var isOpen = curPath === node.path || curPath.indexOf(node.path + '/') === 0;
        if (isOpen) arrow.classList.add('open');
    }

    var icon = document.createElement('span');
    icon.className = 'tree-icon';
    icon.textContent = '\uD83D\uDCC1';

    var label = document.createElement('span');
    label.className = 'tree-label';
    label.textContent = node.name || '/';

    row.appendChild(arrow);
    row.appendChild(icon);
    row.appendChild(label);
    parent.appendChild(row);

    if (hasChildren) {
        var childContainer = document.createElement('div');
        childContainer.className = 'tree-children';
        var isOpen = curPath === node.path || curPath.indexOf(node.path + '/') === 0;
        if (isOpen) childContainer.classList.add('open');

        for (var i = 0; i < node.children.length; i++) {
            renderTreeNode(childContainer, node.children[i], depth + 1);
        }
        parent.appendChild(childContainer);

        (function(a, cc, p) {
            row.onclick = function(e) {
                e.stopPropagation();
                var isExpanded = cc.classList.contains('open');
                if (isExpanded) {
                    a.classList.remove('open');
                    cc.classList.remove('open');
                } else {
                    a.classList.add('open');
                    cc.classList.add('open');
                }
                navigateTo(p);
            };
        })(arrow, childContainer, node.path);
    } else {
        row.onclick = function(e) {
            e.stopPropagation();
            navigateTo(node.path);
        };
    }
}

function setLang(lang) {
    curLang = lang;
    localStorage.setItem('upftp-lang', lang);
    document.getElementById('langEN').className = 'lang-btn' + (lang==='en'?' active':'');
    document.getElementById('langZH').className = 'lang-btn' + (lang==='zh'?' active':'');
    document.getElementById('sLangEN').className = 'lang-btn' + (lang==='en'?' active':'');
    document.getElementById('sLangZH').className = 'lang-btn' + (lang==='zh'?' active':'');
    document.querySelectorAll('[data-i18n]').forEach(function(el) {
        var k = el.getAttribute('data-i18n');
        if (LANG[lang][k]) el.textContent = LANG[lang][k];
    });
    document.querySelectorAll('[data-i18n-placeholder]').forEach(function(el) {
        var k = el.getAttribute('data-i18n-placeholder');
        if (LANG[lang][k]) el.placeholder = LANG[lang][k];
    });
    updateBreadcrumb();
    render(allFiles);
    fetch('/api/settings/language', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({language: lang})
    });
}

function updateBreadcrumb() {
    var el = document.getElementById('breadcrumb');
    var titleEl = document.getElementById('pageTitle');
    if (curPath === '/') {
        el.innerHTML = '<a href="javascript:navigateTo(\'/\')">/' + T('root') + '</a>';
        titleEl.textContent = T('all_files');
    } else {
        var parts = curPath.split('/').filter(Boolean);
        var html = '<a href="javascript:navigateTo(\'/\')">/</a>';
        var acc = '';
        for (var i = 0; i < parts.length; i++) {
            acc += '/' + parts[i];
            if (i < parts.length - 1) {
                html += ' <a href="javascript:navigateTo(\'' + acc + '\')">' + parts[i] + '</a> / ';
            } else {
                html += ' ' + parts[i];
            }
        }
        el.innerHTML = html;
        titleEl.textContent = parts[parts.length - 1];
    }
}

function getIcon(name, isDir) {
    if (isDir) return '&#128193;';
    var ext = (name.split('.').pop() || '').toLowerCase();
    var map = {jpg:'🖼️',jpeg:'🖼️',png:'🖼️',gif:'🖼️',svg:'🖼️',webp:'🖼️',
        mp4:'🎬',avi:'🎬',mov:'🎬',mkv:'🎬',webm:'🎬',
        mp3:'🎵',wav:'🎵',flac:'🎵',ogg:'🎵',m4a:'🎵',
        pdf:'📄',doc:'📄',docx:'📄',xls:'📄',xlsx:'📄',ppt:'📄',pptx:'📄',
        zip:'📦',rar:'📦','7z':'📦',tar:'📦',gz:'📦',
        js:'💻',ts:'💻',py:'💻',go:'💻',html:'💻',css:'💻',sh:'💻',java:'💻',c:'💻',cpp:'💻',rs:'💻',
        txt:'📝',md:'📝',json:'📝',yaml:'📝',yml:'📝',xml:'📝',csv:'📝',log:'📝',toml:'📝'};
    return map[ext] || '📄';
}

function render(files) {
    allFiles = files;
    var body = document.getElementById('fileListBody');
    if (!files || files.length === 0) {
        body.innerHTML = '<div class="empty"><p>' + T('empty') + '</p></div>';
        return;
    }
    var h = '';
    for (var i = 0; i < files.length; i++) {
        var f = files[i];
        var icon = getIcon(f.Name, f.IsDir);
        var click = f.IsDir
            ? 'navigateTo(\'' + escHtml(f.Path) + '\')'
            : '';
        h += '<div class="file-row">';
        h += '<div class="file-name-cell">';
        h += '<span class="file-icon">' + icon + '</span>';
        h += '<span class="file-name" onclick="' + click + '">' + escHtml(f.Name) + '</span>';
        h += '</div>';
        h += '<div class="file-size">' + escHtml(f.Size || '') + '</div>';
        h += '<div class="file-date">' + escHtml(f.ModTime || '') + '</div>';
        h += '<div class="file-actions">';
        h += '<button class="act-btn" onclick="copyLink(&#39;' + escHtml(f.Path) + '&#39;)">' + T('copy_link') + '</button>';
        if (f.IsDir) {
            h += '<button class="act-btn" onclick="dlFile(&#39;' + escHtml(f.Path) + '&#39;)">' + T('zip_download') + '</button>';
        } else {
            h += '<button class="act-btn" onclick="dlFile(&#39;' + escHtml(f.Path) + '&#39;)">' + T('download') + '</button>';
        }
        h += '<button class="act-btn" onclick="showQR(&#39;' + escHtml(f.Path) + '&#39;)">' + T('qr_download') + '</button>';
        if (f.CanPreview && !f.IsDir) {
            h += '<button class="act-btn" onclick="showPreview(&#39;' + escHtml(f.Path) + '&#39;,&#39;' + escHtml(f.Name) + '&#39;,' + f.FileType + ')">' + T('preview') + '</button>';
        }
        h += '</div></div>';
    }
    body.innerHTML = h;
}

function escHtml(s) {
    return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;').replace(/'/g,'&#39;');
}

function copyLink(path) {
    var url = window.location.origin + '/download/' + path;
    if (httpAuthOn && httpAuthUser && httpAuthPass) {
        url = window.location.origin.replace('://', '://' + encodeURIComponent(httpAuthUser) + ':' + encodeURIComponent(httpAuthPass) + '@') + '/download/' + path;
    }
    if (navigator.clipboard) {
        navigator.clipboard.writeText(url).then(function(){ toast(T('copied')); });
    } else {
        var ta = document.createElement('textarea'); ta.value = url;
        document.body.appendChild(ta); ta.select(); document.execCommand('copy');
        document.body.removeChild(ta); toast(T('copied'));
    }
}

function dlFile(path) {
    var a = document.createElement('a');
    a.href = '/download/' + path; a.click();
}

function showQR(path) {
    var url = window.location.origin + '/download/' + path;
    document.getElementById('qrTitle').textContent = T('qr_title');
    document.getElementById('qrImage').src = '/api/qrcode?url=' + encodeURIComponent(url);
    document.getElementById('qrLink').textContent = url;
    document.getElementById('qrModal').classList.add('show');
}

function closeQR() {
    document.getElementById('qrModal').classList.remove('show');
}

function showPreview(path, name, fileType) {
    var modal = document.getElementById('previewModal');
    var body = document.getElementById('previewBody');
    var title = document.getElementById('previewTitle');
    title.textContent = name;
    body.innerHTML = '<div class="loading">' + T('loading') + '</div>';
    modal.classList.add('show');

    if (fileType === 1) {
        body.innerHTML = '<img src="' + path + '" alt="" onerror="previewErr()">';
    } else if (fileType === 4) {
        body.innerHTML = '<video controls src="' + path + '" onerror="previewErr()"></video>';
    } else if (fileType === 5) {
        body.innerHTML = '<audio controls src="' + path + '" onerror="previewErr()" style="width:100%"></audio>';
    } else {
        fetch(path).then(function(r) {
            if (!r.ok) throw new Error('fetch failed');
            return r.text();
        }).then(function(text) {
            body.innerHTML = '<pre>' + escHtml(text) + '</pre>';
        }).catch(function() {
            previewErr();
        });
    }
}

function previewErr() {
    document.getElementById('previewBody').innerHTML = '<div class="error">' + T('preview_err') + '</div>';
}

function closePreview() {
    document.getElementById('previewModal').classList.remove('show');
    document.getElementById('previewBody').innerHTML = '';
}

function doUpload() {
    var input = document.createElement('input');
    input.type = 'file'; input.multiple = true;
    input.onchange = function(e) {
        var fd = new FormData();
        fd.append('path', curPath);
        for (var i = 0; i < e.target.files.length; i++) fd.append('files', e.target.files[i]);
        fetch('/api/upload', {method:'POST', body:fd}).then(function(r) {
            if (r.ok) navigateTo(curPath);
        });
    };
    input.click();
}

function onSearch(q) {
    var rows = document.querySelectorAll('.file-row');
    var ql = q.toLowerCase();
    for (var i = 0; i < rows.length; i++) {
        var name = rows[i].querySelector('.file-name').textContent.toLowerCase();
        rows[i].style.display = name.indexOf(ql) >= 0 ? 'grid' : 'none';
    }
}

function toast(msg) {
    var el = document.getElementById('toast');
    el.textContent = msg; el.classList.add('show');
    setTimeout(function(){ el.classList.remove('show'); }, 2500);
}

function openSettings() {
    document.getElementById('httpAuthUser').value = httpAuthUser;
    document.getElementById('httpAuthPass').value = httpAuthPass;
    var toggle = document.getElementById('httpAuthToggle');
    var fields = document.getElementById('httpAuthFields');
    if (httpAuthOn) {
        toggle.classList.add('on');
        fields.style.display = 'block';
    } else {
        toggle.classList.remove('on');
        fields.style.display = 'none';
    }
    fetch('/api/settings').then(function(r){ return r.json(); }).then(function(s) {
        document.getElementById('ftpUser').value = s.ftpUser || '';
        document.getElementById('ftpPass').value = s.ftpPass || '';
    });
    document.getElementById('settingsModal').classList.add('show');
}

function closeSettings() {
    document.getElementById('settingsModal').classList.remove('show');
}

function toggleHttpAuth() {
    var toggle = document.getElementById('httpAuthToggle');
    var fields = document.getElementById('httpAuthFields');
    toggle.classList.toggle('on');
    if (toggle.classList.contains('on')) {
        fields.style.display = 'block';
    } else {
        fields.style.display = 'none';
    }
}

function saveHttpAuth() {
    var toggle = document.getElementById('httpAuthToggle');
    var enabled = toggle.classList.contains('on');
    var user = document.getElementById('httpAuthUser').value;
    var pass = document.getElementById('httpAuthPass').value;
    var statusEl = document.getElementById('httpAuthStatus');
    statusEl.textContent = '';
    fetch('/api/settings/http-auth', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({enabled: enabled, username: user, password: pass})
    }).then(function(r) { return r.json(); }).then(function(d) {
        if (d.success) {
            httpAuthOn = enabled;
            httpAuthUser = user;
            httpAuthPass = pass;
            statusEl.textContent = T('save_success');
            statusEl.className = 'settings-status success';
        } else {
            statusEl.textContent = T('save_error');
            statusEl.className = 'settings-status error';
        }
    }).catch(function() {
        statusEl.textContent = T('save_error');
        statusEl.className = 'settings-status error';
    });
}

function saveFTP() {
    var user = document.getElementById('ftpUser').value;
    var pass = document.getElementById('ftpPass').value;
    var statusEl = document.getElementById('ftpStatus');
    statusEl.textContent = '';
    fetch('/api/settings/ftp', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({username: user, password: pass})
    }).then(function(r) { return r.json(); }).then(function(d) {
        if (d.success) {
            statusEl.textContent = T('save_success');
            statusEl.className = 'settings-status success';
        } else {
            statusEl.textContent = T('save_error');
            statusEl.className = 'settings-status error';
        }
    }).catch(function() {
        statusEl.textContent = T('save_error');
        statusEl.className = 'settings-status error';
    });
}

document.addEventListener('DOMContentLoaded', function() {
    allFiles = window.__UPFTP_CONFIG__.files;
    curPath = window.__UPFTP_CONFIG__.currentPath;
    setLang(curLang);
    updateBreadcrumb();
    render(allFiles);
    loadTree();
});

window.addEventListener('popstate', function(e) {
    if (e.state && e.state.path) {
        navigateTo(e.state.path);
    }
});
