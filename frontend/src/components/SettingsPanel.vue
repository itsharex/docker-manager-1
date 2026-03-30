<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { appSettings } from '../ui/settings';
import { updates } from '../ui/updates';
import { feedback } from '../ui/feedback';

const props = defineProps<{
  systemInfo?: any;
}>();

const { t } = useI18n();

const apiHint = computed(() => `${appSettings.runtime.apiBaseUrl.replace(/\/+$/, '')}/api`);
const languageLabelKey = computed(() => {
  switch (appSettings.general.language) {
    case 'vi':
      return 'settings.vietnamese';
    case 'zh':
      return 'settings.chinese';
    default:
      return 'settings.english';
  }
});
const summaryCards = computed(() => [
  { label: t('settings.activeLanguage'), value: t(languageLabelKey.value) },
  { label: t('settings.appVersion'), value: `v${appSettings.about.appVersion}` },
  { label: t('settings.buildDate'), value: appSettings.about.buildDate },
]);

const resetUI = () => {
  appSettings.ui.theme = 'dark';
  appSettings.ui.density = 'comfortable';
  appSettings.ui.fontScale = 1;
  appSettings.ui.showSidebarStats = true;
};

const updateState = updates.state;
const updateStatusTone = computed(() => {
  switch (updateState.status) {
    case 'available':
      return 'border-color: var(--warning-soft-border); color: var(--warning-soft-text); background: var(--warning-soft-bg);';
    case 'up-to-date':
      return 'border-color: rgba(61, 220, 132, 0.4); color: #bbf7d0; background: rgba(61, 220, 132, 0.12);';
    case 'error':
      return 'border-color: var(--danger-soft-border); color: var(--danger-soft-text); background: var(--danger-soft-bg);';
    default:
      return 'border-color: var(--glass-border); color: var(--text-muted); background: var(--glass);';
  }
});

const checkedAtLabel = computed(() => {
  if (!updateState.checkedAt) return t('settings.neverChecked');
  return new Date(updateState.checkedAt).toLocaleString();
});

const releaseDateLabel = computed(() => {
  if (!updateState.releaseDate) return t('common.notAvailable');
  return new Date(updateState.releaseDate).toLocaleString();
});

const statusLabel = computed(() => {
  switch (updateState.status) {
    case 'checking':
      return t('settings.updateChecking');
    case 'available':
      return t('settings.updateAvailable');
    case 'up-to-date':
      return t('settings.updateUpToDate');
    case 'error':
      return t('settings.updateCheckFailed');
    default:
      return t('settings.updateIdle');
  }
});

const checkUpdates = async (silent = false) => {
  try {
    await updates.refresh({ silent });
    if (!silent) {
      if (updateState.status === 'available') feedback.info(updateState.message);
      else if (updateState.status === 'up-to-date') feedback.success(updateState.message);
    }
  } catch {
    if (!silent) feedback.error(updateState.message);
  }
};

const openUpdatePage = () => {
  updates.openUpdateUrl();
};

const applyUpdate = async () => {
  const accepted = await feedback.confirmAction({
    title: t('common.pleaseConfirm'),
    message: t('settings.updateConfirm', { version: updateState.latestVersion || 'latest' }),
    confirmText: t('settings.updateNow'),
    cancelText: t('common.cancel'),
  });
  if (!accepted) return;

  try {
    const result = await updates.apply();
    feedback.info(result.message || t('settings.updateStarted'));
    window.setTimeout(() => {
      window.location.reload();
    }, 12000);
  } catch {
    feedback.error(updateState.message || t('common.actionFailed'));
  }
};

onMounted(() => {
  if (appSettings.updates.autoCheck && !updateState.checkedAt && updateState.status === 'idle') {
    void checkUpdates(true);
  }
});
</script>

