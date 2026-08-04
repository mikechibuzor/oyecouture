<script setup lang="ts">
import {
  ADMIN_COPY,
  ADMIN_ROUTES,
} from '~/constants/admin'
import type { AdminProductDraft } from '~/types/admin'
import { ProductSection } from '~/types/product'

definePageMeta({ layout: 'admin', middleware: 'admin' })

useHead({ title: `${ADMIN_COPY.productNewHeading} · ${ADMIN_COPY.brand}` })

const { createProduct } = useAdminProducts()
const { session } = useAdminSession()
const router = useRouter()

const draft = ref<AdminProductDraft>({
  slug: '',
  name: '',
  section: ProductSection.CasualsMen,
  category: '',
  priceKobo: 0,
  description: '',
  tag: '',
  hasSizes: true,
  active: true,
  sortOrder: 0,
})
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

const submit = async (): Promise<void> => {
  isSaving.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const product = await createProduct(draft.value)
    successMessage.value = ADMIN_COPY.productSaved
    await router.push(ADMIN_ROUTES.productEdit(product.slug))
  } catch (saveError) {
    if (!(await isExpiredSession(saveError))) {
      errorMessage.value = productErrorMessage(saveError)
    }
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div class="admin-page">
    <NuxtLink class="admin-back-link" :to="ADMIN_ROUTES.products">
      {{ ADMIN_COPY.backToProducts }}
    </NuxtLink>

    <header class="admin-detail-head">
      <div>
        <p class="admin-detail-eyebrow">{{ ADMIN_COPY.productsHeading }}</p>
        <h1>{{ ADMIN_COPY.productNewHeading }}</h1>
      </div>
    </header>

    <AdminProductForm
      v-model="draft"
      :is-saving="isSaving"
      :error-message="errorMessage"
      :success-message="successMessage"
      @submit="submit"
    />
  </div>
</template>
