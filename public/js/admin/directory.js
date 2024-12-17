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

// 处理文件点击事件
function handleFileClick(filePath) {
    fetch(`/api/browse/file?path=${encodeURIComponent(filePath)}`)
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                if (data.data.isText) {
                    // 显示文本内容
                    showTextContent(data.data.content, filePath);
                } else {
                    // 下载文件
                    window.location.href = data.data.downloadUrl;
                }
            } else {
                showError('获取文件内容失败');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            showError('获取文件内容失败');
        });
}

// 显示文本内容
function showTextContent(content, filePath) {
    const modal = document.createElement('div');
    modal.className = 'fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50';
    
    const modalContent = document.createElement('div');
    modalContent.className = 'bg-white dark:bg-gray-800 rounded-lg shadow-xl w-4/5 h-4/5 flex flex-col';
    
    const header = document.createElement('div');
    header.className = 'flex justify-between items-center p-4 border-b dark:border-gray-700';
    
    const title = document.createElement('h3');
    title.className = 'text-lg font-semibold dark:text-gray-200';
    const fileName = filePath.split('/').pop(); // 获取文件名
    title.textContent = `文件内容: ${fileName}`;
    
    const closeButton = document.createElement('button');
    closeButton.className = 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200';
    closeButton.innerHTML = '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>';
    closeButton.onclick = () => modal.remove();
    
    header.appendChild(title);
    header.appendChild(closeButton);
    
    const contentWrapper = document.createElement('div');
    contentWrapper.className = 'flex-1 p-4 overflow-auto';
    
    const pre = document.createElement('pre');
    pre.className = 'whitespace-pre-wrap break-words text-sm font-mono dark:text-gray-300';
    pre.textContent = content;
    
    contentWrapper.appendChild(pre);
    modalContent.appendChild(header);
    modalContent.appendChild(contentWrapper);
    modal.appendChild(modalContent);
    
    document.body.appendChild(modal);
    
    // 点击模态框外部关闭
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.remove();
        }
    });
    
    // ESC键关闭
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') {
            modal.remove();
        }
    });
}

// 显示错误信息
function showError(message) {
    const toast = document.createElement('div');
    toast.className = 'fixed bottom-4 right-4 bg-red-500 text-white px-6 py-3 rounded-lg shadow-lg z-50';
    toast.textContent = message;
    
    document.body.appendChild(toast);
    
    setTimeout(() => {
        toast.remove();
    }, 3000);
}

// 初始化文件点击事件
document.addEventListener('DOMContentLoaded', function() {
    // 为所有文件链接添加点击事件
    const fileLinks = document.querySelectorAll('[data-file-path]');
    fileLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            const filePath = this.getAttribute('data-file-path');
            handleFileClick(filePath);
        });
    });
});