<script setup lang="ts">
import { ArrowLeft, Check, Copy, Printer } from 'lucide'
import {
  MEASUREMENT_STEPS,
  MEASUREMENT_WORKSHEETS,
} from '~/constants/sizeGuide'
import type { MeasurementProfile } from '~/types/order'
import { StorefrontRoute } from '~/types/storefront'
import { fieldsForMeasurementProfile, serializeMeasurementClipboard } from '~/utils/measurements'

const route = useRoute()
const requestedProfile = computed<MeasurementProfile>(() => {
  const profile = Array.isArray(route.query.profile) ? route.query.profile[0] : route.query.profile
  return profile === 'suit' || profile === 'native' || profile === 'mixed' ? profile : 'native'
})
const activeWorksheetId = ref<'suit' | 'native'>(requestedProfile.value === 'suit' ? 'suit' : 'native')
const measurementValues = reactive<Record<string, string>>({})
const copyStatus = ref<{ message: string, status: 'success' | 'error' } | null>(null)

const activeWorksheet = computed(() =>
  MEASUREMENT_WORKSHEETS.find((worksheet) => worksheet.id === activeWorksheetId.value)
  ?? MEASUREMENT_WORKSHEETS[0],
)

const selectWorksheet = (worksheetId: 'suit' | 'native'): void => {
  activeWorksheetId.value = worksheetId
  copyStatus.value = null
}

const printWorksheet = (): void => {
  window.print()
}

const copyProfile = computed<MeasurementProfile>(() =>
  requestedProfile.value === 'mixed' ? 'mixed' : activeWorksheetId.value,
)
const returnTo = computed(() => {
  const value = Array.isArray(route.query.returnTo) ? route.query.returnTo[0] : route.query.returnTo
  return typeof value === 'string' && value.startsWith('/') ? `${value}#measurements` : ''
})

const copyMeasurements = async (): Promise<void> => {
  const values = Object.fromEntries(
    fieldsForMeasurementProfile(copyProfile.value)
      .map((field) => [field.id, Number(measurementValues[field.id])] as const)
      .filter((entry) => Number.isFinite(entry[1]) && entry[1] > 0),
  )

  if (Object.keys(values).length === 0) {
    copyStatus.value = { message: 'Add at least one measurement before copying.', status: 'error' }
    return
  }

  try {
    await navigator.clipboard.writeText(serializeMeasurementClipboard(copyProfile.value, values))
    copyStatus.value = {
      message: 'Measurements copied. Return to your order and choose Paste measurements.',
      status: 'success',
    }
  } catch {
    copyStatus.value = {
      message: 'Clipboard access is unavailable. Please try again in a secure browser window.',
      status: 'error',
    }
  }
}

useSeoMeta({
  title: 'Measurement Guide — Oyecouture Bespoke',
  description: 'Follow Oyecouture’s visual measurement guide and complete a private worksheet.',
})
</script>

<template>
  <div class="measurement-page">
    <header class="measurement-hero">
      <div class="measurement-hero-content">
        <nav class="measurement-breadcrumbs" aria-label="Breadcrumb">
          <NuxtLink :to="StorefrontRoute.SizeGuide">Size guide</NuxtLink>
          <span aria-hidden="true">/</span>
          <span>Measurements</span>
        </nav>
        <p class="eyebrow">Cut to a single measure</p>
        <h1>How to measure</h1>
        <p>Use a flexible tape, stand naturally and ask someone to help where possible. Record every value in inches.</p>
      </div>
    </header>

    <section class="measurement-steps" aria-labelledby="measurement-steps-title">
      <div class="measurement-section-heading">
        <p class="eyebrow">01 — The guide</p>
        <h2 id="measurement-steps-title">Take each measurement once, carefully</h2>
      </div>

      <ol class="measurement-step-grid">
        <li v-for="(step, stepIndex) in MEASUREMENT_STEPS" :key="step.id">
          <img
            :src="step.imagePath"
            :alt="step.imageAlt"
            width="960"
            height="960"
            loading="lazy"
          >
          <div>
            <span>{{ String(stepIndex + 1).padStart(2, '0') }}</span>
            <h3>{{ step.title }}</h3>
            <p>{{ step.instruction }}</p>
            <small v-if="step.allowance">{{ step.allowance }}</small>
          </div>
        </li>
      </ol>
    </section>

    <section id="worksheet" class="measurement-worksheet" aria-labelledby="measurement-worksheet-title">
      <div class="measurement-section-heading">
        <p class="eyebrow">02 — Your worksheet</p>
        <h2 id="measurement-worksheet-title">Record and keep</h2>
        <p>Values are in inches. The worksheet remains in this browser and is not submitted anywhere.</p>
      </div>

      <div class="measurement-workbook">
        <div class="measurement-workbook-tabs" role="tablist" aria-label="Garment type">
          <button
            v-for="worksheet in MEASUREMENT_WORKSHEETS"
            :key="worksheet.id"
            type="button"
            role="tab"
            :aria-selected="activeWorksheetId === worksheet.id"
            :class="{ 'is-active': activeWorksheetId === worksheet.id }"
            @click="selectWorksheet(worksheet.id)"
          >
            {{ worksheet.label }}
          </button>
        </div>

        <div class="measurement-input-grid">
          <label v-for="field in activeWorksheet.fields" :key="field.id">
            <span>{{ field.label }}</span>
            <span class="measurement-input-shell">
              <input
                v-model.trim="measurementValues[field.id]"
                type="number"
                min="0"
                step="0.25"
                inputmode="decimal"
                :aria-label="`${field.label} in inches`"
              >
              <small>in</small>
            </span>
          </label>
        </div>

        <div class="measurement-workbook-actions">
          <button class="measurement-copy-action" type="button" @click="copyMeasurements">
            <StorefrontIcon :icon-node="copyStatus?.status === 'success' ? Check : Copy" :size="16" :stroke-width="1.6" />
            Copy all measurements
          </button>
          <button class="measurement-print-action" type="button" @click="printWorksheet">
            <StorefrontIcon :icon-node="Printer" :size="16" :stroke-width="1.6" />
            Print or save worksheet
          </button>
          <NuxtLink v-if="returnTo" class="measurement-return-action" :to="returnTo">
            <StorefrontIcon :icon-node="ArrowLeft" :size="16" :stroke-width="1.6" />
            Return to confirmation
          </NuxtLink>
        </div>

        <p
          v-if="copyStatus"
          class="measurement-copy-status"
          :data-status="copyStatus.status"
          role="status"
        >
          {{ copyStatus.message }}
        </p>
      </div>
    </section>
  </div>
</template>
