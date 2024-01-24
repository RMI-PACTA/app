<template>
  <div
    class="standard-full-width"
  >
    <slot />
  </div>
</template>

<style lang="scss">
.standard-full-width {
  width: calc(100vw - 1rem);
  margin: 0 calc(50% - 50vw + .5rem);
  display: flex;
  align-items: center;
  justify-content: center;
}

.standard-full-width > * {
  width: 100%;
}
</style>
