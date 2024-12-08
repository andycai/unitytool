// Stats management functionality
function statsManagement() {
    return {
        stats: [],
        showModal: false,
        currentStat: null,
        currentStatIndex: -1,
        memoryChart: null,
        fpsChart: null,
        detailMemoryChart: null,
        detailFpsChart: null,
        textureChart: null,
        meshChart: null,
        animationChart: null,
        audioChart: null,
        fontChart: null,
        textAssetChart: null,
        shaderChart: null,
        chartInstances: {},
        currentPointIndex: {},
        total: 0,
        filters: {
            startDate: '',
            endDate: '',
            search: '',
            memoryThreshold: ''
        },

        init() {
            this.initDateRange();
            this.initCharts();
            this.fetchStats();
            this.watchFilters();
        },

        initDateRange() {
            // 默认显示最近7天
            const end = new Date();
            const start = new Date();
            start.setDate(start.getDate() - 7);
            
            this.filters.startDate = start.toISOString().split('T')[0];
            this.filters.endDate = end.toISOString().split('T')[0];
        },

        initCharts() {
            const chartConfigs = [
                { id: 'memoryChart', label: 'Memory', dataKeys: ['total_mem', 'used_mem', 'mono_used_mem'], colors: ['rgb(59, 130, 246)', 'rgb(239, 68, 68)', 'rgb(16, 185, 129)'] },
                { id: 'fpsChart', label: 'FPS', dataKey: 'fps', color: 'rgb(99, 102, 241)' },
                { id: 'textureChart', label: 'Texture', dataKey: 'texture', color: 'rgb(255, 159, 64)' },
                { id: 'meshChart', label: 'Mesh', dataKey: 'mesh', color: 'rgb(255, 99, 71)' },
                { id: 'animationChart', label: 'Animation', dataKey: 'animation', color: 'rgb(50, 205, 50)' },
                { id: 'audioChart', label: 'Audio', dataKey: 'audio', color: 'rgb(0, 191, 255)' },
                { id: 'fontChart', label: 'Font', dataKey: 'font', color: 'rgb(255, 140, 0)' },
                { id: 'textAssetChart', label: 'Text Asset', dataKey: 'text_asset', color: 'rgb(186, 85, 211)' },
                { id: 'shaderChart', label: 'Shader', dataKey: 'shader', color: 'rgb(0, 128, 128)' }
            ];

            chartConfigs.forEach(config => {
                const canvas = document.getElementById(config.id);
                if (!canvas) return;

                const ctx = canvas.getContext('2d');
                const options = {
                    responsive: true,
                    maintainAspectRatio: false,
                    interaction: {
                        mode: 'index',
                        intersect: false,
                    },
                    plugins: {
                        zoom: {
                            zoom: {
                                wheel: { enabled: true },
                                pinch: { enabled: true },
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
                                text: config.label
                            }
                        },
                        x: {
                            type: 'time',
                            time: {
                                unit: 'minute',
                                displayFormats: {
                                    minute: 'HH:mm'
                                }
                            },
                            title: {
                                display: true,
                                text: '时间'
                            }
                        }
                    },
                    onClick: (evt, elements) => {
                        if (elements.length > 0) {
                            const index = elements[0].index;
                            const datasetIndex = elements[0].datasetIndex;
                            const dataset = evt.chart.data.datasets[datasetIndex];
                            const pic = dataset.pic[index];
                            const process = dataset.process[index];
                            this.updateClickedPointInfo(pic, process);
                            
                            this.currentPointIndex[config.id] = index;
                            
                            evt.chart.setActiveElements([{
                                datasetIndex: datasetIndex,
                                index: index
                            }]);
                            evt.chart.update();
                        }
                    }
                };

                if (config.id === 'memoryChart' || config.id === 'detailMemoryChart') {
                    const datasets = config.dataKeys.map((key, index) => ({
                        label: key.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' '),
                        data: [],
                        borderColor: config.colors[index],
                        tension: 0.1,
                        pointStyle: 'circle',
                        pointRadius: 5,
                        pointHoverRadius: 7,
                        pic: [],
                        process: []
                    }));

                    this.chartInstances[config.id] = new Chart(ctx, {
                        type: 'line',
                        data: { datasets },
                        options
                    });
                } else {
                    this.chartInstances[config.id] = new Chart(ctx, {
                        type: 'line',
                        data: {
                            datasets: [{
                                label: config.label,
                                data: [],
                                borderColor: config.color,
                                tension: 0.1,
                                pointStyle: 'circle',
                                pointRadius: 5,
                                pointHoverRadius: 7,
                                pic: [],
                                process: []
                            }]
                        },
                        options
                    });
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
            const texture = this.stats.map(stat => stat.texture);
            const mesh = this.stats.map(stat => stat.mesh);
            const animation = this.stats.map(stat => stat.animation);
            const audio = this.stats.map(stat => stat.audio);
            const font = this.stats.map(stat => stat.font);
            const textAsset = this.stats.map(stat => stat.text_asset);
            const shader = this.stats.map(stat => stat.shader);
            const pics = this.stats.map(stat => stat.pic || '');
            const processes = this.stats.map(stat => stat.process || '');

            // 更新内存图表
            const memoryChart = this.chartInstances['memoryChart'];
            if (memoryChart) {
                memoryChart.data.labels = timestamps;
                memoryChart.data.datasets[0].data = totalMem;
                memoryChart.data.datasets[1].data = usedMem;
                memoryChart.data.datasets[2].data = monoMem;
                memoryChart.data.datasets.forEach(dataset => {
                    dataset.pic = pics;
                    dataset.process = processes;
                });
                memoryChart.update();
            }

            // 更新其他图表
            const chartData = {
                'fpsChart': fps,
                'textureChart': texture,
                'meshChart': mesh,
                'animationChart': animation,
                'audioChart': audio,
                'fontChart': font,
                'textAssetChart': textAsset,
                'shaderChart': shader
            };

            Object.entries(chartData).forEach(([chartId, data]) => {
                const chart = this.chartInstances[chartId];
                if (chart) {
                    chart.data.labels = timestamps;
                    chart.data.datasets[0].data = data;
                    chart.data.datasets[0].pic = pics;
                    chart.data.datasets[0].process = processes;
                    chart.update();
                }
            });
        },

        async fetchStats() {
            try {
                const params = new URLSearchParams({
                    ...this.filters
                });

                const response = await fetch(`/api/stats?${params}`);
                if (!response.ok) throw new Error('获取统计列表失败');
                const data = await response.json();
                this.stats = data.stats;
                this.total = data.total;
                this.updateCharts();
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        viewDetails(stat) {
            this.currentStat = stat;
            this.currentStatIndex = this.stats.findIndex(s => s.id === stat.id);
            this.showModal = true;
            this.updateClickedPointInfo(stat.pic, stat.process);
        },

        closeModal() {
            this.showModal = false;
            this.currentStat = null;
            this.currentStatIndex = -1;
        },

        hasPreviousStat() {
            return this.currentStatIndex > 0;
        },

        hasNextStat() {
            return this.currentStatIndex < this.stats.length - 1;
        },

        previousStat() {
            if (this.hasPreviousStat()) {
                this.currentStatIndex--;
                this.currentStat = this.stats[this.currentStatIndex];
                this.updateClickedPointInfo(this.currentStat.pic, this.currentStat.process);
            }
        },

        nextStat() {
            if (this.hasNextStat()) {
                this.currentStatIndex++;
                this.currentStat = this.stats[this.currentStatIndex];
                this.updateClickedPointInfo(this.currentStat.pic, this.currentStat.process);
            }
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
            if (!confirm('确定要清理旧数据吗？此操作不可恢复。')) return;

            try {
                const response = await fetch(`/api/stats/before?date=${this.filters.endDate}`, {
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

        formatMemory(bytes) {
            if (!bytes) return '0 B';
            
            const units = ['B', 'KB', 'MB', 'GB', 'TB'];
            const k = 1024;
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            
            if (i === 0) return bytes + ' ' + units[i];
            
            return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + units[i];
        },

        formatProcess(process) {
            if (typeof process === 'string') {
                try {
                    process = JSON.parse(process);
                } catch (e) {
                    return '解析执行统计数据的 JSON 失败：' + process;
                }
            }

            if (!Array.isArray(process)) return '';

            return process.map(item => {
                if (!item) return '';

                let result = '[' + item.name + ']<br>';
                if (Array.isArray(item.list)) {
                    return result + item.list.map(subItem => {
                        return subItem.replace(/\\n/g, '<br>').replace(/ /g, '&nbsp;')
                    }).join('<br>');
                }
                return result.replace(/\\n/g, '<br>').replace(/ /g, '&nbsp;')
            }).join('<br><br>');
        },

        updateClickedPointInfo(pic, process) {
            const processElement = document.getElementById('processInfo');
            processElement.innerHTML = this.formatProcess(process);
            
            const screenshotElement = document.getElementById('screenshot');
            const picPath = pic ? pic.replace(/\\/g, '/') : '';
            let picHTML = '';
            
            if (picPath) {
                picHTML = `
                    <img src="${picPath}" alt="Stats Image" class="h-40 cursor-zoom-in stats-thumbnail" 
                         @mouseenter="showEnlarged('${picPath}')"
                         @mouseleave="hideEnlarged()">
                `;
            }
            screenshotElement.innerHTML = picHTML;
        },

        showEnlarged(imgSrc) {
            const enlargedContainer = document.getElementById('enlargedImageContainer');
            const enlargedImage = document.getElementById('enlargedImage');
            enlargedImage.src = imgSrc;
            enlargedContainer.classList.remove('hidden');
        },

        hideEnlarged() {
            const enlargedContainer = document.getElementById('enlargedImageContainer');
            setTimeout(() => {
                enlargedContainer.classList.add('hidden');
            }, 300);
        },

        navigateChartPoint(direction, chartId) {
            const chart = this.chartInstances[chartId];
            if (!chart) return;

            const dataset = chart.data.datasets[0];
            if (!dataset || !dataset.data || dataset.data.length === 0) return;

            if (this.currentPointIndex[chartId] === undefined) {
                this.currentPointIndex[chartId] = 0;
            }

            let newIndex;
            if (direction === 'nextPoint') {
                newIndex = this.currentPointIndex[chartId] - 1;
                if (newIndex < 0) {
                    Alpine.store('notification').show('已经是最后一个数据点', 'warning');
                    return;
                }
            } else {
                newIndex = this.currentPointIndex[chartId] + 1;
                if (newIndex >= dataset.data.length) {
                    Alpine.store('notification').show('已经是第一个数据点', 'warning');
                    return;
                }
            }

            this.currentPointIndex[chartId] = newIndex;

            const pic = dataset.pic[this.currentPointIndex[chartId]];
            const process = dataset.process[this.currentPointIndex[chartId]];

            this.updateClickedPointInfo(pic, process);

            chart.setActiveElements([{
                datasetIndex: 0,
                index: this.currentPointIndex[chartId]
            }]);
            chart.update();
        },

        resetZoom(chartId) {
            const chart = this.chartInstances[chartId];
            if (chart && chart.resetZoom) {
                chart.resetZoom();
            }
        },

        watchFilters() {
            this.$watch('filters', () => {
                this.fetchStats();
            });
        }
    }
} 