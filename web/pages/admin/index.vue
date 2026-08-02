<script setup lang="ts">
import { Search } from 'lucide'
import {
  ADMIN_COPY,
  ADMIN_DEFAULT_PERIOD,
  ADMIN_ORDERS_PAGE_SIZE,
  ADMIN_ORDER_FILTER_OPTIONS,
  ADMIN_ROUTES,
  ADMIN_SEARCH_DEBOUNCE_MS,
  ORDER_KIND_LABELS,
  ORDER_STATUS_LABELS,
} from '~/constants/admin'
import type {
  AdminOrderMetrics,
  AdminOrderSummary,
  MetricsPeriod,
} from '~/types/admin'
import { type AdminOrderFilterKey, type AdminOrderFilterValues, OrderStatus } from '~/types/admin'
import { OrderKind } from '~/types/order'
import { formatOrderDate } from '~/utils/admin'
import { formatProductPrice } from '~/utils/product'

definePageMeta({ layout: 'admin', middleware: 'admin' })

useHead({ title: `${ADMIN_COPY.ordersHeading} · ${ADMIN_COPY.brand}` })

const { fetchOrders, fetchMetrics } = useAdminOrders()
const { session } = useAdminSession()
const router = useRouter()

const searchTerm = ref('')
const statusFilter = ref<OrderStatus | ''>('')
const kindFilter = ref<OrderKind | ''>('')
const fromFilter = ref('')
const toFilter = ref('')
const offset = ref(0)

const orders = ref<readonly AdminOrderSummary[]>([])
const metrics = ref<AdminOrderMetrics | null>(null)
const metricsPeriod = ref<MetricsPeriod>(ADMIN_DEFAULT_PERIOD)
const total = ref(0)
const isLoading = ref(true)
const hasFailed = ref(false)

const hasFilters = computed(
  () =>
    searchTerm.value !== '' ||
    statusFilter.value !== '' ||
    kindFilter.value !== '' ||
    fromFilter.value !== '' ||
    toFilter.value !== '',
)

// The filter popover consumes strings so each filter category can share the
// same component contract.
const filterValues = computed<AdminOrderFilterValues>(() => ({
  status: statusFilter.value,
  kind: kindFilter.value,
  from: fromFilter.value,
  to: toFilter.value,
}))

// The visible row range reads as "1-10 of 30" in the table footer.
const rangeStart = computed(() => (total.value === 0 ? 0 : offset.value + 1))

// The visible range ends at the number of rows actually loaded on this page.
const rangeEnd = computed(() => offset.value + orders.value.length)

// Pagination controls disable themselves at the list boundaries.
const canGoBack = computed(() => offset.value > 0)
const canGoForward = computed(() => rangeEnd.value < total.value)

// Current page derives from offset so direct page jumps and next/previous agree.
const currentPage = computed(() => Math.floor(offset.value / ADMIN_ORDERS_PAGE_SIZE) + 1)

// Total pages bottoms out at one so the pagination area never divides by zero.
const totalPages = computed(() =>
  Math.max(1, Math.ceil(total.value / ADMIN_ORDERS_PAGE_SIZE)),
)

// The owner gets nearby pages plus first/last; ellipses keep large lists tidy.
const paginationItems = computed<(number | 'ellipsis')[]>(() => {
  const pages = new Set<number>([1, totalPages.value])

  for (let page = currentPage.value - 1; page <= currentPage.value + 1; page += 1) {
    if (page >= 1 && page <= totalPages.value) {
      pages.add(page)
    }
  }

  const sortedPages = [...pages].sort((firstPage, secondPage) => firstPage - secondPage)
  return sortedPages.reduce<(number | 'ellipsis')[]>((items, page, index) => {
    const previousPage = sortedPages[index - 1]

    if (previousPage !== undefined && page - previousPage > 1) {
      items.push('ellipsis')
    }

    items.push(page)
    return items
  }, [])
})

// A request counter, not a boolean: a slow early response must not overwrite a
// fresher one when the owner types quickly.
let latestRequestId = 0

// A session that expired while the page was open should return the owner to
// sign-in rather than showing a load failure they cannot act on.
const isExpiredSession = async (requestError: unknown): Promise<boolean> => {
  if ((requestError as { status?: number } | null)?.status !== 401) {
    return false
  }

  session.value = null
  await router.push(ADMIN_ROUTES.login)

  return true
}

