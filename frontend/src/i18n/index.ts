import { watch } from 'vue';
import { createI18n } from 'vue-i18n';
import { appSettings } from '../ui/settings';
import vi from './vi';
import en from './en';

export const i18n = createI18n({
  legacy: false,
  locale: appSettings.general.language,
  fallbackLocale: 'en',
  messages: {
    vi,
    en,
  },
});

watch(
  () => appSettings.general.language,
  (language) => {
    i18n.global.locale.value = language;
    document.documentElement.lang = language;
  },
  { immediate: true }
);
