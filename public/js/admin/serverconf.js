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
            pfid: 0,
            pfname: '',
            child: 0,
            adKey: '',
            entryURL: '',
            cdnURL: '',
            cdnVersion: '',
            loginAPI: '',
            loginURL: '',
            serverListURL: '',
            version: '',
            time: '',
            serverZoneURL: '',
            lastServerListURL: '',
            noticeNumURL: '',
            noticeURL: '',
            pkgVersion: '',
            doPatcher: 0
        },
        noticeList: [],
        noticeNum: {
            noticenum: 0,
            eject: 0
        },

        init() {
            this.loadServerList();
            this.loadLastServer();
            this.loadServerInfo();
            this.loadNoticeList();
            this.loadNoticeNum();
        },

        async loadServerList() {
            try {
                const response = await fetch('/open/game/serverlist');
                if (!response.ok) throw new Error('加载失败');
                this.serverList = await response.json();
            } catch (error) {
                console.error('加载服务器列表失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async loadLastServer() {
            try {
                const response = await fetch('/open/game/lastserver');
                if (!response.ok) throw new Error('加载失败');
                this.lastServer = await response.json();
            } catch (error) {
                console.error('加载最后登录服务器失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async loadServerInfo() {
            try {
                const response = await fetch('/open/game/serverinfo');
                if (!response.ok) throw new Error('加载失败');
                this.serverInfo = await response.json();
            } catch (error) {
                console.error('加载服务器信息失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async loadNoticeList() {
            try {
                const response = await fetch('/open/game/noticelist');
                if (!response.ok) throw new Error('加载失败');
                this.noticeList = await response.json();
            } catch (error) {
                console.error('加载公告列表失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async loadNoticeNum() {
            try {
                const response = await fetch('/open/game/noticenum');
                if (!response.ok) throw new Error('加载失败');
                this.noticeNum = await response.json();
            } catch (error) {
                console.error('加载公告数量失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        addServer() {
            this.serverList.serverlist.push({
                server_id: '',
                name: '',
                server_status: '',
                available: '1',
                mergeid: '0',
                online: String(Math.floor(Date.now() / 1000)),
                server_port: '',
                server_ip: ''
            });
        },

        removeServer(index) {
            this.serverList.serverlist.splice(index, 1);
        },

        addLastServerItem() {
            this.lastServer.lastserver.last_server.push({
                server_id: '',
                name: '',
                server_status: '',
                server_port: '',
                server_ip: ''
            });
        },

        removeLastServerItem(index) {
            this.lastServer.lastserver.last_server.splice(index, 1);
        },

        addNotice() {
            this.noticeList.push({
                title: '',
                content: ''
            });
        },

        removeNotice(index) {
            this.noticeList.splice(index, 1);
        },

        async saveServerList() {
            try {
                const response = await fetch('/api/game/serverlist', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.serverList),
                });
                if (!response.ok) throw new Error('保存失败');
                Alpine.store('notification').show('服务器列表保存成功', 'success');
            } catch (error) {
                console.error('保存服务器列表失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async saveLastServer() {
            try {
                const response = await fetch('/api/game/lastserver', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.lastServer),
                });
                if (!response.ok) throw new Error('保存失败');
                Alpine.store('notification').show('最后登录服务器保存成功', 'success');
            } catch (error) {
                console.error('保存最后登录服务器失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async saveServerInfo() {
            try {
                const response = await fetch('/api/game/serverinfo', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.serverInfo),
                });
                if (!response.ok) throw new Error('保存失败');
                Alpine.store('notification').show('服务器信息保存成功', 'success');
            } catch (error) {
                console.error('保存服务器信息失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async saveNoticeList() {
            try {
                const response = await fetch('/api/game/noticelist', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.noticeList),
                });
                if (!response.ok) throw new Error('保存失败');
                Alpine.store('notification').show('公告列表保存成功', 'success');
            } catch (error) {
                console.error('保存公告列表失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async saveNoticeNum() {
            try {
                const response = await fetch('/api/game/noticenum', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.noticeNum),
                });
                if (!response.ok) throw new Error('保存失败');
                Alpine.store('notification').show('公告数量保存成功', 'success');
            } catch (error) {
                console.error('保存公告数量失败:', error);
                Alpine.store('notification').show(error.message, 'error');
            }
        }
    }
}