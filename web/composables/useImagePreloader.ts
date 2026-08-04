const preloadedImageSources = new Set<string>()

const preloadImageSource = async (imageSource: string): Promise<void> => {
  if (preloadedImageSources.has(imageSource)) {
    return
  }

  preloadedImageSources.add(imageSource)

  const image = new Image()
  image.decoding = 'async'
  image.src = imageSource

  try {
    await image.decode()
  } catch {
    // Failed preloads should not block the visible image's normal browser path.
  }
}

export const useImagePreloader = () => {
  const preloadImage = (imageSource: string | null | undefined): void => {
    if (!import.meta.client || !imageSource) {
      return
    }

    void preloadImageSource(imageSource)
  }

  const preloadImages = (
    imageSources: readonly (string | null | undefined)[],
  ): void => {
    for (const imageSource of imageSources) {
      preloadImage(imageSource)
    }
  }

  return {
    preloadImage,
    preloadImages,
  }
}
