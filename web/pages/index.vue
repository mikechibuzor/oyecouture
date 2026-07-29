<script setup lang="ts">
import {
  ASYNC_REQUEST_STATUS,
  FOUNDATION_PAGE_COPY,
} from '~/constants/foundation'
import { ServiceHealthStatus } from '~/types/health'

const { data, error, refresh, status } = await useApiHealth()

const apiConnected = computed(
  () => data.value?.services.api === ServiceHealthStatus.Available,
)
const databaseConnected = computed(
  () => data.value?.services.database === ServiceHealthStatus.Available,
)
const isChecking = computed(
  () => status.value === ASYNC_REQUEST_STATUS.pending,
)
const apiStatusLabel = computed(() =>
  apiConnected.value
    ? FOUNDATION_PAGE_COPY.connected
    : FOUNDATION_PAGE_COPY.unavailable,
)
const databaseStatusLabel = computed(() =>
  databaseConnected.value
    ? FOUNDATION_PAGE_COPY.connected
    : FOUNDATION_PAGE_COPY.unavailable,
)

const refreshHealth = async (): Promise<void> => {
  await refresh()
}
</script>

<template>
  <main class="foundation-page">
    <header class="foundation-header">
      <span class="brand-mark">{{ FOUNDATION_PAGE_COPY.brandMark }}</span>
      <div>
        <p class="eyebrow">{{ FOUNDATION_PAGE_COPY.brand }}</p>
        <h1>{{ FOUNDATION_PAGE_COPY.heading }}</h1>
      </div>
    </header>

    <section class="service-status" aria-labelledby="service-status-title">
      <h2 id="service-status-title">{{ FOUNDATION_PAGE_COPY.servicesHeading }}</h2>

      <dl>
        <div>
          <dt>{{ FOUNDATION_PAGE_COPY.webService }}</dt>
          <dd class="status-ok">{{ FOUNDATION_PAGE_COPY.connected }}</dd>
        </div>
        <div>
          <dt>{{ FOUNDATION_PAGE_COPY.apiService }}</dt>
          <dd :class="apiConnected ? 'status-ok' : 'status-error'">
            {{ apiStatusLabel }}
          </dd>
        </div>
        <div>
          <dt>{{ FOUNDATION_PAGE_COPY.databaseService }}</dt>
          <dd :class="databaseConnected ? 'status-ok' : 'status-error'">
            {{ databaseStatusLabel }}
          </dd>
        </div>
      </dl>

      <p v-if="error" class="connection-error" role="status">
        {{ error.message }}
      </p>

      <button type="button" :disabled="isChecking" @click="refreshHealth">
        {{ isChecking ? FOUNDATION_PAGE_COPY.checking : FOUNDATION_PAGE_COPY.checkAgain }}
      </button>
    </section>
  </main>
</template>
