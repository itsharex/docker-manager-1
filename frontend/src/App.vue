<script setup lang="ts">
import { ref, onMounted } from 'vue';
import {
  LayoutDashboard,
  Container,
  Box,
  HardDrive,
  Network,
  Activity,
  Cpu,
  Database
} from 'lucide-vue-next';
import { dockerApi } from './api';
import ContainerList from './components/ContainerList.vue';
import ImageList from './components/ImageList.vue';
import VolumeList from './components/VolumeList.vue';
import NetworkList from './components/NetworkList.vue';

const activeTab = ref('containers');
const systemInfo = ref<any>(null);
const containers = ref<any[]>([]);

const tabs = [
  { id: 'dashboard', name: 'Dashboard', icon: LayoutDashboard },
  { id: 'containers', name: 'Containers', icon: Container },
  { id: 'images', name: 'Images', icon: Box },
  { id: 'volumes', name: 'Volumes', icon: HardDrive },
  { id: 'networks', name: 'Networks', icon: Network },
];

const fetchStats = async () => {
  try {
    const { data: info } = await dockerApi.getSystemInfo();
    systemInfo.value = info;
    const { data: containerList } = await dockerApi.getContainers();
    containers.value = containerList;
  } catch (err) {
    console.error('Failed to fetch stats:', err);
  }
};

onMounted(() => {
  fetchStats();
  setInterval(fetchStats, 5000);
});
</script>

<template>
  <div class="app-container">
    <!-- Sidebar -->
    <aside class="sidebar glass-panel">
      <div class="logo">
        <Activity class="icon-primary" :size="32" />
        <span>Docker Hub</span>
      </div>

      <nav class="nav-links">
        <button v-for="tab in tabs" :key="tab.id" class="nav-item" :class="{ active: activeTab === tab.id }"
          @click="activeTab = tab.id">
          <component :is="tab.icon" :size="20" />
          {{ tab.name }}
        </button>
      </nav>

      <div class="sidebar-footer">
        <div class="system-stats" v-if="systemInfo">
          <div class="stat-item">
            <Cpu :size="16" />
            <span>{{ systemInfo.NCPU }} CPUs</span>
          </div>
          <div class="stat-item">
            <Database :size="16" />
            <span>{{ (systemInfo.MemTotal / 1024 / 1024 / 1024).toFixed(1) }} GB</span>
          </div>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="main-content">
      <header class="content-header">
        <h1>{{tabs.find(t => t.id === activeTab)?.name}}</h1>
        <div class="header-actions">
          <div class="status-badge" v-if="systemInfo">
            <span class="pulse"></span>
            Docker {{ systemInfo.ServerVersion }}
          </div>
        </div>
      </header>

      <section class="content-area animate-fade-in">
        <!-- Dashboard overflow -->
        <div v-if="activeTab === 'dashboard'" class="dashboard-grid">
          <div class="stat-card glass-panel">
            <div class="card-header">
              <Container :size="20" class="icon-indigo" />
              <h3>Containers</h3>
            </div>
            <div class="value">{{ systemInfo?.ContainersRunning || 0 }} / {{ systemInfo?.Containers || 0 }}</div>
            <div class="label">Running / Total</div>
          </div>
          <div class="stat-card glass-panel">
            <div class="card-header">
              <Box :size="20" class="icon-indigo" />
              <h3>Images</h3>
            </div>
            <div class="value">{{ systemInfo?.Images || 0 }}</div>
            <div class="label">Total Images</div>
          </div>
          <div class="stat-card glass-panel">
            <div class="card-header">
              <Cpu :size="20" class="icon-indigo" />
              <h3>System</h3>
            </div>
            <div class="value">{{ systemInfo?.NCPU || 0 }}</div>
            <div class="label">Core Count</div>
          </div>
        </div>

        <!-- Components -->
        <ContainerList v-else-if="activeTab === 'containers'" />
        <ImageList v-else-if="activeTab === 'images'" />
        <VolumeList v-else-if="activeTab === 'volumes'" />
        <NetworkList v-else-if="activeTab === 'networks'" />
      </section>
    </main>
  </div>
</template>

<style scoped>
.app-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  background: radial-gradient(circle at top right, rgba(99, 102, 241, 0.05), transparent),
    radial-gradient(circle at bottom left, rgba(79, 70, 229, 0.05), transparent);
}

.sidebar {
  width: 260px;
  height: calc(100vh - 32px);
  margin: 16px;
  display: flex;
  flex-direction: column;
  padding: 24px 16px;
  flex-shrink: 0;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 40px;
  padding: 0 8px;
  color: white;
}

.icon-primary {
  color: var(--primary);
  filter: drop-shadow(0 0 8px rgba(99, 102, 241, 0.5));
}

.nav-links {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex-grow: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 12px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
}

.nav-item:hover {
  background: var(--glass);
  color: var(--text-main);
}

.nav-item.active {
  background: var(--primary);
  color: white;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}

.sidebar-footer {
  margin-top: auto;
  padding-top: 24px;
  border-top: 1px solid var(--glass-border);
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.system-stats {
  padding: 0 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.85rem;
  color: var(--text-muted);
}

.main-content {
  flex-grow: 1;
  padding: 40px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.content-header h1 {
  font-size: 2rem;
  font-weight: 700;
}

.status-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 16px;
  background: var(--glass);
  border: 1px solid var(--glass-border);
  border-radius: 20px;
  font-size: 0.85rem;
  color: var(--text-muted);
}

.pulse {
  width: 8px;
  height: 8px;
  background: var(--success);
  border-radius: 50%;
  box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.4);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7);
  }

  70% {
    box-shadow: 0 0 0 10px rgba(16, 185, 129, 0);
  }

  100% {
    box-shadow: 0 0 0 0 rgba(16, 185, 129, 0);
  }
}

.content-area {
  flex-grow: 1;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
}

.stat-card {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.stat-card h3 {
  font-size: 1rem;
  color: var(--text-muted);
  font-weight: 500;
}

.stat-card .value {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--text-main);
}

.stat-card .label {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.icon-indigo {
  color: var(--primary);
}

.stat-card h3 {
  font-size: 1rem;
  color: var(--text-muted);
  font-weight: 500;
  margin: 0;
}
</style>
