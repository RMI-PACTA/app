<script setup lang="ts" />

<template>
  <div>
    <ModalError />
    <ModalLoading />
    <ModalPermissionDenied />
    <ModalFakeUsers />
    <!-- TODO(grady) Remove This in Production Environments -->
    <ModalMissingTranslations />
    <PVToast />
    <PVConfirmDialog
      class="w-30rem"
      style="max-width: calc(100vw - 1rem);"
    />
  </div>
</template>
