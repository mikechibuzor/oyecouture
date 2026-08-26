<script setup lang="ts">
import { Search } from 'lucide'
import {
  ADMIN_COPY,
  ADMIN_DEFAULT_PERIOD,
  ADMIN_ORDERS_PAGE_SIZE,
  ADMIN_ORDER_FILTER_OPTIONS,
  ADMIN_ROWS_PER_PAGE_OPTIONS,
  ADMIN_ROUTES,
  ADMIN_SEARCH_DEBOUNCE_MS,
  ORDER_KIND_LABELS,
  ORDER_STATUS_LABELS,
} from '~/constants/admin'
import type {
  AdminFilterKey,
  AdminOrderMetrics,
  AdminOrderSummary,
} from '~/types/admin'
import {
  type AdminOrderFilterKey,
  type AdminOrderFilterValues,
  MetricsPeriod,
  OrderStatus,
} from '~/types/admin'
import { OrderKind } from '~/types/order'
import { formatOrderDate } from '~/utils/admin'
import { formatProductPrice } from '~/utils/product'

definePageMeta({ layout: 'admin', middleware: 'admin' })

useHead({ title: `${ADMIN_COPY.ordersHeading} · ${ADMIN_COPY.brand}` })

const { fetchOrders, fetchMetrics } = useAdminOrders()
const { session } = useAdminSession()
const router = useRouter()
const route = useRoute()

// Normalises a query value because Vue Router may expose repeated keys as arrays.
const queryString = (value: unknown): string =>
  Array.isArray(value) ? String(value[0] ?? '') : typeof value === 'string' ? value : ''

// Keeps rows-per-page inside the options the UI can actually render.
const pageSizeFromQuery = (value: unknown): string => {
  const queryValue = queryString(value)
  return ADMIN_ROWS_PER_PAGE_OPTIONS.some((option) => option.value === queryValue)
    ? queryValue
    : String(ADMIN_ORDERS_PAGE_SIZE)
}

// Page numbers are one-based in the URL and converted to an offset internally.
const pageFromQuery = (value: unknown): number => {
  const page = Number(queryString(value))
  return Number.isInteger(page) && page > 0 ? page : 1
}

// Query status is accepted only when it matches the order lifecycle enum.
const statusFromQuery = (value: unknown): OrderStatus | '' => {
  const queryValue = queryString(value)
  return Object.values(OrderStatus).includes(queryValue as OrderStatus)
    ? (queryValue as OrderStatus)
    : ''
}

// Query order type is accepted only when it matches the order kind enum.
const kindFromQuery = (value: unknown): OrderKind | '' => {
  const queryValue = queryString(value)
  return Object.values(OrderKind).includes(queryValue as OrderKind)
    ? (queryValue as OrderKind)
    : ''
}

// Metrics period is URL-backed too so the summary strip survives a reload.
const periodFromQuery = (value: unknown): MetricsPeriod => {
  const queryValue = queryString(value)
  return Object.values(MetricsPeriod).includes(queryValue as MetricsPeriod)
    ? (queryValue as MetricsPeriod)
    : ADMIN_DEFAULT_PERIOD
}

const initialPageSize = pageSizeFromQuery(route.query.limit)

const searchTerm = ref(queryString(route.query.q))
const statusFilter = ref<OrderStatus | ''>(statusFromQuery(route.query.status))
const kindFilter = ref<OrderKind | ''>(kindFromQuery(route.query.kind))
const fromFilter = ref(queryString(route.query.from))
const toFilter = ref(queryString(route.query.to))
const offset = ref((pageFromQuery(route.query.page) - 1) * Number(initialPageSize))
const pageSize = ref(initialPageSize)

const orders = ref<readonly AdminOrderSummary[]>([])
const metrics = ref<AdminOrderMetrics | null>(null)
const metricsPeriod = ref<MetricsPeriod>(periodFromQuery(route.query.period))
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

// Rows per page is stored as a string for UiSelect and converted at the API edge.
const pageSizeValue = computed(() => Number(pageSize.value) || ADMIN_ORDERS_PAGE_SIZE)

// Pagination controls disable themselves at the list boundaries.
const canGoBack = computed(() => offset.value > 0)
const canGoForward = computed(() => rangeEnd.value < total.value)

