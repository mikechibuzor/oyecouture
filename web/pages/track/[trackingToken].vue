<script setup lang="ts">
import { ExternalLink, PackageCheck } from 'lucide'
import {
  FULFILLMENT_PATHS,
  FULFILLMENT_STAGE_LABELS,
  TRACKING_COPY,
} from '~/constants/tracking'
import { FulfillmentStage } from '~/types/tracking'
import type { FulfillmentEvent } from '~/types/tracking'
import { formatOrderDate } from '~/utils/admin'

const route = useRoute()
const { fetchOrderTracking } = useOrderTracking()

const trackingToken = computed(() => {
  const rawToken = route.params.trackingToken
  return Array.isArray(rawToken) ? rawToken[0] ?? '' : rawToken ?? ''
})

const {
  data: tracking,
  error: trackingError,
  pending: isLoading,
  refresh,
} = await useAsyncData(
  `order-tracking-${trackingToken.value}`,
  () => fetchOrderTracking(trackingToken.value),
)

const stages = computed(() => {
  if (!tracking.value) {
    return []
  }

  const path = FULFILLMENT_PATHS[tracking.value.kind]
  // Out for delivery is an optional local-rider step that is deliberately not
  // on the customer timeline, so it is read as dispatched. Without this,
  // indexOf returns -1 and every step renders incomplete, showing a paid and
  // shipped order as though nothing had happened to it yet.
  const timelineStage =
    tracking.value.currentStage === FulfillmentStage.OutForDelivery
      ? FulfillmentStage.Dispatched
      : tracking.value.currentStage
  const currentIndex = path.indexOf(timelineStage)
  const eventsByStage = new Map(
    tracking.value.events.map((event) => [event.stage, event] as const),
  )

  return path.map((stage, stageIndex) => ({
    stage,
    label: FULFILLMENT_STAGE_LABELS[stage],
    event: eventsByStage.get(stage),
    isComplete: stageIndex <= currentIndex,
    isCurrent: stage === tracking.value?.currentStage,
  }))
})

const hasShipmentDetails = computed(() =>
  Boolean(
    tracking.value?.shipment.carrierName ||
      tracking.value?.shipment.carrierTrackingNumber ||
      tracking.value?.shipment.carrierTrackingUrl,
  ),
)

useSeoMeta({
  title: () => tracking.value
    ? `${tracking.value.orderNo} progress | Oyecouture`
    : `${TRACKING_COPY.heading} | Oyecouture`,
  robots: 'noindex, nofollow',
})

const retryLoad = (): void => {
  void refresh()
}

const eventDate = (event?: FulfillmentEvent): string =>
  event ? formatOrderDate(event.occurredAt) : ''
</script>

<template>
  <div class="tracking-page">
    <header class="tracking-header">
      <div>
        <p class="tracking-eyebrow">{{ TRACKING_COPY.eyebrow }}</p>
        <h1>{{ tracking?.orderNo ?? TRACKING_COPY.heading }}</h1>
      </div>
      <StorefrontIcon :icon-node="PackageCheck" :size="34" :stroke-width="1.35" />
    </header>

    <p v-if="isLoading" class="tracking-message">{{ TRACKING_COPY.loading }}</p>

    <section v-else-if="trackingError || !tracking" class="tracking-message" role="alert">
      <p>{{ TRACKING_COPY.loadFailed }}</p>
      <button type="button" class="tracking-retry" @click="retryLoad">
        {{ TRACKING_COPY.retry }}
      </button>
    </section>

    <template v-else>
      <section class="tracking-summary" aria-label="Order summary">
        <div>
          <span>{{ TRACKING_COPY.currentStatus }}</span>
          <strong>{{ tracking.currentStageLabel }}</strong>
        </div>
        <div v-if="tracking.estimatedDeliveryDate">
          <span>{{ TRACKING_COPY.estimatedDelivery }}</span>
          <strong>
            {{
              new Intl.DateTimeFormat('en-NG', {
                day: 'numeric',
                month: 'long',
                year: 'numeric',
              }).format(new Date(tracking.estimatedDeliveryDate))
            }}
          </strong>
        </div>
      </section>

      <div class="tracking-layout">
        <section class="tracking-timeline-section">
          <p class="tracking-section-label">{{ TRACKING_COPY.eyebrow }}</p>
          <ol class="tracking-timeline">
            <li
              v-for="stage in stages"
              :key="stage.stage"
              :class="{ 'is-complete': stage.isComplete, 'is-current': stage.isCurrent }"
            >
              <span class="tracking-marker" aria-hidden="true" />
              <div>
                <strong>{{ stage.label }}</strong>
                <time v-if="stage.event" :datetime="stage.event.occurredAt">
                  {{ eventDate(stage.event) }}
                </time>
                <p v-if="stage.event?.publicNote">{{ stage.event.publicNote }}</p>
                <p v-else-if="stage.isCurrent">{{ TRACKING_COPY.awaitingUpdate }}</p>
              </div>
            </li>
          </ol>
        </section>

        <aside class="tracking-details">
          <section>
            <p class="tracking-section-label">{{ TRACKING_COPY.items }}</p>
            <ul class="tracking-items">
              <li v-for="item in tracking.items" :key="`${item.name}-${item.size}`">
                <span>{{ item.name }}<small v-if="item.size">{{ item.size }}</small></span>
                <strong>x{{ item.quantity }}</strong>
              </li>
            </ul>
          </section>

          <section v-if="hasShipmentDetails">
            <p class="tracking-section-label">{{ TRACKING_COPY.shipment }}</p>
            <dl class="tracking-shipment">
              <div v-if="tracking.shipment.carrierName">
                <dt>{{ TRACKING_COPY.carrier }}</dt>
                <dd>{{ tracking.shipment.carrierName }}</dd>
              </div>
              <div v-if="tracking.shipment.carrierTrackingNumber">
                <dt>{{ TRACKING_COPY.trackingNumber }}</dt>
                <dd>{{ tracking.shipment.carrierTrackingNumber }}</dd>
              </div>
            </dl>
            <a
              v-if="tracking.shipment.carrierTrackingUrl"
              class="tracking-courier-link"
              :href="tracking.shipment.carrierTrackingUrl"
              target="_blank"
              rel="noreferrer"
            >
              {{ TRACKING_COPY.courierLink }}
              <StorefrontIcon :icon-node="ExternalLink" :size="14" :stroke-width="1.75" />
            </a>
          </section>
        </aside>
      </div>
    </template>
  </div>
</template>
