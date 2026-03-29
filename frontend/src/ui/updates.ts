import { reactive } from 'vue';
import { dockerApi } from '../api';
import { appSettings } from './settings';

type UpdateStatus = 'idle' | 'checking' | 'available' | 'up-to-date' | 'error';

type UpdateCheckResult = {
  latestVersion: string | null;
  updateUrl: string;
  checkedAt: string;
  message: string;
  releaseDate: string | null;
  hasUpdate: boolean;
};

const state = reactive({
  status: 'idle' as UpdateStatus,
  currentVersion: appSettings.about.appVersion,
  latestVersion: null as string | null,
  message: '',
  checkedAt: '' as string,
  releaseDate: null as string | null,
  updateUrl: '',
});

const getUpdateUrl = () => {
  const namespace = encodeURIComponent(appSettings.updates.dockerHubNamespace.trim());
  const repoPrefix = encodeURIComponent(appSettings.updates.dockerHubRepoPrefix.trim());
  return `https://hub.docker.com/r/${namespace}/${repoPrefix}-frontend/tags`;
};

const formatError = (error: unknown) => {
  if (error instanceof Error && error.message) return error.message;
  return 'Unable to reach Docker Hub right now.';
};

const checkForUpdates = async (): Promise<UpdateCheckResult> => {
  const namespace = appSettings.updates.dockerHubNamespace.trim();
  const repoPrefix = appSettings.updates.dockerHubRepoPrefix.trim();

  if (!namespace || !repoPrefix) {
    throw new Error('Docker Hub namespace and repository prefix are required.');
  }

  const response = await dockerApi.checkAppUpdates({
    currentVersion: appSettings.about.appVersion,
    namespace,
    repoPrefix,
  });
  const payload = response.data as {
    latestVersion?: string | null;
    updateUrl?: string;
    checkedAt?: string;
    message?: string;
    releaseDate?: string | null;
    hasUpdate?: boolean;
  };

  return {
    latestVersion: payload.latestVersion || null,
    updateUrl: payload.updateUrl || getUpdateUrl(),
    checkedAt: payload.checkedAt || new Date().toISOString(),
    releaseDate: payload.releaseDate || null,
    message: payload.message || 'Unable to determine update status.',
    hasUpdate: !!payload.hasUpdate,
  };
};

const refresh = async (opts?: { silent?: boolean }) => {
  state.currentVersion = appSettings.about.appVersion;
  state.status = 'checking';
  if (!opts?.silent) state.message = 'Checking Docker Hub for a newer frontend image...';

  try {
    const result = await checkForUpdates();
    state.latestVersion = result.latestVersion;
    state.checkedAt = result.checkedAt;
    state.releaseDate = result.releaseDate;
    state.updateUrl = result.updateUrl;
    state.message = result.message;
    state.status = result.hasUpdate
      ? 'available'
      : 'up-to-date';
    return result;
  } catch (error) {
    state.status = 'error';
    state.message = formatError(error);
    state.checkedAt = new Date().toISOString();
    state.updateUrl = getUpdateUrl();
    throw error;
  }
};

const openUpdateUrl = () => {
  const target = state.updateUrl || getUpdateUrl();
  window.open(target, '_blank', 'noopener,noreferrer');
};

export const updates = {
  state,
  refresh,
  openUpdateUrl,
};
