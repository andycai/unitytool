function deleteFile(path, event) {
    event.preventDefault();
    event.stopPropagation();

    if (!confirm('确定要删除这个文件吗？')) {
        return;
    }

    fetch('/admin/browse/' + encodeURIComponent(path), {
        method: 'DELETE'
    })
        .then(response => {
            if (response.ok) {
                window.location.reload();
            } else {
                alert('删除失败: ' + response.statusText);
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('删除失败: ' + error.message);
        });
}

function uploadToFTP(path, type, event) {
    event.preventDefault();
    event.stopPropagation();

    if (!confirm('确定要上传这个' + type.toUpperCase() + '文件吗？')) {
        return;
    }

    fetch('/admin/ftp/upload?file=' + encodeURIComponent(path) + '&type=' + type, {
        method: 'POST'
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                alert('上传成功');
            } else {
                alert('上传失败: ' + data.error);
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('上传失败: ' + error.message);
        });
}