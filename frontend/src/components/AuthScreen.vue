<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { LockKeyhole, ShieldCheck, UserRound } from 'lucide-vue-next';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  setupRequired: boolean;
  loading?: boolean;
}>();

const emit = defineEmits<{
  login: [payload: { username: string; password: string }];
  setup: [payload: { username: string; password: string }];
}>();
const { t } = useI18n();

const username = ref('');
const password = ref('');
const confirmPassword = ref('');
const localError = ref('');

const title = computed(() =>
  props.setupRequired ? t('auth.titleSetup') : t('auth.titleLogin')
);

const description = computed(() =>
  props.setupRequired
    ? t('auth.descriptionSetup')
    : t('auth.descriptionLogin')
);

const submitLabel = computed(() => {
  if (props.loading) return props.setupRequired ? t('auth.creatingAccount') : t('auth.signingIn');
  return props.setupRequired ? t('auth.createAccount') : t('auth.signIn');
});

const submit = () => {
  localError.value = '';

  if (username.value.trim().length < 3) {
    localError.value = t('auth.usernameError');
    return;
  }

  if (password.value.trim().length < 8) {
    localError.value = t('auth.passwordError');
    return;
  }

  if (props.setupRequired && password.value !== confirmPassword.value) {
    localError.value = t('auth.confirmPasswordError');
    return;
  }

  const payload = {
    username: username.value.trim(),
    password: password.value,
  };

  if (props.setupRequired) {
    emit('setup', payload);
    return;
  }

  emit('login', payload);
};

watch(
  () => props.setupRequired,
  () => {
    password.value = '';
    confirmPassword.value = '';
    localError.value = '';
  }
);
</script>

<template>
  <div class="flex min-h-screen min-h-dvh items-center justify-center p-4 sm:p-6">
      <div class="grid w-full max-w-5xl overflow-hidden border lg:grid-cols-[1.15fr_0.85fr]" style="border-color: var(--glass-border);">
        <section class="relative hidden min-h-[620px] overflow-hidden border-r p-10 lg:block" style="border-color: var(--glass-border); background: linear-gradient(180deg, rgba(29,78,216,0.08), rgba(0,0,0,0.08));">
          <div class="absolute inset-0 opacity-50" style="background:
            linear-gradient(180deg, rgba(255,255,255,0.04) 1px, transparent 1px),
            linear-gradient(90deg, rgba(255,255,255,0.04) 1px, transparent 1px); background-size: 28px 28px;">
          </div>
          <div class="relative flex h-full flex-col justify-between">
            <div>
              <p class="mb-3 text-[11px] uppercase tracking-[0.26em]" style="color: var(--text-muted);">{{ t('auth.accessNode') }}</p>
              <h1 class="max-w-md text-5xl font-bold leading-none tracking-tight">
                {{ t('auth.leftTitle') }}
              </h1>
              <p class="mt-6 max-w-lg text-base leading-7" style="color: var(--text-muted);">
                {{ t('auth.leftDescription') }}
              </p>
            </div>

            <div class="grid gap-4">
              <div class="flex items-start gap-4 border p-4" style="border-color: var(--glass-border); background: var(--glass);">
                <ShieldCheck class="mt-0.5" :size="18" style="color: var(--success);" />
                <div>
                  <p class="font-semibold">{{ t('auth.protectedFlow') }}</p>
                  <p class="mt-1 text-sm leading-6" style="color: var(--text-muted);">
                    {{ t('auth.protectedFlowDesc') }}
                  </p>
                </div>
              </div>
              <div class="grid grid-cols-2 gap-4">
                <div class="border p-4" style="border-color: var(--glass-border); background: var(--glass);">
                  <p class="text-[11px] uppercase tracking-[0.22em]" style="color: var(--text-muted);">{{ t('auth.surface') }}</p>
                  <p class="mt-2 text-2xl font-bold">{{ t('auth.cliLike') }}</p>
                </div>
                <div class="border p-4" style="border-color: var(--glass-border); background: var(--glass);">
                  <p class="text-[11px] uppercase tracking-[0.22em]" style="color: var(--text-muted);">{{ t('auth.theme') }}</p>
                  <p class="mt-2 text-2xl font-bold">{{ t('auth.rigidUi') }}</p>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section class="glass-panel min-h-[620px] p-6 sm:p-8 lg:p-10">
          <div class="mx-auto flex h-full max-w-md flex-col justify-center">
            <div class="mb-8">
              <div class="mb-5 flex items-center gap-4">
                <div class="grid h-14 w-14 place-items-center border" style="border-color: var(--primary); background: rgba(29, 78, 216, 0.12); color: var(--primary);">
                  <ShieldCheck :size="22" />
                </div>
                <div>
                  <p class="text-[11px] uppercase tracking-[0.24em]" style="color: var(--text-muted);">{{ t('auth.secureAccess') }}</p>
                  <h2 class="text-3xl font-bold tracking-tight">{{ title }}</h2>
                </div>
              </div>
              <p class="text-sm leading-6" style="color: var(--text-muted);">{{ description }}</p>
            </div>

            <form class="space-y-4" @submit.prevent="submit">
              <label class="block">
                <span class="mb-2 block text-sm font-semibold">{{ t('auth.username') }}</span>
                <div class="flex items-center gap-3 border px-4 py-3" style="border-color: var(--glass-border); background: var(--input-bg);">
                  <UserRound :size="16" style="color: var(--text-muted);" />
                  <input v-model="username" class="min-w-0 flex-1 bg-transparent outline-none" type="text" autocomplete="username" :placeholder="t('auth.usernamePlaceholder')" :disabled="loading" />
                </div>
              </label>

              <label class="block">
                <span class="mb-2 block text-sm font-semibold">{{ t('auth.password') }}</span>
                <div class="flex items-center gap-3 border px-4 py-3" style="border-color: var(--glass-border); background: var(--input-bg);">
                  <LockKeyhole :size="16" style="color: var(--text-muted);" />
                  <input v-model="password" class="min-w-0 flex-1 bg-transparent outline-none" type="password" autocomplete="current-password" :placeholder="t('auth.passwordPlaceholder')" :disabled="loading" />
                </div>
              </label>

              <label v-if="setupRequired" class="block">
                <span class="mb-2 block text-sm font-semibold">{{ t('auth.confirmPassword') }}</span>
                <div class="flex items-center gap-3 border px-4 py-3" style="border-color: var(--glass-border); background: var(--input-bg);">
                  <LockKeyhole :size="16" style="color: var(--text-muted);" />
                  <input v-model="confirmPassword" class="min-w-0 flex-1 bg-transparent outline-none" type="password" autocomplete="new-password" :placeholder="t('auth.repeatPasswordPlaceholder')" :disabled="loading" />
                </div>
              </label>

              <p v-if="localError" class="border px-4 py-3 text-sm font-medium" style="border-color: rgba(255,95,86,0.5); color: #fecaca; background: rgba(255,95,86,0.08);">
                {{ localError }}
              </p>

              <button class="btn btn-primary mt-3 w-full" type="submit" :disabled="loading">
                {{ submitLabel }}
              </button>
            </form>
          </div>
        </section>
      </div>
  </div>
</template>
