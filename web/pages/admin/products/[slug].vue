<script setup lang="ts">
import {
  ADMIN_COPY,
  ADMIN_ROUTES,
} from '~/constants/admin'
import type { AdminProduct, AdminProductDraft } from '~/types/admin'
import { koboToNaira } from '~/utils/product'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const route = useRoute()
const router = useRouter()
const { fetchProduct, updateProduct, deleteProduct } = useAdminProducts()
const { notifyError, notifySuccess } = useAppNotifications()
const { session } = useAdminSession()

const productSlug = computed(() => String(route.params.slug ?? ''))

useHead({ title: `${ADMIN_COPY.productEditHeading} · ${ADMIN_COPY.brand}` })

const product = ref<AdminProduct | null>(null)
const isDeleting = ref(false)

// Deletion is permanent and the confirmation says so, including what survives:
// past orders keep the name and price they were sold at. Archiving is offered
// as the reversible alternative in the same sentence.
const removeProduct = async (): Promise<void> => {
  // eslint-disable-next-line no-alert
  if (!window.confirm(ADMIN_COPY.productDeleteConfirm)) {
    return
  }

  isDeleting.value = true

  try {
    await deleteProduct(productSlug.value)
    notifySuccess(ADMIN_COPY.productDeletedNotice)
    await router.push(ADMIN_ROUTES.products)
  } catch {
    notifyError(ADMIN_COPY.productDeleteFailed)
  } finally {
    isDeleting.value = false
  }
}
const draft = ref<AdminProductDraft | null>(null)
const isLoading = ref(true)
const hasFailed = ref(false)
const isSaving = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const isExpiredSession = async (requestError: unknown): Promise<boolean> => {
  if ((requestError as { status?: number } | null)?.status !== 401) {
    return false
  }

  session.value = null
  await router.push(ADMIN_ROUTES.login)

  return true
}

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

const load = async (): Promise<void> => {
  isLoading.value = true
  hasFailed.value = false
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const loadedProduct = await fetchProduct(productSlug.value)
    product.value = loadedProduct
    draft.value = toDraft(loadedProduct)
  } catch (loadError) {
    if (!(await isExpiredSession(loadError))) {
      hasFailed.value = true
      product.value = null
      draft.value = null
    }
  } finally {
    isLoading.value = false
  }
}

const submit = async (): Promise<void> => {
  if (!draft.value) {
    return
  }

  isSaving.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const savedProduct = await updateProduct(productSlug.value, draft.value)
    product.value = savedProduct
    draft.value = toDraft(savedProduct)
    successMessage.value = ADMIN_COPY.productSaved

    if (savedProduct.slug !== productSlug.value) {
      await router.replace(ADMIN_ROUTES.productEdit(savedProduct.slug))
    }
  } catch (saveError) {
    if (!(await isExpiredSession(saveError))) {
      errorMessage.value = productErrorMessage(saveError)
    }
  } finally {
    isSaving.value = false
  }
}

const handleProductImagesUpdated = (updatedProduct: AdminProduct): void => {
  product.value = updatedProduct
  draft.value = toDraft(updatedProduct)
}

onMounted(() => {
  void load()
})
</script>

<template>
  <div class="admin-page">
    <NuxtLink class="admin-back-link" :to="ADMIN_ROUTES.products">
      {{ ADMIN_COPY.backToProducts }}
    </NuxtLink>

    <p v-if="isLoading" class="admin-orders-message">
      {{ ADMIN_COPY.productDetailLoading }}
    </p>

    <div v-else-if="hasFailed" class="admin-orders-message" role="alert">
      <p>{{ ADMIN_COPY.productDetailFailed }}</p>
      <button type="button" class="admin-inline-action" @click="load">
        {{ ADMIN_COPY.retry }}
      </button>
    </div>

    <template v-else-if="draft">
      <header class="admin-detail-head">
        <div>
          <p class="admin-detail-eyebrow">{{ ADMIN_COPY.productsHeading }}</p>
          <h1>{{ ADMIN_COPY.productEditHeading }}</h1>
          <p v-if="product" class="admin-page-subheading">{{ product.name }}</p>
        </div>
      </header>

      <AdminProductForm
        v-model="draft"
        :is-saving="isSaving"
        :error-message="errorMessage"
        :success-message="successMessage"
        @submit="submit"
      />

      <AdminProductImageManager
        v-if="product"
        :product="product"
        @updated="handleProductImagesUpdated"
      />

      <!-- Kept apart from the save actions and at the end of the page, so it is
           never the button reached for by habit. -->
      <section class="admin-form-section admin-product-danger">
        <header class="admin-form-section-head">
          <h2>{{ ADMIN_COPY.deleteProduct }}</h2>
          <p>{{ ADMIN_COPY.productDeleteConfirm }}</p>
        </header>

        <div class="admin-form-actions">
          <button
            type="button"
            class="admin-inline-action is-danger"
            :disabled="isDeleting"
            @click="removeProduct"
          >
            {{ isDeleting ? ADMIN_COPY.deletingProduct : ADMIN_COPY.deleteProduct }}
          </button>
        </div>
      </section>
    </template>
  </div>
</template>
