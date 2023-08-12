import { defineNuxtPlugin } from '#app'

import PrimeVue from 'primevue/config'

import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import FileUpload from 'primevue/fileupload'
import Textarea from 'primevue/textarea'
import Message from 'primevue/message'
import ProgressSpinner from 'primevue/progressspinner'

import ToastService from 'primevue/toastservice'

export default defineNuxtPlugin(({ vueApp }) => {
  vueApp.use(PrimeVue)

  vueApp.use(ToastService)

  vueApp.component('PVButton', Button)
  vueApp.component('PVFileUpload', FileUpload)
  vueApp.component('PVInputText', InputText)
  vueApp.component('PVMessage', Message)
  vueApp.component('PVProgressSpinner', ProgressSpinner)
  vueApp.component('PVTextarea', Textarea)
})
