import { defineNuxtPlugin } from '#app'

import PrimeVue from 'primevue/config'

import Accordion from 'primevue/accordion'
import AccordionTab from 'primevue/accordiontab'
import Button from 'primevue/button'
import Calendar from 'primevue/calendar'
import Card from 'primevue/card'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import Dropdown from 'primevue/dropdown'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import InputSwitch from 'primevue/inputswitch'
import FileUpload from 'primevue/fileupload'
import Textarea from 'primevue/textarea'
import Tooltip from 'primevue/tooltip'
import Message from 'primevue/message'
import OverlayPanel from 'primevue/overlaypanel'
import ProgressSpinner from 'primevue/progressspinner'
import ToastService from 'primevue/toastservice'
import Toast from 'primevue/toast'

export default defineNuxtPlugin(({ vueApp }) => {
  vueApp.use(PrimeVue)

  vueApp.use(ToastService)

  vueApp.component('PVAccordion', Accordion)
  vueApp.component('PVAccordionTab', AccordionTab)
  vueApp.component('PVButton', Button)
  vueApp.component('PVCalendar', Calendar)
  vueApp.component('PVCard', Card)
  vueApp.component('PVColumn', Column)
  vueApp.component('PVDataTable', DataTable)
  vueApp.component('PVDialog', Dialog)
  vueApp.component('PVDropdown', Dropdown)
  vueApp.component('PVFileUpload', FileUpload)
  vueApp.component('PVInputNumber', InputNumber)
  vueApp.component('PVInputText', InputText)
  vueApp.component('PVInputSwitch', InputSwitch)
  vueApp.component('PVMessage', Message)
  vueApp.component('PVOverlayPanel', OverlayPanel)
  vueApp.component('PVProgressSpinner', ProgressSpinner)
  vueApp.component('PVTextarea', Textarea)
  vueApp.component('PVToast', Toast)

  vueApp.directive('tooltip', Tooltip)
})
