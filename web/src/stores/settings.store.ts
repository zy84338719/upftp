import { defineStore } from 'pinia'
import { ref } from 'vue'
import i18n from '@/i18n'
import type { LocaleType, AppConfig } from '@/types'
import { DEFAULT_LANGUAGE, LOCAL_STORAGE_KEYS } from '@/constants'
import {
  setLang as apiSetLang,
} from '@/api'
import { copyLink } from '@/utils'

export const useSettingsStore = defineStore('settings', () => {
  const language = ref<LocaleType>((localStorage.getItem(LOCAL_STORAGE_KEYS.language) as LocaleType) || DEFAULT_LANGUAGE)
  const httpAuthOn = ref(false)
  const httpAuthUser = ref('')
  const httpAuthPass = ref('')
  const securityKey = ref('')

  async function setLang(lang: string) {
    const locale = lang as LocaleType
    language.value = locale
    localStorage.setItem(LOCAL_STORAGE_KEYS.language, lang)
    i18n.global.locale.value = locale
    await apiSetLang(lang)
  }

  function initConfig(config: AppConfig) {
    language.value = (config.language as LocaleType) || DEFAULT_LANGUAGE
    httpAuthOn.value = config.httpAuthOn || false
    httpAuthUser.value = config.httpAuthUser || ''
    httpAuthPass.value = config.httpAuthPass || ''
    i18n.global.locale.value = language.value
  }

  return {
    language,
    httpAuthOn,
    httpAuthUser,
    httpAuthPass,
    securityKey,
    setLang,
    copyLink,
    initConfig
  }
})