const load = async (): Promise<void> => {
  const requestId = ++latestRequestId
  isLoading.value = true
  hasFailed.value = false

  try {
    const page = await fetchOrders({
      status: statusFilter.value,
      kind: kindFilter.value,
      q: searchTerm.value,
      from: fromFilter.value,
      to: toFilter.value,
      limit: ADMIN_ORDERS_PAGE_SIZE,
      offset: offset.value,
    })

    if (requestId !== latestRequestId) {
      return
    }

    orders.value = page.orders
    total.value = page.total
  } catch (loadError) {
    if (requestId !== latestRequestId || (await isExpiredSession(loadError))) {
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

// The summary strip is independent of the filters, so a failure there must not
// blank the list the owner came for.
const loadMetrics = async (): Promise<void> => {
  try {
    metrics.value = await fetchMetrics(metricsPeriod.value)
  } catch (metricsError) {
    if (!(await isExpiredSession(metricsError))) {
      metrics.value = null
    }
  }
}

// Changing a filter always returns to the first page: otherwise a narrower
// result set would leave the owner stranded on an empty page.
let searchTimer: ReturnType<typeof setTimeout> | undefined

watch([statusFilter, kindFilter, fromFilter, toFilter], () => {
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

// The period narrows the summary only; the list below is unaffected.
watch(metricsPeriod, () => {
  void loadMetrics()
})

onBeforeUnmount(() => clearTimeout(searchTimer))

const goToPreviousPage = (): void => {
  offset.value = Math.max(0, offset.value - ADMIN_ORDERS_PAGE_SIZE)
  void load()
}

// Moves directly to one numbered page.
const goToPage = (page: number): void => {
  if (page === currentPage.value) {
    return
  }

  offset.value = (page - 1) * ADMIN_ORDERS_PAGE_SIZE
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
  fromFilter.value = ''
  toFilter.value = ''
}

// Updates one filter category from the table filter popover.
const updateFilter = (key: AdminOrderFilterKey, value: string): void => {
  if (key === 'status') {
    statusFilter.value = value as OrderStatus | ''
  } else if (key === 'kind') {
    kindFilter.value = value as OrderKind | ''
  } else if (key === 'from') {
    fromFilter.value = value
  } else {
    toFilter.value = value
  }
}

// Removes one applied filter chip without disturbing the other active filters.
const removeFilter = (key: AdminOrderFilterKey): void => {
  updateFilter(key, '')
}

// Clicking a metric card filters the list by that status.
const selectMetric = (nextStatus: OrderStatus | null): void => {
  statusFilter.value = nextStatus ?? ''
}

onMounted(() => {
  void load()
  void loadMetrics()
})
</script>

<template>
  <div class="admin-page">
    <header class="admin-page-head">
      <div>
        <h1>{{ ADMIN_COPY.ordersHeading }}</h1>
        <p class="admin-page-subheading">{{ ADMIN_COPY.ordersSubheading }}</p>
      </div>
    </header>

    <AdminMetricCards
      v-model:period="metricsPeriod"
      :metrics="metrics"
      :active-status="statusFilter"
      @select="selectMetric"
    />

    <div class="admin-toolbar">
      <div class="admin-toolbar-filters">
        <AdminTableFilter
          :options="ADMIN_ORDER_FILTER_OPTIONS"
          :values="filterValues"
          @update="updateFilter"
        />
        <AdminAppliedFilters
          :options="ADMIN_ORDER_FILTER_OPTIONS"
          :values="filterValues"
          @remove="removeFilter"
          @clear="clearFilters"
        />
      </div>

      <label class="admin-search">
        <span class="visually-hidden">{{ ADMIN_COPY.searchLabel }}</span>
        <StorefrontIcon
          class="admin-search-icon"
          :icon-node="Search"
          :size="16"
          :stroke-width="1.75"
        />
        <input
          v-model.trim="searchTerm"
          type="search"
          :placeholder="ADMIN_COPY.searchPlaceholder"
        >
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
              <td class="admin-table-order-no">
                <NuxtLink :to="ADMIN_ROUTES.orderDetail(order.orderNo)">
                  {{ order.orderNo }}
                </NuxtLink>
              </td>
              <td>
                <NuxtLink
                  class="admin-customer"
                  :to="ADMIN_ROUTES.orderDetail(order.orderNo)"
                >
                  <span class="admin-customer-avatar" aria-hidden="true">
                    {{ order.customerName.slice(0, 1).toUpperCase() }}
                  </span>
                  <span class="admin-customer-detail">
                    <span class="admin-customer-name">{{ order.customerName }}</span>
                    <span class="admin-customer-contact">{{ order.phone }}</span>
                  </span>
                </NuxtLink>
              </td>
              <td>
                <span class="admin-kind" :data-kind="order.kind">
                  {{ ORDER_KIND_LABELS[order.kind] }}
                </span>
              </td>
              <td>
                <span class="admin-status" :data-status="order.status">
                  {{ ORDER_STATUS_LABELS[order.status] }}
                </span>
              </td>
              <td class="admin-table-number">{{ order.itemCount }}</td>
              <td class="admin-table-number admin-table-total">
                {{ formatProductPrice(order.totalKobo) }}
              </td>
              <td class="admin-table-placed">{{ formatOrderDate(order.createdAt) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <nav class="admin-pagination">
        <p class="admin-pagination-count">
          {{ rangeStart }}–{{ rangeEnd }} of {{ total }}
        </p>

        <div class="admin-pagination-controls">
          <button
            type="button"
            class="admin-pagination-button"
            :disabled="!canGoBack"
            @click="goToPreviousPage"
          >
            {{ ADMIN_COPY.previous }}
          </button>

          <template v-for="(item, index) in paginationItems" :key="`${item}-${index}`">
            <span v-if="item === 'ellipsis'" class="admin-pagination-ellipsis">…</span>
            <button
              v-else
              type="button"
              class="admin-pagination-button is-number"
              :class="{ 'is-current': item === currentPage }"
              :aria-current="item === currentPage ? 'page' : undefined"
              @click="goToPage(item)"
            >
              {{ item }}
            </button>
          </template>

          <button
            type="button"
            class="admin-pagination-button"
            :disabled="!canGoForward"
            @click="goToNextPage"
          >
            {{ ADMIN_COPY.next }}
          </button>
        </div>
      </nav>
    </template>
  </div>
</template>
