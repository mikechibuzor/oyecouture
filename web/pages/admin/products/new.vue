<script setup lang="ts">
import {
  ADMIN_COPY,
  ADMIN_ROUTES,
} from '~/constants/admin'
import type { AdminProductDraft } from '~/types/admin'
import { ProductSection } from '~/types/product'

definePageMeta({ layout: 'admin', middleware: 'admin' })

useHead({ title: `${ADMIN_COPY.productNewHeading} · ${ADMIN_COPY.brand}` })

const { createProduct, uploadProductImage } = useAdminProducts()
const { session } = useAdminSession()
const router = useRouter()

const draft = ref<AdminProductDraft>({
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
const isSaving = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

// Images cannot be attached until the product exists, because the upload
// endpoint is addressed by slug. They are held here and sent immediately after
// the product is created, so the owner still adds everything in one pass.
const pendingImages = ref<File[]>([])

const selectPendingImages = (changeEvent: Event): void => {
  const input = changeEvent.target as HTMLInputElement
  pendingImages.value = Array.from(input.files ?? [])
}

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

    // Uploaded one at a time and in order, because position is assigned on
    // arrival: firing them in parallel would shuffle the gallery.
    let uploadFailed = false
    for (const imageFile of pendingImages.value) {
      try {
        await uploadProductImage(product.slug, imageFile)
      } catch {
        uploadFailed = true
      }
    }

    // The product exists either way, so the owner is sent to its editor rather
    // than left on a create form. A partial upload is reported there.
    successMessage.value = uploadFailed
      ? ADMIN_COPY.newProductImagesFailed
      : ADMIN_COPY.productSaved
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
    >
      <template #extra>
        <label class="admin-field admin-field-wide">
          <span>{{ ADMIN_COPY.newProductImagesLabel }}</span>
          <input type="file" accept="image/*" multiple @change="selectPendingImages">
          <small>
            {{
              pendingImages.length > 0
                ? ADMIN_COPY.newProductImagesPending(pendingImages.length)
                : ADMIN_COPY.newProductImagesHint
            }}
          </small>
        </label>
      </template>
    </AdminProductForm>
  </div>
</template>
