<script setup lang="ts">
const { $axios, $missingTranslations } = useNuxtApp()
const { missingTranslations: { missingTranslationsVisible } } = useModal()
const { t } = useI18n()

const prefix = 'components/modal/MissingTranslations'
const existing = useState<Map<string, string>>(`${prefix}.existing`, () => new Map())
const languages = ['en', 'es', 'fr', 'de']
const tt = (key: string) => t(`${prefix}.${key}`)

onMounted(async () => {
  for (const lang of languages) {
    const langData = await $axios({
      method: 'get',
      url: `/_nuxt/lang/${lang}.json`,
    }).then((r) => {
      return JSON.stringify(r.data)
    })
    existing.value.set(lang, langData)
  }
})

const constructIdealMap = (lang: string): Map<string, Map<string, string>> => {
  const ideal = new Map<string, Map<string, string>>()

  const e = existing.value.get(lang)
  if (e) {
    const asObj = JSON.parse(e) as Map<string, Map<string, string>>
    for (const [prefix, values] of Object.entries(asObj)) {
      let m = ideal.get(prefix)
      if (!m) {
        m = new Map<string, string>()
        ideal.set(prefix, m)
      }
      for (const [key, value] of Object.entries(values as Map<string, string>)) {
        if (m.has(key)) {
          throw new Error(`Duplicate key ${key} in ${prefix}`)
        }
        m.set(key, value)
      }
    }
  }
  const missing = $missingTranslations.values.value.get(lang) ?? new Set<string>()
  for (const key of missing) {
    const splits = key.split('.')
    if (splits.length < 2) {
      throw new Error(`Invalid key structure '${key}'`)
    }
    const file = splits[0]
    const actualKey = splits.slice(1).join('.')
    let m = ideal.get(file)
    if (!m) {
      m = new Map<string, string>()
      ideal.set(file, m)
    }
    if (m.has(actualKey)) {
      throw new Error(`Duplicate key ${actualKey} in ${file}`)
    }
    ideal.set(file, m.set(actualKey, `TODO - ${actualKey}`))
  }

  return ideal
}
const mapToJson = (map: Map<string, Map<string, string>>): string => {
  const obj = Object.fromEntries(
    Array.from(map).map(([key, value]) => {
      const o = Object.fromEntries(value)
      const result = Object.keys(o).sort().reduce<Record<string, string>>((acc, key) => {
        acc[key] = o[key]
        return acc
      }, {})
      return [key, result]
    }),
  )
  const sortedObj = Object.keys(obj).sort().reduce<Record<string, Record<string, string>>>((acc, key) => {
    acc[key] = obj[key]
    return acc
  }, {})
  return JSON.stringify(sortedObj, null, 2)
}
interface TabValue {
  language: string
  ideal: string
  numMissing: number
}
const tabs = computed<TabValue[]>(() => {
  const result: TabValue[] = []
  for (const lang of languages) {
    result.push({
      language: lang,
      ideal: mapToJson(constructIdealMap(lang)),
      numMissing: ($missingTranslations.values.value.get(lang) ?? new Set()).size,
    })
  }
  return result
})
</script>

<template>
  <StandardModal
    v-model:visible="missingTranslationsVisible"
    header="Missing Translations"
    sub-header="A set of tools to help you find and fix missing translations"
  >
    <PVTabView>
      <PVTabPanel
        v-for="tab in tabs"
        :key="tab.language"
        :header="tab.language"
      >
        <div class="flex flex-column gap-2">
          <div class="flex gap-2">
            <CopyToClipboardButton
              :value="tab.ideal"
              :cta="tt('Copy to Clipboard')"
            />
            <DownloadButton
              :value="tab.ideal"
              :cta="tt('Download File')"
              :file-name="`${tab.language}-${new Date().getTime()}.json`"
            />
          </div>
          <div class="code">
            {{ tab.ideal }}
          </div>
        </div>
      </PVTabPanel>
    </PVTabView>
  </StandardModal>
</template>

<style scoped lang="scss">
.code {
    font-family: monospace;
    font-size: 0.9em;
    font-weight: 444;
    white-space: pre-wrap;
    background: #eee;
    padding: 0.5em;
}
</style>
