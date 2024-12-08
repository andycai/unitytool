// Stats management functionality
function statsManagement() {
    return {
        stats: [],
        showModal: false,
        currentStat: null,
        memoryChart: null,
        fpsChart: null,

        init() {
            this.fetchStats();
            this.initCharts();
        },

        async fetchStats() {
            try {
                const response = await fetch('/api/stats');
                if (!response.ok) throw new Error('获取统计列表失败');
                this.stats = await response.json();
                this.updateCharts();
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        initCharts() {
            // 初始化内存使用趋势图
            const memoryCtx = document.getElementById('memoryChart').getContext('2d');
            this.memoryChart = new Chart(memoryCtx, {
                type: 'line',
                data: {
                    labels: [],
                    datasets: [
                        {
                            label: '总内存',
                            data: [],
                            borderColor: 'rgb(59, 130, 246)',
                            tension: 0.1
                        },
                        {
                            label: '已用内存',
                            data: [],
                            borderColor: 'rgb(239, 68, 68)',
                            tension: 0.1
                        },
                        {
                            label: 'Mono内存',
                            data: [],
                            borderColor: 'rgb(16, 185, 129)',
                            tension: 0.1
                        }
                    ]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    interaction: {
                        mode: 'index',
                        intersect: false,
                    },
                    plugins: {
                        zoom: {
                            zoom: {
                                wheel: {
                                    enabled: true,
                                },
                                pinch: {
                                    enabled: true
                                },
                                mode: 'x',
                            },
                            pan: {
                                enabled: true,
                                mode: 'x',
                            }
                        }
                    },
                    scales: {
                        y: {
                            beginAtZero: true,
                            title: {
                                display: true,
                                text: '内存 (MB)'
                            }
                        },
                        x: {
                            type: 'time',
                            time: {
                                unit: 'minute'
                            },
                            title: {
                                display: true,
                                text: '时间'
                            }
                        }
                    }
                }
            });

            // 初始化FPS趋势图
            const fpsCtx = document.getElementById('fpsChart').getContext('2d');
            this.fpsChart = new Chart(fpsCtx, {
                type: 'line',
                data: {
                    labels: [],
                    datasets: [{
                        label: 'FPS',
                        data: [],
                        borderColor: 'rgb(99, 102, 241)',
                        tension: 0.1
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    interaction: {
                        mode: 'index',
                        intersect: false,
                    },
                    plugins: {
                        zoom: {
                            zoom: {
                                wheel: {
                                    enabled: true,
                                },
                                pinch: {
                                    enabled: true
                                },
                                mode: 'x',
                            },
                            pan: {
                                enabled: true,
                                mode: 'x',
                            }
                        }
                    },
                    scales: {
                        y: {
                            beginAtZero: true,
                            title: {
                                display: true,
                                text: 'FPS'
                            }
                        },
                        x: {
                            type: 'time',
                            time: {
                                unit: 'minute'
                            },
                            title: {
                                display: true,
                                text: '时间'
                            }
                        }
                    }
                }
            });
        },

        updateCharts() {
            if (!this.stats.length) return;

            const timestamps = this.stats.map(stat => new Date(stat.mtime * 1000));
            const totalMem = this.stats.map(stat => stat.system_mem / (1024 * 1024));
            const usedMem = this.stats.map(stat => stat.used_mem / (1024 * 1024));
            const monoMem = this.stats.map(stat => stat.mono_used_mem / (1024 * 1024));
            const fps = this.stats.map(stat => stat.fps);

            // 更新内存图表
            this.memoryChart.data.labels = timestamps;
            this.memoryChart.data.datasets[0].data = totalMem;
            this.memoryChart.data.datasets[1].data = usedMem;
            this.memoryChart.data.datasets[2].data = monoMem;
            this.memoryChart.update();

            // 更新FPS图表
            this.fpsChart.data.labels = timestamps;
            this.fpsChart.data.datasets[0].data = fps;
            this.fpsChart.update();
        },

        async viewDetails(stat) {
            try {
                const response = await fetch(`/api/stats/details?login_id=${stat.login_id}&mtime=${stat.mtime}`);
                if (!response.ok) throw new Error('获取统计详情失败');
                this.currentStat = await response.json();
                this.showModal = true;
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        closeModal() {
            this.showModal = false;
            this.currentStat = null;
        },

        async deleteStat(id) {
            if (!confirm('确定要删除这条统计记录吗？')) return;

            try {
                const response = await fetch(`/api/stats/${id}`, {
                    method: 'DELETE'
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '删除失败');
                }

                Alpine.store('notification').show('统计记录删除成功', 'success');
                this.fetchStats();
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async clearOldData() {
            const date = new Date();
            date.setDate(date.getDate() - 7); // 默认清理7天前的数据
            
            if (!confirm('确定要清理7天前的数据吗？此操作不可恢复。')) return;

            try {
                const response = await fetch(`/api/stats/before?date=${Math.floor(date.getTime() / 1000)}`, {
                    method: 'DELETE'
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '清理失败');
                }

                Alpine.store('notification').show('旧数据清理成功', 'success');
                this.fetchStats();
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        refreshData() {
            this.fetchStats();
        },

        formatDate(timestamp) {
            if (!timestamp) return '';
            return new Date(timestamp * 1000).toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
                hour12: false
            });
        },

        formatMemory(bytes) {
            if (!bytes) return '0 MB';
            const mb = bytes / (1024 * 1024);
            return mb.toFixed(2) + ' MB';
        }
    }
} 