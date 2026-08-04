// Resolves which API origin to call.
//
// In production the browser and the SSR process must not use the same value:
// the browser calls the public path (/api) which Caddy proxies, while SSR runs
// inside the Docker network and can reach the API container directly — avoiding
// a pointless round trip out through the proxy and back in. In development both
// resolve to the same localhost origin.
export const useApiBaseUrl = (): string => {
  const runtimeConfig = useRuntimeConfig()

  if (import.meta.server) {
    return runtimeConfig.apiInternalBase
  }

  if (runtimeConfig.public.apiBase) {
    return runtimeConfig.public.apiBase
  }

  return `${window.location.protocol}//${window.location.hostname}:8080`
}
