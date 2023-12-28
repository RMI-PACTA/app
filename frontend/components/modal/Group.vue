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
  </div>
</template>