// Current page derives from offset so direct page jumps and next/previous agree.
const currentPage = computed(() => Math.floor(offset.value / pageSizeValue.value) + 1)

// Total pages bottoms out at one so the pagination area never divides by zero.
const totalPages = computed(() =>
  Math.max(1, Math.ceil(total.value / pageSizeValue.value)),
)

// The summary mirrors the compact pagination inspiration while keeping range detail.
const paginationSummary = computed(() =>
  `Page ${currentPage.value} of ${totalPages.value}`,
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
let isWritingRouteQuery = false
let isReadingRouteQuery = false

// The table state is encoded into readable URL query params.
const adminQuery = computed<Record<string, string>>(() => {
  const query: Record<string, string> = {
    page: String(currentPage.value),
    limit: pageSize.value,
    period: metricsPeriod.value,
  }

  if (searchTerm.value) {
    query.q = searchTerm.value
  }

  if (statusFilter.value) {
    query.status = statusFilter.value
  }

  if (kindFilter.value) {
    query.kind = kindFilter.value
  }

  if (fromFilter.value) {
    query.from = fromFilter.value
  }

  if (toFilter.value) {
    query.to = toFilter.value
  }

  return query
})

// Prevents replace() churn when the current URL already represents the state.
const queryMatchesRoute = (query: Record<string, string>): boolean => {
  const routeKeys = Object.keys(route.query)
  const queryKeys = Object.keys(query)

  if (routeKeys.length !== queryKeys.length) {
    return false
  }

  return queryKeys.every((key) => queryString(route.query[key]) === query[key])
}

// Writes the current table state into the URL without adding history entries.
const replaceRouteQuery = (): void => {
  if (isReadingRouteQuery || queryMatchesRoute(adminQuery.value)) {
    return
  }

  isWritingRouteQuery = true
  void router.replace({ query: adminQuery.value }).finally(() => {
    isWritingRouteQuery = false
  })
}

// Pulls browser back/forward URL state back into the table controls.
const syncFromRouteQuery = async (): Promise<void> => {
  isReadingRouteQuery = true
  const nextPageSize = pageSizeFromQuery(route.query.limit)
  const nextPage = pageFromQuery(route.query.page)

  searchTerm.value = queryString(route.query.q)
  statusFilter.value = statusFromQuery(route.query.status)
  kindFilter.value = kindFromQuery(route.query.kind)
  fromFilter.value = queryString(route.query.from)
  toFilter.value = queryString(route.query.to)
  pageSize.value = nextPageSize
  metricsPeriod.value = periodFromQuery(route.query.period)
  offset.value = (nextPage - 1) * Number(nextPageSize)

  await nextTick()
  isReadingRouteQuery = false
}

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
      limit: pageSizeValue.value,
      offset: offset.value,
    })

    if (requestId !== latestRequestId) {
      return
    }

    if (page.total > 0 && page.orders.length === 0 && offset.value >= page.total) {
      offset.value = Math.floor((page.total - 1) / pageSizeValue.value) * pageSizeValue.value
      replaceRouteQuery()
      void load()
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
  if (isReadingRouteQuery) {
    return
  }

  offset.value = 0
  replaceRouteQuery()
  void load()
})

watch(searchTerm, () => {
  if (isReadingRouteQuery) {
    return
  }

  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    offset.value = 0
    replaceRouteQuery()
    void load()
  }, ADMIN_SEARCH_DEBOUNCE_MS)
})

// The period narrows the summary only; the list below is unaffected.
watch(metricsPeriod, () => {
  if (isReadingRouteQuery) {
    return
  }

  replaceRouteQuery()
  void loadMetrics()
})

watch(pageSize, () => {
  if (isReadingRouteQuery) {
    return
  }

  offset.value = 0
  replaceRouteQuery()
  void load()
})

watch(
  () => route.query,
  async () => {
    if (isWritingRouteQuery) {
      return
    }

    const previousPeriod = metricsPeriod.value
    await syncFromRouteQuery()
    void load()

    if (metricsPeriod.value !== previousPeriod) {
      void loadMetrics()
    }
  },
)

onBeforeUnmount(() => clearTimeout(searchTimer))

const goToPreviousPage = (): void => {
  offset.value = Math.max(0, offset.value - pageSizeValue.value)
  replaceRouteQuery()
  void load()
}

