<script setup lang="ts">
import { Search } from 'lucide'
import {
  ADMIN_COPY,
  ADMIN_PRODUCT_FILTER_OPTIONS,
  ADMIN_PRODUCTS_PAGE_SIZE,
  ADMIN_ROUTES,
  ADMIN_SEARCH_DEBOUNCE_MS,
  PRODUCT_SECTION_LABELS,
} from '~/constants/admin'
import type {
  AdminFilterKey,
  AdminProduct,
  AdminProductActiveFilter,
  AdminProductFilterKey,
  AdminProductFilterValues,
} from '~/types/admin'
import type { ProductSection } from '~/types/product'
import { formatOrderDate } from '~/utils/admin'
import { formatProductPrice } from '~/utils/product'

definePageMeta({ layout: 'admin', middleware: 'admin' })

useHead({ title: `${ADMIN_COPY.productsHeading} · ${ADMIN_COPY.brand}` })

const { fetchProducts, updateProductActive } = useAdminProducts()
const { session } = useAdminSession()
const router = useRouter()

const searchTerm = ref('')
const sectionFilter = ref<ProductSection | ''>('')
const activeFilter = ref<AdminProductActiveFilter>('')
const offset = ref(0)

const products = ref<readonly AdminProduct[]>([])
const total = ref(0)
const isLoading = ref(true)
const hasFailed = ref(false)
const updateError = ref('')
const updatingSlug = ref<string | null>(null)

// Any active filter changes the empty state from "no products" to "no matches".
const hasFilters = computed(
  () => searchTerm.value !== '' || sectionFilter.value !== '' || activeFilter.value !== '',
)

const filterValues = computed<AdminProductFilterValues>(() => ({
  section: sectionFilter.value,
  active: activeFilter.value,
}))

// The visible row range reads as "1-25 of 92" in the table footer.
const rangeStart = computed(() => (total.value === 0 ? 0 : offset.value + 1))

// The end of the visible range tracks the loaded page, not just the page size.
const rangeEnd = computed(() => offset.value + products.value.length)

// Pagination controls disable themselves at the list boundaries.
const canGoBack = computed(() => offset.value > 0)
const canGoForward = computed(() => rangeEnd.value < total.value)

// Current page derives from offset so direct page jumps and next/previous agree.
const currentPage = computed(() => Math.floor(offset.value / ADMIN_PRODUCTS_PAGE_SIZE) + 1)

// Total pages bottoms out at one so the pagination area never divides by zero.
const totalPages = computed(() =>
  Math.max(1, Math.ceil(total.value / ADMIN_PRODUCTS_PAGE_SIZE)),
)

// Products now use the same numbered pagination model as orders.
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

// A request counter prevents stale responses from replacing fresh filter results.
let latestRequestId = 0

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
  updateError.value = ''

  try {
    const page = await fetchProducts({
      section: sectionFilter.value,
      active: activeFilter.value,
      q: searchTerm.value,
      limit: ADMIN_PRODUCTS_PAGE_SIZE,
      offset: offset.value,
    })

    if (requestId !== latestRequestId) {
      return
    }

    products.value = page.products
    total.value = page.total
  } catch (loadError) {
    if (requestId !== latestRequestId || (await isExpiredSession(loadError))) {
      return
    }

    hasFailed.value = true
    products.value = []
    total.value = 0
  } finally {
    if (requestId === latestRequestId) {
      isLoading.value = false
    }
  }
}

let searchTimer: ReturnType<typeof setTimeout> | undefined

