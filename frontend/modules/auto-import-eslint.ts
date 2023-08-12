// See https://github.com/nuxt/eslint-plugin-nuxt/issues/173#issuecomment-1552676382

import type { Import, Unimport } from 'unimport'
import { addTemplate, defineNuxtModule } from '@nuxt/kit'

const padding = ' '.repeat(4)
const autoImportEslint = defineNuxtModule({
  setup (_, nuxt) {
    const autoImports: Record<string, string[]> = {
      // global imports
      global: [
        '$fetch',
        'useCloneDeep',
        'defineNuxtConfig',
        'definePageMeta',
        'defineI18nConfig'
      ]
    }

    nuxt.hook('imports:context', async (context: Unimport) => {
      const imports = await context.getImports()
      imports.forEach(autoImport => {
        const list = autoImports[autoImport.from] ?? []
        const name = autoImport.as ?? autoImport.name
        autoImports[autoImport.from] = list.concat([name])
      })
    })

    nuxt.hook('imports:extend', (composableImport: Import[]) => {
      autoImports.composables = composableImport.map(autoImport => {
        if (autoImport.as !== undefined) {
          return autoImport.as.toString()
        }
        return autoImport.name.toString()
      })
    })

    nuxt.hook('modules:done', () => {
      const filename = '.eslint.globals.json'

      const getContents = (): string => {
        // To prevent formatter accidentally fix padding of template string
        let contents = ''
        contents += '{\n'
        contents += '  "globals": {'
        for (const autoImport in autoImports) {
          contents += `\n${padding}// ${autoImport}`
          autoImports[autoImport].forEach(imp => {
            contents += '\n'
            contents += `${padding}"${imp}": "readonly",`
          })
        }
        contents = `${contents.slice(0, -1)}\n`
        contents += '  }\n'
        contents += '}\n'

        return contents
      }

      addTemplate({
        filename,
        getContents,
        write: true
      })

      // console.log(`globals file is generated at ${fullPath}`)
    })
  }
})

export default autoImportEslint
