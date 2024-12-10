function serverConfig() {
    return {
        currentTab: 'serverList',
        serverList: { serverlist: [] },
        lastServer: {
            lastserver: {
                default_server: {
                    server_id: '',
                    name: '',
                    server_status: '',
                    server_port: '',
                    server_ip: ''
                },
                last_server: []
            },
            params: '',
            sdkParams: ''
        },
        serverInfo: {
            fields: []
        },

        init() {
            this.loadServerList();
            this.loadLastServer();
            this.loadServerInfo();
            this.loadNoticeList();
            this.loadNoticeNum();
        },

        // 字段管理方法
        addField() {
            this.serverInfo.fields.push({
                key: '',
                value: '',
                type: 'string'
            });
        },

        removeField(index) {
            this.serverInfo.fields.splice(index, 1);
        },

        // 转换服务器信息为API格式
        prepareServerInfoForAPI() {
            const result = {
                fields: []
            };

            // 处理所有字段
            for (const field of this.serverInfo.fields) {
                if (field && field.key && field.value !== undefined) {
                    let value = field.value;
                    switch (field.type) {
                        case 'number':
                            value = parseFloat(value) || 0;
                            break;
                        case 'boolean':
                            value = value === 'true' || value === true;
                            break;
                        default:
                            value = String(value);
                    }
                    result[field.key] = value;
                    result.fields.push({
                        key: field.key,
                        value: value,
                        type: field.type
                    });
                }
            }

            return result;
        },

        // 从API响应解析服务器信息
        parseServerInfoFromAPI(data) {
            this.serverInfo.fields = [];

            if (data.fields && Array.isArray(data.fields)) {
                this.serverInfo.fields = data.fields.map(field => ({
                    key: field.key,
                    value: field.value !== null ? field.value.toString() : '',
                    type: field.type || 'string'
                }));
            } else {
                for (const [key, value] of Object.entries(data)) {
                    if (key !== 'fields') {
                        let type = 'string';
                        if (typeof value === 'number') {
                            type = 'number';
                        } else if (typeof value === 'boolean') {
                            type = 'boolean';
                        }
                        
                        this.serverInfo.fields.push({
                            key: key,
                            value: value !== null ? value.toString() : '',
                            type: type
                        });
                    }
                }
            }
        },

        async loadServerList() {
            try {
                const response = await fetch('/api/game/serverlist');
                if (!response.ok) throw new Error('加载失败');
                this.serverList = await response.json();
            } catch (error) {
                console.error('加载服务器列表失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async loadLastServer() {
            try {
                const response = await fetch('/api/game/lastserver');
                if (!response.ok) throw new Error('加载失败');
                this.lastServer = await response.json();
            } catch (error) {
                console.error('加载最后登录服务器失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async loadServerInfo() {
            try {
                const response = await fetch('/api/game/serverinfo');
                if (!response.ok) throw new Error('加载失败');
                const data = await response.json();
                this.parseServerInfoFromAPI(data);
            } catch (error) {
                console.error('加载服务器信息失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async loadNoticeList() {
            try {
                const response = await fetch('/api/game/noticelist');
                if (!response.ok) throw new Error('加载失败');
                this.noticeList = await response.json();
            } catch (error) {
                console.error('加载公告列表失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async loadNoticeNum() {
            try {
                const response = await fetch('/api/game/noticenum');
                if (!response.ok) throw new Error('加载失败');
                this.noticeNum = await response.json();
            } catch (error) {
                console.error('加载公告数量失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async saveServerInfo() {
            try {
                const data = this.prepareServerInfoForAPI();
                console.log('Saving server info:', data);

                const response = await fetch('/api/game/serverinfo', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(data),
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '保存失败');
                }

                Alpine.store('notification').show('服务器信息保存成功', 'success');
                
                await this.loadServerInfo();
            } catch (error) {
                console.error('保存服务器信息失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        }
    }
}