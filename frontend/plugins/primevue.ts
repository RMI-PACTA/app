import { defineNuxtPlugin } from '#app'

import PrimeVue from 'primevue/config'

import Accordion from 'primevue/accordion'
import AccordionTab from 'primevue/accordiontab'
import Button from 'primevue/button'
import Calendar from 'primevue/calendar'
import Card from 'primevue/card'
import Chips from 'primevue/chips'
import Column from 'primevue/column'
import ConfirmDialog from 'primevue/confirmdialog'
import ConfirmationService from 'primevue/confirmationservice'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import Dropdown from 'primevue/dropdown'
import FileUpload from 'primevue/fileupload'
import Image from 'primevue/image'
import InlineMessage from 'primevue/inlinemessage'
import InputNumber from 'primevue/inputnumber'
import InputSwitch from 'primevue/inputswitch'
import InputText from 'primevue/inputtext'
import Message from 'primevue/message'
import Menu from 'primevue/menu'
import MultiSelect from 'primevue/multiselect'
import OverlayPanel from 'primevue/overlaypanel'
import ProgressSpinner from 'primevue/progressspinner'
import TabPanel from 'primevue/tabpanel'
import TabView from 'primevue/tabview'
import Textarea from 'primevue/textarea'
import Toast from 'primevue/toast'
import ToastService from 'primevue/toastservice'
import Tooltip from 'primevue/tooltip'
import TriStateCheckbox from 'primevue/tristatecheckbox'

export default defineNuxtPlugin(({ vueApp }) => {
  vueApp.use(PrimeVue)

  vueApp.use(ToastService)

  vueApp.component('PVAccordion', Accordion)
  vueApp.component('PVAccordionTab', AccordionTab)
  vueApp.component('PVButton', Button)
  vueApp.component('PVCalendar', Calendar)
  vueApp.component('PVCard', Card)
  vueApp.component('PVChips', Chips)
  vueApp.component('PVColumn', Column)
  vueApp.component('PVConfirmDialog', ConfirmDialog)
  vueApp.component('PVDataTable', DataTable)
  vueApp.component('PVDialog', Dialog)
  vueApp.component('PVDropdown', Dropdown)
  vueApp.component('PVFileUpload', FileUpload)
  vueApp.component('PVImage', Image)
  vueApp.component('PVInlineMessage', InlineMessage)
  vueApp.component('PVInputNumber', InputNumber)
  vueApp.component('PVInputSwitch', InputSwitch)
  vueApp.component('PVInputText', InputText)
  vueApp.component('PVMessage', Message)
  vueApp.component('PVMenu', Menu)
  vueApp.component('PVMultiSelect', MultiSelect)
  vueApp.component('PVOverlayPanel', OverlayPanel)
  vueApp.component('PVProgressSpinner', ProgressSpinner)
  vueApp.component('PVTabPanel', TabPanel)
  vueApp.component('PVTabView', TabView)
  vueApp.component('PVTextarea', Textarea)
  vueApp.component('PVToast', Toast)
  vueApp.component('PVTriStateCheckbox', TriStateCheckbox)

  vueApp.directive('tooltip', Tooltip)
  vueApp.use(ConfirmationService)
})
