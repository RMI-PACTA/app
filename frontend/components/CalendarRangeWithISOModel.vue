<script setup lang="ts">
interface Props {
  value: Array<string | undefined> | undefined
}
const props = defineProps<Props>()
interface Emits {
  (event: 'update:value', value: Array<string | undefined> | undefined): void
}
const emit = defineEmits<Emits>()

const decode = (value: string | null | undefined): Date | null => {
  if (value === null || value === undefined || value === '') {
    return null
  }
  return new Date(value)
}
const encode = (value: Date | null | undefined): string => {
  if (value === null || value === undefined) {
    return ''
  }
  return value.toISOString()
}
const addDay = (date: Date | null): Date | null => {
  if (date === null) {
    return null
  }
  return new Date(date.getTime() + 24 * 60 * 60 * 1000)
}
const subDay = (date: Date | null): Date | null => {
  if (date === null) {
    return null
  }
  return new Date(date.getTime() - 24 * 60 * 60 * 1000)
}

// We use the addDay/subDay to make sure end-dates selected on the calendar are
// included in the range for search. Otherwise, the filter will exclude the last day,
// since the Calendar component selects the first millisecond of the date given.
const model = computed<Array<Date | null>>({
  get: () => props.value
    ? [
        decode(props.value[0]),
        subDay(decode(props.value[1])),
      ]
    : [new Date(), new Date()],
  set: (value: Array<Date | null>) => {
    emit('update:value', [
      encode(value[0]),
      encode(addDay(value[1])),
    ])
  },
})
</script>

<template>
  <PVCalendar
    v-model="model"
    inline
    selection-mode="range"
    :manual-input="false"
  />
</template>
