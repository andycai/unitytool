function taskManagement() {
    return {
        tasks: [],
        taskLogs: [],
        currentTask: null,
        currentTaskLog: null,
        showTaskModal: false,
        showLogsModal: false,
        showLogDetailModal: false,
        showProgressModal: false,
        showRunningTasksModal: false,
        editMode: false,
        progressInterval: null,
        runningTasks: [],
        runningTasksInterval: null,
        currentPage: 1,
        pageSize: 10,
        totalPages: 1,
        form: {
            name: '',
            description: '',
            type: 'script',
            script: '',
            url: '',
            method: 'GET',
            headers: '',
            body: '',
            timeout: 300,
            status: 'active'
        },
        init() {
            this.fetchTasks();
            this.startRunningTasksPolling();
        },
        async fetchTasks() {
            try {
                const response = await fetch('/api/citask');
                if (!response.ok) throw new Error('获取任务列表失败');
                this.tasks = await response.json();
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },
        createTask() {
            this.editMode = false;
            this.form = {
                name: '',
                description: '',
                type: 'script',
                script: '',
                url: '',
                method: 'GET',
                headers: '',
                body: '',
                timeout: 300,
                status: 'active'
            };
            this.showTaskModal = true;
        },
        editTask(task) {
            this.editMode = true;
            this.form = {
                ...task
            };
            this.showTaskModal = true;
        },
        async submitTask() {
            try {
                const url = this.editMode ? `/api/citask/${this.form.id}` : '/api/citask';
                const method = this.editMode ? 'PUT' : 'POST';
                
                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.form)
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '操作失败');
                }

                await this.fetchTasks();
                this.showTaskModal = false;
                Alpine.store('notification').show(this.editMode ? '任务更新成功' : '任务创建成功', 'success');
            } catch (error) {
                console.error('保存任务失败:', error);
                Alpine.store('notification').show('保存任务失败', 'error');
            }
        },
        async deleteTask(id) {
            if (!confirm('确定要删除这个任务吗？')) return;

            try {
                const response = await fetch(`/api/citask/${id}`, {
                    method: 'DELETE',
                });

                if (!response.ok) throw new Error('删除任务失败');

                await this.fetchTasks();
                Alpine.store('notification').show('任务删除成功', 'success');
            } catch (error) {
                console.error('删除任务失败:', error);
                Alpine.store('notification').show('删除任务失败', 'error');
            }
        },
        async runTask(task) {
            try {
                const response = await fetch(`/api/citask/run/${task.id}`, {
                    method: 'POST'
                });
                if (!response.ok) throw new Error('启动任务失败');
                const taskLog = await response.json();
                
                // 显示进度模态框
                this.currentTask = task;
                this.currentTaskLog = taskLog;
                this.showProgressModal = true;
                
                // 开始轮询进度
                this.startProgressPolling(taskLog.id);
                
                // 自动滚动到底部
                this.$nextTick(() => {
                    this.scrollOutputToBottom();
                });

                Alpine.store('notification').show('任务已开始执行', 'success');
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },
        async stopTask() {
            if (!this.currentTaskLog) return;
            
            try {
                const response = await fetch(`/api/citask/stop/${this.currentTaskLog.id}`, {
                    method: 'POST'
                });
                if (!response.ok) throw new Error('停止任务失败');
                
                // 停止进度轮询
                this.stopProgressPolling();
                
                Alpine.store('notification').show('任务已停止', 'success');
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },
        async viewLogs(task) {
            try {
                const response = await fetch(`/api/citask/logs/${task.id}`);
                if (!response.ok) throw new Error('获取任务日志失败');
                this.taskLogs = await response.json();
                this.currentPage = 1;
                this.updatePaginatedLogs();
                this.showLogsModal = true;
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },
        viewLogDetail(log) {
            this.currentTaskLog = log;
            this.showLogDetailModal = true;
        },
        startProgressPolling(logId) {
            // 清除现有的轮询
            this.stopProgressPolling();
            
            // 开始新的轮询
            this.progressInterval = setInterval(async () => {
                try {
                    const response = await fetch(`/api/citask/progress/${logId}`);
                    if (!response.ok) throw new Error('获取进度失败');
                    const progress = await response.json();
                    
                    // 更新进度信息
                    this.currentTaskLog = progress;
                    
                    // 自动滚动到底部
                    this.$nextTick(() => {
                        this.scrollOutputToBottom();
                    });
                    
                    // 如果任务已结束，停止轮询
                    if (progress.status !== 'running') {
                        this.stopProgressPolling();
                    }
                } catch (error) {
                    console.error('Progress polling error:', error);
                    this.stopProgressPolling();
                }
            }, 1000); // 每秒轮询一次
        },
        stopProgressPolling() {
            if (this.progressInterval) {
                clearInterval(this.progressInterval);
                this.progressInterval = null;
            }
        },
        closeProgress() {
            this.stopProgressPolling();
            this.showProgressModal = false;
            this.currentTask = null;
            this.currentTaskLog = null;
            this.fetchTasks(); // 刷新任务列表
        },
        getProgressWidth() {
            if (!this.currentTaskLog) return '0%';
            if (this.currentTaskLog.status === 'success') return '100%';
            if (this.currentTaskLog.status === 'failed') return '100%';
            return this.currentTaskLog.progress + '%';
        },
        getProgressText() {
            if (!this.currentTaskLog) return '准备中...';
            if (this.currentTaskLog.status === 'success') return '完成';
            if (this.currentTaskLog.status === 'failed') return '失败';
            if (this.currentTaskLog.status === 'running') return '执行中...';
            return '准备中...';
        },
        getRunningTime() {
            if (!this.currentTaskLog?.start_time) return '0秒';
            const start = new Date(this.currentTaskLog.start_time);
            const end = this.currentTaskLog.end_time > this.currentTaskLog.start_time ? new Date(this.currentTaskLog.end_time) : new Date();
            const seconds = Math.floor((end - start) / 1000);
            
            if (seconds < 60) return `${seconds}秒`;
            if (seconds < 3600) return `${Math.floor(seconds / 60)}分${seconds % 60}秒`;
            return `${Math.floor(seconds / 3600)}时${Math.floor((seconds % 3600) / 60)}分${seconds % 60}秒`;
        },
        scrollOutputToBottom() {
            const outputLog = this.$refs.outputLog;
            if (outputLog) {
                outputLog.scrollTop = outputLog.scrollHeight;
            }
        },
        formatDate(date) {
            if (!date) return '';
            return new Date(date).toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
                hour12: false
            });
        },
        formatDateTime(timestamp) {
            if (!timestamp) return '';
            const date = new Date(timestamp);
            const year = date.getFullYear();
            const month = String(date.getMonth() + 1).padStart(2, '0');
            const day = String(date.getDate()).padStart(2, '0');
            const hours = String(date.getHours()).padStart(2, '0');
            const minutes = String(date.getMinutes()).padStart(2, '0');
            const seconds = String(date.getSeconds()).padStart(2, '0');
            return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
        },
        firstPage() {
            if (this.currentPage !== 1) {
                this.currentPage = 1;
                this.updatePaginatedLogs();
            }
        },
        prevPage() {
            if (this.currentPage > 1) {
                this.currentPage--;
                this.updatePaginatedLogs();
            }
        },
        nextPage() {
            if (this.currentPage < this.totalPages) {
                this.currentPage++;
                this.updatePaginatedLogs();
            }
        },
        updatePaginatedLogs() {
            this.totalPages = Math.ceil(this.taskLogs.length / this.pageSize);
        },
        get paginatedLogs() {
            const startIndex = (this.currentPage - 1) * this.pageSize;
            const endIndex = startIndex + this.pageSize;
            return this.taskLogs.slice(startIndex, endIndex);
        },
        // 开始轮询正在执行的任务
        startRunningTasksPolling() {
            this.fetchRunningTasks();
            // this.runningTasksInterval = setInterval(() => {
            //     this.fetchRunningTasks();
            // }, 5000); // 每5秒更新一次
        },
        // 停止轮询
        stopRunningTasksPolling() {
            if (this.runningTasksInterval) {
                clearInterval(this.runningTasksInterval);
                this.runningTasksInterval = null;
            }
        },
        // 获取正在执行的任务
        async fetchRunningTasks() {
            try {
                const response = await fetch('/api/citask/running');
                if (!response.ok) throw new Error('获取正在执行的任务失败');
                const tasks = await response.json();
                const data = tasks.data == null ? [] : tasks.data;
                this.runningTasks = data;
            } catch (error) {
                console.error('获取正在执行的任务失败:', error);
            }
        },
        // 显示正在执行的任务列表
        showRunningTasks() {
            this.showRunningTasksModal = true;
            this.fetchRunningTasks();
        },
        // 查看任务进度
        viewTaskProgress(task) {
            this.currentTask = task;
            this.showRunningTasksModal = false;
            this.showProgressModal = true;
            this.startProgressPolling(task.id);
        },
        // 在组件销毁时清理
        destroy() {
            this.stopRunningTasksPolling();
            this.stopProgressPolling();
        },
        // 查看日志详情
        async viewLog(logId) {
            try {
                const response = await fetch(`/api/citask/progress/${logId}`);
                if (!response.ok) throw new Error('获取任务进度失败');
                const progress = await response.json();
                
                // 如果任务正在执行中，显示进度界面
                if (progress.status === 'running') {
                    this.currentTask = { id: logId };
                    this.currentTaskLog = progress;
                    this.showProgressModal = true;
                    this.startProgressPolling(logId);
                } else {
                    // 否则显示结果界面
                    this.currentTaskLog = progress;
                    this.showLogDetailModal = true;
                }
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },
        // 获取状态显示样式
        getStatusBadge(status) {
            const statusClasses = {
                'success': 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
                'failed': 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
                'running': 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
                'pending': 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200'
            };
            
            const statusText = {
                'success': '成功',
                'failed': '失败',
                'running': '执行中',
                'pending': '等待中'
            };

            const classes = statusClasses[status] || statusClasses['pending'];
            const text = statusText[status] || status;

            return `<span class="px-2 py-1 text-xs font-medium rounded-full ${classes}">${text}</span>`;
        },
    }
} 