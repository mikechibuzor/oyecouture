<script setup lang="ts">
import {
  ADMIN_AUDIT_RESOURCE,
  ADMIN_AUDIT_SEGMENT,
  ADMIN_COPY,
  ADMIN_PILL_TONE,
  ADMIN_ROUTES,
} from '~/constants/admin'
import { HTTP_STATUS, HttpMethod } from '~/types/http'
import type { AuditEntry } from '~/types/team'
import { formatOrderDate } from '~/utils/admin'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const router = useRouter()
const { session } = useAdminSession()
const { fetchAudit } = useAdminAudit()

const entries = ref<readonly AuditEntry[]>([])
const total = ref(0)
const isLoading = ref(true)
const hasFailed = ref(false)
const pageSize = 50
const offset = ref(0)

useHead({ title: `${ADMIN_COPY.auditHeading} · ${ADMIN_COPY.brand}` })

const canGoBack = computed(() => offset.value > 0)
const canGoForward = computed(() => offset.value + pageSize < total.value)

// The log stores the request that was made, which is precise but reads as
// plumbing. The owner is not here to audit HTTP verbs, so each entry is
// described in the terms they think in. Anything unrecognised falls back to the
// raw method and path rather than being hidden, since an entry nobody can
// account for is exactly the one worth seeing.
const describe = (entry: AuditEntry): string => {
  const path = entry.path.replace(/^\/admin\//, '')
  const [resource, first, second] = path.split('/')
  const method = entry.method

  if (resource === ADMIN_AUDIT_RESOURCE.team) {
    if (method === HttpMethod.Post) return 'Added an admin'
    if (method === HttpMethod.Patch) return 'Changed an admin\u2019s role or access'
    if (method === HttpMethod.Delete) return 'Deleted an admin'
  }
  if (resource === ADMIN_AUDIT_RESOURCE.password) return 'Changed their own password'
  if (resource === ADMIN_AUDIT_RESOURCE.products) {
    if (second === ADMIN_AUDIT_SEGMENT.images) {
      if (method === HttpMethod.Post) return `Uploaded an image to ${first}`
      if (method === HttpMethod.Delete) return `Removed an image from ${first}`
      if (method === HttpMethod.Patch) return `Reordered images on ${first}`
    }
    if (second === ADMIN_AUDIT_SEGMENT.active) return `Archived or restored ${first}`
    if (method === HttpMethod.Post) return 'Created a product'
    if (method === HttpMethod.Put) return `Edited ${first}`
    if (method === HttpMethod.Delete) return `Deleted ${first}`
  }
  if (resource === ADMIN_AUDIT_RESOURCE.orders) {
    if (second === ADMIN_AUDIT_SEGMENT.status) return `Changed the status of ${first}`
    if (second === ADMIN_AUDIT_SEGMENT.tracking) return `Moved ${first} to its next stage`
  }
  if (resource === ADMIN_AUDIT_RESOURCE.deliveryZones) {
    if (method === HttpMethod.Post) return 'Added a delivery area'
    if (method === HttpMethod.Patch) return `Edited the ${first} delivery area`
  }
  if (resource === ADMIN_AUDIT_RESOURCE.settings) return 'Changed storefront settings'
  if (resource === ADMIN_AUDIT_RESOURCE.homeContent) return 'Edited the homepage content'

  return `${method} ${entry.path}`
}

// Grouped rather than shown as a number: the owner cares whether it took
// effect, not which status code said so.
const resultOf = (entry: AuditEntry): { label: string, tone: string } => {
  if (entry.statusCode >= HTTP_STATUS.ok && entry.statusCode < HTTP_STATUS.badRequest) {
    return { label: ADMIN_COPY.auditResultOk, tone: ADMIN_PILL_TONE.positive }
  }

  const refusals: readonly number[] = [
    HTTP_STATUS.unauthorized,
    HTTP_STATUS.forbidden,
    HTTP_STATUS.conflict,
  ]
  if (refusals.includes(entry.statusCode)) {
    return { label: ADMIN_COPY.auditResultRefused, tone: ADMIN_PILL_TONE.muted }
  }

  return { label: ADMIN_COPY.auditResultFailed, tone: ADMIN_PILL_TONE.muted }
}

const load = async (): Promise<void> => {
  isLoading.value = true
  hasFailed.value = false

  try {
    const page = await fetchAudit(pageSize, offset.value)
    entries.value = page.entries
    total.value = page.total
  } catch (loadError) {
    if ((loadError as { status?: number } | null)?.status === HTTP_STATUS.unauthorized) {
      session.value = null
      await router.push(ADMIN_ROUTES.login)
      return
    }

    hasFailed.value = true
  } finally {
    isLoading.value = false
  }
}

const goBack = async (): Promise<void> => {
  offset.value = Math.max(0, offset.value - pageSize)
  await load()
}

const goForward = async (): Promise<void> => {
  offset.value += pageSize
  await load()
}

onMounted(load)
</script>

<template>
  <div class="admin-page">
    <header class="admin-page-head">
      <div>
        <h1>{{ ADMIN_COPY.auditHeading }}</h1>
        <p class="admin-page-subheading">{{ ADMIN_COPY.auditSubheading }}</p>
      </div>
    </header>

    <p class="admin-audit-note">{{ ADMIN_COPY.auditNote }}</p>

    <p v-if="isLoading" class="admin-orders-message">{{ ADMIN_COPY.auditLoading }}</p>
    <p v-else-if="hasFailed" class="admin-orders-message" role="alert">
      {{ ADMIN_COPY.auditLoadFailed }}
    </p>
    <p v-else-if="entries.length === 0" class="admin-orders-message">
      {{ ADMIN_COPY.auditEmpty }}
    </p>

    <template v-else>
      <div class="admin-table-scroll">
        <table class="admin-table">
          <thead>
            <tr>
              <th scope="col">{{ ADMIN_COPY.auditWhenLabel }}</th>
              <th scope="col">{{ ADMIN_COPY.auditWhoLabel }}</th>
              <th scope="col">{{ ADMIN_COPY.auditActionLabel }}</th>
              <th scope="col">{{ ADMIN_COPY.auditResultLabel }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="entry in entries" :key="entry.id">
              <td :data-label="ADMIN_COPY.auditWhenLabel">
                <time :datetime="entry.occurredAt">{{ formatOrderDate(entry.occurredAt) }}</time>
              </td>
              <td :data-label="ADMIN_COPY.auditWhoLabel">{{ entry.actorEmail }}</td>
              <td :data-label="ADMIN_COPY.auditActionLabel">{{ describe(entry) }}</td>
              <td :data-label="ADMIN_COPY.auditResultLabel">
                <span class="admin-team-status" :class="resultOf(entry).tone">
                  {{ resultOf(entry).label }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <nav v-if="canGoBack || canGoForward" class="admin-pagination">
        <button
          type="button"
          class="admin-pagination-button"
          :disabled="!canGoBack"
          @click="goBack"
        >
          Previous
        </button>
        <button
          type="button"
          class="admin-pagination-button"
          :disabled="!canGoForward"
          @click="goForward"
        >
          Next
        </button>
      </nav>
    </template>
  </div>
</template>
