function taskManagement() {
    return {
        tasks: [],
        taskLogs: [],
        currentLog: null,
        showTaskModal: false,
        showLogsModal: false,
        showLogDetailModal: false,
        editMode: false,
        form: {
            id: null,
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
            this.loadTasks();
        },

        // 加载任务列表
        async loadTasks() {
            try {
                const response = await fetch('/api/citask');
                if (!response.ok) throw new Error('加载任务列表失败');
                this.tasks = await response.json();
            } catch (error) {
                console.error('加载任务列表失败:', error);
                showToast('error', '加载任务列表失败');
            }
        },

        // 创建新任务
        createTask() {
            this.editMode = false;
            this.form = {
                id: null,
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

        // 编辑任务
        editTask(task) {
            this.editMode = true;
            this.form = { ...task };
            this.showTaskModal = true;
        },

        // 提交任务表单
        async submitTask() {
            try {
                const url = this.editMode ? `/api/citask/${this.form.id}` : '/api/citask';
                const method = this.editMode ? 'PUT' : 'POST';
                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.form),
                });

                if (!response.ok) throw new Error('保存任务失败');

                await this.loadTasks();
                this.showTaskModal = false;
                showToast('success', this.editMode ? '任务更新成功' : '任务创建成功');
            } catch (error) {
                console.error('保存任务失败:', error);
                showToast('error', '保存任务失败');
            }
        },

        // 删除任务
        async deleteTask(id) {
            if (!confirm('确定要删除这个任务吗？')) return;

            try {
                const response = await fetch(`/api/citask/${id}`, {
                    method: 'DELETE',
                });

                if (!response.ok) throw new Error('删除任务失败');

                await this.loadTasks();
                showToast('success', '任务删除成功');
            } catch (error) {
                console.error('删除任务失败:', error);
                showToast('error', '删除任务失败');
            }
        },

        // 执行任务
        async runTask(task) {
            try {
                const response = await fetch(`/api/citask/${task.id}/run`, {
                    method: 'POST',
                });

                if (!response.ok) throw new Error('执行任务失败');

                const result = await response.json();
                showToast('success', '任务已开始执行');

                // 如果需要实时更新任务状态，可以启动轮询
                this.pollTaskStatus(result.log_id);
            } catch (error) {
                console.error('执行任务失败:', error);
                showToast('error', '执行任务失败');
            }
        },

        // 轮询任务状态
        async pollTaskStatus(logId) {
            const poll = async () => {
                try {
                    const response = await fetch(`/api/citask/${logId}/status?log_id=${logId}`);
                    if (!response.ok) throw new Error('获取任务状态失败');

                    const log = await response.json();
                    if (log.status !== 'running') {
                        // 任务完成，更新任务列表
                        await this.loadTasks();
                        return;
                    }

                    // 继续轮询
                    setTimeout(poll, 2000);
                } catch (error) {
                    console.error('获取任务状态失败:', error);
                }
            };

            poll();
        },

        // 查看任务日志
        async viewLogs(task) {
            try {
                const response = await fetch(`/api/citask/${task.id}/logs`);
                if (!response.ok) throw new Error('加载任务日志失败');

                this.taskLogs = await response.json();
                this.showLogsModal = true;
            } catch (error) {
                console.error('加载任务日志失败:', error);
                showToast('error', '加载任务日志失败');
            }
        },

        // 查看日志详情
        viewLogDetail(log) {
            this.currentLog = log;
            this.showLogDetailModal = true;
        },

        // 格式化日期
        formatDate(date) {
            if (!date) return '';
            return new Date(date).toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
            });
        }
    };
} 