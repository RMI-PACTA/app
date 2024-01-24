<script setup lang="ts" />

<template>
  <div class="modal-group">
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
