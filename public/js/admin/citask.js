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
        editMode: false,
        progressInterval: null,
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
                    if (this.$refs.outputLog) {
                        this.$refs.outputLog.scrollTop = this.$refs.outputLog.scrollHeight;
                    }
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
                        if (this.$refs.outputLog) {
                            this.$refs.outputLog.scrollTop = this.$refs.outputLog.scrollHeight;
                        }
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
            const end = this.currentTaskLog.end_time ? new Date(this.currentTaskLog.end_time) : new Date();
            const seconds = Math.floor((end - start) / 1000);
            
            if (seconds < 60) return `${seconds}秒`;
            if (seconds < 3600) return `${Math.floor(seconds / 60)}分${seconds % 60}秒`;
            return `${Math.floor(seconds / 3600)}时${Math.floor((seconds % 3600) / 60)}分${seconds % 60}秒`;
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
        }
    }
} 