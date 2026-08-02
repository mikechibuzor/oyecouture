<script setup lang="ts">
import {
  ADMIN_COPY,
  ADMIN_ORDERS_PAGE_SIZE,
  ADMIN_ROUTES,
  ADMIN_SEARCH_DEBOUNCE_MS,
  ORDER_KIND_LABELS,
  ORDER_KIND_OPTIONS,
  ORDER_STATUS_LABELS,
  ORDER_STATUS_OPTIONS,
} from '~/constants/admin'
import type { AdminOrderSummary } from '~/types/admin'
import { OrderStatus } from '~/types/admin'
import { OrderKind } from '~/types/order'
import { formatOrderDate } from '~/utils/admin'
import { formatProductPrice } from '~/utils/product'

definePageMeta({ layout: 'admin', middleware: 'admin' })

useHead({ title: `${ADMIN_COPY.ordersHeading} · ${ADMIN_COPY.brand}` })

const { fetchOrders } = useAdminOrders()
const { session } = useAdminSession()
const router = useRouter()

const searchTerm = ref('')
const statusFilter = ref<OrderStatus | ''>('')
const kindFilter = ref<OrderKind | ''>('')
const offset = ref(0)

const orders = ref<readonly AdminOrderSummary[]>([])
const total = ref(0)
const isLoading = ref(true)
const hasFailed = ref(false)

const hasFilters = computed(
  () => searchTerm.value !== '' || statusFilter.value !== '' || kindFilter.value !== '',
)

const rangeStart = computed(() => (total.value === 0 ? 0 : offset.value + 1))
const rangeEnd = computed(() => offset.value + orders.value.length)
const canGoBack = computed(() => offset.value > 0)
const canGoForward = computed(() => rangeEnd.value < total.value)

// A request counter, not a boolean: a slow early response must not overwrite a
// fresher one when the owner types quickly.
let latestRequestId = 0

const load = async (): Promise<void> => {
  const requestId = ++latestRequestId
  isLoading.value = true
  hasFailed.value = false

  try {
    const page = await fetchOrders({
      status: statusFilter.value,
      kind: kindFilter.value,
      q: searchTerm.value,
      limit: ADMIN_ORDERS_PAGE_SIZE,
      offset: offset.value,
    })

    if (requestId !== latestRequestId) {
      return
    }

    orders.value = page.orders
    total.value = page.total
  } catch (loadError) {
    if (requestId !== latestRequestId) {
      return
    }

    // A session that expired while the page was open should return the owner
    // to sign-in rather than showing a load failure they cannot act on.
    if ((loadError as { status?: number } | null)?.status === 401) {
      session.value = null
      await router.push(ADMIN_ROUTES.login)
      return
    }

    hasFailed.value = true
    orders.value = []
    total.value = 0
  } finally {
    if (requestId === latestRequestId) {
      isLoading.value = false
    }
  }
}

// Changing a filter always returns to the first page: otherwise a narrower
// result set would leave the owner stranded on an empty page.
let searchTimer: ReturnType<typeof setTimeout> | undefined

watch([statusFilter, kindFilter], () => {
  offset.value = 0
  void load()
})

watch(searchTerm, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    offset.value = 0
    void load()
  }, ADMIN_SEARCH_DEBOUNCE_MS)
})

onBeforeUnmount(() => clearTimeout(searchTimer))

const goToPreviousPage = (): void => {
  offset.value = Math.max(0, offset.value - ADMIN_ORDERS_PAGE_SIZE)
  void load()
}

const goToNextPage = (): void => {
  offset.value += ADMIN_ORDERS_PAGE_SIZE
  void load()
}

const clearFilters = (): void => {
  searchTerm.value = ''
  statusFilter.value = ''
  kindFilter.value = ''
}

onMounted(load)
</script>

<template>
  <section class="admin-orders">
    <header class="admin-orders-head">
      <h1>{{ ADMIN_COPY.ordersHeading }}</h1>
      <p v-if="total > 0" class="admin-orders-count">
        {{ rangeStart }}–{{ rangeEnd }} of {{ total }}
      </p>
    </header>

    <div class="admin-filters">
      <label class="admin-field admin-field-search">
        <span class="visually-hidden">{{ ADMIN_COPY.searchLabel }}</span>
        <input
          v-model.trim="searchTerm"
          type="search"
          :placeholder="ADMIN_COPY.searchPlaceholder"
        >
      </label>

      <label class="admin-field">
        <span class="visually-hidden">{{ ADMIN_COPY.statusLabel }}</span>
        <select v-model="statusFilter">
          <option value="">{{ ADMIN_COPY.anyStatus }}</option>
          <option v-for="status in ORDER_STATUS_OPTIONS" :key="status" :value="status">
            {{ ORDER_STATUS_LABELS[status] }}
          </option>
        </select>
      </label>

      <label class="admin-field">
        <span class="visually-hidden">{{ ADMIN_COPY.kindLabel }}</span>
        <select v-model="kindFilter">
          <option value="">{{ ADMIN_COPY.anyKind }}</option>
          <option v-for="kind in ORDER_KIND_OPTIONS" :key="kind" :value="kind">
            {{ ORDER_KIND_LABELS[kind] }}
          </option>
        </select>
      </label>
    </div>

    <p v-if="isLoading" class="admin-orders-message">{{ ADMIN_COPY.loading }}</p>

    <div v-else-if="hasFailed" class="admin-orders-message" role="alert">
      <p>{{ ADMIN_COPY.loadFailed }}</p>
      <button type="button" class="admin-inline-action" @click="load">
        {{ ADMIN_COPY.retry }}
      </button>
    </div>

    <div v-else-if="orders.length === 0" class="admin-orders-message">
      <p>{{ hasFilters ? ADMIN_COPY.emptyFiltered : ADMIN_COPY.empty }}</p>
      <button
        v-if="hasFilters"
        type="button"
        class="admin-inline-action"
        @click="clearFilters"
      >
        {{ ADMIN_COPY.clearFilters }}
      </button>
    </div>

    <template v-else>
      <div class="admin-table-scroll">
        <table class="admin-table">
          <thead>
            <tr>
              <th scope="col">Order</th>
              <th scope="col">Customer</th>
              <th scope="col">Type</th>
              <th scope="col">Status</th>
              <th scope="col" class="admin-table-number">Items</th>
              <th scope="col" class="admin-table-number">Total</th>
              <th scope="col">Placed</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="order in orders" :key="order.orderNo">
              <td class="admin-table-order-no">{{ order.orderNo }}</td>
              <td>
                <span class="admin-table-name">{{ order.customerName }}</span>
                <span class="admin-table-contact">{{ order.phone }}</span>
              </td>
              <td>{{ ORDER_KIND_LABELS[order.kind] }}</td>
              <td>
                <span class="admin-status" :data-status="order.status">
                  {{ ORDER_STATUS_LABELS[order.status] }}
                </span>
              </td>
              <td class="admin-table-number">{{ order.itemCount }}</td>
              <td class="admin-table-number">{{ formatProductPrice(order.totalKobo) }}</td>
              <td class="admin-table-placed">{{ formatOrderDate(order.createdAt) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <nav v-if="canGoBack || canGoForward" class="admin-pagination">
        <button
          type="button"
          class="admin-inline-action"
          :disabled="!canGoBack"
          @click="goToPreviousPage"
        >
          {{ ADMIN_COPY.previous }}
        </button>
        <button
          type="button"
          class="admin-inline-action"
          :disabled="!canGoForward"
          @click="goToNextPage"
        >
          {{ ADMIN_COPY.next }}
        </button>
      </nav>
    </template>
  </section>
</template>