// Moves directly to one numbered page.
const goToPage = (page: number): void => {
  if (page === currentPage.value) {
    return
  }

  offset.value = (page - 1) * pageSizeValue.value
  replaceRouteQuery()
  void load()
}

const goToNextPage = (): void => {
  offset.value += pageSizeValue.value
  replaceRouteQuery()
  void load()
}

const clearFilters = (): void => {
  searchTerm.value = ''
  statusFilter.value = ''
  kindFilter.value = ''
  fromFilter.value = ''
  toFilter.value = ''
}

const isOrderFilterKey = (key: AdminFilterKey): key is AdminOrderFilterKey =>
  key === 'status' || key === 'kind' || key === 'from' || key === 'to'

// Updates one filter category from the table filter popover.
const updateFilter = (key: AdminFilterKey, value: string): void => {
  if (!isOrderFilterKey(key)) {
    return
  }

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
const removeFilter = (key: AdminFilterKey): void => {
  updateFilter(key, '')
}

// Clicking a metric card filters the list by that status.
const selectMetric = (nextStatus: OrderStatus | null): void => {
  statusFilter.value = nextStatus ?? ''
}

// Nested links keep their native behavior; the rest of the row opens details.
const isInteractiveRowTarget = (eventTarget: EventTarget | null): boolean =>
  eventTarget instanceof Element &&
  Boolean(eventTarget.closest('a, button, input, select, textarea, summary'))

// Opens the same detail route exposed by the order/customer links.
const openOrderDetail = (orderNo: string): void => {
  void router.push(ADMIN_ROUTES.orderDetail(orderNo))
}

// Row clicks make the table behave like an order desk, not a static report.
const handleOrderRowClick = (mouseEvent: MouseEvent, orderNo: string): void => {
  if (!isInteractiveRowTarget(mouseEvent.target)) {
    openOrderDetail(orderNo)
  }
}

// Keyboard users get the same row-level affordance with Enter or Space.
const handleOrderRowKeydown = (keyboardEvent: KeyboardEvent, orderNo: string): void => {
  if (
    isInteractiveRowTarget(keyboardEvent.target) ||
    !['Enter', ' '].includes(keyboardEvent.key)
  ) {
    return
  }

  keyboardEvent.preventDefault()
  openOrderDetail(orderNo)
}

onMounted(() => {
  replaceRouteQuery()
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
            <tr
              v-for="order in orders"
              :key="order.orderNo"
              class="admin-table-row"
              tabindex="0"
              role="link"
              :aria-label="`Open ${order.orderNo} details`"
              @click="handleOrderRowClick($event, order.orderNo)"
              @keydown="handleOrderRowKeydown($event, order.orderNo)"
            >
              <td data-label="Order" class="admin-table-order-no">
                <NuxtLink :to="ADMIN_ROUTES.orderDetail(order.orderNo)">
                  {{ order.orderNo }}
                </NuxtLink>
              </td>
              <td data-label="Customer">
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
              <td data-label="Type">
                <span class="admin-kind" :data-kind="order.kind">
                  {{ ORDER_KIND_LABELS[order.kind] }}
                </span>
              </td>
              <td data-label="Status">
                <span class="admin-status" :data-status="order.status">
                  {{ ORDER_STATUS_LABELS[order.status] }}
                </span>
              </td>
              <td data-label="Items" class="admin-table-number">{{ order.itemCount }}</td>
              <td data-label="Total" class="admin-table-number admin-table-total">
                {{ formatProductPrice(order.totalKobo) }}
              </td>
              <td data-label="Placed" class="admin-table-placed">
                {{ formatOrderDate(order.createdAt) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <nav class="admin-pagination">
        <div class="admin-pagination-meta">
          <p class="admin-pagination-count">
            {{ paginationSummary }} · {{ rangeStart }}–{{ rangeEnd }} of {{ total }}
          </p>
          <div class="admin-pagination-size">
            <span>{{ ADMIN_COPY.rowsPerPage }}</span>
            <UiSelect
              v-model="pageSize"
              :options="ADMIN_ROWS_PER_PAGE_OPTIONS"
              :label="ADMIN_COPY.rowsPerPage"
            />
          </div>
        </div>

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
