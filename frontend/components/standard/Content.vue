<template>
  <div class="standard-content flex flex-column gap-2 py-3 md:py-5 align-items-start">
    <slot />
  </div>
</template>

<style lang="scss">
.standard-content {
  h1,
  h2,
  h3,
  h4,
  h5,
  h6 {
    color: #003b63;
  }

  h1,
  h1.title-bar {
    margin-bottom: 1rem;
  }

  p {
    margin: 0;
    line-height: 1.5rem;
    color: rgb(0 0 0 / 85%);
  }

  a {
    font-weight: 600;
    color: rgb(0 0 0 / 85%);
  }

  & > .p-button {
    margin: 0.5rem 0;
  }
}
</style>
