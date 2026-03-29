<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { feedback } from './ui/feedback';
import { authState, clearAuthSession } from './ui/auth';
import { setUnauthorizedHandler } from './api';
import UiFeedback from './components/UiFeedback.vue';

const router = useRouter();
const { t } = useI18n();

onMounted(() => {
  setUnauthorizedHandler(() => {
    if (!authState.user && !authState.token) return;
    clearAuthSession();
    feedback.warning(t('app.sessionExpired'));
    if (router.currentRoute.value.name !== 'auth') {
      void router.replace({ name: 'auth' });
    }
  });
});

onUnmounted(() => {
  setUnauthorizedHandler(null);
});
</script>

<template>
  <router-view />
  <UiFeedback />
</template>
