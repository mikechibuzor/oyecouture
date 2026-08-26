<script setup lang="ts">
import { Search, X } from 'lucide'
import {
  ADMIN_COPY,
  ADMIN_PRODUCT_FILTER_OPTIONS,
  ADMIN_PRODUCTS_PAGE_SIZE,
  ADMIN_ROWS_PER_PAGE_OPTIONS,
  ADMIN_ROUTES,
  ADMIN_SEARCH_DEBOUNCE_MS,
  PRODUCT_SECTION_LABELS,
} from '~/constants/admin'
import type {
  AdminFilterKey,
  AdminProduct,
  AdminProductActiveFilter,
  AdminProductDraft,
  AdminProductFilterKey,
  AdminProductFilterValues,
} from '~/types/admin'
import { ProductSection } from '~/types/product'
import { formatOrderDate } from '~/utils/admin'
import { formatProductPrice, koboToNaira } from '~/utils/product'

definePageMeta({ layout: 'admin', middleware: 'admin' })

useHead({ title: `${ADMIN_COPY.productsHeading} · ${ADMIN_COPY.brand}` })

type ProductEditorMode = 'create' | 'edit'

const {
  createProduct,
  fetchProducts,
  fetchProduct,
  updateProduct,
  updateProductActive,
} = useAdminProducts()
const { session } = useAdminSession()
const router = useRouter()
const route = useRoute()
const { notifySuccess } = useAppNotifications()
const { withLoading } = useAppLoading()

// Normalises a query value because Vue Router may expose repeated keys as arrays.
const queryString = (value: unknown): string =>
  Array.isArray(value) ? String(value[0] ?? '') : typeof value === 'string' ? value : ''

// Keeps rows-per-page inside the options the UI can actually render.
const pageSizeFromQuery = (value: unknown): string => {
  const queryValue = queryString(value)
  return ADMIN_ROWS_PER_PAGE_OPTIONS.some((option) => option.value === queryValue)
    ? queryValue
    : String(ADMIN_PRODUCTS_PAGE_SIZE)
}

// Page numbers are one-based in the URL and converted to an offset internally.
const pageFromQuery = (value: unknown): number => {
  const page = Number(queryString(value))
  return Number.isInteger(page) && page > 0 ? page : 1
}

// Query section is accepted only when it matches the product section enum.
const sectionFromQuery = (value: unknown): ProductSection | '' => {
  const queryValue = queryString(value)
  return Object.values(ProductSection).includes(queryValue as ProductSection)
    ? (queryValue as ProductSection)
    : ''
}

// Query active state is accepted only when it matches the availability filter.
const activeFromQuery = (value: unknown): AdminProductActiveFilter => {
  const queryValue = queryString(value)
  return queryValue === 'true' || queryValue === 'false' ? queryValue : ''
}

const initialPageSize = pageSizeFromQuery(route.query.limit)

const searchTerm = ref(queryString(route.query.q))
const sectionFilter = ref<ProductSection | ''>(sectionFromQuery(route.query.section))
const activeFilter = ref<AdminProductActiveFilter>(activeFromQuery(route.query.active))
const offset = ref((pageFromQuery(route.query.page) - 1) * Number(initialPageSize))
const pageSize = ref(initialPageSize)

const products = ref<readonly AdminProduct[]>([])
const total = ref(0)
const isLoading = ref(true)
const hasFailed = ref(false)
const updateError = ref('')
const updatingSlug = ref<string | null>(null)
const productEditorMode = ref<ProductEditorMode | null>(null)
const selectedProductSlug = ref<string | null>(null)
const selectedProduct = ref<AdminProduct | null>(null)
const selectedDraft = ref<AdminProductDraft | null>(null)
const isEditorLoading = ref(false)
const hasEditorFailed = ref(false)
const isEditorSaving = ref(false)
const editorErrorMessage = ref('')
const editorSuccessMessage = ref('')
const editorDrawer = ref<HTMLElement | null>(null)
const editorCloseButton = ref<HTMLButtonElement | null>(null)

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

// Rows per page is stored as a string for UiSelect and converted at the API edge.
const pageSizeValue = computed(() => Number(pageSize.value) || ADMIN_PRODUCTS_PAGE_SIZE)

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

const isEditorOpen = computed(() => productEditorMode.value !== null)

