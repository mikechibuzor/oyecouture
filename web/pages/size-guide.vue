<script setup lang="ts">
import { SIZE_CHARTS } from '~/constants/sizeGuide'
import { StorefrontRoute } from '~/types/storefront'

const { appointmentWhatsAppUrl, opensWhatsApp } = await useAppointmentWhatsApp()

useStorefrontSeo({
  title: 'Size Guide',
  description: 'Oyecouture ready-to-wear size charts and made-to-measure guidance.',
})
</script>

<template>
  <div class="size-guide-page">
    <header class="size-guide-hero">
      <div class="size-guide-hero-inner">
        <div class="size-guide-hero-copy">
          <p class="eyebrow">The right fit</p>
          <h1>Custom size guide</h1>
          <p>Find your ready-to-wear size or prepare measurements for a piece cut entirely to you.</p>
        </div>
        <img
          src="/images/measurements/measuring-tape.webp"
          alt="Flexible measuring tape"
          width="950"
          height="865"
        >
      </div>
    </header>

    <section class="size-guide-intro">
      <div>
        <p class="eyebrow">How to read the chart</p>
        <h2>Measure your body, not your clothes</h2>
        <p>
          Keep the tape level and comfortably close to the body. All chart measurements
          are in centimetres. When you sit between sizes, choose the larger size.
        </p>
      </div>
      <img
        src="/images/measurements/size-figure.webp"
        alt="Front view of the body for taking clothing measurements"
        width="382"
        height="1200"
        loading="lazy"
      >
    </section>

    <section class="size-chart-section" aria-labelledby="ready-to-wear-title">
      <div class="size-chart-heading">
        <p class="eyebrow">Ready to wear</p>
        <h2 id="ready-to-wear-title">Choose your size</h2>
        <span>All body measurements in cm</span>
      </div>

      <article v-for="chart in SIZE_CHARTS" :key="chart.id" class="size-chart-block">
        <h3>{{ chart.title }}</h3>
        <div class="size-chart-scroll">
          <table>
            <thead>
              <tr>
                <th scope="col">Size</th>
                <th scope="col">UK size</th>
                <th scope="col">{{ chart.primaryLabel }}</th>
                <th scope="col">{{ chart.secondaryLabel }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in chart.rows" :key="row.size">
                <th scope="row">{{ row.size }}</th>
                <td>{{ row.ukSize }}</td>
                <td>{{ row.primary }}</td>
                <td>{{ row.secondary }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </article>
    </section>

    <section class="size-guide-next">
      <div>
        <p class="eyebrow">Made to measure</p>
        <h2>A single measurement — Yours</h2>
        <p>Follow our visual guide, record each measurement, then keep or print your completed worksheet.</p>
      </div>
      <div class="size-guide-actions">
        <NuxtLink :to="StorefrontRoute.Measurements">Open measurement guide</NuxtLink>
        <a
          :href="appointmentWhatsAppUrl"
          :target="opensWhatsApp ? '_blank' : undefined"
          :rel="opensWhatsApp ? 'noopener' : undefined"
        >Book a fitting</a>
      </div>
    </section>
  </div>
</template>
