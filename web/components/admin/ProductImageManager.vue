<script setup lang="ts">
import { ArrowDown, ArrowUp, Trash2, Upload } from 'lucide'
import { ADMIN_COPY, ADMIN_ROUTES } from '~/constants/admin'
import type { AdminProduct } from '~/types/admin'
import type { ProductImage } from '~/types/product'

interface AdminProductImageManagerProps {
  product: AdminProduct
}

const props = defineProps<AdminProductImageManagerProps>()

const emit = defineEmits<{
  updated: [product: AdminProduct]
}>()

const router = useRouter()
const { session } = useAdminSession()
const { uploadProductImage, reorderProductImages, deleteProductImage } = useAdminProducts()
const { notifySuccess } = useAppNotifications()
const { withLoading } = useAppLoading()

const selectedFile = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
const isUploading = ref(false)
const movingImageID = ref<string | null>(null)
const deletingImageID = ref<string | null>(null)
const errorMessage = ref('')

const images = computed(() => props.product.images)
const selectedFileName = computed(() => selectedFile.value?.name ?? '')
const canUpload = computed(() => Boolean(selectedFile.value) && !isUploading.value)

const isExpiredSession = async (requestError: unknown): Promise<boolean> => {
  if ((requestError as { status?: number } | null)?.status !== 401) {
    return false
  }

  session.value = null
  await router.push(ADMIN_ROUTES.login)

  return true
}

const updateProduct = (updatedProduct: AdminProduct): void => {
  errorMessage.value = ''
  emit('updated', updatedProduct)
}

const handleFileSelection = (inputEvent: Event): void => {
  selectedFile.value = (inputEvent.target as HTMLInputElement).files?.[0] ?? null
}

const uploadSelectedImage = async (): Promise<void> => {
  if (!selectedFile.value) {
    return
  }

  const imageFile = selectedFile.value
  isUploading.value = true
  errorMessage.value = ''

  try {
    const updatedProduct = await withLoading(
      () => uploadProductImage(props.product.slug, imageFile),
      ADMIN_COPY.uploadingProductImage,
    )
    selectedFile.value = null
    if (fileInput.value) {
      fileInput.value.value = ''
    }
    updateProduct(updatedProduct)
    notifySuccess(ADMIN_COPY.productImageUploadedNotice, props.product.name)
  } catch (uploadError) {
    if (!(await isExpiredSession(uploadError))) {
      errorMessage.value = ADMIN_COPY.productImageUploadFailed
    }
  } finally {
    isUploading.value = false
  }
}

const reorderedImageIDs = (imageIndex: number, direction: -1 | 1): readonly string[] => {
  const nextImageIDs = images.value.map((image) => image.id)
  const targetIndex = imageIndex + direction
  const movedImageID = nextImageIDs[imageIndex]
  const targetImageID = nextImageIDs[targetIndex]

  if (!movedImageID || !targetImageID) {
    return nextImageIDs
  }

  nextImageIDs[imageIndex] = targetImageID
  nextImageIDs[targetIndex] = movedImageID

  return nextImageIDs
}

const moveImage = async (image: ProductImage, imageIndex: number, direction: -1 | 1): Promise<void> => {
  movingImageID.value = image.id
  errorMessage.value = ''

  try {
    const updatedProduct = await withLoading(
      () => reorderProductImages(
        props.product.slug,
        reorderedImageIDs(imageIndex, direction),
      ),
      ADMIN_COPY.productImageUpdatedNotice,
    )
    updateProduct(updatedProduct)
    notifySuccess(ADMIN_COPY.productImageUpdatedNotice, props.product.name)
  } catch (reorderError) {
    if (!(await isExpiredSession(reorderError))) {
      errorMessage.value = ADMIN_COPY.productImageUpdateFailed
    }
  } finally {
    movingImageID.value = null
  }
}

const removeImage = async (image: ProductImage): Promise<void> => {
  deletingImageID.value = image.id
  errorMessage.value = ''

  try {
    const updatedProduct = await withLoading(
      () => deleteProductImage(props.product.slug, image.id),
      ADMIN_COPY.productImageDeletedNotice,
    )
    updateProduct(updatedProduct)
    notifySuccess(ADMIN_COPY.productImageDeletedNotice, props.product.name)
  } catch (deleteError) {
    if (!(await isExpiredSession(deleteError))) {
      errorMessage.value = ADMIN_COPY.productImageDeleteFailed
    }
  } finally {
    deletingImageID.value = null
  }
}
</script>

<template>
  <section class="admin-detail-panel admin-product-images-panel">
    <div class="admin-product-images-head">
      <div>
        <h2>{{ ADMIN_COPY.productImagesHeading }}</h2>
        <p>{{ ADMIN_COPY.productImagesIntro }}</p>
      </div>

      <div class="admin-product-upload">
        <label class="admin-inline-action">
          <StorefrontIcon :icon-node="Upload" :size="14" :stroke-width="1.75" />
          {{ selectedFileName || ADMIN_COPY.uploadProductImage }}
          <input
            ref="fileInput"
            type="file"
            accept="image/jpeg,image/png,image/webp"
            @change="handleFileSelection"
          >
        </label>
        <button
          type="button"
          class="admin-inline-action"
          :disabled="!canUpload"
          @click="uploadSelectedImage"
        >
          {{ isUploading ? ADMIN_COPY.uploadingProductImage : ADMIN_COPY.uploadProductImage }}
        </button>
      </div>
    </div>

    <p v-if="errorMessage" class="admin-detail-error" role="alert">
      {{ errorMessage }}
    </p>

    <p v-if="images.length === 0" class="admin-detail-muted">
      {{ ADMIN_COPY.productImagesEmpty }}
    </p>

    <ol v-else class="admin-product-images">
      <li v-for="(image, imageIndex) in images" :key="image.id" class="admin-product-image">
        <img :src="image.url" :alt="`${product.name} image ${imageIndex + 1}`" loading="lazy">

        <div class="admin-product-image-detail">
          <strong>Image {{ imageIndex + 1 }}</strong>
          <small>{{ image.id }}</small>
        </div>

        <div class="admin-product-image-actions">
          <button
            type="button"
            class="admin-pagination-button is-number"
            :aria-label="ADMIN_COPY.productImageMoveUp"
            :disabled="imageIndex === 0 || movingImageID === image.id"
            @click="moveImage(image, imageIndex, -1)"
          >
            <StorefrontIcon :icon-node="ArrowUp" :size="14" :stroke-width="1.75" />
          </button>
          <button
            type="button"
            class="admin-pagination-button is-number"
            :aria-label="ADMIN_COPY.productImageMoveDown"
            :disabled="imageIndex === images.length - 1 || movingImageID === image.id"
            @click="moveImage(image, imageIndex, 1)"
          >
            <StorefrontIcon :icon-node="ArrowDown" :size="14" :stroke-width="1.75" />
          </button>
          <button
            type="button"
            class="admin-pagination-button is-number"
            :aria-label="ADMIN_COPY.productImageRemove"
            :disabled="deletingImageID === image.id"
            @click="removeImage(image)"
          >
            <StorefrontIcon :icon-node="Trash2" :size="14" :stroke-width="1.75" />
          </button>
        </div>
      </li>
    </ol>
  </section>
</template>
