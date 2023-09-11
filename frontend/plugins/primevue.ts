import { defineNuxtPlugin } from '#app'

import PrimeVue from 'primevue/config'

import Accordion from 'primevue/accordion'
import AccordionTab from 'primevue/accordiontab'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import FileUpload from 'primevue/fileupload'
import Textarea from 'primevue/textarea'
import Tooltip from 'primevue/tooltip'
import Message from 'primevue/message'
import ProgressSpinner from 'primevue/progressspinner'

import ToastService from 'primevue/toastservice'

export default defineNuxtPlugin(({ vueApp }) => {
  vueApp.use(PrimeVue)

  vueApp.use(ToastService)

  vueApp.component('PVAccordion', Accordion)
  vueApp.component('PVAccordionTab', AccordionTab)
  vueApp.component('PVButton', Button)
  vueApp.component('PVCard', Card)
  vueApp.component('PVColumn', Column)
  vueApp.component('PVDataTable', DataTable)
  vueApp.component('PVDialog', Dialog)
  vueApp.component('PVFileUpload', FileUpload)
  vueApp.component('PVInputText', InputText)
  vueApp.component('PVMessage', Message)
  vueApp.component('PVProgressSpinner', ProgressSpinner)
  vueApp.component('PVTextarea', Textarea)

  vueApp.directive('tooltip', Tooltip)
})
