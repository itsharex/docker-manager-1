<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import {
    Play,
    Square,
    Trash2,
    Terminal,
    FileText,
    Search,
    RefreshCw
} from 'lucide-vue-next';
import { dockerApi } from '../api';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

dayjs.extend(relativeTime);

const containers = ref<any[]>([]);
const loading = ref(true);
const searchQuery = ref('');

const fetchContainers = async () => {
    try {
        const { data } = await dockerApi.getContainers();
        containers.value = data;
    } catch (err) {
        console.error('Failed to fetch containers:', err);
    } finally {
        loading.value = false;
    }
};

const handleAction = async (action: string, id: string) => {
    try {
        if (action === 'start') await dockerApi.startContainer(id);
        else if (action === 'stop') await dockerApi.stopContainer(id);
        else if (action === 'remove') {
            if (confirm('Are you sure you want to remove this container?')) {
                await dockerApi.removeContainer(id);
            } else return;
        }
        await fetchContainers();
    } catch (err) {
        alert(`Action failed: ${err}`);
    }
};

const getStatusColor = (status: string) => {
    if (status.includes('Up')) return 'var(--success)';
    if (status.includes('Exited')) return 'var(--danger)';
    return 'var(--warning)';
};

let interval: any;
onMounted(() => {
    fetchContainers();
    interval = setInterval(fetchContainers, 5000);
});

onUnmounted(() => {
    clearInterval(interval);
});
</script>

<template>
    <div class="container-list-view">
        <div class="toolbar glass-panel">
            <div class="search-box">
                <Search :size="18" />
                <input v-model="searchQuery" type="text" placeholder="Search containers..." />
            </div>
            <button class="btn btn-ghost" @click="fetchContainers">
                <RefreshCw :size="18" :class="{ 'animate-spin': loading }" />
                Refresh
            </button>
        </div>

        <div class="table-container glass-panel">
            <table class="docker-table">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Image</th>
                        <th>Status</th>
                        <th>Ports</th>
                        <th>Created</th>
                        <th class="actions-cell">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="container in containers" :key="container.Id">
                        <td class="name-cell">
                            <div class="container-name">
                                {{ container.Names[0].replace('/', '') }}
                                <span class="id-short">{{ container.Id.substring(0, 12) }}</span>
                            </div>
                        </td>
                        <td><code class="image-name">{{ container.Image }}</code></td>
                        <td>
                            <div class="status-pill" :style="{ '--color': getStatusColor(container.Status) }">
                                <span class="dot"></span>
                                {{ container.Status }}
                            </div>
                        </td>
                        <td>
                            <div class="ports">
                                <span v-for="port in container.Ports" :key="port.PublicPort" class="port-tag">
                                    {{ port.PublicPort }}:{{ port.PrivatePort }}
                                </span>
                            </div>
                        </td>
                        <td class="time-cell">{{ dayjs.unix(container.Created).fromNow() }}</td>
                        <td class="actions-cell">
                            <div class="action-group">
                                <button v-if="!container.Status.includes('Up')" class="btn-icon btn-ghost" title="Start"
                                    @click="handleAction('start', container.Id)">
                                    <Play :size="16" />
                                </button>
                                <button v-else class="btn-icon btn-ghost" title="Stop"
                                    @click="handleAction('stop', container.Id)">
                                    <Square :size="16" />
                                </button>
                                <button class="btn-icon btn-ghost" title="Logs">
                                    <FileText :size="16" />
                                </button>
                                <button class="btn-icon btn-ghost" title="Terminal">
                                    <Terminal :size="16" />
                                </button>
                                <button class="btn-icon btn-ghost text-danger" title="Remove"
                                    @click="handleAction('remove', container.Id)">
                                    <Trash2 :size="16" />
                                </button>
                            </div>
                        </td>
                    </tr>
                    <tr v-if="containers.length === 0 && !loading">
                        <td colspan="6" class="empty-state">No containers found</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.container-list-view {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.toolbar {
    padding: 12px 24px;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.search-box {
    display: flex;
    align-items: center;
    gap: 12px;
    background: var(--glass);
    padding: 8px 16px;
    border-radius: 10px;
    border: 1px solid var(--glass-border);
    width: 300px;
}

.search-box input {
    background: transparent;
    border: none;
    color: white;
    outline: none;
    font-size: 0.9rem;
    width: 100%;
}

.table-container {
    overflow: hidden;
}

.docker-table {
    width: 100%;
    border-collapse: collapse;
    text-align: left;
}

.docker-table th {
    padding: 16px 24px;
    font-family: 'Outfit', sans-serif;
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--text-muted);
    border-bottom: 1px solid var(--glass-border);
}

.docker-table td {
    padding: 16px 24px;
    font-size: 0.9rem;
    border-bottom: 1px solid var(--glass-border);
}

.docker-table tr:last-child td {
    border-bottom: none;
}

.docker-table tr:hover {
    background: var(--glass);
}

.container-name {
    font-weight: 600;
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.id-short {
    font-size: 0.75rem;
    color: var(--text-muted);
    font-weight: 400;
}

.image-name {
    background: var(--glass);
    padding: 4px 8px;
    border-radius: 6px;
    color: var(--primary);
}

.status-pill {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 4px 12px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid var(--glass-border);
    border-radius: 20px;
    font-size: 0.8rem;
    color: var(--color);
}

.dot {
    width: 6px;
    height: 6px;
    background: var(--color);
    border-radius: 50%;
}

.ports {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
}

.port-tag {
    background: var(--glass);
    color: var(--text-muted);
    font-size: 0.75rem;
    padding: 2px 6px;
    border-radius: 4px;
}

.action-group {
    display: flex;
    gap: 4px;
}

.actions-cell {
    text-align: right;
    width: 200px;
}

.text-danger {
    color: var(--danger) !important;
}

.animate-spin {
    animation: spin 1s linear infinite;
}

@keyframes spin {
    from {
        transform: rotate(0deg);
    }

    to {
        transform: rotate(360deg);
    }
}

.empty-state {
    padding: 80px 0;
    text-align: center;
    color: var(--text-muted);
}
</style>