<template>
  <div class="flex flex-col gap-6">
    <section class="grid gap-5 xl:grid-cols-[1.1fr_0.9fr]">
      <div class="glass-panel p-6">
        <p class="section-heading">{{ t('settings.title') }}</p>
        <h2 class="max-w-3xl text-3xl font-bold tracking-tight">{{ t('settings.title') }}</h2>
        <p class="mt-3 max-w-2xl text-sm leading-7" style="color: var(--text-muted);">
          {{ t('settings.subtitle') }}
        </p>

        <div class="mt-6 grid gap-3 sm:grid-cols-3">
          <div v-for="card in summaryCards" :key="card.label" class="border p-4" style="border-color: var(--glass-border); background: var(--glass);">
            <p class="text-[11px] uppercase tracking-[0.22em]" style="color: var(--text-muted);">{{ card.label }}</p>
            <p class="mt-2 text-2xl font-bold">{{ card.value }}</p>
          </div>
        </div>
      </div>

      <div class="glass-panel p-6">
        <div class="mb-4 flex items-center justify-between gap-4">
          <p class="section-heading mb-0">{{ t('settings.quickActions') }}</p>
          <span class="border px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em]" style="border-color: var(--primary); color: var(--primary);">
            {{ t('settings.statusLive') }}
          </span>
        </div>

        <div class="grid gap-4">
          <label class="block">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.language') }}</span>
            <select v-model="appSettings.general.language" class="app-select">
              <option value="vi">{{ t('settings.vietnamese') }}</option>
              <option value="en">{{ t('settings.english') }}</option>
              <option value="zh">{{ t('settings.chinese') }}</option>
            </select>
          </label>

          <label class="block">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.theme') }}</span>
            <select v-model="appSettings.ui.theme" class="app-select">
              <option value="dark">{{ t('settings.dark') }}</option>
              <option value="light">{{ t('settings.light') }}</option>
            </select>
          </label>

          <button class="btn btn-ghost w-full" @click="resetUI">{{ t('settings.resetUi') }}</button>
        </div>
      </div>
    </section>

    <div class="grid gap-6 2xl:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
      <section class="glass-panel p-5">
        <p class="section-heading">{{ t('settings.general') }}</p>
        <div class="grid gap-4 lg:grid-cols-2">
          <label class="block lg:col-span-2">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.autoRefresh') }}</span>
            <select v-model.number="appSettings.general.autoRefreshMs" class="app-select">
              <option :value="0">{{ t('settings.off') }}</option>
              <option :value="2000">2s</option>
              <option :value="5000">5s</option>
              <option :value="10000">10s</option>
            </select>
            <small class="mt-2 block text-xs" style="color: var(--text-muted);">{{ t('settings.autoRefreshHelp') }}</small>
          </label>

          <label class="flex items-center justify-between border px-4 py-3" style="border-color: var(--glass-border); background: var(--glass);">
            <span class="text-sm font-semibold">{{ t('settings.confirmDestructive') }}</span>
            <input v-model="appSettings.general.confirmDestructive" type="checkbox" class="h-5 w-5 accent-blue-600" />
          </label>

          <label class="block">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.timeFormat') }}</span>
            <select v-model="appSettings.general.timeFormat" class="app-select">
              <option value="24h">{{ t('settings.hour24') }}</option>
              <option value="12h">{{ t('settings.hour12') }}</option>
            </select>
          </label>
        </div>
      </section>

      <section class="glass-panel p-5">
        <p class="section-heading">{{ t('settings.interface') }}</p>
        <div class="grid gap-4 lg:grid-cols-2">
          <label class="block">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.theme') }}</span>
            <select v-model="appSettings.ui.theme" class="app-select">
              <option value="dark">{{ t('settings.dark') }}</option>
              <option value="light">{{ t('settings.light') }}</option>
            </select>
          </label>

          <label class="block">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.density') }}</span>
            <select v-model="appSettings.ui.density" class="app-select">
              <option value="comfortable">{{ t('settings.comfortable') }}</option>
              <option value="compact">{{ t('settings.compact') }}</option>
            </select>
          </label>

          <label class="block lg:col-span-2">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.fontScale') }} ({{ appSettings.ui.fontScale.toFixed(2) }})</span>
            <div class="border px-4 py-4" style="border-color: var(--glass-border); background: var(--glass);">
              <input v-model.number="appSettings.ui.fontScale" type="range" min="0.9" max="1.15" step="0.01" class="w-full accent-blue-600" />
            </div>
          </label>

          <label class="flex items-center justify-between border px-4 py-3 lg:col-span-2" style="border-color: var(--glass-border); background: var(--glass);">
            <span class="text-sm font-semibold">{{ t('settings.showSidebarStats') }}</span>
            <input v-model="appSettings.ui.showSidebarStats" type="checkbox" class="h-5 w-5 accent-blue-600" />
          </label>
        </div>
      </section>

      <section class="glass-panel p-5">
        <p class="section-heading">{{ t('settings.runtime') }}</p>
        <div class="grid gap-4 lg:grid-cols-2">
          <label class="block lg:col-span-2">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.dockerApiEndpoint') }}</span>
            <input v-model.trim="appSettings.runtime.apiBaseUrl" type="text" placeholder="http://localhost:8080" class="app-input" />
            <small class="mt-2 block text-xs" style="color: var(--text-muted);">
              {{ t('settings.dockerApiHelp', { value: apiHint }) }}
            </small>
          </label>

          <label class="block">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.defaultLogTail') }}</span>
            <input v-model.number="appSettings.runtime.defaultLogTail" type="number" min="50" max="5000" step="50" class="app-input" />
          </label>

          <label class="block">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.terminalShell') }}</span>
            <select v-model="appSettings.runtime.terminalShell" class="app-select">
              <option value="/bin/sh">/bin/sh</option>
              <option value="/bin/bash">/bin/bash</option>
            </select>
          </label>

          <label class="block">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.terminalTheme') }}</span>
            <select v-model="appSettings.runtime.terminalTheme" class="app-select">
              <option value="ocean">{{ t('settings.themeOcean') }}</option>
              <option value="matrix">{{ t('settings.themeMatrix') }}</option>
              <option value="amber">{{ t('settings.themeAmber') }}</option>
            </select>
          </label>

          <label class="block">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.terminalFontSize') }}</span>
            <input v-model.number="appSettings.runtime.terminalFontSize" type="number" min="11" max="20" step="1" class="app-input" />
          </label>

          <label class="block lg:col-span-2">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.composeRefresh') }}</span>
            <select v-model.number="appSettings.runtime.composeRefreshMs" class="app-select">
              <option :value="0">{{ t('settings.off') }}</option>
              <option :value="2000">2s</option>
              <option :value="5000">5s</option>
              <option :value="10000">10s</option>
            </select>
          </label>
        </div>
      </section>

      <section class="glass-panel p-5">
        <div class="mb-4 flex items-center justify-between gap-4">
          <p class="section-heading mb-0">{{ t('settings.updates') }}</p>
          <span class="border px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em]" :style="updateStatusTone">
            {{ statusLabel }}
          </span>
        </div>

        <div class="grid gap-4 lg:grid-cols-2">
          <label class="flex items-center justify-between border px-4 py-3 lg:col-span-2" style="border-color: var(--glass-border); background: var(--glass);">
            <span class="text-sm font-semibold">{{ t('settings.autoCheckUpdates') }}</span>
            <input v-model="appSettings.updates.autoCheck" type="checkbox" class="h-5 w-5 accent-blue-600" />
          </label>

          <label class="block">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.dockerHubNamespace') }}</span>
            <input v-model.trim="appSettings.updates.dockerHubNamespace" type="text" class="app-input" placeholder="ngthanhvu" />
          </label>

          <label class="block">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.dockerHubRepoPrefix') }}</span>
            <input v-model.trim="appSettings.updates.dockerHubRepoPrefix" type="text" class="app-input" placeholder="docker-manager" />
          </label>

          <div class="border p-4 lg:col-span-2" style="border-color: var(--glass-border); background: var(--glass);">
            <div class="grid gap-3 sm:grid-cols-2">
              <div>
                <p class="text-[11px] uppercase tracking-[0.22em]" style="color: var(--text-muted);">{{ t('settings.currentVersion') }}</p>
                <p class="mt-2 text-xl font-bold">v{{ updateState.currentVersion }}</p>
              </div>
              <div>
                <p class="text-[11px] uppercase tracking-[0.22em]" style="color: var(--text-muted);">{{ t('settings.latestVersion') }}</p>
                <p class="mt-2 text-xl font-bold">{{ updateState.latestVersion ? `v${updateState.latestVersion}` : t('common.notAvailable') }}</p>
              </div>
              <div>
                <p class="text-[11px] uppercase tracking-[0.22em]" style="color: var(--text-muted);">{{ t('settings.lastChecked') }}</p>
                <p class="mt-2 text-sm font-medium">{{ checkedAtLabel }}</p>
              </div>
              <div>
                <p class="text-[11px] uppercase tracking-[0.22em]" style="color: var(--text-muted);">{{ t('settings.latestPublished') }}</p>
                <p class="mt-2 text-sm font-medium">{{ releaseDateLabel }}</p>
              </div>
            </div>

            <p class="mt-4 text-sm leading-6" style="color: var(--text-muted);">
              {{ updateState.message || t('settings.updateHelp') }}
            </p>
          </div>

          <div class="flex flex-wrap gap-3 lg:col-span-2">
            <button class="btn btn-ghost" type="button" :disabled="updateState.status === 'checking'" @click="checkUpdates()">
              {{ updateState.status === 'checking' ? t('settings.updateChecking') : t('settings.checkUpdates') }}
            </button>
            <button
              class="btn btn-primary"
              type="button"
              :disabled="updateState.status !== 'available' || updateState.applying"
              @click="applyUpdate"
            >
              {{ t('settings.updateNow') }}
            </button>
            <button
              class="btn btn-ghost"
              type="button"
              :disabled="updateState.status === 'checking'"
              @click="openUpdatePage"
            >
              {{ t('settings.openUpdatePage') }}
            </button>
          </div>
        </div>
      </section>

      <section class="glass-panel p-5">
        <p class="section-heading">{{ t('settings.notifications') }}</p>
        <div class="grid gap-4 lg:grid-cols-2">
          <label class="block lg:col-span-2">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.toastDuration') }}</span>
            <input v-model.number="appSettings.notifications.toastDurationMs" type="number" min="1000" max="10000" step="100" class="app-input" />
          </label>

          <label class="flex items-center justify-between border px-4 py-3" style="border-color: var(--glass-border); background: var(--glass);">
            <span class="text-sm font-semibold">{{ t('settings.showSuccessToasts') }}</span>
            <input v-model="appSettings.notifications.showSuccessToast" type="checkbox" class="h-5 w-5 accent-blue-600" />
          </label>

          <label class="flex items-center justify-between border px-4 py-3" style="border-color: var(--glass-border); background: var(--glass);">
            <span class="text-sm font-semibold">{{ t('settings.showDetailedErrors') }}</span>
            <input v-model="appSettings.notifications.showDetailedErrors" type="checkbox" class="h-5 w-5 accent-blue-600" />
          </label>
        </div>
      </section>

      <section class="glass-panel p-5">
        <p class="section-heading">{{ t('settings.safety') }}</p>
        <div class="grid gap-4">
          <label class="flex items-center justify-between border px-4 py-3" style="border-color: var(--glass-border); background: var(--glass);">
            <span class="text-sm font-semibold">{{ t('settings.requireDeleteTyping') }}</span>
            <input v-model="appSettings.safety.softDeleteRequireTyping" type="checkbox" class="h-5 w-5 accent-blue-600" />
          </label>

          <label class="block">
            <span class="mb-2 block text-sm font-semibold">{{ t('settings.protectedResources') }}</span>
            <input
              v-model="appSettings.safety.protectedResources"
              type="text"
              :placeholder="t('settings.protectedResourcesPlaceholder')"
              class="app-input"
            />
          </label>
        </div>
      </section>

      <section class="glass-panel p-5">
        <p class="section-heading">{{ t('settings.about') }}</p>
        <div class="grid gap-px border text-sm sm:grid-cols-[220px_minmax(0,1fr)]" style="border-color: var(--glass-border); background: var(--glass-border);">
          <div class="px-4 py-3 font-semibold" style="background: var(--table-header-bg);">{{ t('settings.appVersion') }}</div>
          <div class="px-4 py-3" style="background: var(--bg-card);">v{{ appSettings.about.appVersion }}</div>
          <div class="px-4 py-3 font-semibold" style="background: var(--table-header-bg);">{{ t('settings.buildDate') }}</div>
          <div class="px-4 py-3" style="background: var(--bg-card);">{{ appSettings.about.buildDate }}</div>
          <div class="px-4 py-3 font-semibold" style="background: var(--table-header-bg);">{{ t('settings.engine') }}</div>
          <div class="px-4 py-3" style="background: var(--bg-card);">{{ props.systemInfo?.ServerVersion || t('common.notAvailable') }}</div>
          <div class="px-4 py-3 font-semibold" style="background: var(--table-header-bg);">{{ t('settings.os') }}</div>
          <div class="px-4 py-3" style="background: var(--bg-card);">{{ props.systemInfo?.OperatingSystem || t('common.notAvailable') }}</div>
        </div>
      </section>
    </div>
  </div>
</template>