const selectedProductTitle = computed(() =>
  productEditorMode.value === 'create'
    ? ADMIN_COPY.productNewHeading
    : selectedProduct.value?.name ?? ADMIN_COPY.productEditHeading,
)

const canManageSelectedProduct = computed(
  () => productEditorMode.value === 'edit' && selectedProduct.value !== null,
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
let isWritingRouteQuery = false
let isReadingRouteQuery = false

// The table state is encoded into readable URL query params.
const adminQuery = computed<Record<string, string>>(() => {
  const query: Record<string, string> = {
    page: String(currentPage.value),
    limit: pageSize.value,
  }

  if (searchTerm.value) {
    query.q = searchTerm.value
  }

  if (sectionFilter.value) {
    query.section = sectionFilter.value
  }

  if (activeFilter.value) {
    query.active = activeFilter.value
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
  sectionFilter.value = sectionFromQuery(route.query.section)
  activeFilter.value = activeFromQuery(route.query.active)
  pageSize.value = nextPageSize
  offset.value = (nextPage - 1) * Number(nextPageSize)

  await nextTick()
  isReadingRouteQuery = false
}

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
    const page = await withLoading(
      () => fetchProducts({
        section: sectionFilter.value,
        active: activeFilter.value,
        q: searchTerm.value,
        limit: pageSizeValue.value,
        offset: offset.value,
      }),
      ADMIN_COPY.productLoading,
    )

    if (requestId !== latestRequestId) {
      return
    }

    if (page.total > 0 && page.products.length === 0 && offset.value >= page.total) {
      offset.value = Math.floor((page.total - 1) / pageSizeValue.value) * pageSizeValue.value
      replaceRouteQuery()
      void load()
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

    await syncFromRouteQuery()
    void load()
  },
)

onBeforeUnmount(() => clearTimeout(searchTimer))

const goToPreviousPage = (): void => {
  offset.value = Math.max(0, offset.value - pageSizeValue.value)
  replaceRouteQuery()
  void load()
}

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

const productErrorMessage = (requestError: unknown): string => {
  const errorCode = (requestError as { data?: { error?: { code?: string } } } | null)
    ?.data?.error?.code

  return errorCode === 'duplicate_product_slug'
    ? ADMIN_COPY.productDuplicateSlug
    : ADMIN_COPY.productSaveFailed
}

const toDraft = (adminProduct: AdminProduct): AdminProductDraft => ({
  slug: adminProduct.slug,
  name: adminProduct.name,
  section: adminProduct.section,
  category: adminProduct.category,
  priceNaira: koboToNaira(adminProduct.priceKobo),
  description: adminProduct.description,
  tag: adminProduct.tag,
  hasSizes: adminProduct.hasSizes,
  active: adminProduct.active,
  sortOrder: adminProduct.sortOrder,
})

const createEmptyDraft = (): AdminProductDraft => ({
  slug: '',
  name: '',
  section: ProductSection.CasualsMen,
  category: '',
  priceNaira: 0,
  description: '',
  tag: '',
  hasSizes: true,
  active: true,
  sortOrder: 0,
})

const syncProductInList = (updatedProduct: AdminProduct): void => {
  products.value = products.value.map((product) =>
    product.id === updatedProduct.id ? updatedProduct : product,
  )
}

const toggleProductActive = async (
  product: AdminProduct,
  onError: (message: string) => void = (message) => {
    updateError.value = message
  },
): Promise<AdminProduct | null> => {
  updatingSlug.value = product.slug
  onError('')

  try {
    await withLoading(
      () => updateProductActive(product.slug, !product.active),
      productPendingLabel(product),
    )
    const updatedProduct: AdminProduct = { ...product, active: !product.active }
    syncProductInList(updatedProduct)

    if (selectedProduct.value?.id === updatedProduct.id) {
      selectedProduct.value = updatedProduct
      selectedDraft.value = toDraft(updatedProduct)
    }

    return updatedProduct
  } catch (updateActiveError) {
    if (!(await isExpiredSession(updateActiveError))) {
      onError(ADMIN_COPY.productUpdateFailed)
    }

    return null
  } finally {
    updatingSlug.value = null
  }
}

// Nested controls keep their native behavior; the rest of the row opens edit.
const isInteractiveRowTarget = (eventTarget: EventTarget | null): boolean =>
  eventTarget instanceof Element &&
  Boolean(eventTarget.closest('a, button, input, select, textarea, summary'))

const openProductEditor = async (productSlug: string): Promise<void> => {
  productEditorMode.value = 'edit'
  selectedProductSlug.value = productSlug
  selectedProduct.value = products.value.find((product) => product.slug === productSlug) ?? null
  selectedDraft.value = selectedProduct.value ? toDraft(selectedProduct.value) : null
  isEditorLoading.value = true
  hasEditorFailed.value = false
  editorErrorMessage.value = ''
  editorSuccessMessage.value = ''

  try {
    const loadedProduct = await withLoading(
      () => fetchProduct(productSlug),
      ADMIN_COPY.productDetailLoading,
    )
    selectedProduct.value = loadedProduct
    selectedDraft.value = toDraft(loadedProduct)
    syncProductInList(loadedProduct)
  } catch (loadError) {
    if (!(await isExpiredSession(loadError))) {
      hasEditorFailed.value = true
      selectedProduct.value = null
      selectedDraft.value = null
    }
  } finally {
    isEditorLoading.value = false
  }
}

const openNewProductEditor = (): void => {
  productEditorMode.value = 'create'
  selectedProductSlug.value = null
  selectedProduct.value = null
  selectedDraft.value = createEmptyDraft()
  isEditorLoading.value = false
  hasEditorFailed.value = false
  editorErrorMessage.value = ''
  editorSuccessMessage.value = ''
}

const closeProductEditor = (): void => {
  productEditorMode.value = null
  selectedProductSlug.value = null
  selectedProduct.value = null
  selectedDraft.value = null
  hasEditorFailed.value = false
  editorErrorMessage.value = ''
  editorSuccessMessage.value = ''
}

const saveSelectedProduct = async (): Promise<void> => {
  if (!selectedDraft.value) {
    return
  }

  isEditorSaving.value = true
  editorErrorMessage.value = ''
  editorSuccessMessage.value = ''

  try {
    const productDraft = selectedDraft.value
    const wasCreatingProduct = productEditorMode.value === 'create'
    const savedProduct = await withLoading(
      () => wasCreatingProduct
        ? createProduct(productDraft)
        : selectedProductSlug.value
          ? updateProduct(selectedProductSlug.value, productDraft)
          : Promise.resolve(null),
      ADMIN_COPY.savingProduct,
    )

    if (!savedProduct) {
      return
    }

    productEditorMode.value = 'edit'
    selectedProductSlug.value = savedProduct.slug
    selectedProduct.value = savedProduct
    selectedDraft.value = toDraft(savedProduct)
    syncProductInList(savedProduct)

    await load()
    closeProductEditor()
    notifySuccess(
      wasCreatingProduct ? ADMIN_COPY.productCreatedNotice : ADMIN_COPY.productUpdatedNotice,
      savedProduct.name,
    )
  } catch (saveError) {
    if (!(await isExpiredSession(saveError))) {
      editorErrorMessage.value = productErrorMessage(saveError)
    }
  } finally {
    isEditorSaving.value = false
  }
}

const toggleSelectedProductActive = async (): Promise<void> => {
  if (!canManageSelectedProduct.value || !selectedProduct.value) {
    return
  }

  const wasActive = selectedProduct.value.active
  const updatedProduct = await toggleProductActive(selectedProduct.value, (message) => {
    editorErrorMessage.value = message
  })

  if (!updatedProduct) {
    return
  }

  await load()
  closeProductEditor()
  notifySuccess(
    wasActive ? ADMIN_COPY.productArchivedNotice : ADMIN_COPY.productRestoredNotice,
    updatedProduct.name,
  )
}

const handleProductImagesUpdated = (updatedProduct: AdminProduct): void => {
  selectedProduct.value = updatedProduct
  selectedDraft.value = toDraft(updatedProduct)
  syncProductInList(updatedProduct)
}

const retrySelectedProductEditor = (): void => {
  if (selectedProductSlug.value) {
    void openProductEditor(selectedProductSlug.value)
  }
}

const handleProductRowClick = (mouseEvent: MouseEvent, productSlug: string): void => {
  if (!isInteractiveRowTarget(mouseEvent.target)) {
    void openProductEditor(productSlug)
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
  void openProductEditor(productSlug)
}

useModalFocus({
  isOpen: isEditorOpen,
  container: editorDrawer,
  initialFocus: editorCloseButton,
  onEscape: closeProductEditor,
  lockScroll: false,
})

watch(isEditorOpen, (isOpen) => {
  if (!import.meta.client) {
    return
  }

  document.body.classList.toggle('has-admin-drawer', isOpen)
})

onMounted(() => {
  replaceRouteQuery()
  void load()
})

onBeforeUnmount(() => {
  document.body.classList.remove('has-admin-drawer')
})
</script>

<template>
  <div class="admin-page">
    <header class="admin-page-head">
      <div>
        <h1>{{ ADMIN_COPY.productsHeading }}</h1>
        <p class="admin-page-subheading">{{ ADMIN_COPY.productsSubheading }}</p>
      </div>

      <button type="button" class="admin-inline-action" @click="openNewProductEditor">
        {{ ADMIN_COPY.newProduct }}
      </button>
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
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="product in products"
              :key="product.id"
              class="admin-table-row"
              tabindex="0"
              role="button"
              :aria-label="`Edit ${product.name}`"
              @click="handleProductRowClick($event, product.slug)"
              @keydown="handleProductRowKeydown($event, product.slug)"
            >
              <td data-label="Product">
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
              <td data-label="Section">{{ PRODUCT_SECTION_LABELS[product.section] }}</td>
              <td data-label="Category">{{ product.category }}</td>
              <td data-label="Price" class="admin-table-number admin-table-total">
                {{ formatProductPrice(product.priceKobo) }}
              </td>
              <td data-label="Status">
                <span class="admin-status" :data-status="productStatusTone(product)">
                  {{ productAvailabilityLabel(product) }}
                </span>
              </td>
              <td data-label="Updated" class="admin-table-placed">
                {{ formatOrderDate(product.createdAt) }}
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

    <Teleport to="body">
      <div
        v-if="isEditorOpen"
        class="admin-drawer-backdrop"
        role="presentation"
        @click.self="closeProductEditor"
      >
        <aside
          ref="editorDrawer"
          class="admin-drawer"
          role="dialog"
          aria-modal="true"
          aria-labelledby="admin-product-editor-title"
          tabindex="-1"
        >
          <header class="admin-drawer-head">
            <div>
              <p class="admin-detail-eyebrow">{{ ADMIN_COPY.productsHeading }}</p>
              <h2 id="admin-product-editor-title">{{ selectedProductTitle }}</h2>
            </div>

            <div class="admin-drawer-actions">
              <button
                v-if="canManageSelectedProduct && selectedProduct"
                type="button"
                class="admin-inline-action"
                :disabled="isUpdatingProduct(selectedProduct)"
                @click="toggleSelectedProductActive"
              >
                {{ productActionButtonLabel(selectedProduct) }}
              </button>
              <button
                ref="editorCloseButton"
                type="button"
                class="admin-drawer-close"
                :aria-label="ADMIN_COPY.closeProductEditor"
                @click="closeProductEditor"
              >
                <StorefrontIcon :icon-node="X" :size="18" :stroke-width="1.75" />
              </button>
            </div>
          </header>

          <p v-if="isEditorLoading" class="admin-orders-message">
            {{ ADMIN_COPY.productDetailLoading }}
          </p>

          <div v-else-if="hasEditorFailed" class="admin-orders-message" role="alert">
            <p>{{ ADMIN_COPY.productDetailFailed }}</p>
            <button
              type="button"
              class="admin-inline-action"
              @click="retrySelectedProductEditor"
            >
              {{ ADMIN_COPY.retry }}
            </button>
          </div>

          <template v-else-if="selectedDraft">
            <AdminProductForm
              v-model="selectedDraft"
              :is-saving="isEditorSaving"
              :error-message="editorErrorMessage"
              :success-message="editorSuccessMessage"
              @submit="saveSelectedProduct"
            />

            <AdminProductImageManager
              v-if="selectedProduct"
              :product="selectedProduct"
              @updated="handleProductImagesUpdated"
            />
          </template>
        </aside>
      </div>
    </Teleport>
  </div>
</template>
