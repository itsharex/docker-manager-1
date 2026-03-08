<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Network, Trash2, RefreshCw, Globe, ShieldCheck } from 'lucide-vue-next';
import { dockerApi } from '../api';

const networks = ref<any[]>([]);
const loading = ref(true);

const fetchNetworks = async () => {
    try {
        loading.value = true;
        const { data } = await dockerApi.getNetworks();
        networks.value = data || [];
    } catch (err) {
        console.error('Failed to fetch networks:', err);
    } finally {
        loading.value = false;
    }
};

const removeNetwork = async (id: string) => {
    if (confirm('Are you sure you want to remove this network?')) {
        try {
            await dockerApi.removeNetwork(id);
            await fetchNetworks();
        } catch (err) {
            alert(`Failed to remove network: ${err}`);
        }
    }
};

onMounted(fetchNetworks);
</script>

<template>
    <div class="network-list-view">
        <div class="toolbar glass-panel">
            <div class="title-with-icon">
                <Network :size="20" class="icon-indigo" />
                <h2>Networks</h2>
            </div>
            <button class="btn btn-ghost" @click="fetchNetworks">
                <RefreshCw :size="18" :class="{ 'animate-spin': loading }" />
                Refresh
            </button>
        </div>

        <div class="grid-container">
            <div v-for="net in networks" :key="net.Id" class="network-card glass-panel animate-fade-in">
                <div class="network-header">
                    <div class="network-name">{{ net.Name }}</div>
                    <button class="btn-icon btn-ghost text-danger" @click="removeNetwork(net.Id)">
                        <Trash2 :size="16" />
                    </button>
                </div>
                <div class="network-details">
                    <div class="detail-item">
                        <Globe :size="14" />
                        <span>Driver: {{ net.Driver }}</span>
                    </div>
                    <div class="detail-item">
                        <ShieldCheck :size="14" />
                        <span>Scope: {{ net.Scope }}</span>
                    </div>
                </div>
                <div class="network-id">{{ net.Id.substring(0, 12) }}</div>
            </div>
            <div v-if="networks.length === 0 && !loading" class="empty-state glass-panel">
                No networks found
            </div>
        </div>
    </div>
</template>

<style scoped>
.network-list-view {
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

.network-card {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 12px;
    position: relative;
}

.network-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.network-name {
    font-family: 'Outfit', sans-serif;
    font-weight: 600;
    color: var(--text-main);
}

.network-details {
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

.network-id {
    font-size: 0.7rem;
    color: var(--text-muted);
    opacity: 0.5;
    margin-top: 4px;
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
