<script setup lang="ts">
import {
  ADMIN_COPY,
  ADMIN_ROUTES,
} from '~/constants/admin'
import type { AdminProduct, AdminProductDraft } from '~/types/admin'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const route = useRoute()
const router = useRouter()
const { fetchProduct, updateProduct } = useAdminProducts()
const { session } = useAdminSession()

const productSlug = computed(() => String(route.params.slug ?? ''))

useHead({ title: `${ADMIN_COPY.productEditHeading} · ${ADMIN_COPY.brand}` })

const product = ref<AdminProduct | null>(null)
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
  priceKobo: adminProduct.priceKobo,
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
        </div>
      </header>

      <AdminProductForm
        v-model="draft"
        :is-saving="isSaving"
        :error-message="errorMessage"
        :success-message="successMessage"
        @submit="submit"
      />
    </template>
  </div>
</template>
