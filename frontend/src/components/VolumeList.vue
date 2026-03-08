<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { HardDrive, Trash2, RefreshCw, Calendar, Server } from 'lucide-vue-next';
import { dockerApi } from '../api';
import dayjs from 'dayjs';

const volumes = ref<any[]>([]);
const loading = ref(true);

const fetchVolumes = async () => {
    try {
        loading.value = true;
        const { data } = await dockerApi.getVolumes();
        volumes.value = data || [];
    } catch (err) {
        console.error('Failed to fetch volumes:', err);
    } finally {
        loading.value = false;
    }
};

const removeVolume = async (id: string) => {
    if (confirm('Are you sure you want to remove this volume?')) {
        try {
            await dockerApi.removeVolume(id);
            await fetchVolumes();
        } catch (err) {
            alert(`Failed to remove volume: ${err}`);
        }
    }
};

onMounted(fetchVolumes);
</script>

<template>
    <div class="volume-list-view">
        <div class="toolbar glass-panel">
            <div class="title-with-icon">
                <HardDrive :size="20" class="icon-indigo" />
                <h2>Volumes</h2>
            </div>
            <button class="btn btn-ghost" @click="fetchVolumes">
                <RefreshCw :size="18" :class="{ 'animate-spin': loading }" />
                Refresh
            </button>
        </div>

        <div class="grid-container">
            <div v-for="vol in volumes" :key="vol.Name" class="volume-card glass-panel animate-fade-in">
                <div class="volume-header">
                    <div class="volume-name">
                        {{ vol.Name.substring(0, 20) }}{{ vol.Name.length > 20 ? '...' : '' }}
                    </div>
                    <button class="btn-icon btn-ghost text-danger" @click="removeVolume(vol.Name)">
                        <Trash2 :size="16" />
                    </button>
                </div>
                <div class="volume-details">
                    <div class="detail-item">
                        <Server :size="14" />
                        <span>Driver: {{ vol.Driver }}</span>
                    </div>
                    <div class="detail-item">
                        <Calendar :size="14" />
                        <span>{{ dayjs(vol.CreatedAt).format('MMM D, YYYY') }}</span>
                    </div>
                </div>
            </div>
            <div v-if="volumes.length === 0 && !loading" class="empty-state glass-panel">
                No volumes found
            </div>
        </div>
    </div>
</template>

<style scoped>
.volume-list-view {
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

.title-with-icon {
    display: flex;
    align-items: center;
    gap: 12px;
}

.title-with-icon h2 {
    font-size: 1.2rem;
    margin: 0;
}

.grid-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
}

.volume-card {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.volume-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.volume-name {
    font-family: 'Outfit', sans-serif;
    font-weight: 600;
    color: var(--text-main);
    word-break: break-all;
}

.volume-details {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.detail-item {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.85rem;
    color: var(--text-muted);
}

.icon-indigo {
    color: var(--primary);
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
    grid-column: 1 / -1;
    padding: 60px;
    text-align: center;
    color: var(--text-muted);
}
</style>