watch([sectionFilter, activeFilter], () => {
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
  offset.value = Math.max(0, offset.value - ADMIN_PRODUCTS_PAGE_SIZE)
  void load()
}

const goToPage = (page: number): void => {
  if (page === currentPage.value) {
    return
  }

  offset.value = (page - 1) * ADMIN_PRODUCTS_PAGE_SIZE
  void load()
}

const goToNextPage = (): void => {
  offset.value += ADMIN_PRODUCTS_PAGE_SIZE
  void load()
}

const clearFilters = (): void => {
  searchTerm.value = ''
  sectionFilter.value = ''
  activeFilter.value = ''
}

const isProductFilterKey = (key: AdminFilterKey): key is AdminProductFilterKey =>
  key === 'section' || key === 'active'

const updateFilter = (key: AdminFilterKey, value: string): void => {
  if (!isProductFilterKey(key)) {
    return
  }

  if (key === 'section') {
    sectionFilter.value = value as ProductSection | ''
    return
  }

  activeFilter.value = value as AdminProductActiveFilter
}

const removeFilter = (key: AdminFilterKey): void => {
  updateFilter(key, '')
}

const productImageUrl = (product: AdminProduct): string | null =>
  product.images[0]?.url || null

const productAvailabilityLabel = (product: AdminProduct): string =>
  product.active ? 'Visible' : 'Archived'

const productActionLabel = (product: AdminProduct): string =>
  product.active ? ADMIN_COPY.archiveProduct : ADMIN_COPY.restoreProduct

const productPendingLabel = (product: AdminProduct): string =>
  product.active ? ADMIN_COPY.archivingProduct : ADMIN_COPY.restoringProduct

const productActionButtonLabel = (product: AdminProduct): string =>
  isUpdatingProduct(product) ? productPendingLabel(product) : productActionLabel(product)

const productStatusTone = (product: AdminProduct): string =>
  product.active ? 'paid' : 'cancelled'

const isUpdatingProduct = (product: AdminProduct): boolean =>
  updatingSlug.value === product.slug

const toggleProductActive = async (product: AdminProduct): Promise<void> => {
  updatingSlug.value = product.slug
  updateError.value = ''

  try {
    await updateProductActive(product.slug, !product.active)
    await load()
  } catch (updateActiveError) {
    if (!(await isExpiredSession(updateActiveError))) {
      updateError.value = ADMIN_COPY.productUpdateFailed
    }
  } finally {
    updatingSlug.value = null
  }
}

// Nested controls keep their native behavior; the rest of the row opens edit.
const isInteractiveRowTarget = (eventTarget: EventTarget | null): boolean =>
  eventTarget instanceof Element &&
  Boolean(eventTarget.closest('a, button, input, select, textarea, summary'))

const openProductEdit = (productSlug: string): void => {
  void router.push(ADMIN_ROUTES.productEdit(productSlug))
}

const handleProductRowClick = (mouseEvent: MouseEvent, productSlug: string): void => {
  if (!isInteractiveRowTarget(mouseEvent.target)) {
    openProductEdit(productSlug)
  }
}

const handleProductRowKeydown = (keyboardEvent: KeyboardEvent, productSlug: string): void => {
  if (
    isInteractiveRowTarget(keyboardEvent.target) ||
    !['Enter', ' '].includes(keyboardEvent.key)
  ) {
    return
  }

  keyboardEvent.preventDefault()
  openProductEdit(productSlug)
}

onMounted(() => {
  void load()
})
</script>

<template>
  <div class="admin-page">
    <header class="admin-page-head">
      <div>
        <h1>{{ ADMIN_COPY.productsHeading }}</h1>
        <p class="admin-page-subheading">{{ ADMIN_COPY.productsSubheading }}</p>
      </div>

      <NuxtLink class="admin-inline-action" :to="ADMIN_ROUTES.productNew">
        {{ ADMIN_COPY.newProduct }}
      </NuxtLink>
    </header>

    <div class="admin-toolbar">
      <div class="admin-toolbar-filters">
        <AdminTableFilter
          :options="ADMIN_PRODUCT_FILTER_OPTIONS"
          :values="filterValues"
          @update="updateFilter"
        />
        <AdminAppliedFilters
          :options="ADMIN_PRODUCT_FILTER_OPTIONS"
          :values="filterValues"
          @remove="removeFilter"
          @clear="clearFilters"
        />
      </div>

      <label class="admin-search">
        <span class="visually-hidden">{{ ADMIN_COPY.productSearchLabel }}</span>
        <StorefrontIcon
          class="admin-search-icon"
          :icon-node="Search"
          :size="16"
          :stroke-width="1.75"
        />
        <input
          v-model.trim="searchTerm"
          type="search"
          :placeholder="ADMIN_COPY.productSearchPlaceholder"
        >
      </label>
    </div>

    <p v-if="updateError" class="admin-detail-error" role="alert">
      {{ updateError }}
    </p>

    <p v-if="isLoading" class="admin-orders-message">{{ ADMIN_COPY.productLoading }}</p>

    <div v-else-if="hasFailed" class="admin-orders-message" role="alert">
      <p>{{ ADMIN_COPY.productLoadFailed }}</p>
      <button type="button" class="admin-inline-action" @click="load">
        {{ ADMIN_COPY.retry }}
      </button>
    </div>

    <div v-else-if="products.length === 0" class="admin-orders-message">
      <p>{{ hasFilters ? ADMIN_COPY.productEmptyFiltered : ADMIN_COPY.productEmpty }}</p>
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
              <th scope="col">Product</th>
              <th scope="col">Section</th>
              <th scope="col">Category</th>
              <th scope="col" class="admin-table-number">Price</th>
              <th scope="col">Status</th>
              <th scope="col">Updated</th>
              <th scope="col">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="product in products"
              :key="product.id"
              class="admin-table-row"
              tabindex="0"
              role="link"
              :aria-label="`Edit ${product.name}`"
              @click="handleProductRowClick($event, product.slug)"
              @keydown="handleProductRowKeydown($event, product.slug)"
            >
              <td>
                <div class="admin-product">
                  <img
                    v-if="productImageUrl(product)"
                    :src="productImageUrl(product) ?? ''"
                    :alt="product.name"
                    loading="lazy"
                  >
                  <span v-else class="admin-product-placeholder" aria-hidden="true">
                    {{ product.name.slice(0, 1).toUpperCase() }}
                  </span>
                  <span class="admin-product-detail">
                    <strong>{{ product.name }}</strong>
                    <small>{{ product.slug }}</small>
                  </span>
                </div>
              </td>
              <td>{{ PRODUCT_SECTION_LABELS[product.section] }}</td>
              <td>{{ product.category }}</td>
              <td class="admin-table-number admin-table-total">
                {{ formatProductPrice(product.priceKobo) }}
              </td>
              <td>
                <span class="admin-status" :data-status="productStatusTone(product)">
                  {{ productAvailabilityLabel(product) }}
                </span>
              </td>
              <td class="admin-table-placed">{{ formatOrderDate(product.createdAt) }}</td>
              <td>
                <div class="admin-product-actions">
                  <NuxtLink
                    class="admin-inline-action"
                    :to="ADMIN_ROUTES.productEdit(product.slug)"
                  >
                    {{ ADMIN_COPY.editProduct }}
                  </NuxtLink>
                  <button
                    type="button"
                    class="admin-inline-action"
                    :disabled="isUpdatingProduct(product)"
                    @click="toggleProductActive(product)"
                  >
                    {{ productActionButtonLabel(product) }}
                  </button>
                </div>
              </td>
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
