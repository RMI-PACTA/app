<script setup lang="ts" />

<template>
  <div class="standard-content flex flex-column gap-3 pt-3 pb-4 md:pt-5 md:pb-6 align-items-start">
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

  a & :not(.p-button) {
    font-weight: 600;
    color: rgb(0 0 0 / 85%);
  }

  & > .p-button {
    margin: 0.5rem 0;
  }

  min-height: calc(100vh - 9rem - 4px);
  justify-content: center;
}
</style>
