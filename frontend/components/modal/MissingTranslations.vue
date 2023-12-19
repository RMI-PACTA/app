<script setup lang="ts">
const { $axios, $missingTranslations } = useNuxtApp()
const { loading: { withLoading }, missingTranslations: { missingTranslationsVisible } } = useModal()
const { t } = useI18n()

const prefix = 'components/modal/MissingTranslations'
const existing = useState<Map<string, string>>(`${prefix}.existing`, () => new Map())
const languages = ['en', 'es', 'fr', 'de']
const tt = (key: string) => t(`${prefix}.${key}`)

const onOpen = async () => {
  await withLoading(async () => {
    await Promise.all(languages.map((lang) => $axios({
      method: 'get',
      url: `/_nuxt/lang/${lang}.json`,
      transformResponse: (res) => res,
      responseType: 'json',
    }).then((r) => {
      existing.value.set(lang, r.data)
    })))
  }, `${prefix}.onOpen`)
}

const constructIdealMap = (lang: string): { map: Map<string, Map<string, string>>, errors: string[], missing: Set<string> } => {
  const errors: string[] = []
  const ideal = new Map<string, Map<string, string>>()

  const e = existing.value.get(lang)
  if (e) {
    let asObj: Map<string, Map<string, string>>
    try {
      asObj = JSON.parse(e) as Map<string, Map<string, string>>
    } catch (e: any) {
      errors.push(`Error parsing ${lang} JSON: ${e}`)
      return { map: ideal, errors, missing: new Set<string>() }
    }
    for (const [prefix, values] of Object.entries(asObj)) {
      let m = ideal.get(prefix)
      if (!m) {
        m = new Map<string, string>()
        ideal.set(prefix, m)
      }
      for (const [key, value] of Object.entries(values as Map<string, string>)) {
        if (m.has(key)) {
          errors.push(`Duplicate key ${key} in ${prefix}`)
        }
        m.set(key, value)
      }
    }
  }
  const missing = $missingTranslations.values.value.get(lang) ?? new Set<string>()
  for (const key of missing) {
    const splits = key.split('.')
    if (splits.length < 2) {
      errors.push(`Invalid key structure '${key}'`)
    }
    const file = splits[0]
    const actualKey = splits.slice(1).join('.')
    let m = ideal.get(file)
    if (!m) {
      m = new Map<string, string>()
      ideal.set(file, m)
    }
    if (m.has(actualKey)) {
      errors.push(`Duplicate key ${actualKey} in ${file}`)
    }
    ideal.set(file, m.set(actualKey, `TODO - ${actualKey}`))
  }

  return { map: ideal, errors, missing }
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
  errors: string[]
  missing: Set<string>
}
const tabs = computed<TabValue[]>(() => {
  const result: TabValue[] = []
  for (const lang of languages) {
    const { map, errors, missing } = constructIdealMap(lang)
    result.push({
      language: lang,
      ideal: mapToJson(map),
      numMissing: ($missingTranslations.values.value.get(lang) ?? new Set()).size,
      errors,
      missing,
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
    @opened="onOpen"
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
          <StandardDebug
            :label="`Errors (${tab.errors.length})`"
            :value="tab.errors"
            always
          />
          <StandardDebug
            :label="`Missing ${tab.numMissing} Translation Strings`"
            :value="Array.from(tab.missing).sort()"
            always
          />
          <div class="code">
            {{ tab.ideal }}
          </div>
        </div>
      </PVTabPanel>
    </PVTabView>
    <StandardDebug
      :value="languages"
      label="Languages"
    />
    <StandardDebug
      :value="constructIdealMap('en')"
      label="EN Ideal"
    />
    <StandardDebug
      :value="tabs"
      label="Tabs"
    />
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
