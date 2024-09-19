function logSystem() {
    return {
        logs: [],
        stats: [],
        showModal: false,
        showStatModal: false,
        selectedLog: {},
        selectedStat: {},
        page: 1,
        limit: 30,
        total: 0,
        searchQuery: '',
        searchTimeout: null,
        selectedDate: '',
        showConfirmModal: false,
        currentView: 'logs',
        statsPage: 1,
        statsLimit: 10,
        statsTotal: 0,

        init() {
            this.fetchLogs();
            this.fetchStats();
        },

        async fetchLogs() {
            const response = await fetch(`/api/logs?page=${this.page}&limit=${this.limit}&search=${this.searchQuery}`);
            const data = await response.json();
            this.logs = data.logs;
            this.total = data.total;
            this.page = data.page;
            this.limit = data.limit;
        },

        async fetchStats() {
            const response = await fetch(`/api/stats?page=${this.statsPage}&pageSize=${this.statsLimit}`);
            const data = await response.json();
            this.stats = data.stats;
            this.statsTotal = data.total;
        },

        showLogDetails(log) {
            this.selectedLog = log;
            this.showModal = true;
        },

        showStatDetails(stat) {
            this.selectedStat = stat;
            this.showStatModal = true;
            this.renderChart();
        },

        changePage(newPage) {
            if (newPage >= 1 && newPage <= Math.ceil(this.total / this.limit)) {
                this.page = newPage;
                this.fetchLogs();
            }
        },

        changeStatsPage(newPage) {
            if (newPage >= 1 && newPage <= Math.ceil(this.statsTotal / this.statsLimit)) {
                this.statsPage = newPage;
                this.fetchStats();
            }
        },

        debounceSearch() {
            clearTimeout(this.searchTimeout);
            this.searchTimeout = setTimeout(() => {
                this.page = 1;
                this.fetchLogs();
            }, 300);
        },

        deleteLogsBefore() {
            if (!this.selectedDate) {
                alert('Please select a date first.');
                return;
            }
            this.showConfirmModal = true;
        },

        async confirmDelete() {
            this.showConfirmModal = false;
            try {
                const response = await fetch(`/api/logs?date=${this.selectedDate}`, {
                    method: 'DELETE'
                });
                const result = await response.json();
                if (response.ok) {
                    alert(`${result.count} logs deleted successfully.`);
                    this.fetchLogs();
                } else {
                    throw new Error(result.error);
                }
            } catch (error) {
                alert('Failed to delete logs: ' + error.message);
            }
        },

        formatStack(stack) {
            return stack.replace(/\n/g, '<br>').replace(/ /g, '&nbsp;');
        },

        viewLogs() {
            this.currentView = 'logs';
        },

        viewStats() {
            this.currentView = 'stats';
        },

        renderChart() {
            const ctx = document.getElementById('performanceChart').getContext('2d');
            new Chart(ctx, {
                type: 'line',
                data: {
                    labels: this.selectedStat.timestamps,
                    datasets: [
                        {
                            label: 'FPS',
                            data: this.selectedStat.fps,
                            borderColor: 'rgba(75, 192, 192, 1)',
                            borderWidth: 1,
                            fill: false
                        },
                        {
                            label: 'Total Memory',
                            data: this.selectedStat.total_mem,
                            borderColor: 'rgba(255, 99, 132, 1)',
                            borderWidth: 1,
                            fill: false
                        },
                        {
                            label: 'Used Memory',
                            data: this.selectedStat.used_mem,
                            borderColor: 'rgba(54, 162, 235, 1)',
                            borderWidth: 1,
                            fill: false
                        },
                        {
                            label: 'Mono Used Memory',
                            data: this.selectedStat.mono_used_mem,
                            borderColor: 'rgba(255, 206, 86, 1)',
                            borderWidth: 1,
                            fill: false
                        },
                        {
                            label: 'Mono Stack Memory',
                            data: this.selectedStat.mono_stack_mem,
                            borderColor: 'rgba(75, 192, 192, 1)',
                            borderWidth: 1,
                            fill: false
                        },
                        {
                            label: 'Texture',
                            data: this.selectedStat.texture,
                            borderColor: 'rgba(153, 102, 255, 1)',
                            borderWidth: 1,
                            fill: false
                        },
                        {
                            label: 'Audio',
                            data: this.selectedStat.audio,
                            borderColor: 'rgba(255, 159, 64, 1)',
                            borderWidth: 1,
                            fill: false
                        },
                        {
                            label: 'Text Asset',
                            data: this.selectedStat.text_asset,
                            borderColor: 'rgba(255, 99, 132, 1)',
                            borderWidth: 1,
                            fill: false
                        },
                        {
                            label: 'Shader',
                            data: this.selectedStat.shader,
                            borderColor: 'rgba(54, 162, 235, 1)',
                            borderWidth: 1,
                            fill: false
                        }
                    ]
                },
                options: {
                    responsive: true,
                    scales: {
                        x: {
                            type: 'time',
                            time: {
                                unit: 'minute'
                            }
                        }
                    }
                }
            });
        },

        async deleteStatsBefore() {
            if (!this.selectedDate) {
                alert('Please select a date first.');
                return;
            }
            try {
                const response = await fetch(`/api/stats?date=${this.selectedDate}`, {
                    method: 'DELETE'
                });
                const result = await response.json();
                if (response.ok) {
                    alert(`${result.count} stats deleted successfully.`);
                    this.fetchStats();
                } else {
                    throw new Error(result.error);
                }
            } catch (error) {
                alert('Failed to delete stats: ' + error.message);
            }
        }
    }
}
 