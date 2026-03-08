<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Box, Trash2, RefreshCw, Layers } from 'lucide-vue-next';
import { dockerApi } from '../api';
import dayjs from 'dayjs';

const images = ref<any[]>([]);
const loading = ref(true);

const fetchImages = async () => {
    try {
        loading.value = true;
        const { data } = await dockerApi.getImages();
        images.value = data;
    } catch (err) {
        console.error('Failed to fetch images:', err);
    } finally {
        loading.value = false;
    }
};

const removeImage = async (id: string) => {
    if (confirm('Are you sure you want to remove this image?')) {
        try {
            await dockerApi.removeImage(id);
            await fetchImages();
        } catch (err) {
            alert(`Failed to remove image: ${err}`);
        }
    }
};

const formatSize = (bytes: number) => {
    return (bytes / 1024 / 1024).toFixed(1) + ' MB';
};

onMounted(fetchImages);
</script>

<template>
    <div class="image-list-view">
        <div class="toolbar glass-panel">
            <div class="title-with-icon">
                <Box :size="20" class="icon-indigo" />
                <h2>Images</h2>
            </div>
            <button class="btn btn-ghost" @click="fetchImages">
                <RefreshCw :size="18" :class="{ 'animate-spin': loading }" />
                Refresh
            </button>
        </div>

        <div class="grid-container">
            <div v-for="image in images" :key="image.Id" class="image-card glass-panel animate-fade-in">
                <div class="image-header">
                    <div class="image-repo">
                        {{ image.RepoTags?.[0] || '<none>:<none>' }}
                    </div>
                    <button class="btn-icon btn-ghost text-danger" @click="removeImage(image.Id)">
                        <Trash2 :size="16" />
                    </button>
                </div>
                <div class="image-footer">
                    <div class="info-item">
                        <Layers :size="14" />
                        <span>{{ image.Id.substring(7, 19) }}</span>
                    </div>
                    <div class="info-item">
                        <span>{{ formatSize(image.Size) }}</span>
                    </div>
                    <div class="info-item date">
                        {{ dayjs.unix(image.Created).format('MMM D, YYYY') }}
                    </div>
                </div>
            </div>
            <div v-if="images.length === 0 && !loading" class="empty-state glass-panel">
                No images found
            </div>
        </div>
    </div>
</template>

<style scoped>
.image-list-view {
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
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 20px;
}

.image-card {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 16px;
    transition: transform 0.2s;
}

.image-card:hover {
    transform: translateY(-4px);
    border-color: var(--primary);
}

.image-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
}

.image-repo {
    font-family: 'Outfit', sans-serif;
    font-weight: 600;
    font-size: 1rem;
    word-break: break-all;
    color: var(--text-main);
}

.image-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: auto;
    padding-top: 12px;
    border-top: 1px solid var(--glass-border);
}

.info-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 0.8rem;
    color: var(--text-muted);
}

.date {
    font-style: italic;
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
